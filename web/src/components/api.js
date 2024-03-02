import axios from 'axios';
import c from '../config';

const instance = axios.create({
    baseURL: `${c.apiUrl}`,
  });

// Add a request interceptor
instance.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('accessToken');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);


// Add a response interceptor
instance.interceptors.response.use(
    (response) => response,
    async (error) => {
      const originalRequest = error.config;
  
      // If the error status is 401 and there is no originalRequest._retry flag,
      // it means the token has expired and we need to refresh it
      if (error.response.status === 401 && !originalRequest._retry) {
        originalRequest._retry = true;
  
        try {
          const refreshToken = localStorage.getItem('refreshToken');
          const response = await axios.post(`${c.apiUrl}/refresh-token`,{}, {
            headers: {
              'Content-Type': 'application/json',
              'Authorization': `Bearer ${refreshToken}`
            }
          });
          const accessToken = response.data.access_token;
  
          localStorage.setItem('accessToken', accessToken);
  
          // Retry the original request with the new token
          originalRequest.headers.Authorization = `Bearer ${accessToken}`;
          return axios(originalRequest);
        } catch (error) {
            console.log(error)
          // Handle refresh token error or redirect to login
        }
      }
  
      return Promise.reject(error);
    }
  );
  
export default instance
