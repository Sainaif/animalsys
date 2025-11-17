import { mount } from '@vue/test-utils'
import { describe, it, expect, vi, beforeAll } from 'vitest'
import PrimeVue from 'primevue/config'
import CreateInventoryItem from './CreateInventoryItem.vue'
import { inventoryService } from '@/services/inventoryService'
import { useRouter } from 'vue-router'

// Mock the router
vi.mock('vue-router', () => ({
  useRouter: vi.fn(() => ({
    push: vi.fn()
  }))
}))

vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (key) => key
  })
}))

// Mock the inventoryService
vi.mock('@/services/inventoryService', () => ({
  inventoryService: {
    createInventoryItem: vi.fn()
  }
}))

// Mock PrimeVue toast composable
vi.mock('primevue/usetoast', () => ({
  useToast: () => ({
    add: vi.fn()
  })
}))

beforeAll(() => {
  if (!window.matchMedia) {
    window.matchMedia = vi.fn().mockImplementation((query) => ({
      matches: false,
      media: query,
      onchange: null,
      addListener: vi.fn(),
      removeListener: vi.fn(),
      addEventListener: vi.fn(),
      removeEventListener: vi.fn(),
      dispatchEvent: vi.fn()
    }))
  }
})

describe('CreateInventoryItem.vue', () => {
  it('renders the form correctly', () => {
    const wrapper = mount(CreateInventoryItem, {
      global: {
        plugins: [PrimeVue]
      }
    })
    expect(wrapper.find('h2').text()).toBe('Create Inventory Item')
    expect(wrapper.find('#name').exists()).toBe(true)
    expect(wrapper.find('#category').exists()).toBe(true)
    expect(wrapper.find('button[type="submit"]').exists()).toBe(true)
  })

  it('calls createInventoryItem on form submission', async () => {
    const push = vi.fn()
    useRouter.mockImplementationOnce(() => ({
      push
    }))

    const wrapper = mount(CreateInventoryItem, {
      global: {
        plugins: [PrimeVue]
      }
    })

    await wrapper.find('#name').setValue('Test Item')
    await wrapper.find('form').trigger('submit.prevent')

    expect(inventoryService.createInventoryItem).toHaveBeenCalled()
    expect(push).toHaveBeenCalledWith({ name: 'inventory' })
  })
})
