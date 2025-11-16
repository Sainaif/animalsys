import api from './api'
import type { PaginatedResponse } from '@/types/common'
import type {
  User,
  UserFilters,
  CreateUserPayload,
  UpdateUserPayload
} from '@/types/user'

const toPaginatedResponse = (payload: any): PaginatedResponse<User> => {
  const list = Array.isArray(payload?.users)
    ? payload.users
    : Array.isArray(payload?.data)
      ? payload.data
      : Array.isArray(payload)
        ? payload
        : []

  return {
    data: list,
    total: typeof payload?.total === 'number' ? payload.total : list.length,
    limit: typeof payload?.limit === 'number' ? payload.limit : list.length,
    offset: typeof payload?.offset === 'number' ? payload.offset : 0
  }
}

export const userService = {
  async getUsers(params?: UserFilters): Promise<PaginatedResponse<User>> {
    const response = await api.get('/users', { params })
    return toPaginatedResponse(response.data)
  },

  async getUser(id: string): Promise<User> {
    const response = await api.get(`/users/${id}`)
    return response.data
  },

  async createUser(payload: CreateUserPayload): Promise<User> {
    const response = await api.post('/users', payload)
    return response.data
  },

  async updateUser(id: string, payload: UpdateUserPayload): Promise<User> {
    const response = await api.put(`/users/${id}`, payload)
    return response.data
  },

  async deleteUser(id: string): Promise<void> {
    await api.delete(`/users/${id}`)
  },

  async resetPassword(id: string, newPassword: string): Promise<void> {
    await api.put(`/users/${id}/reset-password`, { new_password: newPassword })
  }
}
