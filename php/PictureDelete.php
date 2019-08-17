<?php

require 'main.php';

require 'middleware/AuthToken.php';


$data = json_decode(file_get_contents('php://input'), true);

$filePath = $data['file_path'];

if ($filePath == ""){
    http_response_code(404);
    exit;
}

if(file_exists($filePath)){
        //將檔案刪除
    if (unlink($filePath)){
        $resp = array(
            "data" => null,
            "message" => "刪除圖片成功",
            "code" => Code::DeletePicture_Success,
            "status" => "Y",
        );
        
        echo (json_encode($resp));
    }else{
        $resp = array(
            "data" => null,
            "message" => "刪除圖片失敗",
            "code" => Code::DeleteContent_Fail,
            "status" => "N",
        );
    }
    
}else{
    $resp = array(
        "data" => null,
        "message" => "刪除圖片失敗,檔案不存在",
        "code" => Code::DeleteContent_Fail,
        "status" => "N",
    );
}