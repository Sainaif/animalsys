import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

// Layouts
const PublicLayout = () => import('../layouts/PublicLayout.vue')
const AuthenticatedLayout = () => import('../layouts/AuthenticatedLayout.vue')

// Public Pages
const Home = () => import('../views/public/Home.vue')
const Login = () => import('../views/public/Login.vue')
const Register = () => import('../views/public/Register.vue')
const AnimalsPublic = () => import('../views/public/AnimalsPublic.vue')
const AnimalDetails = () => import('../views/public/AnimalDetails.vue')
const CampaignsPublic = () => import('../views/public/CampaignsPublic.vue')

// Authenticated Pages
const Dashboard = () => import('../views/Dashboard.vue')
const Profile = () => import('../views/Profile.vue')

// Admin Pages
const Users = () => import('../views/admin/Users.vue')

// Animals Module
const Animals = () => import('../views/animals/Animals.vue')
const AnimalForm = () => import('../views/animals/AnimalForm.vue')
const AnimalView = () => import('../views/animals/AnimalView.vue')

// Adoptions Module
const Adoptions = () => import('../views/adoptions/Adoptions.vue')
const AdoptionForm = () => import('../views/adoptions/AdoptionForm.vue')
const AdoptionView = () => import('../views/adoptions/AdoptionView.vue')

// Volunteers Module
const Volunteers = () => import('../views/volunteers/Volunteers.vue')
const VolunteerForm = () => import('../views/volunteers/VolunteerForm.vue')
const VolunteerView = () => import('../views/volunteers/VolunteerView.vue')

// Schedule Module
const Schedules = () => import('../views/schedules/Schedules.vue')

// Documents Module
const Documents = () => import('../views/documents/Documents.vue')

// Finance Module
const Finance = () => import('../views/finance/Finance.vue')

// Donors Module
const Donors = () => import('../views/donors/Donors.vue')

// Inventory Module
const Inventory = () => import('../views/inventory/Inventory.vue')

// Veterinary Module
const Veterinary = () => import('../views/veterinary/Veterinary.vue')

// Campaigns Module
const Campaigns = () => import('../views/campaigns/Campaigns.vue')

// Partners Module
const Partners = () => import('../views/partners/Partners.vue')

// Communications Module
const Communications = () => import('../views/communications/Communications.vue')

// Error Pages
const NotFound = () => import('../views/errors/NotFound.vue')
const Unauthorized = () => import('../views/errors/Unauthorized.vue')

const routes = [
  // Public routes
  {
    path: '/',
    component: PublicLayout,
    children: [
      {
        path: '',
        name: 'home',
        component: Home,
        meta: { public: true }
      },
      {
        path: 'login',
        name: 'login',
        component: Login,
        meta: { public: true, guest: true }
      },
      {
        path: 'register',
        name: 'register',
        component: Register,
        meta: { public: true, guest: true }
      },
      {
        path: 'animals',
        name: 'animals-public',
        component: AnimalsPublic,
        meta: { public: true }
      },
      {
        path: 'animals/:id',
        name: 'animal-details',
        component: AnimalDetails,
        meta: { public: true }
      },
      {
        path: 'campaigns',
        name: 'campaigns-public',
        component: CampaignsPublic,
        meta: { public: true }
      },
    ]
  },

  // Authenticated routes
  {
    path: '/app',
    component: AuthenticatedLayout,
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'dashboard',
        component: Dashboard,
      },
      {
        path: 'profile',
        name: 'profile',
        component: Profile,
      },

      // Users (Admin only)
      {
        path: 'users',
        name: 'users',
        component: Users,
        meta: { requiresRole: 'admin' }
      },

      // Animals (Employee+)
      {
        path: 'animals',
        meta: { requiresRole: 'employee' },
        children: [
          {
            path: '',
            name: 'animals-list',
            component: Animals,
          },
          {
            path: 'create',
            name: 'animal-create',
            component: AnimalForm,
          },
          {
            path: ':id',
            name: 'animal-view',
            component: AnimalView,
          },
          {
            path: ':id/edit',
            name: 'animal-edit',
            component: AnimalForm,
          },
        ]
      },

      // Adoptions
      {
        path: 'adoptions',
        children: [
          {
            path: '',
            name: 'adoptions-list',
            component: Adoptions,
          },
          {
            path: 'apply',
            name: 'adoption-apply',
            component: AdoptionForm,
          },
          {
            path: ':id',
            name: 'adoption-view',
            component: AdoptionView,
          },
        ]
      },

      // Volunteers (Employee+)
      {
        path: 'volunteers',
        meta: { requiresRole: 'employee' },
        children: [
          {
            path: '',
            name: 'volunteers-list',
            component: Volunteers,
          },
          {
            path: 'create',
            name: 'volunteer-create',
            component: VolunteerForm,
          },
          {
            path: ':id',
            name: 'volunteer-view',
            component: VolunteerView,
          },
          {
            path: ':id/edit',
            name: 'volunteer-edit',
            component: VolunteerForm,
          },
        ]
      },

      // Schedules (Volunteer+)
      {
        path: 'schedules',
        name: 'schedules',
        component: Schedules,
        meta: { requiresRole: 'volunteer' }
      },

      // Documents (Employee+)
      {
        path: 'documents',
        name: 'documents',
        component: Documents,
        meta: { requiresRole: 'employee' }
      },

      // Finance (Employee+)
      {
        path: 'finance',
        name: 'finance',
        component: Finance,
        meta: { requiresRole: 'employee' }
      },

      // Donors (Employee+)
      {
        path: 'donors',
        name: 'donors',
        component: Donors,
        meta: { requiresRole: 'employee' }
      },

      // Inventory (Employee+)
      {
        path: 'inventory',
        name: 'inventory',
        component: Inventory,
        meta: { requiresRole: 'employee' }
      },

      // Veterinary (Employee+)
      {
        path: 'veterinary',
        name: 'veterinary',
        component: Veterinary,
        meta: { requiresRole: 'employee' }
      },

      // Campaigns (Employee+)
      {
        path: 'campaigns',
        name: 'campaigns',
        component: Campaigns,
        meta: { requiresRole: 'employee' }
      },

      // Partners (Employee+)
      {
        path: 'partners',
        name: 'partners',
        component: Partners,
        meta: { requiresRole: 'employee' }
      },

      // Communications (Employee+)
      {
        path: 'communications',
        name: 'communications',
        component: Communications,
        meta: { requiresRole: 'employee' }
      },
    ]
  },

  // Error routes
  {
    path: '/unauthorized',
    name: 'unauthorized',
    component: Unauthorized,
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'not-found',
    component: NotFound,
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else {
      return { top: 0 }
    }
  },
})

// Navigation guards
router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()

  // Check if route requires authentication
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next({ name: 'login', query: { redirect: to.fullPath } })
    return
  }

  // Check if route is for guests only (logged out users)
  if (to.meta.guest && authStore.isAuthenticated) {
    next({ name: 'dashboard' })
    return
  }

  // Check role requirements
  if (to.meta.requiresRole) {
    const hasRole = authStore.hasRole(to.meta.requiresRole)
    if (!hasRole) {
      next({ name: 'unauthorized' })
      return
    }
  }

  next()
})

export default router
