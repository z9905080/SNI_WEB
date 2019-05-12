<?php

require 'main.php';

$sid = $_COOKIE['sid'];

$db = new MysqliDb($dbCofig);

$token_data = $db->where("token", $sid)->getOne('user_token');

$db->disconnect();

if ($token_data != null) {
    $nowTime = date("Y-m-d H:i:s");
    if (strtotime($token_data['expire_time']) < strtotime($nowTimeFormat)) {
        $respInst = new APIResponse(
            null,
            "登入令牌過期，請重新登入",
            "10001",
            "N");
        $resp = $respInst->getAPIResponse();
        echo (json_encode($resp));
        exit;
    }
} else {
    $respInst = new APIResponse(
        null,
        "驗證錯誤，請重新登入",
        "10002",
        "N");
    $resp = $respInst->getAPIResponse();
    echo (json_encode($resp));
    exit;
}
