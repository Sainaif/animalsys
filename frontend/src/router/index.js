import { createRouter, createWebHistory } from 'vue-router'
import store from '../store'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: () => import('../views/Home.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue')
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('../views/Register.vue')
  },
  {
    path: '/animals',
    name: 'Animals',
    component: () => import('../views/Animals.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/adoptions',
    name: 'Adoptions',
    component: () => import('../views/Adoptions.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/schedules',
    name: 'Schedules',
    component: () => import('../views/Schedules.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/documents',
    name: 'Documents',
    component: () => import('../views/Documents.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/finances',
    name: 'Finances',
    component: () => import('../views/Finances.vue'),
    meta: { requiresAuth: true, requiresRole: ['admin', 'employee'] }
  },
  {
    path: '/users',
    name: 'Users',
    component: () => import('../views/Users.vue'),
    meta: { requiresAuth: true, requiresRole: ['admin'] }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const isAuthenticated = store.getters['auth/isAuthenticated']
  const user = store.getters['auth/user']

  if (to.meta.requiresAuth && !isAuthenticated) {
    next('/login')
  } else if (to.meta.requiresRole && user) {
    if (to.meta.requiresRole.includes(user.role)) {
      next()
    } else {
      next('/')
    }
  } else {
    next()
  }
})

export default router
