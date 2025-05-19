import axios from 'axios';
import { API_BASE_URL, API_ENDPOINTS } from '../config/api';

const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
    'Accept': 'application/json',
  },
  withCredentials: true,
});

// Function to fetch CSRF token
const fetchCsrfToken = async () => {
  try {
    const response = await axios.get(API_ENDPOINTS.CSRF.TOKEN, {
      baseURL: API_BASE_URL,
      withCredentials: true,
    });
    return response.data.csrfToken;
  } catch (error) {
    console.error('Failed to fetch CSRF token:', error);
    return null;
  }
};

api.interceptors.request.use(
  async (config) => {
    // Skip auth token for CSRF endpoint
    const isCsrfEndpoint = config.url === '/api/csrf-token';
    
    // Add auth token if available and not CSRF endpoint
    if (!isCsrfEndpoint) {
      const token = localStorage.getItem('token');
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
    }

    // Add CSRF token for non-GET requests
    if (config.method !== 'get') {
      const csrfToken = await fetchCsrfToken();
      if (csrfToken) {
        config.headers['X-Csrf-Token'] = csrfToken;
      }
    }

    return config;
  },
  (error) => {
    console.error('Request error:', error);
    return Promise.reject(error);
  }
);

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response) {
      const errorMessage = error.response.data.error || 'An unexpected error occurred';
      
      switch (error.response.status) {
        case 400:
          console.error('Bad Request:', errorMessage);
          break;
        case 401:
          localStorage.removeItem('token');
          window.location.href = '/login';
          break;
        case 403:
          console.error('Forbidden:', errorMessage);
          window.location.href = '/';
          break;
        case 404:
          console.error('Resource not found:', errorMessage);
          break;
        case 500:
          console.error('Server error:', errorMessage);
          break;
        default:
          console.error('API Error:', error.response.status, errorMessage);
      }
    } else if (error.request) {
      console.error('Network Error:', error.request);
      if (error.message === 'Network Error') {
        console.error('CORS Error: The request was blocked due to CORS policy');
      }
    } else {
      console.error('Error:', error.message);
    }
    
    return Promise.reject({
      error: error.response?.data?.error || error.message || 'An unexpected error occurred'
    });
  }
);

export default api;
