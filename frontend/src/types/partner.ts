export interface Partner {
  id?: number
  organization_name: string
  partner_type: 'shelter' | 'rescue' | 'veterinary' | 'foster' | 'other'
  contact_person: string
  email: string
  phone?: string
  address?: {
    street: string
    city: string
    state: string
    postal_code: string
    country: string
  }
  website?: string
  license_number?: string
  capacity?: number
  specialization?: string[]
  status: 'active' | 'inactive' | 'pending'
  partnership_start_date?: string
  partnership_end_date?: string
  notes?: string
  created_at?: string
  updated_at?: string
}

export interface AnimalTransfer {
  id?: number
  animal_id: number
  animal?: any
  from_organization_id?: number
  from_organization?: Partner
  to_organization_id?: number
  to_organization?: Partner
  transfer_date: string
  transfer_reason: string
  transfer_type: 'incoming' | 'outgoing' | 'temporary'
  expected_return_date?: string
  actual_return_date?: string
  transport_method?: string
  transport_cost?: number
  health_certificate_provided: boolean
  microchip_transferred: boolean
  status: 'pending' | 'in_transit' | 'completed' | 'cancelled'
  notes?: string
  documents?: string[]
  created_by?: number
  created_at?: string
  updated_at?: string
}

export interface PartnerAgreement {
  id?: number
  partner_id: number
  partner?: Partner
  agreement_type: 'mou' | 'contract' | 'transfer_agreement' | 'foster_agreement'
  start_date: string
  end_date?: string
  terms: string
  signed_by_partner: boolean
  signed_by_us: boolean
  signed_date?: string
  document_url?: string
  status: 'draft' | 'active' | 'expired' | 'terminated'
  renewal_date?: string
  created_at?: string
  updated_at?: string
}
