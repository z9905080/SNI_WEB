<?php

require 'main.php';

require 'middleware/AuthToken.php';

if ($_SERVER['REQUEST_METHOD'] !== 'POST') {
    // http_response_code(404);
    // exit;
}

$data = json_decode(file_get_contents('php://input'), true);

$pageID = $data['page_id'];
$pageName = $data['page_name'];
$htmlContext = $data['html_context'];
$pageGroupID = $data['page_group_id'];


if ($pageID == "") {
    http_response_code(404);
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
        "code" => "20004",
        "status" => "Y",
    );
    
    echo (json_encode($resp));

} else {

    $resp = array(
        "data" => $null,
        "message" => "新增內文失敗". $db->getLastError(),
        "code" => "10004",
        "status" => "Y",
    );
    echo (json_encode($resp));
}
