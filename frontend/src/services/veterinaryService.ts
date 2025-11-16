import api from './api'
import type { PaginatedResponse, QueryParams } from '@/types/common'
import type {
  VeterinaryVisit,
  Vaccination,
  Medication,
  TreatmentPlan,
  MedicalCondition,
  VeterinaryStatistics
} from '@/types/veterinary'

const toPaginatedResponse = <T>(
  payload: any,
  key: string,
  params?: QueryParams
): PaginatedResponse<T> => {
  const items = Array.isArray(payload?.[key]) ? payload[key] : []
  return {
    data: items,
    total: typeof payload?.total === 'number' ? payload.total : items.length,
    limit: typeof payload?.limit === 'number'
      ? payload.limit
      : (typeof params?.limit === 'number' ? params.limit : items.length),
    offset: typeof payload?.offset === 'number'
      ? payload.offset
      : (typeof params?.offset === 'number' ? params.offset : 0)
  }
}

const toList = <T>(payload: any, key: string): T[] => {
  if (Array.isArray(payload?.[key])) {
    return payload[key]
  }
  if (Array.isArray(payload)) {
    return payload
  }
  return []
}

const extractItem = <T>(payload: any, key: string): T => {
  if (payload && typeof payload === 'object' && key in payload) {
    return payload[key]
  }
  return payload as T
}

