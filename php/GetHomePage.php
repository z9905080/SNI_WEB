<?php

require 'main.php';

$db = new MysqliDb(array(
    'host' => "localhost",
    'username' => 'vhost24073',
    'password' => 'tw28150945*sni',
    'db' => 'vhost24073',
    'port' => 3306,
    'charset' => 'utf8'));

$data = $db->where("page_group_id", 1)->getOne('page_content');

echo (json_encode($data));
