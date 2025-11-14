export const exportService = {
  /**
   * Export data to CSV format
   */
  exportToCSV(data: any[], filename: string, columns?: { field: string; header: string }[]) {
    if (!data || data.length === 0) {
      throw new Error('No data to export')
    }

    // If columns not specified, use all keys from first object
    const cols = columns || Object.keys(data[0]).map(key => ({ field: key, header: key }))

    // Create CSV header
    const headers = cols.map(col => col.header).join(',')

    // Create CSV rows
    const rows = data.map(item => {
      return cols.map(col => {
        const value = this.getNestedValue(item, col.field)
        // Escape and quote values that contain commas or quotes
        if (value === null || value === undefined) return ''
        const stringValue = String(value)
        if (stringValue.includes(',') || stringValue.includes('"') || stringValue.includes('\n')) {
          return `"${stringValue.replace(/"/g, '""')}"`
        }
        return stringValue
      }).join(',')
    })

    // Combine header and rows
    const csv = [headers, ...rows].join('\n')

    // Create and download file
    this.downloadFile(csv, filename, 'text/csv')
  },

  /**
   * Get nested object value by dot notation
   */
  getNestedValue(obj: any, path: string): any {
    return path.split('.').reduce((current, key) => current?.[key], obj)
  },

  /**
   * Download file to browser
   */
  downloadFile(content: string, filename: string, mimeType: string) {
    const blob = new Blob([content], { type: mimeType })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = filename
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(url)
  },

  /**
   * Export data to JSON format
   */
  exportToJSON(data: any[], filename: string) {
    const json = JSON.stringify(data, null, 2)
    this.downloadFile(json, filename, 'application/json')
  }
}
