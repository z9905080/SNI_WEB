package storage

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/testcontainers/testcontainers-go/modules/minio"
)

func TestS3(t *testing.T) {
	if testing.Short() {
		t.Skip("需要 Docker")
	}
	ctx := context.Background()
	ctr, err := minio.Run(ctx, "quay.io/minio/minio:RELEASE.2025-04-22T22-12-26Z",
		minio.WithUsername("minioadmin"), minio.WithPassword("minioadmin"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = ctr.Terminate(context.Background()) })
	hostPort, err := ctr.ConnectionString(ctx)
	if err != nil {
		t.Fatal(err)
	}

	st := NewS3(S3Options{
		Endpoint:      "http://" + hostPort,
		Region:        "us-east-1",
		Bucket:        "sniweb",
		AccessKey:     "minioadmin",
		SecretKey:     "minioadmin",
		PublicBaseURL: "https://cdn.example.com/sniweb",
	})
	if _, err := st.client.CreateBucket(ctx, &s3.CreateBucketInput{Bucket: aws.String("sniweb")}); err != nil {
		t.Fatal(err)
	}

	put := func(name, body string) error {
		return st.Put(ctx, name, strings.NewReader(body), int64(len(body)), "image/png")
	}
	if err := put("2026-01-01_00-00-00.png", "aaa"); err != nil {
		t.Fatal(err)
	}
	if err := put("2026-01-02_00-00-00.png", "bb"); err != nil {
		t.Fatal(err)
	}
	if err := put("2026-01-01_00-00-00.png", "zzz"); !errors.Is(err, ErrExists) {
		t.Fatalf("重複應回 ErrExists：%v", err)
	}
	if err := put("../x.png", "x"); !errors.Is(err, ErrInvalidName) {
		t.Fatalf("got %v", err)
	}

	head, err := st.client.HeadObject(ctx, &s3.HeadObjectInput{Bucket: aws.String("sniweb"), Key: aws.String("picture/2026-01-01_00-00-00.png")})
	if err != nil {
		t.Fatal(err)
	}
	if aws.ToString(head.ContentType) != "image/png" || aws.ToString(head.CacheControl) != immutableCache {
		t.Fatalf("物件 metadata 錯誤：%v %v", aws.ToString(head.ContentType), aws.ToString(head.CacheControl))
	}

	objs, err := st.List(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(objs) != 2 || objs[0].Name != "2026-01-02_00-00-00.png" || objs[1].Size != 3 {
		t.Fatalf("List = %+v", objs)
	}

	w := httptest.NewRecorder()
	st.Serve(w, httptest.NewRequest(http.MethodGet, "/", nil), "2026-01-01_00-00-00.png")
	if w.Code != http.StatusFound || w.Header().Get("Location") != "https://cdn.example.com/sniweb/picture/2026-01-01_00-00-00.png" {
		t.Fatalf("Serve: %d %v", w.Code, w.Header())
	}
	w = httptest.NewRecorder()
	st.Serve(w, httptest.NewRequest(http.MethodGet, "/", nil), "..")
	if w.Code != 404 {
		t.Fatalf("不合法檔名應 404：%d", w.Code)
	}

	if err := st.Delete(ctx, "2026-01-01_00-00-00.png"); err != nil {
		t.Fatal(err)
	}
	if err := st.Delete(ctx, "2026-01-01_00-00-00.png"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v", err)
	}
}
