package storage

import (
	"context"
	"errors"
	"io"
	"net/http"
	"slices"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/z9905080/SNI_WEB/backend/internal/content"
)

const s3Prefix = "picture/"

type S3Options struct {
	Endpoint      string
	Region        string
	Bucket        string
	AccessKey     string
	SecretKey     string
	PublicBaseURL string
}

type S3 struct {
	client     *s3.Client
	bucket     string
	publicBase string
}

func NewS3(o S3Options) *S3 {
	client := s3.New(s3.Options{
		Region:                     o.Region,
		BaseEndpoint:               aws.String(o.Endpoint),
		UsePathStyle:               true,
		Credentials:                aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(o.AccessKey, o.SecretKey, "")),
		RequestChecksumCalculation: aws.RequestChecksumCalculationWhenRequired,
		ResponseChecksumValidation: aws.ResponseChecksumValidationWhenRequired,
	})
	return &S3{client: client, bucket: o.Bucket, publicBase: strings.TrimRight(o.PublicBaseURL, "/")}
}

func httpStatus(err error) int {
	var re *awshttp.ResponseError
	if errors.As(err, &re) {
		return re.HTTPStatusCode()
	}
	return 0
}

func (s *S3) Put(ctx context.Context, name string, r io.ReadSeeker, size int64, contentType string) error {
	if !content.ValidImageName(name) {
		return ErrInvalidName
	}
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(s.bucket),
		Key:           aws.String(s3Prefix + name),
		Body:          r,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(contentType),
		CacheControl:  aws.String(immutableCache),
		IfNoneMatch:   aws.String("*"),
	})
	if httpStatus(err) == http.StatusPreconditionFailed {
		return ErrExists
	}
	return err
}

func (s *S3) Delete(ctx context.Context, name string) error {
	if !content.ValidImageName(name) {
		return ErrInvalidName
	}
	key := aws.String(s3Prefix + name)
	if _, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{Bucket: aws.String(s.bucket), Key: key}); err != nil {
		if httpStatus(err) == http.StatusNotFound {
			return ErrNotFound
		}
		return err
	}
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{Bucket: aws.String(s.bucket), Key: key})
	return err
}

func (s *S3) List(ctx context.Context) ([]Object, error) {
	out := []Object{}
	p := s3.NewListObjectsV2Paginator(s.client, &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(s3Prefix),
	})
	for p.HasMorePages() {
		page, err := p.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, o := range page.Contents {
			name := strings.TrimPrefix(aws.ToString(o.Key), s3Prefix)
			if !content.ValidImageName(name) {
				continue
			}
			out = append(out, Object{Name: name, Size: aws.ToInt64(o.Size), ModTime: aws.ToTime(o.LastModified)})
		}
	}
	slices.SortFunc(out, func(a, b Object) int { return strings.Compare(b.Name, a.Name) })
	return out, nil
}

func (s *S3) Serve(w http.ResponseWriter, r *http.Request, name string) {
	if !content.ValidImageName(name) {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Cache-Control", "public, max-age=3600")
	http.Redirect(w, r, s.publicBase+"/"+s3Prefix+name, http.StatusFound)
}
