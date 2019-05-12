<?php

require 'main.php';

/**
 * 產生隨機字串
 *
 * @return string
 */
function RandomString()
{
    $characters = '0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ';
    $randstring = '';
    for ($i = 0; $i < 10; $i++) {
        $randstring = $characters[rand(0, strlen($characters))];
    }
    return $randstring;
}

if ($_SERVER['REQUEST_METHOD'] !== 'POST') {
    // http_response_code(404);
    // exit;
}

$data = json_decode(file_get_contents('php://input'), true);

$account = $data['account'];
$pwd = $data['password'];

if ($account == "" || $pwd == "") {
    http_response_code(404);
    exit;
}

$db = new MysqliDb($dbCofig);

$userData = $db->where("account", $account)->getOne('user');

if ($userData != null) {

    if ($pwd == $userData["pwd"]) {
        $token = md5(RandomString());
        $db->where("user_id", $userData["id"])->delete("user_token");
        $now = new DateTime();
        $nowAfter1H = $now->add(new DateInterval("PT1H"));
        
        $tokenData = array(
            'user_id' => $userData["id"],
            'token' => $token,
            'expire_time' => $nowAfter1H->format("Y-m-d H:i:s"),
        );
        if ($db->insert("user_token", $tokenData)) {
            echo (json_encode($tokenData));
        } else {
            echo (json_encode(array('message' => "新增Token失敗，請重試")));
        }
    }

} else {

}
