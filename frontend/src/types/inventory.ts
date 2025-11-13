export interface InventoryItem {
  id?: number
  name: string
  category: 'food' | 'medicine' | 'supplies' | 'equipment' | 'other'
  description?: string
  sku?: string
  unit: string
  quantity_in_stock: number
  minimum_quantity: number
  maximum_quantity?: number
  unit_cost?: number
  total_value?: number
  supplier?: string
  location?: string
  expiration_date?: string
  last_restocked_date?: string
  status: 'in_stock' | 'low_stock' | 'out_of_stock' | 'expired'
  created_at?: string
  updated_at?: string
}

export interface StockTransaction {
  id?: number
  item_id: number
  item?: InventoryItem
  transaction_type: 'purchase' | 'usage' | 'donation' | 'disposal' | 'adjustment'
  quantity: number
  unit_cost?: number
  total_cost?: number
  transaction_date: string
  reference_number?: string
  notes?: string
  created_by?: number
  created_at?: string
}
