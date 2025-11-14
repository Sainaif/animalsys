export interface PaginatedResponse<T> {
  data: T[]
  total: number
  limit: number
  offset: number
}

export interface ApiError {
  error: string
}

export interface QueryParams {
  limit?: number
  offset?: number
  sort_by?: string
  sort_order?: 'asc' | 'desc'
  search?: string
  [key: string]: any
}

export type Status = 'active' | 'inactive' | 'pending' | 'completed' | 'cancelled'
