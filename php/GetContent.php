<?php

require 'main.php';

$pageID = $_REQUEST["page_id"];

if ($pageID == "") {
    http_response_code(404);
    exit;
}


$db = new MysqliDb($dbCofig);

$data = $db->where("id", $pageID)->getOne('page_content');

echo (json_encode($data));
