import axios from "axios";

const PROTOCOL = window._env_.USE_TLS==="true"?"https":"http";
const WS_PROTOCOL = window._env_.USE_TLS==="true"?"wss":"ws";

const API_URL = window._env_.API_URL;

export const groupsService = PROTOCOL+'://'+API_URL+'/groups';
export const messageService = PROTOCOL+'://'+API_URL+'/messages';
export const userService = PROTOCOL+'://'+API_URL+'/users';
export const searchService = PROTOCOL+'://'+API_URL+'/search';
export const wsService = PROTOCOL+'://'+API_URL+'/ws'
export const wsServiceWebsocket = WS_PROTOCOL+'://'+API_URL+'/ws';
export const webrtcService = PROTOCOL+'://'+API_URL+'/video-call'
export const webrtcServiceWebsocket = WS_PROTOCOL+'://'+API_URL+'/video-call';

async function refreshAccessToken() {
  let response;
  try {
    response = await axios.post(userService+"/refresh", {}, {
      withCredentials: true,
    })
    window.localStorage.setItem("token", response.data.accessToken);
  } catch (err) {
    window.localStorage.clear();
    window.dispatchEvent(new Event("logout"));
  }
}

let axiosObject = axios.create();
axiosObject.defaults.headers.common['Content-Type'] = "application/json";

// Request interceptor for API calls
axiosObject.interceptors.request.use(
    config => {
        let accessToken = window.localStorage.getItem("token");
        config.headers = { 
            'Authorization': `Bearer ${accessToken}`,
            'Accept': 'application/json',
        };
        return config;
    },
    error => {
      Promise.reject(error);
  });
  
  // Response interceptor for API calls
  axiosObject.interceptors.response.use((response) => {
    return response;
  }, async error => {
    const originalRequest = error.config;

    if (error.response.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;
      
      await refreshAccessToken();

      const access_token = window.localStorage.getItem("token");

      if (access_token === null) return Promise.reject(error);

      axiosObject.defaults.headers.common['Authorization'] = 'Bearer ' + access_token;
      return axiosObject(originalRequest);
    }
    return Promise.reject(error);
  });

export default axiosObject;