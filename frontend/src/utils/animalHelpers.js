export const getLocalizedValue = (value, locale) => {
  if (!value) {
    return ''
  }

  if (typeof value === 'string') {
    return value
  }

  if (typeof value === 'object') {
    if (locale && value[locale]) {
      return value[locale]
    }
    return value.en || value.pl || Object.values(value)[0] || ''
  }

  return ''
}

export const getAnimalImage = (animal) => {
  if (!animal) {
    return ''
  }

  if (animal.images?.primary) {
    return animal.images.primary
  }

  if (animal.photo_url) {
    return animal.photo_url
  }

  if (animal.images?.gallery?.length) {
    return animal.images.gallery[0]
  }

  return 'https://placehold.co/600x400?text=Animal'
}

export const getAnimalBreed = (animal) => {
  if (!animal) {
    return ''
  }

  return animal.breed || animal.species || ''
}

export const getAnimalGender = (animal) => {
  if (!animal) {
    return 'unknown'
  }

  return animal.sex || animal.gender || 'unknown'
}

const normalizeKey = (value) => value.toLowerCase().replace(/[^a-z0-9]+/g, '_')

export const translateValue = (value, t, baseKey) => {
  if (!value || typeof value !== 'string' || !t) {
    return value || ''
  }

  const normalized = normalizeKey(value)
  const translationKey = `${baseKey}.${normalized}`
  const translated = t(translationKey)
  return translated === translationKey ? value : translated
}
