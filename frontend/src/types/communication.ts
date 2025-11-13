export interface EmailTemplate {
  id?: number
  name: string
  subject: string
  body: string
  template_type: 'adoption' | 'donation' | 'event' | 'general' | 'newsletter'
  variables: string[]
  is_active: boolean
  created_at?: string
  updated_at?: string
}

export interface EmailCampaign {
  id?: number
  name: string
  template_id: number
  template?: EmailTemplate
  recipient_type: 'donors' | 'adopters' | 'volunteers' | 'all' | 'custom'
  recipient_filter?: Record<string, any>
  scheduled_date?: string
  sent_date?: string
  status: 'draft' | 'scheduled' | 'sending' | 'sent' | 'failed'
  total_recipients?: number
  sent_count?: number
  opened_count?: number
  clicked_count?: number
  created_by?: number
  created_at?: string
  updated_at?: string
}

export interface CommunicationLog {
  id?: number
  communication_type: 'email' | 'phone' | 'sms' | 'in_person' | 'other'
  subject: string
  message?: string
  recipient_type: 'donor' | 'adopter' | 'volunteer' | 'contact' | 'staff'
  recipient_id?: number
  recipient_name?: string
  sender_name?: string
  communication_date: string
  status: 'sent' | 'delivered' | 'failed' | 'pending'
  notes?: string
  created_by?: number
  created_at?: string
}

export interface SMSTemplate {
  id?: number
  name: string
  message: string
  template_type: 'reminder' | 'notification' | 'alert' | 'general'
  variables: string[]
  is_active: boolean
  created_at?: string
  updated_at?: string
}
