import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/auth/Login.vue'),
      meta: { requiresAuth: false, layout: 'blank' }
    },
    {
      path: '/',
      redirect: '/dashboard'
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: () => import('@/views/dashboard/Dashboard.vue'),
      meta: { requiresAuth: true }
    },
    // Animal routes
    {
      path: '/animals',
      name: 'animals',
      component: () => import('@/views/animal/AnimalList.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/animals/new',
      name: 'animal-create',
      component: () => import('@/views/animal/AnimalForm.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/animals/:id',
      name: 'animal-detail',
      component: () => import('@/views/animal/AnimalDetail.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/animals/:id/edit',
      name: 'animal-edit',
      component: () => import('@/views/animal/AnimalForm.vue'),
      meta: { requiresAuth: true }
    },
    // Veterinary routes
    {
      path: '/veterinary',
      name: 'veterinary',
      component: () => import('@/views/veterinary/VeterinaryList.vue'),
      meta: { requiresAuth: true }
    },
    // Adoption routes
    {
      path: '/adoptions',
      name: 'adoptions',
      component: () => import('@/views/adoption/AdoptionList.vue'),
      meta: { requiresAuth: true }
    },
    // Contact routes
    {
      path: '/contacts',
      name: 'contacts',
      component: () => import('@/views/contact/ContactList.vue'),
      meta: { requiresAuth: true }
    },
    // Settings routes
    {
      path: '/settings',
      name: 'settings',
      component: () => import('@/views/settings/Settings.vue'),
      meta: { requiresAuth: true, requiresAdmin: true }
    },
    // User routes
    {
      path: '/users',
      name: 'users',
      component: () => import('@/views/user/UserList.vue'),
      meta: { requiresAuth: true, requiresAdmin: true }
    },
    {
      path: '/profile',
      name: 'profile',
      component: () => import('@/views/user/Profile.vue'),
      meta: { requiresAuth: true }
    },
    // 404
    {
      path: '/:pathMatch(.*)*',
      name: 'not-found',
      component: () => import('@/views/NotFound.vue')
    }
  ]
})

// Navigation guards
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next({ name: 'login', query: { redirect: to.fullPath } })
  } else if (to.meta.requiresAdmin && !authStore.isAdmin) {
    next({ name: 'dashboard' })
  } else if (to.name === 'login' && authStore.isAuthenticated) {
    next({ name: 'dashboard' })
  } else {
    next()
  }
})

export default router
