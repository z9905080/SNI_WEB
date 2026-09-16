SET NAMES utf8mb4;

INSERT INTO `user` (`id`, `account`, `pwd`, `identity`, `user_name`) VALUES
(1, 'admin', '$2a$10$ns9js0F6P2CaUHD15cOsg.mW1MgBoa30WozPKYoOwLPqZaQfWk63S', 1, '管理者');

INSERT INTO `web_config` (`id`, `data_key`, `data_value`) VALUES
(1, 'web_title', '生長之家'),
(2, 'web_sub_title', '一句簡單的感謝'),
(3, 'page_group_sort', '[1,3,2]'),
(4, 'facebook_url', 'https://www.facebook.com/seichonoie.tw');

INSERT INTO `page_group` (`id`, `group_name`, `page_sort`) VALUES
(1, '首頁', ''),
(2, '主要行事活動', '[4,3]'),
(3, '關於我們', '[]');

INSERT INTO `page_content` (`id`, `page_group_id`, `page_name`, `html_context`) VALUES
(1, 1, '【最新訊息】', '<p><span style="font-size: 14pt; color: #e03e2d; font-family: arial, helvetica, sans-serif;">●官方網站重新改版，歡迎瀏覽！</span></p>\n<table style="border-collapse: collapse; width: 100%;" border="1"><tbody>\n<tr><td style="width: 100%; background-color: #18a085; text-align: center;"><strong><span style="color: #ffffff;">國際本部發行的影片</span></strong></td></tr>\n<tr><td><iframe src="//www.youtube.com/embed/V4LhnGIPZt4" width="100%" height="225" allowfullscreen="allowfullscreen"></iframe></td></tr>\n</tbody></table>\n<p>前往 <a href="/#/3">練成會</a> 了解更多。</p>'),
(2, 3, '生長之家簡介', '<h2>生長之家簡介</h2><p style="text-align: justify;">生長之家是以「感謝」為中心的信仰團體。</p>'),
(3, 2, '練成會', '<p>練成會日程如下：</p><table border="1" style="width: 900px;"><tbody><tr><td>日期</td><td>地點</td><td>主題</td><td>備註</td></tr><tr><td>每月第一週</td><td>台北道場</td><td>一般練成會</td><td>需事先報名</td></tr></tbody></table>'),
(4, 2, '誌友會', '<p><span style="background-color: #fbeeb8;">每週日上午舉行。</span></p>');

INSERT INTO `carousel` (`id`, `image`, `url`) VALUES
(1, '/php/picture/sample-1.jpg', ''),
(2, '/php/picture/sample-2.jpg', 'https://www.facebook.com/seichonoie.tw');

INSERT INTO `marquee` (`id`, `text`, `color`) VALUES
(1, '合掌感謝！', '#1EFF00'),
(2, '官方網站改版上線', '#FFFFFF');
