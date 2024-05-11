import { createRouter, createWebHistory } from 'vue-router'
import ChatView from '../views/ChatView.vue'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: ChatView,
      meta: {
        requiredAuth: true
      }
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/LoginView.vue')
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('@/views/RegisterView.vue')
    }
  ]
})

router.beforeEach((to, from, next) => {
  const { isAuth } = useAuthStore()

  if (!isAuth && to.meta.requiredAuth) {
    next({ name: 'login' })
    return
  }

  next()
})

export default router
