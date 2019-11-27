<?php

require 'main.php';

$db = new MysqliDb($dbCofig);

// $data = $db->join("page_group pg", "pg.id=pc.page_group_id", "LEFT")
//     ->get('page_content pc', null, "pc.id as page_content_id, page_group_id, group_name, page_name");

$webSettingList = $db->get("web_config");

$webSettingMap = array();

foreach ($webSettingList as $key => $item) {
    $webSettingMap[$item['data_key']] = $item['data_value'];
}

$groupSort = json_decode($webSettingMap['page_group_sort']);

$pageGroupList = $db->get("page_group");

$pageContentList = $db->get("page_content");

$dataRlt = array();

foreach ($pageGroupList as $index1 => $pageGroupData) {
    $pgGroupID = $pageGroupData['id'];
    $newGroupData = array();
    $newGroupData['page_group_id'] = $pgGroupID;
    $newGroupData['group_name'] = $pageGroupData['group_name'];

    foreach ($pageContentList as $index2 => $pageContentData) {
        if ($pageContentData['page_group_id'] == $pgGroupID) {
            $newGroupData['page_content'][] = array(
                "page_content_id" => $pageContentData["id"],
                "page_name" => $pageContentData["page_name"],
            );
        }
    }
    $dataRlt[] = $newGroupData;
}

$sortDataRlt = array();
foreach ($groupSort as $sIndex => $sortGroupID) {
    foreach ($dataRlt as $key => $data) {
        if ($data['page_group_id'] == $sortGroupID) {
            $sortDataRlt[] = $data;
        }
    }
}

$sortRightDataRlt = array();
foreach ($dataRlt as $key => $data) {
    if (!in_array($data['page_group_id'], $sortDataRlt)) {
        $sortRightDataRlt[] = $data;
    }
}



array_merge($sortDataRlt, $sortRightDataRlt);

echo (json_encode($sortDataRlt));
