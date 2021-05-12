import axios from 'axios'
import Router from '@/router'

const client = axios.create({
  baseURL: `/api`,
  withCredentials: true
})

client.interceptors.response.use((response) => response, (error) => {
  if (error.response.status === 401) {
    Router.push({name: 'login'})
  }
  return Promise.reject(error)
})

export default client
