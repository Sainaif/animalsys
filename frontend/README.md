# AnimalSys Frontend

Modern Vue.js 3 frontend application for the AnimalSys ERP system.

## Tech Stack

- **Vue 3** - Progressive JavaScript framework (Composition API)
- **Vite** - Next generation frontend tooling
- **Pinia** - State management
- **Vue Router** - Official router
- **Vue I18n** - Internationalization (Polish + English)
- **Axios** - HTTP client
- **Tailwind CSS** - Utility-first CSS framework
- **Vitest** - Unit testing framework
- **Playwright** - E2E testing

## Architecture

```
frontend/
├── src/
│   ├── api/              # API client and modules
│   │   ├── client.js     # Axios instance with interceptors
│   │   └── modules/      # API endpoint modules
│   ├── assets/           # Static assets
│   ├── components/       # Reusable Vue components
│   │   └── common/       # Common components
│   ├── composables/      # Composition API functions
│   ├── layouts/          # Layout components
│   │   ├── PublicLayout.vue
│   │   └── AuthenticatedLayout.vue
│   ├── locales/          # i18n translations
│   │   ├── en.json
│   │   ├── pl.json
│   │   ├── _example.json
│   │   └── index.js
│   ├── router/           # Vue Router configuration
│   ├── stores/           # Pinia stores
│   │   ├── auth.js       # Authentication store
│   │   ├── theme.js      # Theme management
│   │   └── notification.js
│   ├── styles/           # Global styles
│   │   └── main.css      # Tailwind + custom CSS
│   ├── views/            # Page components
│   │   ├── public/       # Public pages
│   │   └── ...           # Module pages
│   ├── App.vue           # Root component
│   └── main.js           # Application entry point
└── tests/                # Test files
```

## Features

### Authentication & Authorization
- JWT-based authentication with access & refresh tokens
- Auto token refresh on expiration
- Role-based access control (RBAC)
- Protected routes with navigation guards
- Session persistence

### State Management
- **Auth Store**: User authentication, session management
- **Theme Store**: Dark/Light mode with localStorage persistence
- **Notification Store**: Toast notifications system

### Internationalization
- Full bilingual support (Polish + English)
- Easy to add new languages using `_example.json` template
- Language switcher in navigation
- Stored preference in localStorage

### UI/UX
- Responsive design (mobile, tablet, desktop)
- Dark/Light theme toggle
- Loading states and error handling
- Toast notifications
- Accessible components

### API Integration
- Axios client with interceptors
- Automatic token injection
- Token refresh handling
- Error handling
- Request/response logging

## Installation

```bash
# Install dependencies
npm install
```

## Development

```bash
# Start dev server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview

# Run tests
npm run test

# Run E2E tests
npm run test:e2e

# Lint code
npm run lint

# Format code
npm run format
```

## Environment Variables

Create `.env` file:

```env
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

## Project Structure

### Layouts

**PublicLayout**: For unauthenticated users
- Header with logo, navigation, language & theme switchers
- Footer
- Used for: Home, Login, Register, Public animal listings

**AuthenticatedLayout**: For logged-in users
- Collapsible sidebar navigation
- Header with page title, user menu
- Role-based menu items
- Used for: Dashboard and all admin modules

### Routing

Routes are organized by access level:
- **Public routes** (`/`): Accessible to everyone
- **Auth routes** (`/app`): Require authentication
- **Role-restricted routes**: Require specific roles (admin, employee, etc.)

Navigation guards automatically:
- Redirect unauthenticated users to login
- Redirect authenticated users away from guest pages
- Check role requirements
- Handle unauthorized access

### State Stores

**Auth Store** ([stores/auth.js](src/stores/auth.js)):
- Login/logout/register
- Token management
- User profile
- Role checking
- Password change

**Theme Store** ([stores/theme.js](src/stores/theme.js)):
- Dark/Light mode
- System preference detection
- Persistence to localStorage

**Notification Store** ([stores/notification.js](src/stores/notification.js)):
- Show success/error/warning/info messages
- Auto-dismiss
- Queue management

### API Client

Configured Axios instance with:
- Base URL from environment
- Auto token injection
- Token refresh on 401
- Request/response interceptors
- Error handling

### Composables

**useApi** ([composables/useApi.js](src/composables/useApi.js)):
Reusable composition function for API calls:
```js
const { data, loading, error, execute } = useApi(apiFunc, options)
```

## Adding New Languages

1. Copy `src/locales/_example.json` to `src/locales/[code].json`
2. Translate all strings
3. Add language to `src/locales/index.js`:
   ```js
   import newLang from './newLang.json'

   export const availableLanguages = ['en', 'pl', 'newLang']
   export const languageInfo = {
     // ...
     newLang: { code: 'newLang', name: 'Language Name', nativeName: 'Native Name', flag: '🏳️' }
   }
   export const messages = { en, pl, newLang }
   ```

## Styling

Using Tailwind CSS with custom CSS variables for theming:

```css
:root {
  --primary-color: #3b82f6;
  --bg-primary: #ffffff;
  --text-primary: #1f2937;
  /* ... */
}

[data-theme="dark"] {
  --bg-primary: #1f2937;
  --text-primary: #f9fafb;
  /* ... */
}
```

Components can use:
- Tailwind utility classes
- CSS variables for theme-aware styling
- Scoped styles in SFC

## Testing

### Unit Tests (Vitest)
```bash
npm run test
```

### E2E Tests (Playwright)
```bash
npm run test:e2e
```

## Build & Deployment

```bash
# Build for production
npm run build

# Output in dist/
```

The built application is a static SPA that can be served by any static file server or CDN.

## Browser Support

- Chrome/Edge (latest 2 versions)
- Firefox (latest 2 versions)
- Safari (latest 2 versions)
- Mobile browsers (iOS Safari, Chrome Mobile)

## License

Proprietary - All rights reserved
