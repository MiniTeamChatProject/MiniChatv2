import { createRouter, createWebHistory } from 'vue-router'
import { authStore } from '@/stores/auth'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/rooms',
    name: 'RoomList',
    component: () => import('@/views/RoomList.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/room/:id',
    name: 'ChatRoom',
    component: () => import('@/views/ChatRoom.vue'),
    meta: { requiresAuth: true },
    props: true
  },
  {
    path: '/',
    redirect: '/rooms'
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Route guard - using global authStore directly
router.beforeEach((to) => {
  if (to.meta.requiresAuth && !authStore.isLoggedIn) {
    return '/login'
  }
  if (to.path === '/login' && authStore.isLoggedIn) {
    return '/rooms'
  }
  return true
})

export default router
