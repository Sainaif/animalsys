export interface AdoptionApplication {
  id: string
  animal_id: string
  animal?: {
    id: string
    name: string
    species: string
    photo_url?: string
  }
  applicant_first_name: string
  applicant_last_name: string
  email: string
  phone: string
  address: {
    street: string
    city: string
    state: string
    postal_code: string
    country: string
  }
  date_of_birth: string
  occupation: string
  household_type: 'house' | 'apartment' | 'condo' | 'other'
  has_yard: boolean
  yard_fenced: boolean
  other_pets: Array<{
    type: string
    breed?: string
    age: number
    spayed_neutered: boolean
  }>
  household_members: Array<{
    name: string
    age: number
    relationship: string
  }>
  experience_with_pets: string
  years_of_experience?: number
  previous_adoptions?: string
  veterinarian_info?: {
    name: string
    clinic: string
    phone: string
  }
  references: Array<{
    name: string
    relationship: string
    phone: string
    email?: string
  }>
  employment_status: string
  work_hours?: string
  who_will_care: string
  reason_for_adoption: string
  expectations: string
  application_date: string
  status: 'pending' | 'under_review' | 'approved' | 'rejected' | 'withdrawn'
  reviewed_by?: string
  review_date?: string
  review_notes?: string
  created_at: string
  updated_at: string
}

export interface Adoption {
  id: string
  animal_id: string
  application_id: string
  adopter_first_name: string
  adopter_last_name: string
  adopter_email: string
  adopter_phone: string
  adoption_date: string
  adoption_fee: number
  payment_status: 'pending' | 'partial' | 'completed' | 'refunded'
  payment_method?: string
  contract_signed: boolean
  contract_signed_date?: string
  microchip_transferred: boolean
  return_policy_explained: boolean
  status: 'active' | 'returned' | 'completed'
  follow_up_required: boolean
  follow_up_schedule?: Array<{
    date: string
    type: string
    completed: boolean
    notes?: string
  }>
  return_date?: string
  return_reason?: string
  notes?: string
  created_at: string
  updated_at: string
}

export interface AdoptionStatistics {
  total_adoptions: number
  this_month: number
  this_year: number
  pending_applications: number
  average_time_to_adopt_days: number
}
