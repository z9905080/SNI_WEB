<?php

require 'main.php';

//require 'middleware/AuthToken.php';

// if ($_SERVER['REQUEST_METHOD'] !== 'POST') {
//     // http_response_code(404);
//     // exit;
// }

$data = json_decode(file_get_contents('php://input'), true);

$pageID = $data['page_id'];
$pageName = $data['page_name'];
$htmlContext = $data['html_context'];
$pageGroupID = $data['page_group_id'];


$db = new MysqliDb($dbCofig);

$db->where("id", $pageID);
if ($db->update("page_content", array(
    "page_name" => $pageName,
    "html_context" => $htmlContext,
    "page_group_id" => $pageGroupID,
))) {
    echo (1);
    exit;
    $respInst = new APIResponse(
        null,
        "編輯成功",
        "20001",
        "Y");
    $resp = $respInst->getAPIResponse();
    echo (json_encode($resp));
}else {
    echo (2);
    exit;
    $respInst = new APIResponse(
        null,
        "編輯成功1",
        "20001",
        "Y");
    $resp = $respInst->getAPIResponse();
    echo (json_encode($resp));
}


