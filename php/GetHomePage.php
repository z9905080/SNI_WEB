<?php

require 'main.php';

$db = new MysqliDb($dbCofig);

$pageData = $db->where("page_group_id", 1)->getOne('page_content');

$carousel = $db->get('carousel');

$webSettingList = $db->get("web_config");

$webSettingMap = array();

foreach ($webSettingList as $key => $item) {
    $webSettingMap[$item['data_key']] = $item['data_value'];
}

$returnData = array(
    "page_data" => $pageData,
    "carousel" => $carousel,
    "web_title" => $webSettingMap['web_title'],
    "web_sub_title" => $webSettingMap['web_sub_title'],
);

echo (json_encode($returnData));
