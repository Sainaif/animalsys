import api from './api'
import type { Animal, AnimalStatistics, VeterinaryVisit, Vaccination } from '@/types/animal'
import type { PaginatedResponse, QueryParams } from '@/types/common'

const normalizeListResponse = (
  payload: any,
  params?: QueryParams
): PaginatedResponse<Animal> => {
  const animals = Array.isArray(payload?.animals)
    ? payload.animals
    : Array.isArray(payload?.data)
      ? payload.data
      : Array.isArray(payload)
        ? payload
        : []

  return {
    data: animals,
    total: typeof payload?.total === 'number' ? payload.total : animals.length,
    limit: typeof payload?.limit === 'number'
      ? payload.limit
      : (typeof params?.limit === 'number' ? params.limit : animals.length),
    offset: typeof payload?.offset === 'number'
      ? payload.offset
      : (typeof params?.offset === 'number' ? params.offset : 0)
  }
}

export const animalService = {
  // List animals
  async getAnimals(params?: QueryParams): Promise<PaginatedResponse<Animal>> {
    const response = await api.get('/animals', { params })
    return normalizeListResponse(response.data, params)
  },

  // Get animal by ID
  async getAnimal(id: string): Promise<Animal> {
    const response = await api.get(`/animals/${id}`)
    return response.data
  },

  // Create animal
  async createAnimal(data: Partial<Animal>): Promise<Animal> {
    const response = await api.post('/animals', data)
    return response.data
  },

  // Update animal
  async updateAnimal(id: string, data: Partial<Animal>): Promise<Animal> {
    const response = await api.put(`/animals/${id}`, data)
    return response.data
  },

  // Delete animal
  async deleteAnimal(id: string): Promise<void> {
    await api.delete(`/animals/${id}`)
  },

  // Get statistics
  async getStatistics(): Promise<AnimalStatistics> {
    const response = await api.get('/animals/statistics')
    return response.data
  },

  // Get veterinary visits
  async getVeterinaryVisits(animalId: string): Promise<VeterinaryVisit[]> {
    const response = await api.get(`/animals/${animalId}/visits`)
    return response.data
  },

  // Get vaccinations
  async getVaccinations(animalId: string): Promise<Vaccination[]> {
    const response = await api.get(`/animals/${animalId}/vaccinations`)
    return response.data
  },

  // Upload photo
  async uploadPhoto(animalId: string, file: File): Promise<{ photo_url: string }> {
    const formData = new FormData()
    formData.append('photo', file)
    const response = await api.post(`/animals/${animalId}/photos`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
    return response.data
  }
}
