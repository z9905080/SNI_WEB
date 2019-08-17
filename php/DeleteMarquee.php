<?php

require 'main.php';

require 'middleware/AuthToken.php';

$data = json_decode(file_get_contents('php://input'), true);

$marqueeID = $data['marquee_id'];

if (empty($marqueeID)) {
    http_response_code(404);
    exit;
}

$db = new MysqliDb($dbCofig);

$db->where("id", $marqueeID);
if ($db->delete("marquee")) {

    $resp = (new APIResponse(null, "刪除跑馬燈成功", Code::DeleteMarquee_Success, "Y"))->GetAPIResponse();
    // $resp = array(
    //     "data" => null,
    //     "message" => "刪除跑馬燈成功",
    //     "code" => Code::DeleteMarquee_Success,
    //     "status" => "Y",
    // );

    echo (json_encode($resp));
} else {

    $resp = array(
        "data" => null,
        "message" => "刪除跑馬燈失敗",
        "code" => Code::DeleteMarquee_Fail,
        "status" => "N",
    );

    echo (json_encode($resp));
}
