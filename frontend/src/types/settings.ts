export interface ContactDetails {
  email: string
  phone?: string
  fax?: string
  mobile?: string
  website?: string
  emergency_phone?: string
}

export interface AddressInfo {
  street?: string
  city?: string
  state?: string
  zip_code?: string
  country?: string
}

export interface EmailSettings {
  smtp_host?: string
  smtp_port?: number
  smtp_username?: string
  smtp_password?: string
  from_email: string
  from_name: string
  reply_to_email?: string
  enable_tls: boolean
  email_signature?: string
}

export interface FoundationSettings {
  id?: string
  name: string
  legal_name?: string
  description?: string
  mission?: string
  vision?: string
  contact_info: ContactDetails
  address?: AddressInfo
  email_settings: EmailSettings
  version?: number
  updated_at?: string
  created_at?: string
}

export interface OrganizationSettingsPayload {
  name: string
  legal_name?: string
  description?: string
  contact_info: ContactDetails
  address?: AddressInfo | null
}

export type EmailSettingsPayload = EmailSettings
