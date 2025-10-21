# AnimalSys

A full-stack, containerized animal welfare management system.

## Overview

AnimalSys is a comprehensive web application designed to streamline operations for animal shelters and welfare organizations. It provides tools for managing animals, adoptions, employee schedules, documents, and finances.

## Features

- **Animal Management**: Track animals, their health history, and adoption status
- **Adoption System**: Manage adoption applications with approval workflow
- **User Management**: Role-based access control (Admin, Employee, Volunteer, User)
- **Schedule Management**: Employee shift scheduling with swap and absence requests
- **Document Management**: Store and organize medical records, contracts, and other files
- **Finance Tracking**: Track income and expenses with categorization
- **Authentication**: JWT-based authentication with optional LDAP support
- **Internationalization**: Multi-language support (currently English)

## Technology Stack

### Backend
- **Language**: Go 1.23
- **Framework**: Gin
- **Database**: MongoDB 6.0
- **Authentication**: JWT, LDAP (optional)

### Frontend
- **Framework**: Vue.js 3
- **Build Tool**: Vite
- **State Management**: Vuex
- **Routing**: Vue Router
- **HTTP Client**: Axios
- **i18n**: Vue I18n

### DevOps
- **Containerization**: Docker
- **Orchestration**: Docker Compose
- **Database UI**: Mongo Express

## Quick Start

### Prerequisites
- Docker and Docker Compose
- Git

### Running the Application

1. Clone the repository:
```sh
git clone https://github.com/Sainaif/animalsys.git
cd animalsys
```

2. Start all services:
```sh
docker-compose up --build
```

3. Access the application:
- **Frontend**: http://localhost:5173
- **Backend API**: http://localhost:8080/api
- **Mongo Express**: http://localhost:8081

### Default Credentials
The application starts with no default users. You'll need to register a new account through the registration page.

## Project Structure

```
animalsys/
├── backend/              # Go backend application
│   ├── config/          # Configuration management
│   ├── controllers/     # API controllers
│   ├── middlewares/     # Authentication & RBAC middlewares
│   ├── models/          # Data models
│   ├── routes/          # Route definitions
│   ├── tests/           # Unit tests
│   ├── utils/           # Utility functions (JWT, password, LDAP)
│   ├── main.go          # Application entry point
│   ├── go.mod           # Go dependencies
│   └── Dockerfile       # Backend container configuration
│
├── frontend/            # Vue.js frontend application
│   ├── public/          # Static assets
│   ├── src/
│   │   ├── components/  # Vue components
│   │   ├── locales/     # i18n translations
│   │   ├── router/      # Route configuration
│   │   ├── store/       # Vuex state management
│   │   ├── utils/       # API client
│   │   ├── views/       # Page components
│   │   ├── App.vue      # Root component
│   │   └── main.js      # Application entry point
│   ├── package.json     # NPM dependencies
│   ├── vite.config.js   # Vite configuration
│   └── Dockerfile       # Frontend container configuration
│
├── mongo/               # MongoDB initialization
│   └── init-mongo.js    # Database schema and validation
│
├── docker-compose.yml   # Docker Compose orchestration
└── README.md           # This file
```

## API Endpoints

### Authentication
- `POST /api/auth/login` - User login
- `POST /api/auth/register` - User registration

### Animals
- `GET /api/animals` - List all animals
- `GET /api/animals/:id` - Get animal by ID
- `POST /api/animals` - Create new animal (Admin/Employee)
- `PUT /api/animals/:id` - Update animal (Admin/Employee)
- `DELETE /api/animals/:id` - Delete animal (Admin)

### Adoptions
- `GET /api/adoptions` - List all adoptions
- `GET /api/adoptions/:id` - Get adoption by ID
- `POST /api/adoptions` - Create adoption application
- `PUT /api/adoptions/:id` - Update adoption (Admin/Employee)
- `DELETE /api/adoptions/:id` - Delete adoption (Admin)

### Schedules
- `GET /api/schedules` - List all schedules
- `GET /api/schedules/:id` - Get schedule by ID
- `POST /api/schedules` - Create schedule (Admin/Employee)
- `PUT /api/schedules/:id` - Update schedule
- `DELETE /api/schedules/:id` - Delete schedule (Admin/Employee)

### Documents
- `GET /api/documents` - List all documents
- `GET /api/documents/:id` - Get document by ID
- `POST /api/documents` - Upload document (Admin/Employee)
- `DELETE /api/documents/:id` - Delete document (Admin/Employee)

### Finances
- `GET /api/finances` - List all finances (Admin/Employee)
- `GET /api/finances/:id` - Get finance by ID (Admin/Employee)
- `POST /api/finances` - Create finance record (Admin/Employee)
- `PUT /api/finances/:id` - Update finance (Admin/Employee)
- `DELETE /api/finances/:id` - Delete finance (Admin)

### Users
- `GET /api/users` - List all users
- `GET /api/users/:id` - Get user by ID
- `PUT /api/users/:id` - Update user
- `DELETE /api/users/:id` - Delete user (Admin)

## Environment Variables

### Backend (.env)
```
MONGO_URI=mongodb://root:example@mongo:27017/animalsys?authSource=admin
MONGO_DB=animalsys
JWT_SECRET=supersecretkey
LDAP_ENABLED=false
LDAP_SERVER=
LDAP_BASE_DN=
LDAP_BIND_DN=
LDAP_BIND_PASSWORD=
PORT=8080
```

### Frontend (.env)
```
VITE_API_URL=http://localhost:8080/api
```

## Development

### Backend Development
```sh
cd backend
go mod tidy
go run main.go
```

### Frontend Development
```sh
cd frontend
npm install
npm run dev
```

### Running Tests
```sh
cd backend
go test ./...
```

## User Roles

- **Admin**: Full access to all features
- **Employee**: Can manage animals, adoptions, schedules, documents, and finances
- **Volunteer**: Can view animals and schedules
- **User**: Can view animals and submit adoption applications

## LDAP Integration

To enable LDAP authentication, set the following environment variables:

```
LDAP_ENABLED=true
LDAP_SERVER=ldap://your-ldap-server:389
LDAP_BASE_DN=dc=example,dc=com
LDAP_BIND_DN=cn=admin,dc=example,dc=com
LDAP_BIND_PASSWORD=yourpassword
```

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Support

For issues and questions, please open an issue on GitHub.
