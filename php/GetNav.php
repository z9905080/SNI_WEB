<?php

require 'main.php';

$db = new MysqliDb($dbCofig);

$data = $db->join("page_group pg", "pg.id=pc.page_group_id", "INNER")
    ->get('page_content pc', null, "pc.id as page_content_id, page_group_id, group_name, page_name");

$dataRlt = array();

foreach ($data as $key => $originItem) {
    foreach ($dataRlt as $key => $newItem) {
        if ($newItem["page_group_id"] != $originItem["page_group_id"]) {
            array_push($dataRlt, array(
                "page_group_id" => $originItem["page_group_id"],
                "group_name" => $originItem["group_name"],
                "page_content" => array(
                    "page_content_id" => $originItem["page_content_id"],
                    "page_name" => $originItem["page_name"],
                ),
            ));
        } else {
            array_push($newItem["page_content"],array(
                "page_content_id" => $originItem["page_content_id"],
                "page_name" => $originItem["page_name"],
            ));
        }
    }
}

echo (json_encode($dataRlt));
