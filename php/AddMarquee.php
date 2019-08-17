<?php

require 'main.php';

require 'middleware/AuthToken.php';

$data = json_decode(file_get_contents('php://input'), true);

$text = $data['text'];
$color = $data['color'];

if (empty($text) || empty($color)) {

    $resp = array(
        "data" => null,
        "message" => "新增跑馬燈失敗,參數不得為空",
        "code" => Code::AddMarquee_Fail,
        "status" => "N",
    );
    echo (json_encode($resp));

    exit;
}

$db = new MysqliDb($dbCofig);

$marqueeData = array(
    "text" => $text,
    "color" => $color,
);

$id = $db->insert('marquee', $marqueeData);
if ($id) {
    $marqueeData["id"] = $id;

    $resp = array(
        "data" => $marqueeData,
        "message" => "新增跑馬燈成功",
        "code" => Code::AddMarquee_Success,
        "status" => "Y",
    );

    echo (json_encode($resp));
} else {

    $resp = array(
        "data" => null,
        "message" => "新增跑馬燈失敗" . $db->getLastError(),
        "code" => Code::AddMarquee_Fail,
        "status" => "N",
    );
    echo (json_encode($resp));
}
