# AnimalSys

A full-stack, containerized animal welfare management system.

## 🐳 Quick Start (Docker Compose)

```sh
# In the project root (E:/animalsys)
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
├── backend
├── frontend
└── mongo
```

## ⚙️ Environment Variables

Set in `docker-compose.yml` for local dev. See `backend/config/config.go` for all options.

## 🏗️ Development

- **Backend**: Go + Gin, hot reload via volume mount
- **Frontend**: Vue 3 + Vite, hot reload via volume mount
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
- Animals, Adoptions, Schedule, Documents, Users, Auth (see `/api` routes)

## 📝 License
MIT 