<?php

// ini_set("display_errors", "On"); // 顯示錯誤是否打開( On=開, Off=關 )

// error_reporting(E_ALL & ~E_NOTICE);

require 'main.php';

$db = new MysqliDb($dbCofig);

$marquee = $db->get('marquee');

if (count($marquee) > 0) {

    $resp = array(
        "data" => $marquee,
        "message" => "取得跑馬燈成功",
        "code" => Code::GetMarquee_Success,
        "status" => "Y",
    );
    
    echo (json_encode($resp));
} else {
    $resp = array(
        "message" => "目前沒有跑馬燈，或者取得失敗",
        "code" => Code::GetMarquee_Fail,
        "status" => "Y",
    );

    echo (json_encode($resp));
}
