import axios from 'axios'
import type { InventoryItem, StockTransaction } from '@/types/inventory'

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:3000/api'

export const inventoryService = {
  // Inventory Items
  getInventoryItems: (params?: any) => axios.get(`${API_URL}/inventory-items`, { params }),
  getInventoryItem: (id: number) => axios.get(`${API_URL}/inventory-items/${id}`),
  createInventoryItem: (data: InventoryItem) => axios.post(`${API_URL}/inventory-items`, data),
  updateInventoryItem: (id: number, data: InventoryItem) => axios.put(`${API_URL}/inventory-items/${id}`, data),
  deleteInventoryItem: (id: number) => axios.delete(`${API_URL}/inventory-items/${id}`),

  // Stock Transactions
  getStockTransactions: (params?: any) => axios.get(`${API_URL}/stock-transactions`, { params }),
  getStockTransaction: (id: number) => axios.get(`${API_URL}/stock-transactions/${id}`),
  createStockTransaction: (data: StockTransaction) => axios.post(`${API_URL}/stock-transactions`, data),
  deleteStockTransaction: (id: number) => axios.delete(`${API_URL}/stock-transactions/${id}`)
}
