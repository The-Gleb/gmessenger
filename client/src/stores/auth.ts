import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import Cookies from 'js-cookie'
import { type User } from '@/types'

export const useAuthStore = defineStore('user', () => {
  const user = ref<User | null>(null)
  const accessToken = ref<string | null>(null)
  const refreshToken = ref<string | null>(null)

  const isAuth = computed(() => {
    return !!user.value
  })

  const setUser = (_user: User) => {
    user.value = _user
    localStorage.setItem('auth.user', JSON.stringify(_user))
  }

  const init = () => {
    const accessToken = Cookies.get('access_token')

    const user = localStorage.getItem('auth.user')

    if (accessToken && user) {
      setAccessToken(accessToken)
      setUser(JSON.parse(user) as User)
    }
  }

  const authTokensSetup = (accessToken: string) => {
    setAccessToken(accessToken)
  }

  const destroy = () => {
    clearUser()
    accessToken.value = null
  }

  const setAccessToken = (_token: string) => {
    Cookies.set('access_token', _token)
    accessToken.value = _token
  }

  const clearUser = () => {
    user.value = null
    localStorage.removeItem('auth.user')
  }

  return {
    user: computed(() => user.value),
    accessToken: computed(() => accessToken.value),
    refreshToken: computed(() => refreshToken.value),
    isAuth,
    setUser,
    clearUser,
    authTokensSetup,
    init,
    destroy,
    setAccessToken
  }
})
