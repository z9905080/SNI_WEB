<?php

require 'main.php';

require 'middleware/AuthToken.php';

if ($_SERVER['REQUEST_METHOD'] !== 'POST') {
    // http_response_code(404);
    // exit;
}

$data = json_decode(file_get_contents('php://input'), true);

$groupName = $data['page_group_name'];

if (empty($groupName)) {
    http_response_code(404);
    exit;
}

$db = new MysqliDb($dbCofig);

$pageGroupData = array(
    "group_name" => $groupName,
);

$id = $db->insert('page_group', $pageGroupData);
if ($id) {
    $pageGroupData["page_group_id"] = $id;

    $resp = array(
        "data" => $pageGroupData,
        "message" => "新增頁籤成功",
        "code" => Code::AddNav_Success,
        "status" => "Y",
    );

    echo (json_encode($resp));
} else {

    $resp = array(
        "data" => null,
        "message" => "新增頁籤失敗" . $db->getLastError(),
        "code" => Code::AddNav_Fail,
        "status" => "N",
    );
    echo (json_encode($resp));
}
