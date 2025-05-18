# AnimalSys

A full-stack, containerized animal welfare management system.

## 🐳 Quick Start (Docker Compose)

```sh
# In the project root
docker-compose up --build
```

- Frontend: http://localhost:5173
- Backend API: http://localhost:8080/api
- Mongo Express: http://localhost:8081

## 🗂️ Project Structure

```
E:/animalsys
│   README.md
│   docker-compose.yml
│
├── backend          # Go API server with MongoDB connectivity
├── frontend         # Vue.js SPA with centralized API client
│   └── .env         # Frontend environment config
└── mongo            # MongoDB initialization scripts
```

## ⚙️ Environment Variables

Set in `docker-compose.yml` for local dev or create standalone `.env` files:

- **Backend**: See `backend/config/config.go` for all options
- **Frontend**: Create `.env` file in frontend folder with `VITE_API_URL=http://localhost:8080`

## 🏗️ Development

- **Backend**: Go + Gin, hot reload via volume mount
- **Frontend**: Vue 3 + Vite + Vuex, hot reload via volume mount
- **Database**: MongoDB with schema validation

### Local Backend Dev
```sh
cd backend
go mod tidy
go run main.go
```

### Local Frontend Dev
```sh
cd frontend
npm install
npm run dev
```

## 🧪 Testing
- Backend: `go test ./...`
- Frontend: `npm run test`

## 📦 Modules
- Animals, Adoptions, Schedule, Documents, Users, Auth, Finances (see `/api` routes)

## 🔌 API Configuration
The frontend connects to the backend API using:
- Centralized Axios instance in `src/utils/api.js`
- Environment variables from `.env` file
- Automatic token-based authentication

## 📝 License
MIT