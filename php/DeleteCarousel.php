<?php

require 'main.php';

require 'middleware/AuthToken.php';

$data = json_decode(file_get_contents('php://input'), true);

$carouselID = $data['carousel_id'];

if (empty($carouselID)) {
    http_response_code(404);
    exit;
}

$db = new MysqliDb($dbCofig);

$db->where("id", $carouselID);
if ($db->delete("carousel")) {
    $resp = array(
        "data" => null,
        "message" => "刪除輪播圖成功",
        "code" => Code::DeleteCarousel_Success,
        "status" => "Y",
    );

    echo (json_encode($resp));
} else {

    $resp = array(
        "data" => null,
        "message" => "刪除輪播圖失敗",
        "code" => Code::DeleteCarousel_Fail,
        "status" => "N",
    );

    echo (json_encode($resp));
}
