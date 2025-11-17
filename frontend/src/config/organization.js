const organization = {
  name: import.meta.env.VITE_ORG_NAME || 'Animal Foundation',
  shortName: import.meta.env.VITE_ORG_SHORT_NAME || import.meta.env.VITE_ORG_NAME || 'Animal Foundation',
  legalName: import.meta.env.VITE_ORG_LEGAL_NAME || import.meta.env.VITE_ORG_NAME || 'Animal Foundation',
  tagline: import.meta.env.VITE_ORG_TAGLINE || 'Saving lives, one paw at a time',
  website: import.meta.env.VITE_ORG_WEBSITE || 'https://animalfoundation.org',
  contact: {
    email: import.meta.env.VITE_ORG_EMAIL || 'info@animalfoundation.org',
    phone: import.meta.env.VITE_ORG_PHONE || '+1 (234) 567-890',
    address: import.meta.env.VITE_ORG_ADDRESS || '123 Animal Street, City, State 12345'
  }
}

organization.currentYear = new Date().getFullYear()

export default organization
