import axios from 'axios'
import qs from 'qs'
import { useAuthStore } from '@/stores/auth'
import { BASE_URL } from '@/utils/constants'

const httpClient = axios.create({
  baseURL: BASE_URL,
  headers: {
    'Content-type': 'application/json'
  },
  paramsSerializer: (params) => {
    return qs.stringify(params, { arrayFormat: 'repeat' })
  }
})

const httpClientWithoutInterceptor = axios.create({
  baseURL: BASE_URL,
  headers: {
    'Content-type': 'application/json'
  },
  paramsSerializer: (params) => {
    return qs.stringify(params, { arrayFormat: 'repeat' })
  }
})

httpClient.interceptors.request.use(async (config) => {
  const { accessToken } = useAuthStore()

  if (accessToken) {
    config.headers.Authorization = `Bearer ${accessToken}`
  }

  return config
})

export { httpClient, httpClientWithoutInterceptor }