export const veterinaryService = {
  // Veterinary Visits
  async getVisits(params?: QueryParams): Promise<PaginatedResponse<VeterinaryVisit>> {
    const response = await api.get('/veterinary/visits', { params })
    return toPaginatedResponse<VeterinaryVisit>(response.data, 'visits', params)
  },

  async getVisit(id: string): Promise<VeterinaryVisit> {
    const response = await api.get(`/veterinary/visits/${id}`)
    return extractItem<VeterinaryVisit>(response.data, 'visit')
  },

  async createVisit(data: Partial<VeterinaryVisit>): Promise<VeterinaryVisit> {
    const response = await api.post('/veterinary/visits', data)
    return response.data
  },

  async updateVisit(id: string, data: Partial<VeterinaryVisit>): Promise<VeterinaryVisit> {
    const response = await api.put(`/veterinary/visits/${id}`, data)
    return response.data
  },

  async deleteVisit(id: string): Promise<void> {
    await api.delete(`/veterinary/visits/${id}`)
  },

  async getAnimalVisits(animalId: string): Promise<VeterinaryVisit[]> {
    const response = await api.get(`/animals/${animalId}/veterinary-visits`)
    return toList<VeterinaryVisit>(response.data, 'visits')
  },

  // Vaccinations
  async getVaccinations(params?: QueryParams): Promise<PaginatedResponse<Vaccination>> {
    const response = await api.get('/veterinary/vaccinations', { params })
    return toPaginatedResponse<Vaccination>(response.data, 'vaccinations', params)
  },

  async getVaccination(id: string): Promise<Vaccination> {
    const response = await api.get(`/veterinary/vaccinations/${id}`)
    return extractItem<Vaccination>(response.data, 'vaccination')
  },

  async createVaccination(data: Partial<Vaccination>): Promise<Vaccination> {
    const response = await api.post('/veterinary/vaccinations', data)
    return response.data
  },

  async updateVaccination(id: string, data: Partial<Vaccination>): Promise<Vaccination> {
    const response = await api.put(`/veterinary/vaccinations/${id}`, data)
    return response.data
  },

  async deleteVaccination(id: string): Promise<void> {
    await api.delete(`/veterinary/vaccinations/${id}`)
  },

  async getAnimalVaccinations(animalId: string): Promise<Vaccination[]> {
    const response = await api.get(`/animals/${animalId}/vaccinations`)
    return toList<Vaccination>(response.data, 'vaccinations')
  },

  // Medications
  async getMedications(params?: QueryParams): Promise<PaginatedResponse<Medication>> {
    const response = await api.get('/veterinary/medications', { params })
    return toPaginatedResponse<Medication>(response.data, 'medications', params)
  },

  async getMedication(id: string): Promise<Medication> {
    const response = await api.get(`/veterinary/medications/${id}`)
    return extractItem<Medication>(response.data, 'medication')
  },

  async createMedication(data: Partial<Medication>): Promise<Medication> {
    const response = await api.post('/veterinary/medications', data)
    return response.data
  },

  async updateMedication(id: string, data: Partial<Medication>): Promise<Medication> {
    const response = await api.put(`/veterinary/medications/${id}`, data)
    return response.data
  },

  async deleteMedication(id: string): Promise<void> {
    await api.delete(`/veterinary/medications/${id}`)
  },

  async getAnimalMedications(animalId: string): Promise<Medication[]> {
    const response = await api.get(`/animals/${animalId}/medications`)
    return toList<Medication>(response.data, 'medications')
  },

  async discontinueMedication(id: string, reason: string): Promise<Medication> {
    const response = await api.post(`/veterinary/medications/${id}/discontinue`, { reason })
    return response.data
  },

  // Treatment Plans
  async getTreatmentPlans(params?: QueryParams): Promise<PaginatedResponse<TreatmentPlan>> {
    const response = await api.get('/veterinary/treatment-plans', { params })
    return toPaginatedResponse<TreatmentPlan>(response.data, 'treatment_plans', params)
  },

  async getTreatmentPlan(id: string): Promise<TreatmentPlan> {
    const response = await api.get(`/veterinary/treatment-plans/${id}`)
    return extractItem<TreatmentPlan>(response.data, 'treatment_plan')
  },

  async createTreatmentPlan(data: Partial<TreatmentPlan>): Promise<TreatmentPlan> {
    const response = await api.post('/veterinary/treatment-plans', data)
    return response.data
  },

  async updateTreatmentPlan(id: string, data: Partial<TreatmentPlan>): Promise<TreatmentPlan> {
    const response = await api.put(`/veterinary/treatment-plans/${id}`, data)
    return response.data
  },

  async deleteTreatmentPlan(id: string): Promise<void> {
    await api.delete(`/veterinary/treatment-plans/${id}`)
  },

  async getAnimalTreatmentPlans(animalId: string): Promise<TreatmentPlan[]> {
    const response = await api.get(`/animals/${animalId}/treatment-plans`)
    return toList<TreatmentPlan>(response.data, 'treatment_plans')
  },

  async addProgressNote(
    planId: string,
    note: { date: string; note: string; recorded_by: string }
  ): Promise<TreatmentPlan> {
    const response = await api.post(`/veterinary/treatment-plans/${planId}/progress-notes`, note)
    return response.data
  },

  // Medical Conditions
  async getMedicalConditions(params?: QueryParams): Promise<PaginatedResponse<MedicalCondition>> {
    const response = await api.get('/veterinary/medical-conditions', { params })
    return toPaginatedResponse<MedicalCondition>(response.data, 'conditions', params)
  },

  async getMedicalCondition(id: string): Promise<MedicalCondition> {
    const response = await api.get(`/veterinary/medical-conditions/${id}`)
    return extractItem<MedicalCondition>(response.data, 'condition')
  },

  async createMedicalCondition(data: Partial<MedicalCondition>): Promise<MedicalCondition> {
    const response = await api.post('/veterinary/medical-conditions', data)
    return response.data
  },

  async updateMedicalCondition(
    id: string,
    data: Partial<MedicalCondition>
  ): Promise<MedicalCondition> {
    const response = await api.put(`/veterinary/medical-conditions/${id}`, data)
    return response.data
  },

  async deleteMedicalCondition(id: string): Promise<void> {
    await api.delete(`/veterinary/medical-conditions/${id}`)
  },

  async getAnimalMedicalConditions(animalId: string): Promise<MedicalCondition[]> {
    const response = await api.get(`/animals/${animalId}/medical-conditions`)
    return toList<MedicalCondition>(response.data, 'conditions')
  },

  // Statistics
  async getStatistics(): Promise<VeterinaryStatistics> {
    const response = await api.get('/veterinary/statistics')
    return response.data
  }
}
