import { createRouter, createWebHistory } from 'vue-router';
import Animals from '../views/Animals.vue';
import Adoptions from '../views/Adoptions.vue';
import Schedule from '../views/Schedule.vue';
import Documents from '../views/Documents.vue';
import Users from '../views/Users.vue';
import Login from '../views/Login.vue';
import Register from '../views/Register.vue';
import Reports from '../views/Reports.vue';
import store from '../store';

const routes = [
  { path: '/', redirect: '/animals' },
  { path: '/animals', component: Animals },
  { path: '/adoptions', component: Adoptions },
  { path: '/schedule', component: Schedule },
  { path: '/documents', component: Documents },
  { path: '/users', component: Users, meta: { requiresAuth: true, requiresRole: 'admin' } },
  { path: '/reports', component: Reports, meta: { requiresAuth: true } },
  { path: '/login', component: Login, meta: { guest: true } },
  { path: '/register', component: Register, meta: { guest: true } }
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

router.beforeEach((to, from, next) => {
  const isAuth = store.getters['auth/isAuthenticated'];
  const role = store.getters['auth/userRole'];
  if (to.meta.requiresAuth && !isAuth) {
    next('/login');
  } else if (to.meta.requiresRole && role !== to.meta.requiresRole) {
    next('/');
  } else if (to.meta.guest && isAuth) {
    next('/');
  } else {
    next();
  }
});

export default router; 