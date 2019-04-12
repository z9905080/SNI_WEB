<?php


// if ( $_SERVER['REQUEST_METHOD'] !== 'POST' )
// {
//     http_response_code(404);
//     exit;
// }

$rws_post = $GLOBALS['HTTP_RAW_POST_DATA'];

$mypost = json_decode($rws_post);

$pageID = (string)$mypost->page_id;

if ($pageID === "") {
    http_response_code(404);
    exit;
}


require 'main.php';

$db = new MysqliDb($dbCofig);

$data = $db->where("id", $pageID)->getOne('page_content');

echo (json_encode($data));
