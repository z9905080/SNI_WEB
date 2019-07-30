<?php

require 'main.php';

require 'middleware/AuthToken.php';


$dir = "picture/";
$a=0;//初始化檔案數

$currentDir = getcwd();
get_dir_list($a,$dir,$currentDir);
//函式功能 列出該路徑下所有的檔案包含子目錄
function get_dir_list(&$a,$dir,$replaceDir){//&$a 檔案加總變數 傳址參數,$dir 資料夾路徑
    if(is_dir($dir)){//如果是資料夾才執行
        $dh = opendir($dir);//打開資料夾
        chdir ($dir);//指定目錄
        while (($file = readdir($dh)) !== false) {//列出該目錄的所有檔案
            if (is_dir($file) && basename($file)!='.' && basename($file)!='..'){//若是資料夾 且非 . .. 就在呼叫自已一次 
                get_dir_list($a,$file,$replaceDir);
            }else if(basename($file) != "." && basename($file) != ".."){//若非 . .. 就列出檔案
                echo getcwd()."\\$file <BR/>";//輸出 完整檔案路徑檔名
                $a+=1;//檔案總數加1
            }
        }
        chdir("../");//回到上一層目錄
        closedir($dh);//關閉資料夾
    }
}