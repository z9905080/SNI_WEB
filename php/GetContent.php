<?php

require 'setting.php';

// if ($_SERVER['REQUEST_METHOD'] !== 'POST') {
//     // http_response_code(404);
//     // exit;
// }

$data = json_decode(file_get_contents('php://input'), true);

$pageID = $data['page_id'];

//$pageID = $_REQUEST["page_id"];

if ($pageID == "") {
    http_response_code(404);
    exit;
}

require 'main.php';

$db = new MysqliDb($dbCofig);

$data = $db->where("id", $pageID)->getOne('page_content');

echo (json_encode($data));
