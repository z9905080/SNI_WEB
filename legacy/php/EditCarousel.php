<?php

require 'main.php';

require 'middleware/AuthToken.php';

// if ($_SERVER['REQUEST_METHOD'] !== 'POST') {
//     // http_response_code(404);
//     // exit;
// }

$data = json_decode(file_get_contents('php://input'), true);

$carouselID = $data['carousel_id'];
$imageUrl = $data['image_url'];
$linkUrl = $data['link_url'];

if ($carouselID == "") {
    http_response_code(404);
    exit;
}

$db = new MysqliDb($dbCofig);

$carouselData = array(
    "image" => $imageUrl,
    "url" => $linkUrl,
);


$db->where("id", $carouselID);
if ($db->update("carousel", $carouselData)) {
    $resp = array(
        "data" => null,
        "message" => "編輯輪播圖成功",
        "code" => Code::EditCarousel_Success,
        "status" => "Y",
    );

    echo (json_encode($resp));
} else {

    $resp = array(
        "data" => null,
        "message" => "編輯輪播圖失敗",
        "code" => Code::EditCarousel_Fail,
        "status" => "N",
    );

    echo (json_encode($resp));
}
