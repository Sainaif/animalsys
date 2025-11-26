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
      path: '/animals',
      name: 'public-animals',
      component: () => import('@/views/home/AnimalGallery.vue'),
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
      path: '/staff/veterinary',
      name: 'veterinary',
      component: () => import('@/views/staff/veterinary/VeterinaryDashboard.vue'),
      meta: { requiresAuth: true, layout: 'staff' }
    },
    {
      path: '/staff/veterinary/visits',
      name: 'veterinary-visits',
      component: () => import('@/views/staff/veterinary/VisitList.vue'),
      meta: { requiresAuth: true, layout: 'staff' }
    },
    {
      path: '/staff/veterinary/visits/new',
      name: 'veterinary-visit-create',
      component: () => import('@/views/staff/veterinary/VisitForm.vue'),
      meta: { requiresAuth: true, layout: 'staff' }
    },
    {
      path: '/staff/veterinary/visits/:id/edit',
      name: 'veterinary-visit-edit',
      component: () => import('@/views/staff/veterinary/VisitForm.vue'),
      meta: { requiresAuth: true, layout: 'staff' }
    },
    {
      path: '/staff/veterinary/vaccinations',
      name: 'veterinary-vaccinations',
      component: () => import('@/views/staff/veterinary/VaccinationList.vue'),
      meta: { requiresAuth: true, layout: 'staff' }
    },
    {
      path: '/staff/veterinary/medications',
      name: 'veterinary-medications',
      component: () => import('@/views/staff/veterinary/MedicationList.vue'),
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
    // Finance routes
    {
      path: '/staff/finance',
      name: 'finance',
      component: () => import('@/views/staff/finance/FinanceDashboard.vue'),
      meta: { requiresAuth: true, layout: 'staff' }
    },
    {
      path: '/staff/finance/donors',
      name: 'finance-donors',
      component: () => import('@/views/staff/finance/DonorList.vue'),
      meta: { requiresAuth: true, layout: 'staff' }
    },
    {
      path: '/staff/finance/donors/new',
      name: 'finance-donor-create',
      component: () => import('@/views/staff/finance/DonorForm.vue'),
      meta: { requiresAuth: true, layout: 'staff' }
    },
    {
      path: '/staff/finance/donors/:id/edit',
      name: 'finance-donor-edit',
      component: () => import('@/views/staff/finance/DonorForm.vue'),
      meta: { requiresAuth: true, layout: 'staff' }
    },
    {
      path: '/staff/finance/donations',
      name: 'finance-donations',
      component: () => import('@/views/staff/finance/DonationList.vue'),
      meta: { requiresAuth: true, layout: 'staff' }
    },
    {
      path: '/staff/finance/donations/new',
      name: 'finance-donation-create',
      component: () => import('@/views/staff/finance/DonationForm.vue'),
      meta: { requiresAuth: true, layout: 'staff' }
    },
    {
      path: '/staff/finance/donations/:id/edit',
      name: 'finance-donation-edit',
      component: () => import('@/views/staff/finance/DonationForm.vue'),
      meta: { requiresAuth: true, layout: 'staff' }
    },
    {
      path: '/staff/finance/campaigns',
      name: 'finance-campaigns',
      component: () => import('@/views/staff/finance/CampaignList.vue'),
      meta: { requiresAuth: true, layout: 'staff' }
    },
    {
      path: '/staff/finance/campaigns/new',
      name: 'finance-campaign-create',
      component: () => import('@/views/staff/finance/CampaignForm.vue'),
      meta: { requiresAuth: true, layout: 'staff' }
    },
    {
      path: '/staff/finance/campaigns/:id/edit',
      name: 'finance-campaign-edit',
      component: () => import('@/views/staff/finance/CampaignForm.vue'),
      meta: { requiresAuth: true, layout: 'staff' }
    },
    // Events routes
    {
      path: '/staff/events',
      name: 'events',
      component: () => import('@/views/staff/events/EventList.vue'),
      meta: { requiresAuth: true, layout: 'staff' }
    },
    {
      path: '/staff/events/new',
      name: 'event-create',
      component: () => import('@/views/staff/events/EventForm.vue'),
      meta: { requiresAuth: true, layout: 'staff' }
    },
    {
      path: '/staff/events/:id',
      name: 'event-detail',
      component: () => import('@/views/staff/events/EventDetail.vue'),
      meta: { requiresAuth: true, layout: 'staff' }
    },
    {
      path: '/staff/events/:id/edit',
      name: 'event-edit',
      component: () => import('@/views/staff/events/EventForm.vue'),
      meta: { requiresAuth: true, layout: 'staff' }
    },
    {
      path: '/staff/volunteers',
      name: 'volunteers',
      component: () => import('@/views/staff/events/VolunteerList.vue'),
      meta: { requiresAuth: true, layout: 'staff' }
    },
    {
      path: '/staff/volunteers/:id',
      name: 'volunteer-detail',
      component: () => import('@/views/staff/events/VolunteerDetail.vue'),
      meta: { requiresAuth: true, layout: 'staff' }
    },
    // Communication routes
    {
      path: '/staff/communication',
      name: 'communication',
      component: () => import('@/views/staff/communication/CommunicationDashboard.vue'),
      meta: { requiresAuth: true, layout: 'staff' }
    },
    {
      path: '/staff/communication/templates',
      name: 'communication-templates',
      component: () => import('@/views/staff/communication/EmailTemplateList.vue'),
      meta: { requiresAuth: true, layout: 'staff' }
    },
    {
      path: '/staff/communication/templates/new',
      name: 'communication-template-create',
      component: () => import('@/views/staff/communication/EmailTemplateForm.vue'),
      meta: { requiresAuth: true, layout: 'staff' }
    },
    {
      path: '/staff/communication/templates/:id/edit',
      name: 'communication-template-edit',
      component: () => import('@/views/staff/communication/EmailTemplateForm.vue'),
      meta: { requiresAuth: true, layout: 'staff' }
    },
    {
      path: '/staff/communication/campaigns',
      name: 'communication-campaigns',
      component: () => import('@/views/staff/communication/EmailCampaignList.vue'),
      meta: { requiresAuth: true, layout: 'staff' }
    },
    {
      path: '/staff/communication/campaigns/new',
      name: 'communication-campaign-create',
      component: () => import('@/views/staff/communication/EmailCampaignForm.vue'),
      meta: { requiresAuth: true, layout: 'staff' }
    },
    {
      path: '/staff/communication/campaigns/:id/edit',
      name: 'communication-campaign-edit',
      component: () => import('@/views/staff/communication/EmailCampaignForm.vue'),
      meta: { requiresAuth: true, layout: 'staff' }
    },
    {
      path: '/staff/communication/logs',
      name: 'communication-logs',
      component: () => import('@/views/staff/communication/CommunicationLogList.vue'),
      meta: { requiresAuth: true, layout: 'staff' }
    },
    {
      path: '/staff/communication/logs/new',
      name: 'communication-log-create',
      component: () => import('@/views/staff/communication/CommunicationLogForm.vue'),
      meta: { requiresAuth: true, layout: 'staff' }
    },
    {
      path: '/staff/communication/logs/:id/edit',
      name: 'communication-log-edit',
      component: () => import('@/views/staff/communication/CommunicationLogForm.vue'),
      meta: { requiresAuth: true, layout: 'staff' }
    },
    // Partner routes
    {
      path: '/staff/partners',
      name: 'partners',
      component: () => import('@/views/staff/partners/PartnerList.vue'),
      meta: { requiresAuth: true, layout: 'staff' }
    },
    {
      path: '/staff/partners/new',
      name: 'partner-create',
      component: () => import('@/views/staff/partners/PartnerForm.vue'),
      meta: { requiresAuth: true, layout: 'staff' }
    },
    {
      path: '/staff/partners/:id/edit',
      name: 'partner-edit',
      component: () => import('@/views/staff/partners/PartnerForm.vue'),
      meta: { requiresAuth: true, layout: 'staff' }
    },
    {
      path: '/staff/partners/transfers',
      name: 'animal-transfers',
      component: () => import('@/views/staff/partners/AnimalTransferList.vue'),
      meta: { requiresAuth: true, layout: 'staff' }
    },
    // Inventory routes
    {
      path: '/staff/inventory',
      name: 'inventory',
      component: () => import('@/views/staff/inventory/InventoryList.vue'),
      meta: { requiresAuth: true, layout: 'staff' }
    },
    {
      path: '/staff/inventory/new',
      name: 'inventory-create',
      component: () => import('@/views/staff/inventory/CreateInventoryItem.vue'),
      meta: { requiresAuth: true, layout: 'staff' }
    },
    // Reports routes
    {
      path: '/staff/reports',
      name: 'reports',
      component: () => import('@/views/staff/reports/ReportsDashboard.vue'),
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
