import axios from 'axios'

let client = axios.create({
  baseURL: `/api`,
  withCredentials: true
})

client.interceptors.response.use((response) => response, (error) => {
  if (error.response.status === 401) {
    window.location.href = '/#/login'
  }
  return Promise.reject(error)
})

export default client
