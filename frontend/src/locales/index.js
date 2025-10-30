/**
 * Translation System - Language Registry
 *
 * HOW TO ADD A NEW LANGUAGE:
 * 1. Copy _example.json to {language_code}.json (e.g., de.json for German)
 * 2. Translate all values in the new file
 * 3. Add the language code to the availableLanguages array below
 * 4. Import the translation file and add it to the messages object
 *
 * Example:
 * import de from './de.json'
 * export const availableLanguages = ['en', 'pl', 'de']
 * export const messages = { en, pl, de }
 */

// Import all available language files
import en from './en.json'
import pl from './pl.json'

/**
 * List of all available languages in the system
 * Add new language codes here when adding new translations
 */
export const availableLanguages = ['en', 'pl']

/**
 * Language metadata for UI display
 */
export const languageInfo = {
  en: {
    code: 'en',
    name: 'English',
    nativeName: 'English',
    flag: '🇬🇧'
  },
  pl: {
    code: 'pl',
    name: 'Polish',
    nativeName: 'Polski',
    flag: '🇵🇱'
  }
  // Add new language metadata here
  // de: {
  //   code: 'de',
  //   name: 'German',
  //   nativeName: 'Deutsch',
  //   flag: '🇩🇪'
  // }
}

/**
 * All translation messages
 * Add new language imports here
 */
export const messages = {
  en,
  pl
  // Add new language here
  // de
}

/**
 * Default language fallback
 */
export const defaultLanguage = 'en'

/**
 * Get browser language or default
 */
export const getDefaultLanguage = () => {
  const browserLang = navigator.language.split('-')[0]
  return availableLanguages.includes(browserLang) ? browserLang : defaultLanguage
}
