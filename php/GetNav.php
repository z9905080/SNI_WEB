<?php

require 'main.php';

$db = new MysqliDb(array(
    'host' => "localhost",
    'username' => 'vhost24073',
    'password' => 'tw28150945*sni',
    'db' => 'vhost24073',
    'port' => 3306,
    'charset' => 'utf8'));

$data = $db->join("page_group pg", "pg.id=pc.page_group_id", "INNER")->get('page_content pc');

echo (json_encode($data));
