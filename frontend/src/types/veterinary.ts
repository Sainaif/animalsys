export interface VeterinaryVisit {
  id: string
  animal_id: string
  animal?: {
    id: string
    name: string
    species: string
  }
  visit_date: string
  visit_type: 'checkup' | 'emergency' | 'surgery' | 'vaccination' | 'follow_up' | 'other'
  veterinarian_name: string
  clinic_name?: string
  reason: string
  diagnosis?: string
  treatment_provided?: string
  medications_prescribed?: string
  follow_up_required: boolean
  follow_up_date?: string
  weight?: number
  temperature?: number
  heart_rate?: number
  notes?: string
  cost?: number
  attachments?: string[]
  created_at?: string
  updated_at?: string
}

export interface Vaccination {
  id: string
  animal_id: string
  animal?: {
    id: string
    name: string
    species: string
  }
  vaccine_name: string
  vaccine_type: string
  manufacturer?: string
  batch_number?: string
  vaccination_date: string
  next_due_date?: string
  veterinarian_name: string
  clinic_name?: string
  site_of_injection?: string
  dosage?: string
  adverse_reactions?: string
  notes?: string
  created_at?: string
  updated_at?: string
}

export interface Medication {
  id: string
  animal_id: string
  animal?: {
    id: string
    name: string
    species: string
  }
  medication_name: string
  medication_type: 'antibiotic' | 'pain_relief' | 'anti_inflammatory' | 'antiparasitic' | 'supplement' | 'other'
  dosage: string
  frequency: string
  route: 'oral' | 'topical' | 'injection' | 'intravenous' | 'other'
  start_date: string
  end_date?: string
  prescribed_by: string
  reason: string
  instructions?: string
  side_effects?: string
  status: 'active' | 'completed' | 'discontinued'
  notes?: string
  created_at?: string
  updated_at?: string
}

export interface TreatmentPlan {
  id: string
  animal_id: string
  animal?: {
    id: string
    name: string
    species: string
  }
  plan_name: string
  condition: string
  start_date: string
  end_date?: string
  veterinarian_name: string
  goals: string
  medications?: string[]
  procedures?: string[]
  dietary_requirements?: string
  exercise_requirements?: string
  monitoring_schedule?: string
  progress_notes?: Array<{
    date: string
    note: string
    recorded_by: string
  }>
  status: 'active' | 'completed' | 'cancelled'
  created_at?: string
  updated_at?: string
}

export interface MedicalCondition {
  id: string
  animal_id: string
  animal?: {
    id: string
    name: string
    species: string
  }
  condition_name: string
  diagnosis_date: string
  diagnosed_by: string
  severity: 'mild' | 'moderate' | 'severe' | 'critical'
  status: 'active' | 'managed' | 'resolved' | 'chronic'
  symptoms?: string
  treatment?: string
  prognosis?: string
  notes?: string
  created_at?: string
  updated_at?: string
}

export interface VeterinaryStatistics {
  total_visits: number
  upcoming_vaccinations: number
  active_medications: number
  active_treatment_plans: number
  active_conditions: number
  visits_this_month: number
  vaccinations_due_soon: number
}
