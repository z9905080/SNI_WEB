<?php
/**
 * APIResponse Class
 *
 * @category  Database Access
 * @package   APIResponse
 * @author    Jeffery Way <jeffrey@jeffrey-way.com>
 * @author    Josh Campbell <jcampbell@ajillion.com>
 * @author    Alexander V. Butenko <a.butenka@gmail.com>
 * @copyright Copyright (c) 2010-2017
 * @license   http://opensource.org/licenses/gpl-3.0.html GNU Public License
 * @link      http://github.com/joshcam/PHP-MySQLi-Database-Class
 * @version   2.9.2
 */

class APIResponse
{

    public $response = array();
    /**
     * @param string $code
     * @param string $status
     * @param string $message
     * @param object $data
     */
    public function __construct($data, $message, $code, $status = 'Y')
    {
        $this->response = array(
            "data" => $data,
            "message" => $message,
            "code" => $code,
            "status" => $status,
        );
    }

    public function getAPIResponse(){
        return $this->response;
    }

}
