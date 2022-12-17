export const groupsService = 'http://localhost:8081/api';
export const messageService = 'http://localhost:8082/api';
export const userService = 'http://localhost:8083/api';
export const wsService = 'http://localhost:8084/api'
export const wsServiceWebsocket = 'ws://localhost:8084';

let axiosObject = require('axios').default;
axiosObject.defaults.headers.common['Content-Type'] = "application/json";


async function refreshAccessToken() {
    let response = await axiosObject.post(userService+"/refresh", {}, {
        withCredentials: true,
    })
    if (response.data.accessToken !== undefined) {
        window.localStorage.setItem("token", response.data.accessToken);
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