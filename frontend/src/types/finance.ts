export interface Donor {
  id: string
  donor_type: 'individual' | 'organization' | 'corporate' | 'foundation'
  first_name?: string
  last_name?: string
  organization_name?: string
  email: string
  phone?: string
  address?: {
    street: string
    city: string
    state: string
    postal_code: string
    country: string
  }
  total_donated?: number
  donation_count?: number
  first_donation_date?: string
  last_donation_date?: string
  donor_status: 'active' | 'inactive' | 'lapsed'
  communication_preferences?: {
    email: boolean
    phone: boolean
    mail: boolean
  }
  tags?: string[]
  notes?: string
  created_at?: string
  updated_at?: string
}

export interface Donation {
  id: string
  donor_id: string
  donor?: {
    id: string
    first_name?: string
    last_name?: string
    organization_name?: string
    email: string
  }
  campaign_id?: string
  campaign?: {
    id: string
    name: string
  }
  amount: number
  currency: string
  donation_date: string
  donation_type: 'one_time' | 'recurring' | 'pledge'
  payment_method: 'cash' | 'credit_card' | 'debit_card' | 'bank_transfer' | 'check' | 'paypal' | 'other'
  payment_status: 'pending' | 'completed' | 'failed' | 'refunded'
  transaction_id?: string
  receipt_sent: boolean
  receipt_number?: string
  tax_deductible: boolean
  purpose?: string
  anonymous: boolean
  recurring_frequency?: 'weekly' | 'monthly' | 'quarterly' | 'annually'
  notes?: string
  created_at?: string
  updated_at?: string
}

type LocalizedText = string | Record<string, string>

export interface Campaign {
  id: string
  name: LocalizedText
  description?: LocalizedText
  campaign_type?: 'fundraising' | 'awareness' | 'adoption' | 'event' | 'other' | 'general' | 'capital' | 'emergency' | 'annual' | 'membership' | 'end_of_year'
  type?: 'fundraising' | 'awareness' | 'adoption' | 'event' | 'other' | 'general' | 'capital' | 'emergency' | 'annual' | 'membership' | 'end_of_year'
  start_date: string
  end_date?: string
  goal_amount?: number
  goal?: number
  current_amount?: number
  currentAmount?: number
  raised_amount?: number
  currency?: string
  status: 'planning' | 'active' | 'completed' | 'cancelled' | 'draft' | 'paused'
  target_audience?: string
  coordinator_name?: string
  coordinator_email?: string
  donation_count?: number
  donor_count?: number
  success_metrics?: string
  notes?: string
  public?: boolean
  featured?: boolean
  manager?: string
  created_at?: string
  updated_at?: string
}

export interface FinanceStatistics {
  total_donors: number
  active_donors: number
  total_donations: number
  total_amount: number
  donations_this_month: number
  amount_this_month: number
  active_campaigns: number
  recurring_donors: number
}
