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
    echo (json_encode($resp));
} else {
    $resp = (new APIResponse(null, "刪除跑馬燈失敗", Code::DeleteMarquee_Fail, "N"))->GetAPIResponse();
    echo (json_encode($resp));
}
