export interface Event {
  id: string
  name: string
  event_type: 'fundraiser' | 'adoption_event' | 'education' | 'volunteer' | 'community' | 'other'
  description?: string
  start_date: string
  end_date?: string
  location?: string
  organizer_name?: string
  organizer_email?: string
  max_participants?: number
  registered_participants?: number
  volunteers_needed?: number
  volunteers_assigned?: number
  status: 'planned' | 'active' | 'completed' | 'cancelled'
  budget?: number
  raised_amount?: number
  notes?: string
  created_at?: string
  updated_at?: string
}

export interface Volunteer {
  id: string
  first_name: string
  last_name: string
  email: string
  phone?: string
  address?: {
    street: string
    city: string
    state: string
    postal_code: string
    country: string
  }
  date_of_birth?: string
  emergency_contact_name?: string
  emergency_contact_phone?: string
  volunteer_status: 'active' | 'inactive' | 'pending'
  skills?: string[]
  availability?: {
    monday: boolean
    tuesday: boolean
    wednesday: boolean
    thursday: boolean
    friday: boolean
    saturday: boolean
    sunday: boolean
  }
  preferred_roles?: string[]
  total_hours?: number
  start_date?: string
  background_check_status?: 'pending' | 'approved' | 'rejected'
  background_check_date?: string
  orientation_completed: boolean
  orientation_date?: string
  notes?: string
  created_at?: string
  updated_at?: string
}

export interface Shift {
  id: string
  event_id?: string
  event?: {
    id: string
    name: string
  }
  volunteer_id?: string
  volunteer?: {
    id: string
    first_name: string
    last_name: string
  }
  shift_date: string
  start_time: string
  end_time: string
  role: string
  hours?: number
  status: 'scheduled' | 'completed' | 'cancelled' | 'no_show'
  notes?: string
  created_at?: string
  updated_at?: string
}

export interface EventStatistics {
  total_events: number
  upcoming_events: number
  active_volunteers: number
  total_volunteer_hours: number
  events_this_month: number
  volunteers_this_month: number
}
