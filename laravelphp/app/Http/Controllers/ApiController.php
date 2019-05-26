<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;

class ApiController extends Controller
{
    /**
     * Display a listing of the resource.
     *
     * @return \Illuminate\Http\Response
     */
    public function index()
    {
        //
    }

    /**
     * Show the form for creating a new resource.
     *
     * @return \Illuminate\Http\Response
     */
    public function create()
    {
        //
    }

    /**
     * Store a newly created resource in storage.
     *
     * @param  \Illuminate\Http\Request  $request
     * @return \Illuminate\Http\Response
     */
    public function store(Request $request)
    {
        //
    }

    /**
     * Display the specified resource.
     *
     * @param  int  $id
     * @return \Illuminate\Http\Response
     */
    public function show($id)
    {
        //
    }

    /**
     * Show the form for editing the specified resource.
     * @param Request $request
     * @return \Illuminate\Http\JsonResponse
     */
    public function edit(Request $request)
    {
        $param = $request->all();

        //驗證參數
        $validator = Validator::make($param, [
            'page_id' => 'required|string',
            'page_name' => 'required|string',
            'html_context' => 'required',
            'page_group_id' => 'required',
        ]);
        if ($validator->fails()) {
            throw new ValidationException($validator);
        }

        $result = DB::update('update page_content
                set page_name = ? , html_context = ? , page_group_id = ?
                where page_id = ?', [
            $param['page_name'],
            $param['html_context'],
            $param['page_group_id'],
            $param['page_id'],
        ]);

        if ($request < 1) {
            $resp = array(
                "data" => null,
                "message" => "編輯失敗",
                "code" => "10003",
                "status" => "N",
            );

            response()->json($resp);
        } else {

            $resp = array(
                "data" => null,
                "message" => "編輯成功",
                "code" => "20003",
                "status" => "Y",
            );

            response()->json($resp);
        }
    }

    /**
     * Update the specified resource in storage.
     *
     * @param  \Illuminate\Http\Request  $request
     * @param  int  $id
     * @return \Illuminate\Http\Response
     */
    public function update(Request $request, $id)
    {
        //
    }

    /**
     * Remove the specified resource from storage.
     *
     * @param  int  $id
     * @return \Illuminate\Http\Response
     */
    public function destroy($id)
    {
        //
    }
}
