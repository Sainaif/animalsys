import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    // Public routes
    {
      path: '/',
      name: 'home',
      component: () => import('@/views/home/Home.vue'),
      meta: { requiresAuth: false, layout: 'public' }
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/auth/Login.vue'),
      meta: { requiresAuth: false, layout: 'blank' }
    },
    // Staff routes
    {
      path: '/dashboard',
      name: 'dashboard',
      component: () => import('@/views/dashboard/Dashboard.vue'),
      meta: { requiresAuth: true, layout: 'staff' }
    },
    // Animal routes
    {
      path: '/staff/animals',
      name: 'animals',
      component: () => import('@/views/staff/animals/AnimalList.vue'),
      meta: { requiresAuth: true, layout: 'staff' }
    },
    {
      path: '/staff/animals/new',
      name: 'animal-create',
      component: () => import('@/views/staff/animals/AnimalForm.vue'),
      meta: { requiresAuth: true, layout: 'staff' }
    },
    {
      path: '/staff/animals/:id',
      name: 'animal-detail',
      component: () => import('@/views/staff/animals/AnimalDetail.vue'),
      meta: { requiresAuth: true, layout: 'staff' }
    },
    {
      path: '/staff/animals/:id/edit',
      name: 'animal-edit',
      component: () => import('@/views/staff/animals/AnimalForm.vue'),
      meta: { requiresAuth: true, layout: 'staff' }
    },
    // Veterinary routes
    {
      path: '/veterinary',
      name: 'veterinary',
      component: () => import('@/views/veterinary/VeterinaryList.vue'),
      meta: { requiresAuth: true, layout: 'staff' }
    },
    // Adoption routes
    {
      path: '/staff/adoptions/applications',
      name: 'adoption-applications',
      component: () => import('@/views/staff/adoptions/ApplicationList.vue'),
      meta: { requiresAuth: true, layout: 'staff' }
    },
    {
      path: '/staff/adoptions/applications/:id',
      name: 'adoption-application-detail',
      component: () => import('@/views/staff/adoptions/ApplicationDetail.vue'),
      meta: { requiresAuth: true, layout: 'staff' }
    },
    {
      path: '/staff/adoptions',
      name: 'adoptions',
      component: () => import('@/views/staff/adoptions/AdoptionList.vue'),
      meta: { requiresAuth: true, layout: 'staff' }
    },
    {
      path: '/staff/adoptions/new',
      name: 'adoption-create',
      component: () => import('@/views/staff/adoptions/AdoptionForm.vue'),
      meta: { requiresAuth: true, layout: 'staff' }
    },
    {
      path: '/staff/adoptions/:id',
      name: 'adoption-detail',
      component: () => import('@/views/staff/adoptions/AdoptionList.vue'),
      meta: { requiresAuth: true, layout: 'staff' }
    },
    // Contact routes
    {
      path: '/contacts',
      name: 'contacts',
      component: () => import('@/views/contact/ContactList.vue'),
      meta: { requiresAuth: true, layout: 'staff' }
    },
    // Settings routes
    {
      path: '/settings',
      name: 'settings',
      component: () => import('@/views/settings/Settings.vue'),
      meta: { requiresAuth: true, requiresAdmin: true, layout: 'staff' }
    },
    // User routes
    {
      path: '/users',
      name: 'users',
      component: () => import('@/views/user/UserList.vue'),
      meta: { requiresAuth: true, requiresAdmin: true, layout: 'staff' }
    },
    {
      path: '/profile',
      name: 'profile',
      component: () => import('@/views/user/Profile.vue'),
      meta: { requiresAuth: true, layout: 'staff' }
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
