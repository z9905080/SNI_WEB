<?php

require 'main.php';

require 'middleware/AuthToken.php';

// if ($_SERVER['REQUEST_METHOD'] !== 'POST') {
//     // http_response_code(404);
//     // exit;
// }

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

$db->where("id", $pageID);
if ($db->update("page_content", array(
    "page_name" => $pageName,
    "html_context" => $htmlContext,
    "page_group_id" => $pageGroupID,
))) {
    $resp = array(
        "data" => null,
        "message" => "編輯成功",
        "code" => "20003",
        "status" => "Y",
    );
    
    echo (json_encode($resp));
} else {

    $resp = array(
        "data" => null,
        "message" => "編輯失敗",
        "code" => "10003",
        "status" => "N",
    );

    echo (json_encode($resp));
}
