<?php

// Cross-Origin Resource Sharing Header
header('Access-Control-Allow-Origin: *');
header('Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS');
header('Access-Control-Allow-Headers: X-Requested-With, Content-Type, Accept');

require_once 'MysqliDb.php';

// $db = new MysqliDb(array(
//     'host' => "localhost",
//     'username' => 'vhost24073',
//     'password' => 'tw28150945*sni',
//     'db' => 'vhost24073',
//     'port' => 3306,
//     'charset' => 'utf8'));

$db = new MysqliDb(array(
    'host' => "23.92.66.131",
    'username' => 'shouting_root',
    'password' => 'U6#j&gu8V2!0',
    'db' => 'shouting_vhost24073',
    'port' => 3306,
    'charset' => 'utf8'));
