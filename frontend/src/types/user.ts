import type { QueryParams } from './common'

export type UserRole = 'super_admin' | 'admin' | 'employee' | 'volunteer' | 'user'
export type UserStatus = 'active' | 'inactive' | 'suspended'

export interface User {
  id: string
  email: string
  first_name: string
  last_name: string
  role: UserRole
  status: UserStatus
  phone?: string
  avatar?: string
  language: 'en' | 'pl'
  theme: 'light' | 'dark'
  last_login?: string
  created_at: string
  updated_at: string
}

export interface UserFilters extends QueryParams {
  role?: UserRole | ''
  status?: UserStatus | ''
  search?: string
}

export interface CreateUserPayload {
  email: string
  password: string
  first_name: string
  last_name: string
  role: UserRole
  status?: UserStatus
  phone?: string
  language: 'en' | 'pl'
  theme: 'light' | 'dark'
}

export interface UpdateUserPayload {
  first_name?: string
  last_name?: string
  role?: UserRole
  status?: UserStatus
  phone?: string
  language?: 'en' | 'pl'
  theme?: 'light' | 'dark'
}
