<?php

require 'main.php';

require 'middleware/AuthToken.php';


$data = json_decode(file_get_contents('php://input'), true);

$filePath = $data['file_path'];

if ($filePath == "") {
    http_response_code(404);
    exit;
}

$resp = array();

if (file_exists($filePath)) {
    chmod($filePath,0777);
    //將檔案刪除
    if (@unlink($filePath)) {
        $resp = (new APIResponse(null, "刪除圖片成功", Code::DeletePicture_Success, "Y"))->GetAPIResponse();
    } else {
        $resp = (new APIResponse(null, "刪除圖片失敗", Code::DeleteContent_Fail, "N"))->GetAPIResponse();
    }
} else {
    $resp = (new APIResponse(null, "刪除圖片失敗,檔案不存在", Code::DeleteContent_Fail, "N"))->GetAPIResponse();
}

echo (json_encode($resp));
