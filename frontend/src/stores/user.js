import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/utils/request'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const user = ref(JSON.parse(localStorage.getItem('user') || 'null'))

  const isLoggedIn = computed(() => !!token.value)
  const isAdmin = computed(() => user.value?.role === 'admin')

  const setToken = (newToken) => {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  const setUser = (newUser) => {
    user.value = newUser
    localStorage.setItem('user', JSON.stringify(newUser))
  }

  const logout = () => {
    token.value = ''
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  const login = async (username, password, rememberMe) => {
    const res = await api.post('/auth/login', { username, password, remember_me: rememberMe })
    setToken(res.data.token)
    setUser(res.data.user)
    return res
  }

  const register = async (data) => {
    return await api.post('/auth/register', data)
  }

  const fetchUser = async () => {
    const res = await api.get('/user/me')
    setUser(res.data)
    return res
  }

  return {
    token,
    user,
    isLoggedIn,
    isAdmin,
    setToken,
    setUser,
    logout,
    login,
    register,
    fetchUser
  }
})
