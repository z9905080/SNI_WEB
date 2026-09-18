<?php

require 'main.php';

require 'middleware/AuthToken.php';

$data = json_decode(file_get_contents('php://input'), true);

$pageGroupID = $data['page_group_id'];

if (empty($pageGroupID)) {
    http_response_code(404);
    exit;
}

$db = new MysqliDb($dbCofig);

$db->where("id", $pageGroupID);
if ($db->delete("page_group")) {
    $resp = array(
        "data" => null,
        "message" => "刪除頁籤成功",
        "code" => Code::DeleteNav_Success,
        "status" => "Y",
    );

    echo (json_encode($resp));
} else {

    $resp = array(
        "data" => null,
        "message" => "刪除頁籤失敗",
        "code" => Code::DeleteNav_Fail,
        "status" => "N",
    );

    echo (json_encode($resp));
}
