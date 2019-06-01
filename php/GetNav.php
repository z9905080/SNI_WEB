<?php

require 'main.php';

$db = new MysqliDb($dbCofig);

$data = $db->join("page_content pc", "pg.id=pc.page_group_id", "LEFT")
    ->get('page_group pg', null, "pc.id as page_content_id, page_group_id, group_name, page_name");

$dataRlt = array();

foreach ($data as $key => $originItem) {
    $isExist = false;
    foreach ($dataRlt as $key => $newItem) {
        if ($newItem["page_group_id"] == $originItem["page_group_id"]) {
            array_push($dataRlt[$key]["page_content"], array(
                "page_content_id" => $originItem["page_content_id"],
                "page_name" => $originItem["page_name"],
            ));
            $isExist = true;
        }
    }

    if (!$isExist) {
        array_push($dataRlt, array(
            "page_group_id" => $originItem["page_group_id"],
            "group_name" => $originItem["group_name"],
            "page_content" => array(
                array(
                    "page_content_id" => $originItem["page_content_id"],
                    "page_name" => $originItem["page_name"],
                ),
            ),
        ));
    }

}

echo (json_encode($dataRlt));
