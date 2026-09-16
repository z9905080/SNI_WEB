<?php

ini_set("display_errors", "Off"); // 顯示錯誤是否打開( On=開, Off=關 )

error_reporting(E_ERROR);

require_once 'Code.php';

require_once 'APIResponse.php';

require_once 'MysqliDb.php';

// $db = new MysqliDb(array(
//     'host' => "localhost",
//     'username' => 'vhost24073',
//     'password' => 'tw28150945*sni',
//     'db' => 'vhost24073',
//     'port' => 3306,
//     'charset' => 'utf8'));


// Cross-Origin Resource Sharing Header
header('Access-Control-Allow-Origin: *');
header('Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS');
header('Access-Control-Allow-Headers: *');
header('Cache-Control: max-age=3600');

date_default_timezone_set("Asia/Taipei");

$dbCofig = array(
    'host' => "localhost",
    'username' => 'vhost24073',
    'password' => 'tw28150945*sni',
    'db' => 'vhost24073',
    'port' => 3306,
    'charset' => 'utf8');