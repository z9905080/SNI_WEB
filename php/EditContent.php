<?php

require 'main.php';

require 'middleware/AuthToken.php';

$db = new MysqliDb($dbCofig);

$data = $db->join("page_group pg", "pg.id=pc.page_group_id", "INNER")
    ->get('page_content pc', null, "pc.id as page_content_id, page_group_id, group_name, page_name");

$dataRlt = array();

echo (json_encode($dataRlt));
