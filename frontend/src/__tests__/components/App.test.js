import { describe, it, expect, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createStore } from 'vuex'
import { createRouter, createMemoryHistory } from 'vue-router'
import App from '../../App.vue'

describe('App Component', () => {
  let store
  let router

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
  })

  it('should render the app container', () => {
    const wrapper = mount(App, {
      global: {
        plugins: [store, router]
      }
    })

    expect(wrapper.find('#app').exists()).toBe(true)
  })

  it('should not show navigation when not authenticated', () => {
    const wrapper = mount(App, {
      global: {
        plugins: [store, router]
      }
    })

    expect(wrapper.find('nav').exists()).toBe(false)
  })

  it('should show navigation when authenticated', async () => {
    store.state.auth.token = 'test-token'
    store.state.auth.user = { id: '1', username: 'test', role: 'user' }

    const wrapper = mount(App, {
      global: {
        plugins: [store, router]
      }
    })

    expect(wrapper.find('nav').exists()).toBe(true)
  })
})
