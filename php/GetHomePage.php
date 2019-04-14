<?php

require 'setting.php';

require 'main.php';

$db = new MysqliDb($dbCofig);

$data = $db->where("page_group_id", 1)->getOne('page_content');

echo (json_encode($data));
