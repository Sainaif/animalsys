import { describe, it, expect, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createStore } from 'vuex'
import { createRouter, createMemoryHistory } from 'vue-router'
import { createI18n } from 'vue-i18n'
import App from '../../App.vue'

describe('App Component', () => {
  let store
  let router
  let i18n

  beforeEach(() => {
    store = createStore({
      modules: {
        auth: {
          namespaced: true,
          state: {
            token: null,
            user: null
          },
          getters: {
            isAuthenticated: state => !!state.token,
            user: state => state.user
          },
          actions: {
            logout: () => {}
          }
        }
      }
    })

    router = createRouter({
      history: createMemoryHistory(),
      routes: [
        { path: '/', component: { template: '<div>Home</div>' } },
        { path: '/login', component: { template: '<div>Login</div>' } }
      ]
    })

    i18n = createI18n({
      legacy: false,
      locale: 'en',
      messages: {
        en: {
          nav: {
            home: 'Home',
            animals: 'Animals',
            finances: 'Finances',
            users: 'Users',
            logout: 'Logout'
          }
        }
      }
    })
  })

  it('should render the app container', () => {
    const wrapper = mount(App, {
      global: {
        plugins: [store, router, i18n]
      }
    })

    expect(wrapper.find('#app').exists()).toBe(true)
  })

  it('should not show navigation when not authenticated', () => {
    const wrapper = mount(App, {
      global: {
        plugins: [store, router, i18n]
      }
    })

    expect(wrapper.find('nav').exists()).toBe(false)
  })

  it('should show navigation when authenticated', async () => {
    store.state.auth.token = 'test-token'
    store.state.auth.user = { id: '1', username: 'test', role: 'user' }

    const wrapper = mount(App, {
      global: {
        plugins: [store, router, i18n]
      }
    })

    expect(wrapper.find('nav').exists()).toBe(true)
  })
})
