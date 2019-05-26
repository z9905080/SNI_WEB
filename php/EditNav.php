<?php

require 'main.php';

require 'middleware/AuthToken.php';

$data = json_decode(file_get_contents('php://input'), true);

$pageGroupID = $data['page_group_id'];
$groupName = $data['page_group_name'];

if (empty($pageGroupID)) {
    http_response_code(404);
    exit;
}

$db = new MysqliDb($dbCofig);

$db->where("id", $pageGroupID);
if ($db->update("page_group", array(
    "group_name" => $groupName,
))) {
    $resp = array(
        "data" => null,
        "message" => "編輯頁籤成功",
        "code" => "20005",
        "status" => "Y",
    );

    echo (json_encode($resp));
} else {

    $resp = array(
        "data" => null,
        "message" => "編輯頁籤失敗",
        "code" => "10005",
        "status" => "N",
    );

    echo (json_encode($resp));
}
