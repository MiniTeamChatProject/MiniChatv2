import axios, { type AxiosInstance, type AxiosRequestConfig } from 'axios'

// User Service HTTP client (using Vite proxy)
export const userHttp: AxiosInstance = axios.create({
  baseURL: '/user',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Room Service HTTP client (using Vite proxy)
export const roomHttp: AxiosInstance = axios.create({
  baseURL: '/room',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor for User Service
userHttp.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

// Request interceptor for Room Service
roomHttp.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    // Add X-User-ID header (required by room service)
    const userId = localStorage.getItem('userId')
    if (userId) {
      config.headers['X-User-ID'] = userId
    }
    return config
  },
  (error) => Promise.reject(error)
)

// Response interceptor for both clients
const createResponseInterceptor = (http: AxiosInstance) => {
  http.interceptors.response.use(
    (response) => response.data,
    (error) => {
      if (error.response?.status === 401) {
        // Token expired, redirect to login
        localStorage.clear()
        window.location.href = '/login'
      }
      // Ensure error has consistent structure
      if (error.response?.data) {
        error.message = error.response.data.message || error.message
      }
      return Promise.reject(error)
    }
  )
}

createResponseInterceptor(userHttp)
createResponseInterceptor(roomHttp)

export default userHttp
