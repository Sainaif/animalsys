# AnimalSys - Foundation ERP System

**AnimalSys** is a comprehensive, secure, and lightweight ERP system designed specifically for animal welfare foundations and shelters.

## 🌟 Features

- **Multilingual Support**: Easy to add new languages (Polish & English included)
- **Secure**: JWT authentication, RBAC, audit trails
- **Responsive**: Mobile-first design with Dark/Light mode
- **Comprehensive**: 14 functional modules covering all foundation operations
- **Well-tested**: 90%+ test coverage
- **Easy Deployment**: One-command setup with Docker

## 🏗️ Architecture

- **Frontend**: Vue.js 3 + Vite + Pinia + Tailwind CSS
- **Backend**: Go + Gin framework + Clean Architecture
- **Database**: MongoDB with schema validation
- **Deployment**: Docker + Docker Compose

## 📦 Modules

1. **Users & Roles** - User management with RBAC
2. **Animals** - Complete animal records and medical history
3. **Adoptions** - Adoption workflow from application to completion
4. **Volunteers** - Volunteer management, training, and hours tracking
5. **Schedules & Shifts** - Shift scheduling with swap and absence requests
6. **Documents** - Secure document storage with GridFS
7. **Finances** - Income/expense tracking and financial reports
8. **Donors & Donations** - Donor database and donation tracking
9. **Inventory** - Stock management for food, medicine, and supplies
10. **Veterinary** - Veterinary visits, vaccinations, and treatments
11. **Campaigns** - Fundraising and adoption campaigns
12. **Partners** - Partner and sponsor management
13. **Communications** - Newsletter and bulk messaging
14. **Reports & Analytics** - Comprehensive reporting and statutory compliance

## 🚀 Quick Start

```bash
# Clone repository
git clone <repository-url>
cd animalsys

# Run setup script
chmod +x deployment/setup.sh
./deployment/setup.sh

# Access the application
# Frontend: http://localhost:3000
# Backend API: http://localhost:8080
```

## 📖 Documentation

- [Architecture Overview](docs/ARCHITECTURE.md)
- [Development Guide](docs/DEVELOPMENT.md)
- [API Documentation](docs/API.md)
- [Translation Guide](docs/TRANSLATIONS.md)
- [Deployment Guide](docs/DEPLOYMENT.md)
- [User Manual (PL)](docs/USER_GUIDE_PL.md)
- [User Manual (EN)](docs/USER_GUIDE_EN.md)

## 🧪 Testing

```bash
# Backend tests
cd backend
go test ./... -v -cover

# Frontend tests
cd frontend
npm run test
npm run test:e2e
```

## 📄 License

[Add your license here]

## 🤝 Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](docs/CONTRIBUTING.md) for details.

## 📧 Support

For support and questions, please [open an issue](https://github.com/your-org/animalsys/issues).
