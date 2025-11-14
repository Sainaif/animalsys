export interface Animal {
  id: string
  name: string
  name_en?: string
  species: string
  breed?: string
  category: string
  sex: 'male' | 'female' | 'unknown'
  date_of_birth?: string
  age_years?: number
  age_months?: number
  color?: string
  size?: 'small' | 'medium' | 'large' | 'extra_large'
  weight?: number
  microchip_id?: string
  status: 'available' | 'adopted' | 'under_treatment' | 'fostered' | 'transferred' | 'deceased'
  intake_date: string
  intake_reason?: string
  description?: string
  description_en?: string
  medical_history?: string
  temperament?: string[]
  good_with_kids?: boolean
  good_with_dogs?: boolean
  good_with_cats?: boolean
  house_trained?: boolean
  spayed_neutered?: boolean
  vaccinated?: boolean
  special_needs?: string
  adoption_fee?: number
  adoption_requirements?: string
  photo_url?: string
  photos?: string[]
  daily_notes?: DailyNote[]
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface DailyNote {
  date: string
  note: string
  staff_member: string
}

export interface AnimalStatistics {
  total_animals: number
  by_status: {
    [key: string]: number
  }
  by_species: {
    [key: string]: number
  }
  by_category: {
    [key: string]: number
  }
}

export interface VeterinaryVisit {
  id: string
  animal_id: string
  visit_date: string
  visit_type: string
  veterinarian_name: string
  reason: string
  diagnosis?: string
  treatment?: string
  medications?: string
  follow_up_required: boolean
  follow_up_date?: string
  cost?: number
  created_at: string
}

export interface Vaccination {
  id: string
  animal_id: string
  vaccine_name: string
  vaccine_type: string
  vaccination_date: string
  next_due_date?: string
  administered_by: string
  batch_number?: string
  created_at: string
}
