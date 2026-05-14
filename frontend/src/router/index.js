import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue')
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/Register.vue')
  },
  {
    path: '/',
    component: () => import('@/views/Layout.vue'),
    children: [
      {
        path: '',
        name: 'Home',
        component: () => import('@/views/Home.vue')
      },
      {
        path: 'events',
        name: 'Events',
        component: () => import('@/views/Events.vue')
      },
      {
        path: 'events/:id',
        name: 'EventDetail',
        component: () => import('@/views/EventDetail.vue')
      },
      {
        path: 'my-registrations',
        name: 'MyRegistrations',
        component: () => import('@/views/MyRegistrations.vue')
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/views/Profile.vue')
      },
      {
        path: 'admin/events',
        name: 'AdminEvents',
        component: () => import('@/views/admin/Events.vue'),
        meta: { requiresAdmin: true }
      },
      {
        path: 'admin/events/create',
        name: 'CreateEvent',
        component: () => import('@/views/admin/EventForm.vue'),
        meta: { requiresAdmin: true }
      },
      {
        path: 'admin/events/:id/edit',
        name: 'EditEvent',
        component: () => import('@/views/admin/EventForm.vue'),
        meta: { requiresAdmin: true }
      },
      {
        path: 'admin/dashboard',
        name: 'Dashboard',
        component: () => import('@/views/admin/Dashboard.vue'),
        meta: { requiresAdmin: true }
      },
      {
        path: 'admin/users',
        name: 'AdminUsers',
        component: () => import('@/views/admin/Users.vue'),
        meta: { requiresAdmin: true }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  
  if (to.meta.requiresAdmin && !userStore.isAdmin) {
    next('/')
    return
  }
  
  if ((to.path === '/login' || to.path === '/register') && userStore.token) {
    next('/')
    return
  }
  
  if (to.path.startsWith('/admin') && !userStore.token) {
    next('/login')
    return
  }
  
  next()
})

export default router
