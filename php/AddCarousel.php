<?php

require 'main.php';

require 'middleware/AuthToken.php';

$protocol = ((!empty($_SERVER['HTTPS']) && $_SERVER['HTTPS'] != 'off') || $_SERVER['SERVER_PORT'] == 443) ? "https://" : "http://";
$http = $_SERVER['HTTP_HOST'];

echo $protocol.$http."/php";
exit;


$data = json_decode(file_get_contents('php://input'), true);

$filePath = $data['file_path'];
$link_url = $data['link_url'];

if (empty($filePath) || empty($link_url)) {

    $resp = array(
        "data" => null,
        "message" => "新增輪播圖失敗,參數不得為空",
        "code" => Code::AddCarousel_Fail,
        "status" => "N",
    );
    echo (json_encode($resp));

    exit;
}

$db = new MysqliDb($dbCofig);



$carouselData = array(
    "image" => "",
    "html_context" => $htmlContext,
);

$id = $db->insert('page_content', $pageContentData);
if ($id) {
    $pageContentData["page_content_id"] = $id;

    $resp = array(
        "data" => $pageContentData,
        "message" => "新增內文成功",
        "code" => Code::AddContent_Succes,
        "status" => "Y",
    );

    echo (json_encode($resp));
} else {

    $resp = array(
        "data" => null,
        "message" => "新增內文失敗" . $db->getLastError(),
        "code" => Code::AddContent_Fail,
        "status" => "N",
    );
    echo (json_encode($resp));
}


if (file_exists($filePath)) {
    //將檔案刪除
    if (unlink($filePath)) {
        $resp = array(
            "data" => null,
            "message" => "刪除圖片成功",
            "code" => Code::DeletePicture_Success,
            "status" => "Y",
        );

        echo (json_encode($resp));
    } else {
        $resp = array(
            "data" => null,
            "message" => "刪除圖片失敗",
            "code" => Code::DeleteContent_Fail,
            "status" => "N",
        );
    }
} else {
    $resp = array(
        "data" => null,
        "message" => "刪除圖片失敗,檔案不存在",
        "code" => Code::DeleteContent_Fail,
        "status" => "N",
    );
}
