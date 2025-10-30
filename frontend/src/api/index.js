import { api } from './client'
import { authApi } from './modules/auth'
import { animalsApi } from './modules/animals'
import { adoptionsApi } from './modules/adoptions'
import { volunteersApi } from './modules/volunteers'
import { financeApi } from './modules/finance'
import { donorsApi } from './modules/donors'
import { inventoryApi } from './modules/inventory'
import { veterinaryApi } from './modules/veterinary'
import { campaignsApi } from './modules/campaigns'
import { partnersApi } from './modules/partners'
import { schedulesApi } from './modules/schedules'
import { documentsApi } from './modules/documents'
import { communicationsApi } from './modules/communications'
import { reportsApi } from './modules/reports'

// Export main client
export { api }

// Export all API modules
export const API = {
  auth: authApi,
  animals: animalsApi,
  adoptions: adoptionsApi,
  volunteers: volunteersApi,
  finance: financeApi,
  donors: donorsApi,
  inventory: inventoryApi,
  veterinary: veterinaryApi,
  campaigns: campaignsApi,
  partners: partnersApi,
  schedules: schedulesApi,
  documents: documentsApi,
  communications: communicationsApi,
  reports: reportsApi,
  // Other modules can be added here as they're implemented
}

export default API
