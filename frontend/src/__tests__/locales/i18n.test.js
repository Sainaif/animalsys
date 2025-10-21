import { describe, it, expect, beforeEach } from 'vitest'
import i18n from '../../locales'

describe('i18n Configuration', () => {
  it('should have polish locale configured as default', () => {
    expect(i18n.global.locale.value).toBe('pl')
  })

  it('should have fallback locale set to english', () => {
    expect(i18n.global.fallbackLocale.value).toBe('en')
  })

  describe('Polish translations', () => {
    beforeEach(() => {
      i18n.global.locale.value = 'pl'
    })

    it('should have navigation translations', () => {
      expect(i18n.global.t('nav.home')).toBe('Strona główna')
      expect(i18n.global.t('nav.animals')).toBe('Zwierzęta')
      expect(i18n.global.t('nav.logout')).toBe('Wyloguj')
    })

    it('should have login translations', () => {
      expect(i18n.global.t('login.title')).toBe('Logowanie')
      expect(i18n.global.t('login.username')).toBe('Nazwa użytkownika')
      expect(i18n.global.t('login.password')).toBe('Hasło')
    })

    it('should have animals translations', () => {
      expect(i18n.global.t('animals.title')).toBe('Zwierzęta')
      expect(i18n.global.t('animals.name')).toBe('Imię')
      expect(i18n.global.t('animals.species')).toBe('Gatunek')
    })
  })

  describe('English translations', () => {
    beforeEach(() => {
      i18n.global.locale.value = 'en'
    })

    it('should have navigation translations', () => {
      expect(i18n.global.t('nav.home')).toBe('Home')
      expect(i18n.global.t('nav.animals')).toBe('Animals')
      expect(i18n.global.t('nav.logout')).toBe('Logout')
    })

    it('should have login translations', () => {
      expect(i18n.global.t('login.title')).toBe('Login')
      expect(i18n.global.t('login.username')).toBe('Username')
      expect(i18n.global.t('login.password')).toBe('Password')
    })

    it('should have animals translations', () => {
      expect(i18n.global.t('animals.title')).toBe('Animals')
      expect(i18n.global.t('animals.name')).toBe('Name')
      expect(i18n.global.t('animals.species')).toBe('Species')
    })
  })
})
