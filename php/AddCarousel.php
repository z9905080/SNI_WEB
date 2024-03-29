<?php

require 'main.php';

require 'middleware/AuthToken.php';

$data = json_decode(file_get_contents('php://input'), true);

$imageUrl = $data['image_url'];
$linkUrl = $data['link_url'];

if (empty($imageUrl) || empty($linkUrl)) {

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


// $protocol = ((!empty($_SERVER['HTTPS']) && $_SERVER['HTTPS'] != 'off') || $_SERVER['SERVER_PORT'] == 443) ? "https://" : "http://";
// $http = $_SERVER['HTTP_HOST'];

$carouselData = array(
    "image" => $imageUrl,
    "url" => $linkUrl,
);



$id = $db->insert('carousel', $carouselData);
if ($id) {
    $carouselData["id"] = $id;

    $resp = array(
        "data" => $carouselData,
        "message" => "新增輪播圖成功",
        "code" => Code::AddCarousel_Success,
        "status" => "Y",
    );

    echo (json_encode($resp));
} else {

    $resp = array(
        "data" => null,
        "message" => "新增輪播圖失敗" . $db->getLastError(),
        "code" => Code::AddCarousel_Fail,
        "status" => "N",
    );
    echo (json_encode($resp));
}
