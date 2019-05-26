<?php

//header('Access-Control-Allow-Origin: *');
if ($_SERVER["REQUEST_METHOD"] == "POST") {

    $DirNowDate = "picture/" . (new \DateTime())->format('Y_m_d');
    if (!is_dir($DirNowDate)) {
        mkdir($DirNowDate, 0777);
    }

    //private var
    $size = $_FILES['file']['size'];
    $type = $_FILES['file']['type'];

    $msg = array();
    $msg["error"] = false;
    $msg['content'] = '';
    $reNameFile = array("image/jpeg" => ".jpg",
        "image/png" => ".png",
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
            $FilePathAndName = $DirNowDate . '/' . $nowTime . $reNameFile[$type];

            if (move_uploaded_file($_FILES['file']['tmp_name'], $FilePathAndName)) {

            }

            $msg["url"] = $FilePathAndName;

            //echo '<a href="file/'.$_FILES['file']['name'].'">file/'.$_FILES['file']['name'].'</a>';//顯示檔案路徑
        } else {
            $msg["error"] = true;
            $msg['content'] = 'upload File type must be gif/jpeg/png';
        }
    }

    echo json_encode($msg);
}

function make_thumb($src, $dest, $desired_width = false, $desired_height = false)
{
    /*If no dimenstion for thumbnail given, return false */
    //if (!$desired_height&&!$desired_width) return false;
    $fparts = pathinfo($src);
    $ext = strtolower($fparts['extension']);
    /* if its not an image return false */
    if (!in_array($ext, array('gif', 'jpg', 'png', 'jpeg'))) {
        return false;
    }

    /* read the source image */
    if ($ext == 'gif') {
        $resource = imagecreatefromgif($src);
    } else if ($ext == 'png') {
        $resource = imagecreatefrompng($src);
    } else if ($ext == 'jpg' || $ext == 'jpeg') {
        $resource = imagecreatefromjpeg($src);
    }

    //var_dump($resource);
    $width = imagesx($resource);
    $height = imagesy($resource);

    if ($width > $height) {
        $desired_width = 256;
        $desired_height = false;
    } else {
        $desired_width = false;
        $desired_height = 256;
    }
    /* find the "desired height" or "desired width" of this thumbnail, relative to each other, if one of them is not given  */
    if (!$desired_height) {
        $desired_height = floor($height * ($desired_width / $width));
    }

    if (!$desired_width) {
        $desired_width = floor($width * ($desired_height / $height));
    }

    /* create a new, "virtual" image */
    $virtual_image = imagecreatetruecolor($desired_width, $desired_height);
    if ($width > $desired_width || $height > $desired_height) {
        /* copy source image at a resized size */
        imagecopyresized($virtual_image, $resource, 0, 0, 0, 0, $desired_width, $desired_height, $width, $height);
    } else {
        $virtual_image = $resource;
    }
    /* create the physical thumbnail image to its destination */
    /* Use correct function based on the desired image type from $dest thumbnail source */
    $fparts = pathinfo($dest);
    $ext = strtolower($fparts['extension']);
    /* if dest is not an image type, default to jpg */
    if (!in_array($ext, array('gif', 'jpg', 'png', 'jpeg'))) {
        $ext = 'jpg';
    }

    $dest = $fparts['dirname'] . '/' . $fparts['filename'] . '.' . $ext;

    if ($ext == 'gif') {
        imagegif($virtual_image, $dest);
    } else if ($ext == 'png') {
        imagepng($virtual_image, $dest, 1);
    } else if ($ext == 'jpg' || $ext == 'jpeg') {
        imagejpeg($virtual_image, $dest, 80);
    }

    return array(
        'width' => $width,
        'height' => $height,
        'new_width' => $desired_width,
        'new_height' => $desired_height,
        'dest' => $dest,
    );
}

/**
 * Optimizes PNG file with pngquant 1.8 or later (reduces file size of 24-bit/32-bit PNG images).
 *
 * You need to install pngquant 1.8 on the server (ancient version 1.0 won't work).
 * There's package for Debian/Ubuntu and RPM for other distributions on http://pngquant.org
 *
 * @param $path_to_png_file string - path to any PNG file, e.g. $_FILE['file']['tmp_name']
 * @param $max_quality int - conversion quality, useful values from 60 to 100 (smaller number = smaller file)
 * @return string - content of PNG file after conversion
 */
function compress_png($path_to_png_file, $max_quality = 80)
{
    if (!file_exists($path_to_png_file)) {
        throw new Exception("File does not exist: $path_to_png_file");
    }

    require_once "PNGQuant.php";

    $instance = new PNGQuant();

    // Set the path to the image to compress
    $exit_code = $instance->setBinaryPath("pngquant\\pngquant.exe")
        ->setImage($path_to_png_file)
    // Set the output filepath
        ->setOutputImage($path_to_png_file)
    // Overwrite output file if exists, otherwise pngquant will generate output ...
        ->overwriteExistingFile()
    // As the quality in pngquant isn't fixed (it uses a range)
    // set the quality between 50 and 80
        ->setQuality(70, 90)
    // Execute the command
        ->execute();

    // guarantee that quality won't be worse than that.
    $min_quality = 60;

    // '-' makes it use stdout, required to save to $compressed_png_content variable
    // '<' makes it read from the given file path
    // escapeshellarg() makes this safe to use with any path
    // $compressed_png_content = shell_exec("pngquant --quality=$min_quality-$max_quality - < " . escapeshellarg($path_to_png_file));
}
