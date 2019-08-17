<?php

ini_set("display_errors", "On"); // 顯示錯誤是否打開( On=開, Off=關 )

error_reporting(E_ALL & ~E_NOTICE);

use Code;

require 'main.php';


$db = new MysqliDb($dbCofig);

$carousel = $db->get('carousel');

if (count($carousel) > 0) {

    $resp = array(
        "data" => $carousel,
        "message" => "取得輪播圖成功",
        "code" => Code::GetCarousel_Success,
        "status" => "Y",
    );

    echo "123:" . count($carousel);
    exit;

    echo (json_encode($resp));
} else {
    $resp = array(
        "message" => "目前沒有輪播圖，或者取得失敗",
        "code" => Code::GetCarousel_Fail,
        "status" => "Y",
    );

    echo (json_encode($resp));
}
