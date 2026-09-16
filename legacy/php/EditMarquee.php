<?php

require 'main.php';

require 'middleware/AuthToken.php';

// if ($_SERVER['REQUEST_METHOD'] !== 'POST') {
//     // http_response_code(404);
//     // exit;
// }

$data = json_decode(file_get_contents('php://input'), true);

$marqueeID = $data['marquee_id'];
$text = $data['text'];
$color = $data['color'];

if ($marqueeID == "") {
    http_response_code(404);
    exit;
}

$db = new MysqliDb($dbCofig);

$db->where("id", $marqueeID);
if ($db->update("marquee", array(
    "text" => $text,
    "color" => $color,
))) {
    $resp = array(
        "data" => null,
        "message" => "編輯跑馬燈成功",
        "code" => Code::EditMarquee_Success,
        "status" => "Y",
    );

    echo (json_encode($resp));
} else {

    $resp = array(
        "data" => null,
        "message" => "編輯跑馬燈失敗",
        "code" => Code::EditMarquee_Fail,
        "status" => "N",
    );

    echo (json_encode($resp));
}
