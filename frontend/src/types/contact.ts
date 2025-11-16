import type { QueryParams } from './common'

export type ContactType = 'adopter' | 'donor' | 'volunteer' | 'partner' | 'vendor' | 'other'
export type ContactStatus = 'active' | 'inactive' | 'prospect' | 'archived'

export interface ContactAddress {
  street?: string
  city?: string
  state?: string
  postal_code?: string
  country?: string
}

export interface ContactOwner {
  id: string
  name: string
}

export interface ContactActivity {
  id: string
  type: 'call' | 'email' | 'meeting' | 'note'
  subject: string
  description?: string
  created_at: string
  created_by?: ContactOwner
  outcome?: string
}

export interface Contact {
  id: string
  first_name: string
  last_name: string
  organization?: string
  email?: string
  phone?: string
  type: ContactType
  status: ContactStatus
  tags: string[]
  owner?: ContactOwner
  preferred_channel?: 'email' | 'phone' | 'sms'
  last_contacted_at?: string
  next_follow_up_at?: string
  notes?: string
  address?: ContactAddress
  activities?: ContactActivity[]
  created_at: string
  updated_at: string
}

export interface ContactFilters extends QueryParams {
  type?: ContactType | ''
  status?: ContactStatus | ''
  owner_id?: string
}

export interface ContactPayload {
  first_name: string
  last_name: string
  organization?: string
  email?: string
  phone?: string
  type: ContactType
  status: ContactStatus
  tags?: string[]
  owner_id?: string
  preferred_channel?: 'email' | 'phone' | 'sms'
  next_follow_up_at?: string | null
  notes?: string
  address?: ContactAddress
}
