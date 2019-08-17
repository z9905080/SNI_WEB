<?php

require 'main.php';

require 'middleware/AuthToken.php';

if ($_SERVER['REQUEST_METHOD'] !== 'POST') {
    // http_response_code(404);
    // exit;
}

$data = json_decode(file_get_contents('php://input'), true);

$pageName = $data['page_name'];
$htmlContext = $data['html_context'];
$pageGroupID = $data['page_group_id'];


if (empty($pageName) || empty($htmlContext) || empty($pageGroupID)) {

    $resp = array(
        "data" => null,
        "message" => "新增內文失敗,參數不得為空",
        "code" => Code::AddContent_Fail,
        "status" => "N",
    );
    echo (json_encode($resp));

    exit;
}

$db = new MysqliDb($dbCofig);

$pageContentData = array(
    "page_group_id" => $pageGroupID,
    "page_name" => $pageName,
    "html_context" => $htmlContext,
);

$id = $db->insert('page_content', $pageContentData);
if ($id) {
    $pageContentData["page_content_id"] = $id;

    $resp = array(
        "data" => $pageContentData,
        "message" => "新增內文成功",
        "code" => Code::AddContent_Succes,
        "status" => "Y",
    );

    echo (json_encode($resp));
} else {

    $resp = array(
        "data" => null,
        "message" => "新增內文失敗" . $db->getLastError(),
        "code" => Code::AddContent_Fail,
        "status" => "N",
    );
    echo (json_encode($resp));
}
