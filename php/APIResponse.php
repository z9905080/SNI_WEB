<?php
/**
 * APIResponse Class
 *
 * @category  API Access
 * @package   APIResponse
 */

class APIResponse
{

    /**
     * Static instance of self
     *
     * @var APIResponse
     */
    protected static $_instance;

    public $response = null;

    /**
     * @param object $data
     * @param string $message
     * @param string $code
     * @param string $status
     */
    public function __construct($data = null, $message = null, $code = null, $status = 'Y')
    {

        $this->response = array(
            "data" => $data,
            "message" => $message,
            "code" => $code,
            "status" => $status,
        );

        self::$_instance = $this;
    }

    public function GetAPIResponse()
    {
        return self::$_instance->response;
    }

}


// END class