import axios from 'axios';

const axiosInstance = axios.create();

axiosInstance.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error != null && error.response.status == 401) {
      return window.location.href = "/login";
    }
    return Promise.reject(error);
  }
);

export default axiosInstance;
