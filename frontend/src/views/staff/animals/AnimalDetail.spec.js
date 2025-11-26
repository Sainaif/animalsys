
import { shallowMount, flushPromises } from '@vue/test-utils'
import { describe, it, expect, vi } from 'vitest'
import AnimalDetail from './AnimalDetail.vue'

// Mock dependencies
vi.mock('vue-router', () => ({
  useRoute: () => ({
    params: { id: '1' }
  }),
  useRouter: () => ({
    back: vi.fn(),
    push: vi.fn()
  })
}))
vi.mock('@/services/animalService', () => ({
  animalService: {
    getAnimal: vi.fn().mockResolvedValue({ id: '1', sex: null })
  }
}))
vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (key) => key,
    locale: 'en'
  })
}))

vi.mock('primevue/usetoast', () => ({
  useToast: () => ({
    add: vi.fn()
  })
}))

vi.mock('primevue/useconfirm', () => ({
  useConfirm: () => ({
    require: vi.fn()
  })
}))

describe('AnimalDetail.vue', () => {
  it('displays "unknown" when gender is null', async () => {
    const wrapper = shallowMount(AnimalDetail, {
      global: {
        mocks: {
          $t: (key) => key
        },
        stubs: {
          LoadingSpinner: true,
          Button: true,
          Badge: true,
          Card: true,
          TabView: true,
          TabPanel: true,
          ConfirmDialog: true
        }
      }
    })
    // Wait for the component to mount and fetch data
    await flushPromises()

    expect(wrapper.html()).toContain('animal.unknown')
  })
})
