<?php

require 'main.php';

require 'middleware/AuthToken.php';

if ($_SERVER["REQUEST_METHOD"] == "POST") {

    $Dir = "picture";
    if (!is_dir($Dir)) {
        mkdir($Dir, 0777);
    }

    //private var
    $size = $_FILES['file']['size'];
    $type = $_FILES['file']['type'];

    $msg = array();
    $msg["error"] = false;
    $msg['content'] = '';
    $reNameFile = array(
        "image/jpeg" => ".jpg",
        "image/png" => ".png",
        "image/gif" => ".gif",
    );
    if ($_FILES['file']['error'] > 0) {

        $msg["error"] = true;
        $msg['content'] = 'file upload Error！';
        echo json_encode($msg);
        exit(""); //如果出現錯誤則停止程式
    }
    if ($size > 2048000) {
        $msg["error"] = true;
        $msg['content'] = 'upload file size than 1M';
    } else {
        if ($type == "image/jpeg" || $type == "image/png" || $type == "image/gif") {

            $nowTime = (new \DateTime())->format('Y-m-d_H-i-s');
            $FilePathAndName = $Dir . '/' . $nowTime . $reNameFile[$type];

            if (file_exists($FilePathAndName)) {
                $rand = substr(md5(rand()), 0, 6);
                $FilePathAndName = $FilePathAndName . $rand;
            }

            if (move_uploaded_file($_FILES['file']['tmp_name'], $FilePathAndName)) { }

            $msg["url"] = $FilePathAndName;

            //echo '<a href="file/'.$_FILES['file']['name'].'">file/'.$_FILES['file']['name'].'</a>';//顯示檔案路徑
        } else {
            $msg["error"] = true;
            $msg['content'] = 'upload File type must be gif/jpeg/png';
        }
    }

    echo json_encode($msg);
}
