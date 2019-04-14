import axios from 'axios';
import { Message, Loading } from 'element-ui';

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
axios.interceptors.request.use(config => {
    //加載動畫
    startLoading();
    return config;
}, error => {
    return Promise.reject(error);
}
)

// 響應攔截
axios.interceptors.response.use(response => {
    //結束加載動畫
    endLoading();
    return response;
}, error => {
    //錯誤提醒
    endLoading();
    Message.error(error.response.data);
    return Promise.reject(error);
})

export default axios;