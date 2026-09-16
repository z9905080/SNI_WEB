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
    $resp = (new APIResponse(null, "刪除輪播圖成功", Code::DeleteCarousel_Success, "Y"))->GetAPIResponse();
    echo (json_encode($resp));
} else {
    $resp = (new APIResponse(null, "刪除輪播圖失敗", Code::DeleteCarousel_Fail, "N"))->GetAPIResponse();
    echo (json_encode($resp));
}
