import axios from 'axios';
import { Message, Loading } from 'element-ui';
import router from './router';

let loading;
function startLoading() {
    loading = Loading.service({
        lock: true,
        text: '拚命加載中...',
        background: 'rgba(0, 0, 0, 0, 7)'
    });
}

function endLoading() {
    loading.close();
}

// 請求攔截
// axios.interceptors.request.use(config => {
//     console.log(config);

//     //加載動畫
//     startLoading();

//     if (window.$cookies.isKey("sid")) {
//         //設定統一的請求header
//         config.headers.Authorization = window.$cookies.isKey("sid");
//     }

//     return config;
// }, error => {
//     return Promise.reject(error);
// })
// console.log(axios.interceptors);


// // 響應攔截
// axios.interceptors.response.use(response => {
//     //結束加載動畫
//     endLoading();
//     return response;
// }, error => {
//     //錯誤提醒
//     endLoading();
//     Message.error(error.response.data);

//     //獲取錯誤狀態碼
//     const { status } = error.response;

//     if (status === 401) {
//         Message.error('token失效,請重新登入!');
//         // 清除token
//         localStorage.removeItem('eleToken');
//         // 跳轉至登入頁面
//         router.push('/login');
//     }

//     return Promise.reject(error);
// })

export default axios;