<?php

require 'main.php';

require 'middleware/AuthToken.php';

$data = json_decode(file_get_contents('php://input'), true);

$pageGroupID = $data['page_group_id'];
$pageSortStr = $data['page_sort'];

if (empty($pageGroupID) || empty($pageSortStr)) {
    http_response_code(404);
    exit;
}


$pageSortList = json_decode($pageSortStr);

foreach ($pageSortList as $index => $sortID) {
    if (!is_int($sortID)) {
        $resp = array(
            "data" => null,
            "message" => "格式錯誤",
            "code" => Code::EditContentSort_Fail,
            "status" => "N",
        );
        echo (json_encode($resp));
        exit;
    }
}

$db = new MysqliDb($dbCofig);

$db->where("id", $pageGroupID);
if ($db->update("page_group", array(
    "page_sort" => $pageSortStr,
))) {
    $resp = array(
        "data" => null,
        "message" => "編輯內頁順序成功",
        "code" => Code::EditContentSort_Success,
        "status" => "Y",
    );

    echo (json_encode($resp));
} else {

    $resp = array(
        "data" => null,
        "message" => "編輯內頁順序失敗",
        "code" => Code::EditContentSort_Fail,
        "status" => "N",
    );

    echo (json_encode($resp));
}
