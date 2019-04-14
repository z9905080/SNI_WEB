<?php

require 'setting.php';

require 'main.php';

$db = new MysqliDb($dbCofig);

$pageData = $db->where("page_group_id", 1)->getOne('page_content');

$carousel = $db->get('carousel');

$returnData = array(
    "page_data" => $pageData,
    "carousel" => $carousel,
);

echo (json_encode($returnData));
