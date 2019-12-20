<?php

ini_set("display_errors", "On"); // 顯示錯誤是否打開( On=開, Off=關 )

error_reporting(E_ALL & ~E_NOTICE);

require 'main.php';

require 'middleware/AuthToken.php';

$data = json_decode(file_get_contents('php://input'), true);

$pageGroupSortStr = $data['page_group_sort'];

if (empty($pageGroupSortStr)) {
    http_response_code(404);
    exit;
}

$pageGroupSortList = json_decode($pageGroupSortStr);

foreach ($pageGroupSortList as $index => $pageGroupSortID) {
    if (!is_int($pageGroupSortID)) {
        $resp = array(
            "data" => null,
            "message" => "格式錯誤",
            "code" => Code::EditNavSort_Fail,
            "status" => "N",
        );
        echo (json_encode($resp));
        exit;
    }
}

$db = new MysqliDb($dbCofig);

$db->where("data_key", "page_group_sort");
if ($db->update("web_config", array(
    "data_value" => $pageGroupSortStr,
))) {
    $resp = array(
        "data" => null,
        "message" => "編輯頁籤順序成功",
        "code" => Code::EditNavSort_Success,
        "status" => "Y",
    );

    echo (json_encode($resp));
} else {

    $resp = array(
        "data" => null,
        "message" => "編輯頁籤順序失敗",
        "code" => Code::EditNavSort_Fail,
        "status" => "N",
    );

    echo (json_encode($resp));
}
