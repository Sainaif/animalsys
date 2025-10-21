import { beforeEach, vi } from 'vitest'

// Mock localStorage
const localStorageMock = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn(),
}

global.localStorage = localStorageMock

// Reset mocks before each test
beforeEach(() => {
  // Clear all stored data
  Object.keys(localStorageMock).forEach((key) => {
    if (key !== 'getItem' && key !== 'setItem' && key !== 'removeItem' && key !== 'clear') {
      delete localStorageMock[key]
    }
  })

  localStorageMock.getItem.mockImplementation((key) => {
    return localStorageMock[key] !== undefined ? localStorageMock[key] : null
  })
  localStorageMock.setItem.mockImplementation((key, value) => {
    localStorageMock[key] = value
  })
  localStorageMock.removeItem.mockImplementation((key) => {
    delete localStorageMock[key]
  })
  localStorageMock.clear.mockImplementation(() => {
    Object.keys(localStorageMock).forEach((key) => {
      if (key !== 'getItem' && key !== 'setItem' && key !== 'removeItem' && key !== 'clear') {
        delete localStorageMock[key]
      }
    })
  })
})
