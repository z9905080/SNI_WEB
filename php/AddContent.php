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

$db = new MysqliDb($dbCofig);

$pageContentData = array(
    "page_group_id" => $pageGroupID,
    "page_name" => $pageName,
    "html_context" => $htmlContext,
);

$id = $db->insert('page_content', $pageContentData);
if ($id) {
    $pageContentData["page_content_id"] = $id;

    $respInst = new APIResponse(
        $pageContentData,
        "新增成功",
        "20002",
        "Y");
    $resp = $respInst->getAPIResponse();
    echo (json_encode($resp));

} else {
    $respInst = new APIResponse(
        null,
        "新增失敗" . $db->getLastError(),
        "10002",
        "N");
    $resp = $respInst->getAPIResponse();
    echo (json_encode($resp));
}
