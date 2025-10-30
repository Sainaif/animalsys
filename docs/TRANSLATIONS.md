# Translation System Guide / Przewodnik Systemów Tłumaczeń

[English](#english) | [Polski](#polski)

---

## English

### Overview

AnimalSys uses a simple, file-based translation system that makes it easy to add new languages. The system is built on Vue-i18n and uses JSON files for translations.

### File Structure

```
frontend/src/locales/
├── index.js           # Language registry (add new languages here)
├── _example.json      # Template for new languages
├── en.json           # English translations
└── pl.json           # Polish translations
```

### How to Add a New Language

Follow these 4 simple steps:

#### Step 1: Copy the Example File

```bash
cd frontend/src/locales
cp _example.json de.json  # Replace 'de' with your language code
```

Language codes follow ISO 639-1 standard:
- `de` - German (Deutsch)
- `fr` - French (Français)
- `es` - Spanish (Español)
- `it` - Italian (Italiano)
- `nl` - Dutch (Nederlands)
- etc.

#### Step 2: Translate All Values

Open your new file (e.g., `de.json`) and translate all the values:

```json
{
  "common": {
    "save": "Speichern",    // Translated from "Save"
    "cancel": "Abbrechen",  // Translated from "Cancel"
    ...
  }
}
```

**Important**: Only translate the values (right side), never change the keys (left side)!

#### Step 3: Register the Language

Edit `index.js` and:

1. Import your translation file:
```javascript
import de from './de.json'
```

2. Add the language code to the `availableLanguages` array:
```javascript
export const availableLanguages = ['en', 'pl', 'de']
```

3. Add language metadata:
```javascript
export const languageInfo = {
  // ... existing languages
  de: {
    code: 'de',
    name: 'German',
    nativeName: 'Deutsch',
    flag: '🇩🇪'
  }
}
```

4. Add to messages object:
```javascript
export const messages = {
  en,
  pl,
  de  // Your new language
}
```

#### Step 4: Test Your Translation

1. Rebuild the frontend:
```bash
npm run build
```

2. Open the application and switch to your new language in Settings

### Translation Structure

The translation files are organized by feature:

- `common` - Common words (save, cancel, delete, etc.)
- `nav` - Navigation menu items
- `auth` - Authentication pages
- `dashboard` - Dashboard
- `animals` - Animal management
- `adoptions` - Adoption management
- `volunteers` - Volunteer management
- `schedules` - Scheduling
- `documents` - Document management
- `finances` - Financial management
- `donors` - Donor management
- `inventory` - Inventory management
- `veterinary` - Veterinary management
- `campaigns` - Campaign management
- `partners` - Partner management
- `communications` - Communications
- `reports` - Reports
- `users` - User management
- `settings` - Settings
- `public` - Public portal
- `validation` - Validation messages
- `messages` - System messages

### Using Translations in Code

#### In Vue Components (Composition API)

```vue
<script setup>
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
</script>

<template>
  <button>{{ t('common.save') }}</button>
  <h1>{{ t('dashboard.welcome', { name: userName }) }}</h1>
</template>
```

#### In JavaScript Code

```javascript
import { i18n } from '@/i18n'

const message = i18n.global.t('messages.saveSuccess')
```

### Best Practices

1. **Never hardcode text** - Always use translation keys
2. **Use descriptive keys** - `animals.deleteConfirm` is better than `confirm1`
3. **Group related translations** - Keep all animal-related texts in `animals`
4. **Test thoroughly** - Check all pages after adding a new language
5. **Keep consistency** - Use the same terminology across the application
6. **Handle plurals** - Use vue-i18n pluralization for countable items

### Troubleshooting

**Problem**: New language doesn't appear in the language switcher
- **Solution**: Make sure you added it to `availableLanguages` and `languageInfo`

**Problem**: Translations show as keys (e.g., "common.save")
- **Solution**: Check that the JSON file is valid and properly imported

**Problem**: Some texts are in English even after switching language
- **Solution**: Check if those keys exist in your translation file

---

## Polski

### Przegląd

AnimalSys używa prostego systemu tłumaczeń opartego na plikach, który ułatwia dodawanie nowych języków. System jest zbudowany na Vue-i18n i używa plików JSON do tłumaczeń.

### Struktura Plików

```
frontend/src/locales/
├── index.js           # Rejestr języków (tutaj dodajemy nowe języki)
├── _example.json      # Szablon dla nowych języków
├── en.json           # Tłumaczenia angielskie
└── pl.json           # Tłumaczenia polskie
```

### Jak Dodać Nowy Język

Wykonaj te 4 proste kroki:

#### Krok 1: Skopiuj Plik Przykładowy

```bash
cd frontend/src/locales
cp _example.json de.json  # Zastąp 'de' kodem swojego języka
```

Kody języków według standardu ISO 639-1:
- `de` - Niemiecki (Deutsch)
- `fr` - Francuski (Français)
- `es` - Hiszpański (Español)
- `it` - Włoski (Italiano)
- `nl` - Holenderski (Nederlands)
- itd.

#### Krok 2: Przetłumacz Wszystkie Wartości

Otwórz nowy plik (np. `de.json`) i przetłumacz wszystkie wartości:

```json
{
  "common": {
    "save": "Speichern",    // Przetłumaczone z "Save"
    "cancel": "Abbrechen",  // Przetłumaczone z "Cancel"
    ...
  }
}
```

**Ważne**: Tłumacz tylko wartości (prawa strona), nigdy nie zmieniaj kluczy (lewa strona)!

#### Krok 3: Zarejestruj Język

Edytuj plik `index.js`:

1. Zaimportuj plik tłumaczenia:
```javascript
import de from './de.json'
```

2. Dodaj kod języka do tablicy `availableLanguages`:
```javascript
export const availableLanguages = ['en', 'pl', 'de']
```

3. Dodaj metadane języka:
```javascript
export const languageInfo = {
  // ... istniejące języki
  de: {
    code: 'de',
    name: 'German',
    nativeName: 'Deutsch',
    flag: '🇩🇪'
  }
}
```

4. Dodaj do obiektu messages:
```javascript
export const messages = {
  en,
  pl,
  de  // Twój nowy język
}
```

#### Krok 4: Przetestuj Tłumaczenie

1. Przebuduj frontend:
```bash
npm run build
```

2. Otwórz aplikację i przełącz na nowy język w Ustawieniach

### Struktura Tłumaczeń

Pliki tłumaczeń są zorganizowane według funkcji:

- `common` - Wspólne słowa (zapisz, anuluj, usuń, itp.)
- `nav` - Elementy menu nawigacji
- `auth` - Strony uwierzytelniania
- `dashboard` - Panel główny
- `animals` - Zarządzanie zwierzętami
- `adoptions` - Zarządzanie adopcjami
- `volunteers` - Zarządzanie wolontariuszami
- `schedules` - Grafiki
- `documents` - Zarządzanie dokumentami
- `finances` - Zarządzanie finansami
- `donors` - Zarządzanie darczyńcami
- `inventory` - Zarządzanie magazynem
- `veterinary` - Zarządzanie weterynarią
- `campaigns` - Zarządzanie kampaniami
- `partners` - Zarządzanie partnerami
- `communications` - Komunikacja
- `reports` - Raporty
- `users` - Zarządzanie użytkownikami
- `settings` - Ustawienia
- `public` - Portal publiczny
- `validation` - Komunikaty walidacji
- `messages` - Komunikaty systemowe

### Używanie Tłumaczeń w Kodzie

#### W Komponentach Vue (Composition API)

```vue
<script setup>
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
</script>

<template>
  <button>{{ t('common.save') }}</button>
  <h1>{{ t('dashboard.welcome', { name: userName }) }}</h1>
</template>
```

#### W Kodzie JavaScript

```javascript
import { i18n } from '@/i18n'

const message = i18n.global.t('messages.saveSuccess')
```

### Najlepsze Praktyki

1. **Nigdy nie koduj tekstu na stałe** - Zawsze używaj kluczy tłumaczeń
2. **Używaj opisowych kluczy** - `animals.deleteConfirm` jest lepsze niż `confirm1`
3. **Grupuj powiązane tłumaczenia** - Trzymaj wszystkie teksty o zwierzętach w `animals`
4. **Testuj dokładnie** - Sprawdź wszystkie strony po dodaniu nowego języka
5. **Zachowaj spójność** - Używaj tej samej terminologii w całej aplikacji
6. **Obsługuj liczby mnogie** - Używaj pluralizacji vue-i18n dla elementów policzalnych

### Rozwiązywanie Problemów

**Problem**: Nowy język nie pojawia się w przełączniku języków
- **Rozwiązanie**: Upewnij się, że dodałeś go do `availableLanguages` i `languageInfo`

**Problem**: Tłumaczenia pokazują się jako klucze (np. "common.save")
- **Rozwiązanie**: Sprawdź, czy plik JSON jest poprawny i prawidłowo zaimportowany

**Problem**: Niektóre teksty są po angielsku nawet po przełączeniu języka
- **Rozwiązanie**: Sprawdź, czy te klucze istnieją w pliku tłumaczenia
