
const PROTOCOL = window._env_.USE_TLS==="true"?"https":"http"
const WS_PROTOCOL = window._env_.USE_TLS==="true"?"wss":"ws"

const API_URL = window._env_.API_URL
const GROUPS_ADDRESS = API_URL!=="" ? API_URL : window._env_.GROUPS_SERVICE
const MESSAGES_ADDRESS = API_URL!=="" ? API_URL : window._env_.MESSAGES_SERVICE
const USERS_ADDRESS = API_URL!=="" ? API_URL : window._env_.USERS_SERVICE
const WS_ADDRESS = API_URL!=="" ? API_URL : window._env_.WS_SERVICE
const SEARCH_ADDRESS = API_URL!=="" ? API_URL : window._env_.SEARCH_SERVICE

export const groupsService = PROTOCOL+'://'+GROUPS_ADDRESS+'/groups';
export const messageService = PROTOCOL+'://'+MESSAGES_ADDRESS+'/messages';
export const userService = PROTOCOL+'://'+USERS_ADDRESS+'/users';
export const wsService = PROTOCOL+'://'+WS_ADDRESS+'/ws'
export const wsServiceWebsocket = WS_PROTOCOL+'://'+WS_ADDRESS+'/ws';
export const searchService = PROTOCOL+'://'+SEARCH_ADDRESS+'/search';

let axiosObject = require('axios').default;
axiosObject.defaults.headers.common['Content-Type'] = "application/json";


async function refreshAccessToken() {
    let response = await axiosObject.post(userService+"/refresh", {}, {
        withCredentials: true,
    })
    if (response.data.accessToken !== undefined) {
        window.localStorage.setItem("token", response.data.accessToken);
    } else {
      window.localStorage.clear();
    }
}

// Request interceptor for API calls
axiosObject.interceptors.request.use(
    async config => {
        let accessToken = window.localStorage.getItem("token")
        config.headers = { 
            'Authorization': `Bearer ${accessToken}`,
            'Accept': 'application/json',
        }
        return config;
    },
    error => {
      Promise.reject(error)
  });
  
  // Response interceptor for API calls
  axiosObject.interceptors.response.use((response) => {
    return response
  }, async function (error) {
    const originalRequest = error.config;
    if (error.response.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;
      await refreshAccessToken();
      let access_token = window.localStorage.getItem("token");
      axiosObject.defaults.headers.common['Authorization'] = 'Bearer ' + access_token;
      return axiosObject(originalRequest);
    }
    return Promise.reject(error);
  });

export default axiosObject;