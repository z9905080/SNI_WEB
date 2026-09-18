<?php

require 'main.php';

require 'middleware/AuthToken.php';

// if ($_SERVER['REQUEST_METHOD'] !== 'POST') {
//     // http_response_code(404);
//     // exit;
// }

$data = json_decode(file_get_contents('php://input'), true);

$pageID = $data['page_id'];

if ($pageID == "") {
    http_response_code(404);
    exit;
}

$db = new MysqliDb($dbCofig);

$db->where("id", $pageID);
if ($db->delete("page_content")) {
    $resp = (new APIResponse(null, "刪除內文成功", Code::DeleteContent_Success, "Y"))->GetAPIResponse();
    echo (json_encode($resp));
} else {
    $resp = (new APIResponse(null, "刪除內文失敗", Code::DeleteContent_Fail, "N"))->GetAPIResponse();
    echo (json_encode($resp));
}
