import { mount } from '@vue/test-utils'
import { describe, it, expect, vi } from 'vitest'
import CreateInventoryItem from './CreateInventoryItem.vue'
import { inventoryService } from '@/services/inventoryService'
import { useRouter } from 'vue-router'

// Mock the router
vi.mock('vue-router', () => ({
  useRouter: vi.fn(() => ({
    push: vi.fn()
  }))
}))

// Mock the inventoryService
vi.mock('@/services/inventoryService', () => ({
  inventoryService: {
    createInventoryItem: vi.fn()
  }
}))

describe('CreateInventoryItem.vue', () => {
  it('renders the form correctly', () => {
    const wrapper = mount(CreateInventoryItem)
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

    const wrapper = mount(CreateInventoryItem)

    await wrapper.find('#name').setValue('Test Item')
    await wrapper.find('form').trigger('submit.prevent')

    expect(inventoryService.createInventoryItem).toHaveBeenCalled()
    expect(push).toHaveBeenCalledWith({ name: 'Inventory' })
  })
})
