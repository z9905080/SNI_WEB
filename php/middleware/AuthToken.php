<?php

require 'main.php';

$sid = $_COOKIE['sid'];

echo $sid;

$db = new MysqliDb($dbCofig);

$token_data = $db->where("token", $sid)->getOne('user_token');

$db->disconnect();

if ($token_data != null) {
    $nowTime = date("Y-m-d H:i:s");
    if (strtotime($token_data['expire_time']) < strtotime($nowTimeFormat)) {
    
        $resp = array(
            "data" => null,
            "message" => "登入令牌過期，請重新登入",
            "code" => "10001",
            "status" => "N",
        );
        
        echo (json_encode($resp));
        exit;
    }
} else {

    $resp = array(
        "data" => null,
        "message" => "驗證錯誤，請重新登入",
        "code" => "10002",
        "status" => "N",
    );

    echo (json_encode($resp));
    exit;
}
