<?php

require 'main.php';

$sid = $_COOKIE['sid'];

$db = new MysqliDb($dbCofig);

$token_data = $db->where("token", $sid)->getOne('user_token');

$db->disconnect();

if ($token_data != null) {
    $nowTime = date("Y-m-d H:i:s");
    if (strtotime($token_data['expire_time']) < strtotime($nowTimeFormat)) {
        echo (json_encode($dataRlt));
    }
}


