# Animal Foundation CRM - API Documentation

**Version:** 1.0
**Base URL:** `http://localhost:8080/api/v1`
**Last Updated:** 2025-11-08

---

## Table of Contents

1. [Authentication](#authentication)
2. [Common Structures](#common-structures)
3. [Error Handling](#error-handling)
4. [Permissions](#permissions)
5. [API Endpoints](#api-endpoints)
   - [Authentication & User Management](#authentication--user-management)
   - [Animal Management](#animal-management)
   - [Veterinary Records](#veterinary-records)
   - [Adoption Management](#adoption-management)
   - [Donor Management](#donor-management)
   - [Donation Management](#donation-management)
   - [Campaign Management](#campaign-management)
   - [Event Management](#event-management)
   - [Volunteer Management](#volunteer-management)
   - [Communication Management](#communication-management)
   - [Notification System](#notification-system)
   - [Report Generation](#report-generation)
   - [Dashboard & Analytics](#dashboard--analytics)
   - [Settings Management](#settings-management)
   - [Task Management](#task-management)
   - [Document Management](#document-management)
   - [Partner Management](#partner-management)
   - [Transfer Management](#transfer-management)
   - [Inventory Management](#inventory-management)
   - [Stock Transactions](#stock-transactions)
   - [Audit Logs](#audit-logs)
   - [System Monitoring](#system-monitoring)
   - [Medical Records & Prescriptions](#medical-records--prescriptions)

---

## Authentication

### Authentication Methods

The API uses **JWT (JSON Web Tokens)** for authentication with two token types:

- **Access Token**: Short-lived token (default: 15 minutes) for API requests
- **Refresh Token**: Long-lived token (default: 7 days) for obtaining new access tokens

### Using Authentication

Include the access token in the `Authorization` header:

```
Authorization: Bearer <access_token>
```

### Login Flow

1. **Login**: POST `/auth/login` with credentials → Receive access token + refresh token
2. **Make Requests**: Include access token in Authorization header
3. **Token Expires**: Use refresh token at POST `/auth/refresh` → Receive new access token
4. **Logout**: POST `/auth/logout` to invalidate refresh token

---

## Common Structures

### MongoDB ObjectID

All IDs in the system use MongoDB ObjectID format:
- **Format**: 24-character hexadecimal string
- **Example**: `"507f1f77bcf86cd799439011"`

### Pagination

List endpoints support pagination using query parameters:

```
?limit=20&offset=0
```

- `limit`: Number of items per page (default: 20, max: 100)
- `offset`: Number of items to skip (default: 0)

**Response Format:**
```json
{
  "data": [...],
  "total": 150,
  "limit": 20,
  "offset": 0
}
```

### Date/Time Format

All timestamps use **ISO 8601** format:
```
2025-11-08T15:30:00Z
```

### Common Query Parameters

- **Sorting**: `?sort_by=created_at&sort_order=desc`
- **Search**: `?search=keyword`
- **Date Range**: `?start_date=2025-01-01&end_date=2025-12-31`
- **Status Filter**: `?status=active`

---

## Error Handling

### Error Response Format

```json
{
  "error": "Error message description"
}
```

### HTTP Status Codes

| Code | Meaning | Description |
|------|---------|-------------|
| 200 | OK | Request successful |
| 201 | Created | Resource created successfully |
| 400 | Bad Request | Invalid request data |
| 401 | Unauthorized | Authentication required or failed |
| 403 | Forbidden | Insufficient permissions |
| 404 | Not Found | Resource not found |
| 409 | Conflict | Resource conflict (e.g., duplicate email) |
| 500 | Internal Server Error | Server error |

---

## Permissions

### User Roles

| Role | Description | Access Level |
|------|-------------|--------------|
| `super_admin` | Super Administrator | Full system access |
| `admin` | Administrator | Full access except system-critical operations |
| `employee` | Employee | Can create/update most resources |
| `volunteer` | Volunteer | Limited access to specific features |
| `user` | Basic User | Read-only access to allowed resources |

### Permission System

Permissions are granular and role-based:

- **View**: Read access to resources
- **Create**: Ability to create new resources
- **Update**: Ability to modify existing resources
- **Delete**: Ability to remove resources

**Permission Format**: `Permission[Action][Resource]`

Examples:
- `PermissionViewAnimals`
- `PermissionCreateVeterinary`
- `PermissionUpdateAdoptions`
- `PermissionDeleteDonors`

---

## API Endpoints

---

## Authentication & User Management

### Authentication Endpoints

#### POST /api/v1/auth/login
**Description**: Authenticate user and receive JWT tokens
**Authentication**: None (Public)
**Permissions**: None

**Request Body:**
```json
{
  "email": "admin@example.com",
  "password": "securePassword123"
}
```

**Response: 200 OK**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "507f1f77bcf86cd799439011",
    "email": "admin@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "role": "admin",
    "status": "active",
    "phone": "+1234567890",
    "language": "en",
    "theme": "light",
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-11-08T15:30:00Z"
  }
}
```

**Errors:**
- `400`: Invalid request data
- `401`: Invalid credentials
- `403`: Account suspended or inactive

---

#### POST /api/v1/auth/refresh
**Description**: Refresh access token using refresh token
**Authentication**: None (Public)
**Permissions**: None

**Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Response: 200 OK**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

---

#### POST /api/v1/auth/logout
**Description**: Logout user and invalidate refresh token
**Authentication**: Required
**Permissions**: Authenticated user

**Response: 200 OK**
```json
{
  "message": "logged out successfully"
}
```

---

#### POST /api/v1/auth/register
**Description**: Register a new user (Admin only)
**Authentication**: Required
**Permissions**: Admin

**Request Body:**
```json
{
  "email": "newuser@example.com",
  "password": "securePassword123",
  "first_name": "Jane",
  "last_name": "Smith",
  "role": "employee",
  "phone": "+1234567890",
  "language": "en",
  "theme": "light"
}
```

**Response: 201 Created**
```json
{
  "id": "507f1f77bcf86cd799439012",
  "email": "newuser@example.com",
  "first_name": "Jane",
  "last_name": "Smith",
  "role": "employee",
  "status": "active",
  "phone": "+1234567890",
  "language": "en",
  "theme": "light",
  "created_at": "2025-11-08T15:30:00Z",
  "updated_at": "2025-11-08T15:30:00Z"
}
```

---

#### GET /api/v1/auth/me
**Description**: Get current authenticated user's information
**Authentication**: Required
**Permissions**: Authenticated user

**Response: 200 OK**
```json
{
  "id": "507f1f77bcf86cd799439011",
  "email": "admin@example.com",
  "first_name": "John",
  "last_name": "Doe",
  "role": "admin",
  "status": "active",
  "phone": "+1234567890",
  "language": "en",
  "theme": "light",
  "last_login": "2025-11-08T15:30:00Z",
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-11-08T15:30:00Z"
}
```

---

#### PUT /api/v1/auth/change-password
**Description**: Change current user's password
**Authentication**: Required
**Permissions**: Authenticated user

**Request Body:**
```json
{
  "current_password": "oldPassword123",
  "new_password": "newSecurePassword456"
}
```

**Response: 200 OK**
```json
{
  "message": "password changed successfully"
}
```

---

### User Management Endpoints

#### GET /api/v1/users
**Description**: List all users with pagination and filtering
**Authentication**: Required
**Permissions**: `PermissionViewUsers`

**Query Parameters:**
- `limit` (int): Items per page (default: 20)
- `offset` (int): Number of items to skip (default: 0)
- `role` (string): Filter by role
- `status` (string): Filter by status
- `search` (string): Search by name or email

**Response: 200 OK**
```json
{
  "data": [
    {
      "id": "507f1f77bcf86cd799439011",
      "email": "admin@example.com",
      "first_name": "John",
      "last_name": "Doe",
      "role": "admin",
      "status": "active",
      "phone": "+1234567890",
      "language": "en",
      "theme": "light",
      "last_login": "2025-11-08T15:30:00Z",
      "created_at": "2025-01-01T00:00:00Z",
      "updated_at": "2025-11-08T15:30:00Z"
    }
  ],
  "total": 50,
  "limit": 20,
  "offset": 0
}
```

---

#### GET /api/v1/users/:id
**Description**: Get user by ID
**Authentication**: Required
**Permissions**: `PermissionViewUsers`

**Path Parameters:**
- `id` (string): User ID

**Response: 200 OK**
```json
{
  "id": "507f1f77bcf86cd799439011",
  "email": "admin@example.com",
  "first_name": "John",
  "last_name": "Doe",
  "role": "admin",
  "status": "active",
  "phone": "+1234567890",
  "language": "en",
  "theme": "light",
  "last_login": "2025-11-08T15:30:00Z",
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-11-08T15:30:00Z"
}
```

---

#### POST /api/v1/users
**Description**: Create a new user
**Authentication**: Required
**Permissions**: Admin

**Request Body:** Same as `/auth/register`

**Response: 201 Created**

---

#### PUT /api/v1/users/:id
**Description**: Update user information
**Authentication**: Required
**Permissions**: Admin

**Request Body:**
```json
{
  "first_name": "John",
  "last_name": "Doe",
  "role": "admin",
  "status": "active",
  "phone": "+1234567890",
  "language": "en",
  "theme": "dark"
}
```

**Response: 200 OK**

---

#### DELETE /api/v1/users/:id
**Description**: Delete a user
**Authentication**: Required
**Permissions**: Super Admin

**Response: 200 OK**
```json
{
  "message": "user deleted successfully"
}
```

---

#### PUT /api/v1/users/:id/reset-password
**Description**: Reset user password (Admin)
**Authentication**: Required
**Permissions**: Admin

**Request Body:**
```json
{
  "new_password": "newSecurePassword456"
}
```

**Response: 200 OK**
```json
{
  "message": "password reset successfully"
}
```

---

## Animal Management

### Animal Data Structure

```json
{
  "id": "507f1f77bcf86cd799439013",
  "name": {
    "en": "Buddy",
    "pl": "Budzik",
    "latin": "Canis lupus familiaris"
  },
  "species": "dog",
  "breed": "Golden Retriever",
  "category": "mammal",
  "sex": "male",
  "date_of_birth": "2022-05-15T00:00:00Z",
  "age_years": 2,
  "age_months": 6,
  "color": "golden",
  "size": "large",
  "weight": 30.5,
  "status": "available",
  "intake_date": "2024-01-15T00:00:00Z",
  "intake_reason": "Owner surrender",
  "images": {
    "primary": "https://example.com/images/buddy.jpg",
    "gallery": [
      "https://example.com/images/buddy-1.jpg",
      "https://example.com/images/buddy-2.jpg"
    ],
    "thumbnails": [
      "https://example.com/images/buddy-thumb.jpg"
    ]
  },
  "medical_info": {
    "vaccinated": true,
    "sterilized": true,
    "microchipped": true,
    "microchip_number": "123456789012345",
    "health_status": "healthy",
    "medications": [],
    "allergies": [],
    "special_needs": "",
    "last_vet_visit": "2025-10-15T00:00:00Z",
    "next_vet_visit": "2025-12-15T00:00:00Z"
  },
  "behavior_info": {
    "temperament": ["friendly", "playful", "energetic"],
    "good_with_kids": true,
    "good_with_dogs": true,
    "good_with_cats": false,
    "house_trained": true,
    "special_requirements": "Needs daily exercise"
  },
  "adoption_info": {
    "adoption_fee": 150.00,
    "special_conditions": "Must have a fenced yard",
    "availability_status": "available"
  },
  "description": {
    "en": "Buddy is a friendly and energetic Golden Retriever...",
    "pl": "Budzik jest przyjaznym i energicznym Golden Retrieverem..."
  },
  "daily_notes": [
    {
      "date": "2025-11-08T10:00:00Z",
      "note": "Played well with other dogs in the yard",
      "added_by": "507f1f77bcf86cd799439011"
    }
  ],
  "created_at": "2024-01-15T00:00:00Z",
  "updated_at": "2025-11-08T15:30:00Z"
}
```

### Animal Endpoints

#### GET /api/v1/animals
**Description**: List all animals with pagination and filtering
**Authentication**: Required
**Permissions**: `PermissionViewAnimals`

**Query Parameters:**
- `limit` (int): Items per page (default: 20)
- `offset` (int): Number of items to skip (default: 0)
- `species` (string): Filter by species (e.g., "dog", "cat") **✅ VERIFIED**
- `status` (string): Filter by status (available, adopted, under_treatment, etc.)
- `category` (string): Filter by category (mammal, reptile, bird, etc.)
- `sex` (string): Filter by sex (male, female, unknown)
- `size` (string): Filter by size (small, medium, large, xlarge)
- `available_only` (bool): Filter only available animals for adoption **✅ VERIFIED**
- `good_with_kids` (bool): Filter animals good with kids
- `good_with_dogs` (bool): Filter animals good with dogs
- `good_with_cats` (bool): Filter animals good with cats
- `assigned_caretaker` (string): Filter by caretaker ID
- `min_age` (float): Minimum age in years
- `max_age` (float): Maximum age in years
- `search` (string): Search in names and descriptions
- `sort_by` (string): Field to sort by (default: created_at)
- `sort_order` (string): asc or desc (default: desc)

**Examples:**
- Get available animals: `GET /api/v1/animals?available_only=true`
- Get dogs: `GET /api/v1/animals?species=dog`
- Get available cats: `GET /api/v1/animals?species=cat&available_only=true`
- Get animals good with kids: `GET /api/v1/animals?good_with_kids=true`

**Response: 200 OK**
```json
{
  "data": [/* Array of animals */],
  "total": 150,
  "limit": 20,
  "offset": 0
}
```

---

#### GET /api/v1/animals/:id
**Description**: Get animal by ID
**Authentication**: Required
**Permissions**: `PermissionViewAnimals`

**Response: 200 OK**
Returns complete animal object (see Animal Data Structure above)

---

#### POST /api/v1/animals
**Description**: Create a new animal record
**Authentication**: Required
**Permissions**: `PermissionCreateAnimals`

**Request Body:**
```json
{
  "name": {
    "en": "Buddy",
    "pl": "Budzik"
  },
  "species": "dog",
  "breed": "Golden Retriever",
  "category": "mammal",
  "sex": "male",
  "date_of_birth": "2022-05-15T00:00:00Z",
  "color": "golden",
  "size": "large",
  "weight": 30.5,
  "status": "available",
  "intake_date": "2024-01-15T00:00:00Z",
  "intake_reason": "Owner surrender",
  "medical_info": {
    "vaccinated": true,
    "sterilized": true,
    "microchipped": true,
    "microchip_number": "123456789012345",
    "health_status": "healthy"
  },
  "behavior_info": {
    "temperament": ["friendly", "playful"],
    "good_with_kids": true,
    "good_with_dogs": true,
    "good_with_cats": false,
    "house_trained": true
  },
  "adoption_info": {
    "adoption_fee": 150.00,
    "availability_status": "available"
  },
  "description": {
    "en": "Buddy is a friendly Golden Retriever..."
  }
}
```

**Response: 201 Created**
Returns complete animal object

---

#### PUT /api/v1/animals/:id
**Description**: Update animal information
**Authentication**: Required
**Permissions**: `PermissionUpdateAnimals`

**Request Body:** Same structure as POST (all fields optional)

**Response: 200 OK**

---

#### DELETE /api/v1/animals/:id
**Description**: Delete an animal record
**Authentication**: Required
**Permissions**: `PermissionDeleteAnimals`

**Response: 200 OK**
```json
{
  "message": "animal deleted successfully"
}
```

---

#### POST /api/v1/animals/:id/images
**Description**: Upload animal images
**Authentication**: Required
**Permissions**: `PermissionUpdateAnimals`

**Content-Type**: `multipart/form-data`

**Form Data:**
- `images[]`: One or more image files
- `primary` (bool): Set as primary image

**Response: 200 OK**
```json
{
  "message": "images uploaded successfully",
  "urls": [
    "https://example.com/images/buddy-new.jpg"
  ]
}
```

---

#### POST /api/v1/animals/:id/notes
**Description**: Add a daily note to animal record
**Authentication**: Required
**Permissions**: `PermissionUpdateAnimals`

**Request Body:**
```json
{
  "note": "Played well with other dogs today"
}
```

**Response: 200 OK**
```json
{
  "message": "note added successfully"
}
```

---

#### GET /api/v1/animals/statistics
**Description**: Get animal statistics
**Authentication**: Required
**Permissions**: `PermissionViewAnimals`

**Response: 200 OK**
```json
{
  "total_animals": 150,
  "by_status": {
    "available": 45,
    "adopted": 80,
    "under_treatment": 15,
    "fostered": 10
  },
  "by_species": {
    "dog": 90,
    "cat": 50,
    "other": 10
  },
  "by_category": {
    "mammal": 140,
    "bird": 8,
    "reptile": 2
  }
}
```

---

#### GET /api/v1/animals/:id/visits
**Description**: Get veterinary visits for an animal
**Authentication**: Required
**Permissions**: `PermissionViewVeterinary`

**Response: 200 OK**
Returns array of veterinary visit objects

---

#### GET /api/v1/animals/:id/vaccinations
**Description**: Get vaccination records for an animal
**Authentication**: Required
**Permissions**: `PermissionViewVeterinary`

**Response: 200 OK**
Returns array of vaccination objects

---

#### GET /api/v1/animals/:id/applications
**Description**: Get adoption applications for an animal
**Authentication**: Required
**Permissions**: `PermissionViewAdoptions`

**Response: 200 OK**
Returns array of adoption application objects

---

#### GET /api/v1/animals/:id/adoption
**Description**: Get adoption record for an animal
**Authentication**: Required
**Permissions**: `PermissionViewAdoptions`

**Response: 200 OK**
Returns adoption object if animal is adopted

---

### Public Animal Endpoints (No Authentication)

#### GET /api/v1/public/animals
**Description**: List available animals (for public adoption page)
**Authentication**: None
**Permissions**: None

Returns animals with status "available" only

---

#### GET /api/v1/public/animals/:id
**Description**: Get animal details (public view)
**Authentication**: None
**Permissions**: None

---

#### GET /api/v1/public/animals/species
**Description**: Get list of animal species
**Authentication**: None
**Permissions**: None

**Response: 200 OK**
```json
{
  "species": ["dog", "cat", "rabbit", "bird", "reptile", "other"]
}
```

---

## Veterinary Records

### Veterinary Visit Structure

```json
{
  "id": "507f1f77bcf86cd799439014",
  "animal_id": "507f1f77bcf86cd799439013",
  "visit_date": "2025-11-08T10:00:00Z",
  "visit_type": "checkup",
  "veterinarian_name": "Dr. Smith",
  "veterinarian_license": "VET-12345",
  "clinic_name": "City Vet Clinic",
  "reason": "Annual checkup",
  "diagnosis": "Healthy, no issues found",
  "treatment": "None required",
  "medications_prescribed": [
    {
      "medication_name": "Heartgard",
      "dosage": "One tablet",
      "frequency": "Monthly",
      "duration": "Ongoing"
    }
  ],
  "tests_performed": ["Blood test", "Physical examination"],
  "test_results": "All results normal",
  "weight": 31.0,
  "temperature": 38.5,
  "heart_rate": 90,
  "notes": "Animal in good health",
  "follow_up_required": false,
  "follow_up_date": null,
  "cost": 75.00,
  "payment_status": "paid",
  "created_by": "507f1f77bcf86cd799439011",
  "created_at": "2025-11-08T10:30:00Z",
  "updated_at": "2025-11-08T10:30:00Z"
}
```

### Vaccination Structure

```json
{
  "id": "507f1f77bcf86cd799439015",
  "animal_id": "507f1f77bcf86cd799439013",
  "vaccine_name": "Rabies",
  "vaccine_type": "core",
  "manufacturer": "Pfizer Animal Health",
  "batch_number": "BATCH-2025-001",
  "vaccination_date": "2025-11-08T10:00:00Z",
  "expiration_date": "2026-11-08T10:00:00Z",
  "administered_by": "Dr. Smith",
  "veterinarian_license": "VET-12345",
  "clinic_name": "City Vet Clinic",
  "next_due_date": "2026-11-08T10:00:00Z",
  "site_of_injection": "Left hind leg",
  "adverse_reactions": "",
  "notes": "No reactions observed",
  "cost": 25.00,
  "certificate_number": "CERT-2025-001",
  "created_by": "507f1f77bcf86cd799439011",
  "created_at": "2025-11-08T10:30:00Z",
  "updated_at": "2025-11-08T10:30:00Z"
}
```

### Veterinary Endpoints

#### GET /api/v1/veterinary/visits
**Description**: List all veterinary visits
**Authentication**: Required
**Permissions**: `PermissionViewVeterinary`

**Query Parameters:**
- `limit`, `offset`: Pagination
- `animal_id` (string): Filter by animal
- `visit_type` (string): Filter by visit type
- `start_date`, `end_date`: Date range

**Response: 200 OK**

---

#### GET /api/v1/veterinary/visits/:id
**Description**: Get veterinary visit by ID
**Authentication**: Required
**Permissions**: `PermissionViewVeterinary`

**Response: 200 OK**

---

#### POST /api/v1/veterinary/visits
**Description**: Create veterinary visit record
**Authentication**: Required
**Permissions**: `PermissionCreateVeterinary`

**Request Body:** (See Veterinary Visit Structure)

**Response: 201 Created**

---

#### PUT /api/v1/veterinary/visits/:id
**Description**: Update veterinary visit
**Authentication**: Required
**Permissions**: `PermissionUpdateVeterinary`

**Response: 200 OK**

---

#### DELETE /api/v1/veterinary/visits/:id
**Description**: Delete veterinary visit
**Authentication**: Required
**Permissions**: `PermissionDeleteVeterinary`

**Response: 200 OK**

---

#### GET /api/v1/veterinary/visits/upcoming
**Description**: Get upcoming veterinary visits
**Authentication**: Required
**Permissions**: `PermissionViewVeterinary`

**Response: 200 OK**

---

#### GET /api/v1/veterinary/vaccinations
**Description**: List all vaccinations
**Authentication**: Required
**Permissions**: `PermissionViewVeterinary`

**Query Parameters:**
- `limit`, `offset`: Pagination
- `animal_id` (string): Filter by animal
- `vaccine_type` (string): Filter by vaccine type

**Response: 200 OK**

---

#### GET /api/v1/veterinary/vaccinations/:id
**Description**: Get vaccination by ID
**Authentication**: Required
**Permissions**: `PermissionViewVeterinary`

**Response: 200 OK**

---

#### POST /api/v1/veterinary/vaccinations
**Description**: Create vaccination record
**Authentication**: Required
**Permissions**: `PermissionCreateVeterinary`

**Request Body:** (See Vaccination Structure)

**Response: 201 Created**

---

#### DELETE /api/v1/veterinary/vaccinations/:id
**Description**: Delete vaccination record
**Authentication**: Required
**Permissions**: `PermissionDeleteVeterinary`

**Response: 200 OK**

---

#### GET /api/v1/veterinary/vaccinations/due
**Description**: Get vaccinations that are due or overdue
**Authentication**: Required
**Permissions**: `PermissionViewVeterinary`

**Response: 200 OK**

---

## Adoption Management

### Adoption Application Structure

```json
{
  "id": "507f1f77bcf86cd799439016",
  "animal_id": "507f1f77bcf86cd799439013",
  "applicant_name": "Jane Smith",
  "applicant_email": "jane@example.com",
  "applicant_phone": "+1234567890",
  "applicant_address": {
    "street": "123 Main St",
    "city": "New York",
    "state": "NY",
    "postal_code": "10001",
    "country": "USA"
  },
  "application_date": "2025-11-01T00:00:00Z",
  "status": "pending",
  "household_info": {
    "type": "house",
    "has_yard": true,
    "is_fenced": true,
    "has_pool": false,
    "other_pets": [
      {
        "type": "dog",
        "breed": "Labrador",
        "age": 5,
        "spayed_neutered": true
      }
    ],
    "household_members": 4,
    "children_ages": [8, 12]
  },
  "experience": {
    "previous_pets": true,
    "years_of_experience": 10,
    "specific_breed_experience": true,
    "training_experience": "Basic obedience training"
  },
  "references": [
    {
      "name": "John Doe",
      "relationship": "Veterinarian",
      "phone": "+1234567891",
      "email": "john@vetclinic.com"
    }
  ],
  "employment_info": {
    "employed": true,
    "employer": "Tech Company Inc",
    "hours_away": "8 hours",
    "who_cares_for_pet": "Pet sitter"
  },
  "reasons_for_adoption": "Looking for a family companion",
  "reviewed_by": null,
  "review_date": null,
  "review_notes": "",
  "approval_status": "pending",
  "created_at": "2025-11-01T00:00:00Z",
  "updated_at": "2025-11-01T00:00:00Z"
}
```

### Adoption Structure

```json
{
  "id": "507f1f77bcf86cd799439017",
  "animal_id": "507f1f77bcf86cd799439013",
  "application_id": "507f1f77bcf86cd799439016",
  "adopter_name": "Jane Smith",
  "adopter_email": "jane@example.com",
  "adopter_phone": "+1234567890",
  "adopter_address": {
    "street": "123 Main St",
    "city": "New York",
    "state": "NY",
    "postal_code": "10001",
    "country": "USA"
  },
  "adoption_date": "2025-11-08T00:00:00Z",
  "adoption_fee": 150.00,
  "payment_status": "paid",
  "payment_method": "credit_card",
  "contract_signed": true,
  "contract_date": "2025-11-08T00:00:00Z",
  "microchip_transferred": true,
  "follow_up_schedule": [
    {
      "scheduled_date": "2025-12-08T00:00:00Z",
      "type": "phone_call",
      "status": "pending",
      "notes": ""
    }
  ],
  "return_policy_explained": true,
  "trial_period_end": "2025-11-22T00:00:00Z",
  "status": "active",
  "notes": "Adoption went smoothly",
  "created_by": "507f1f77bcf86cd799439011",
  "created_at": "2025-11-08T00:00:00Z",
  "updated_at": "2025-11-08T00:00:00Z"
}
```

### Adoption Endpoints

#### GET /api/v1/adoptions/applications
**Description**: List adoption applications
**Authentication**: Required
**Permissions**: `PermissionViewAdoptions`

**Query Parameters:**
- `limit`, `offset`: Pagination
- `status` (string): Filter by status
- `animal_id` (string): Filter by animal

**Response: 200 OK**

---

#### GET /api/v1/adoptions/applications/:id
**Description**: Get application by ID
**Authentication**: Required
**Permissions**: `PermissionViewAdoptions`

**Response: 200 OK**

---

#### POST /api/v1/adoptions/applications
**Description**: Create adoption application (Public - anyone can apply)
**Authentication**: Required
**Permissions**: None (authenticated users)

**Request Body:** (See Adoption Application Structure)

**Response: 201 Created**

---

#### PUT /api/v1/adoptions/applications/:id
**Description**: Update adoption application (for review/approval)
**Authentication**: Required
**Permissions**: `PermissionUpdateAdoptions`

**Response: 200 OK**

---

#### DELETE /api/v1/adoptions/applications/:id
**Description**: Delete adoption application
**Authentication**: Required
**Permissions**: `PermissionDeleteAdoptions`

**Response: 200 OK**

---

#### GET /api/v1/adoptions/applications/pending
**Description**: Get pending applications
**Authentication**: Required
**Permissions**: `PermissionViewAdoptions`

**Response: 200 OK**

---

#### GET /api/v1/adoptions
**Description**: List adoptions
**Authentication**: Required
**Permissions**: `PermissionViewAdoptions`

**Response: 200 OK**

---

#### GET /api/v1/adoptions/:id
**Description**: Get adoption by ID
**Authentication**: Required
**Permissions**: `PermissionViewAdoptions`

**Response: 200 OK**

---

#### POST /api/v1/adoptions
**Description**: Create adoption record
**Authentication**: Required
**Permissions**: `PermissionCreateAdoptions`

**Request Body:** (See Adoption Structure)

**Response: 201 Created**

---

#### PUT /api/v1/adoptions/:id
**Description**: Update adoption record
**Authentication**: Required
**Permissions**: `PermissionUpdateAdoptions`

**Response: 200 OK**

---

#### DELETE /api/v1/adoptions/:id
**Description**: Delete adoption record
**Authentication**: Required
**Permissions**: `PermissionDeleteAdoptions`

**Response: 200 OK**

---

#### GET /api/v1/adoptions/statistics
**Description**: Get adoption statistics
**Authentication**: Required
**Permissions**: `PermissionViewAdoptions`

**Response: 200 OK**
```json
{
  "total_adoptions": 150,
  "this_month": 12,
  "this_year": 85,
  "pending_applications": 25,
  "average_time_to_adopt_days": 15.5
}
```

---

#### GET /api/v1/adoptions/follow-ups/pending
**Description**: Get adoptions with pending follow-ups
**Authentication**: Required
**Permissions**: `PermissionViewAdoptions`

**Response: 200 OK**

---

## Donor Management

### Donor Structure

```json
{
  "id": "507f1f77bcf86cd799439018",
  "type": "individual",
  "first_name": "John",
  "last_name": "Donor",
  "organization_name": "",
  "email": "john.donor@example.com",
  "phone": "+1234567890",
  "address": {
    "street": "456 Charity Lane",
    "city": "Boston",
    "state": "MA",
    "postal_code": "02101",
    "country": "USA"
  },
  "preferred_contact_method": "email",
  "preferred_language": "en",
  "donor_since": "2020-01-01T00:00:00Z",
  "total_donated": 5000.00,
  "last_donation_date": "2025-10-15T00:00:00Z",
  "last_donation_amount": 100.00,
  "donation_frequency": "monthly",
  "preferred_donation_method": "credit_card",
  "tax_id": "123-45-6789",
  "wants_tax_receipt": true,
  "wants_newsletter": true,
  "wants_thank_you": true,
  "recognition_level": "gold",
  "anonymity_preference": "public",
  "interests": ["animals", "education", "healthcare"],
  "engagement_score": 85,
  "engagement_level": "high",
  "last_engagement_date": "2025-11-01T00:00:00Z",
  "notes": "Very engaged donor, prefers email communication",
  "tags": ["major_donor", "monthly_supporter"],
  "created_by": "507f1f77bcf86cd799439011",
  "created_at": "2020-01-01T00:00:00Z",
  "updated_at": "2025-11-08T15:30:00Z"
}
```

### Donor Endpoints

#### GET /api/v1/donors
**Description**: List donors
**Authentication**: Required
**Permissions**: `PermissionViewDonors`

**Query Parameters:**
- `limit`, `offset`: Pagination
- `type` (string): individual, organization, foundation
- `engagement_level` (string): Filter by engagement
- `search` (string): Search by name or email

**Response: 200 OK**

---

#### GET /api/v1/donors/:id
**Description**: Get donor by ID
**Authentication**: Required
**Permissions**: `PermissionViewDonors`

**Response: 200 OK**

---

#### POST /api/v1/donors
**Description**: Create donor
**Authentication**: Required
**Permissions**: `PermissionCreateDonors`

**Request Body:** (See Donor Structure)

**Response: 201 Created**

---

#### PUT /api/v1/donors/:id
**Description**: Update donor
**Authentication**: Required
**Permissions**: `PermissionUpdateDonors`

**Response: 200 OK**

---

#### DELETE /api/v1/donors/:id
**Description**: Delete donor
**Authentication**: Required
**Permissions**: `PermissionDeleteDonors`

**Response: 200 OK**

---

#### GET /api/v1/donors/major
**Description**: Get major donors (high value)
**Authentication**: Required
**Permissions**: `PermissionViewDonors`

**Query Parameters:**
- `min_amount` (float): Minimum total donated

**Response: 200 OK**

---

#### GET /api/v1/donors/lapsed
**Description**: Get lapsed donors (haven't donated recently)
**Authentication**: Required
**Permissions**: `PermissionViewDonors`

**Query Parameters:**
- `months` (int): Months since last donation

**Response: 200 OK**

---

#### GET /api/v1/donors/statistics
**Description**: Get donor statistics
**Authentication**: Required
**Permissions**: `PermissionViewDonors`

**Response: 200 OK**
```json
{
  "total_donors": 500,
  "active_donors": 350,
  "new_this_month": 15,
  "by_type": {
    "individual": 400,
    "organization": 80,
    "foundation": 20
  },
  "average_donation": 125.50,
  "lifetime_value_average": 2500.00
}
```

---

#### POST /api/v1/donors/:id/engagement
**Description**: Update donor engagement metrics
**Authentication**: Required
**Permissions**: `PermissionUpdateDonors`

**Request Body:**
```json
{
  "engagement_type": "event_attended",
  "notes": "Attended fundraising gala"
}
```

**Response: 200 OK**

---

#### GET /api/v1/donors/:id/donations
**Description**: Get donations by donor
**Authentication**: Required
**Permissions**: `PermissionViewDonations`

**Response: 200 OK**

---

## Donation Management

### Donation Structure

```json
{
  "id": "507f1f77bcf86cd799439019",
  "donor_id": "507f1f77bcf86cd799439018",
  "campaign_id": "507f1f77bcf86cd79943901a",
  "amount": 100.00,
  "currency": "USD",
  "donation_date": "2025-11-08T00:00:00Z",
  "donation_type": "monetary",
  "payment_method": "credit_card",
  "payment_status": "completed",
  "transaction_id": "TXN-2025-11-08-001",
  "is_recurring": true,
  "recurring_frequency": "monthly",
  "recurring_day_of_month": 1,
  "next_recurring_date": "2025-12-01T00:00:00Z",
  "is_anonymous": false,
  "designation": "general_fund",
  "in_honor_of": "",
  "in_memory_of": "",
  "dedication_notify_name": "",
  "dedication_notify_email": "",
  "tax_deductible": true,
  "tax_receipt_sent": false,
  "tax_receipt_sent_date": null,
  "tax_receipt_number": "",
  "thank_you_sent": false,
  "thank_you_sent_date": null,
  "notes": "Monthly supporter",
  "source": "website",
  "campaign_name": "Annual Fundraiser 2025",
  "matched_by": "",
  "match_multiplier": 1,
  "fees": 2.50,
  "net_amount": 97.50,
  "refunded": false,
  "refund_date": null,
  "refund_reason": "",
  "created_by": "507f1f77bcf86cd799439011",
  "created_at": "2025-11-08T00:00:00Z",
  "updated_at": "2025-11-08T00:00:00Z"
}
```

### Donation Endpoints

#### GET /api/v1/donations
**Description**: List donations
**Authentication**: Required
**Permissions**: `PermissionViewDonations`

**Query Parameters:**
- `limit`, `offset`: Pagination
- `donor_id` (string): Filter by donor
- `campaign_id` (string): Filter by campaign
- `payment_status` (string): Filter by payment status
- `start_date`, `end_date`: Date range

**Response: 200 OK**

---

#### GET /api/v1/donations/:id
**Description**: Get donation by ID
**Authentication**: Required
**Permissions**: `PermissionViewDonations`

**Response: 200 OK**

---

#### POST /api/v1/donations
**Description**: Create donation
**Authentication**: Required
**Permissions**: `PermissionCreateDonations`

**Request Body:** (See Donation Structure)

**Response: 201 Created**

---

#### PUT /api/v1/donations/:id
**Description**: Update donation
**Authentication**: Required
**Permissions**: `PermissionUpdateDonations`

**Response: 200 OK**

---

#### DELETE /api/v1/donations/:id
**Description**: Delete donation
**Authentication**: Required
**Permissions**: `PermissionDeleteDonations`

**Response: 200 OK**

---

#### GET /api/v1/donations/recurring
**Description**: Get recurring donations
**Authentication**: Required
**Permissions**: `PermissionViewDonations`

**Response: 200 OK**

---

#### GET /api/v1/donations/pending-thank-yous
**Description**: Get donations needing thank you notes
**Authentication**: Required
**Permissions**: `PermissionViewDonations`

**Response: 200 OK**

---

#### GET /api/v1/donations/statistics
**Description**: Get donation statistics
**Authentication**: Required
**Permissions**: `PermissionViewDonations`

**Response: 200 OK**
```json
{
  "total_raised": 50000.00,
  "total_donations": 500,
  "this_month": 5000.00,
  "this_year": 35000.00,
  "average_donation": 100.00,
  "recurring_revenue_monthly": 2500.00,
  "by_payment_method": {
    "credit_card": 30000.00,
    "bank_transfer": 15000.00,
    "cash": 5000.00
  }
}
```

---

#### GET /api/v1/donations/date-range
**Description**: Get donations within date range
**Authentication**: Required
**Permissions**: `PermissionViewDonations`

**Query Parameters:**
- `start_date` (string): ISO date
- `end_date` (string): ISO date

**Response: 200 OK**

---

#### POST /api/v1/donations/:id/process
**Description**: Process a donation payment
**Authentication**: Required
**Permissions**: `PermissionUpdateDonations`

**Request Body:**
```json
{
  "transaction_id": "TXN-2025-11-08-001",
  "payment_method": "credit_card"
}
```

**Response: 200 OK**

---

#### POST /api/v1/donations/:id/refund
**Description**: Refund a donation
**Authentication**: Required
**Permissions**: `PermissionDeleteDonations`

**Request Body:**
```json
{
  "reason": "Requested by donor"
}
```

**Response: 200 OK**

---

#### POST /api/v1/donations/:id/thank-you
**Description**: Mark thank you as sent
**Authentication**: Required
**Permissions**: `PermissionUpdateDonations`

**Response: 200 OK**

---

#### POST /api/v1/donations/:id/tax-receipt
**Description**: Generate and send tax receipt
**Authentication**: Required
**Permissions**: `PermissionUpdateDonations`

**Response: 200 OK**

---

## Campaign Management

### Campaign Structure

**✅ VERIFIED - Tested and working**

```json
{
  "id": "507f1f77bcf86cd79943901a",
  "name": {
    "en": "Winter Care Campaign 2025",
    "pl": "Kampania Zimowej Opieki 2025"
  },
  "description": {
    "en": "Raising funds for winter animal care and shelter heating",
    "pl": "Zbieramy fundusze na zimową opiekę nad zwierzętami"
  },
  "type": "general",
  "status": "active",
  "goal_amount": 5000.00,
  "current_amount": 0.00,
  "donor_count": 0,
  "donation_count": 0,
  "average_donation": 0.00,
  "manager": "507f1f77bcf86cd799439011",
  "start_date": "2025-11-08T20:00:00Z",
  "end_date": "2026-01-08T20:00:00Z",
  "public": true,
  "featured": false,
  "image_url": "",
  "video_url": "",
  "tags": [],
  "contact_email": "",
  "contact_phone": "",
  "notes": "",
  "created_by": "507f1f77bcf86cd799439011",
  "updated_by": "507f1f77bcf86cd799439011",
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-11-08T15:30:00Z"
}
```

**Valid Campaign Types:**
- `general` - General fundraising campaign
- `capital` - Capital/building campaign
- `emergency` - Emergency response campaign
- `annual` - Annual fundraising drive
- `event` - Event-based campaign
- `membership` - Membership drive
- `end_of_year` - Year-end giving campaign

**Key Differences from Original Documentation:**
1. ✅ `name` and `description` are **MultilingualName** objects with `en` and `pl` keys
2. ✅ `campaign_type` → `type` (field renamed)
3. ✅ `manager_id` → `manager` (REQUIRED field, renamed)
4. ✅ `is_public` → `public` (field renamed)
5. ✅ `is_featured` → `featured` (field renamed)
6. ✅ `raised_amount` → `current_amount` (field renamed)
7. ✅ Added: `donor_count`, `donation_count`, `average_donation`, `updated_by`
8. ❌ Removed: `slug`, `currency`, `team_members`, `updates`, `thank_you_message`, `seo_*`, `categories`, `progress_percentage`

### Campaign Endpoints

#### GET /api/v1/campaigns
**Description**: List campaigns
**Authentication**: Required
**Permissions**: `PermissionViewCampaigns`

**Query Parameters:**
- `limit`, `offset`: Pagination
- `status` (string): Filter by status
- `type` (string): Filter by campaign type
- `public` (bool): Filter public campaigns
- `featured` (bool): Filter featured campaigns

**Response: 200 OK**

---

#### GET /api/v1/campaigns/:id
**Description**: Get campaign by ID
**Authentication**: Required
**Permissions**: `PermissionViewCampaigns`

**Response: 200 OK**

---

#### POST /api/v1/campaigns
**Description**: Create campaign
**Authentication**: Required
**Permissions**: `PermissionCreateCampaigns`

**Request Body:** (See Campaign Structure)

**Response: 201 Created**

---

#### PUT /api/v1/campaigns/:id
**Description**: Update campaign
**Authentication**: Required
**Permissions**: `PermissionUpdateCampaigns`

**Response: 200 OK**

---

#### DELETE /api/v1/campaigns/:id
**Description**: Delete campaign
**Authentication**: Required
**Permissions**: `PermissionDeleteCampaigns`

**Response: 200 OK**

---

#### GET /api/v1/campaigns/active
**Description**: Get active campaigns
**Authentication**: Required
**Permissions**: `PermissionViewCampaigns`

**Response: 200 OK**

---

#### GET /api/v1/campaigns/featured
**Description**: Get featured campaigns
**Authentication**: Required
**Permissions**: `PermissionViewCampaigns`

**Response: 200 OK**

---

#### GET /api/v1/campaigns/public
**Description**: Get public campaigns
**Authentication**: Required
**Permissions**: `PermissionViewCampaigns`

**Response: 200 OK**

---

#### GET /api/v1/campaigns/statistics
**Description**: Get campaign statistics
**Authentication**: Required
**Permissions**: `PermissionViewCampaigns`

**Response: 200 OK**

---

#### POST /api/v1/campaigns/:id/activate
**Description**: Activate campaign
**Authentication**: Required
**Permissions**: `PermissionUpdateCampaigns`

**Response: 200 OK**

---

#### POST /api/v1/campaigns/:id/pause
**Description**: Pause campaign
**Authentication**: Required
**Permissions**: `PermissionUpdateCampaigns`

**Response: 200 OK**

---

#### POST /api/v1/campaigns/:id/complete
**Description**: Complete campaign
**Authentication**: Required
**Permissions**: `PermissionUpdateCampaigns`

**Response: 200 OK**

---

#### POST /api/v1/campaigns/:id/cancel
**Description**: Cancel campaign
**Authentication**: Required
**Permissions**: `PermissionUpdateCampaigns`

**Response: 200 OK**

---

#### GET /api/v1/campaigns/:id/donations
**Description**: Get donations for campaign
**Authentication**: Required
**Permissions**: `PermissionViewDonations`

**Response: 200 OK**

---

#### GET /api/v1/campaigns/:id/communications
**Description**: Get communications for campaign
**Authentication**: Required
**Permissions**: `PermissionViewCommunications`

**Response: 200 OK**

---

#### GET /api/v1/users/:id/campaigns
**Description**: Get campaigns managed by user
**Authentication**: Required
**Permissions**: `PermissionViewCampaigns`

**Response: 200 OK**

---

## Event Management

### Event Structure

```json
{
  "id": "507f1f77bcf86cd79943901b",
  "name": "Annual Adoption Fair",
  "slug": "annual-adoption-fair-2025",
  "description": "Meet adoptable animals and learn about our foundation",
  "event_type": "adoption_fair",
  "status": "active",
  "start_date": "2025-12-15T10:00:00Z",
  "end_date": "2025-12-15T16:00:00Z",
  "timezone": "America/New_York",
  "location": {
    "name": "City Park",
    "address": "789 Park Ave, New York, NY 10001",
    "coordinates": {
      "latitude": 40.7128,
      "longitude": -74.0060
    }
  },
  "is_virtual": false,
  "virtual_link": "",
  "is_public": true,
  "is_featured": true,
  "capacity": 200,
  "current_attendance": 45,
  "registration_required": true,
  "registration_deadline": "2025-12-14T23:59:59Z",
  "registration_fee": 0.00,
  "volunteers_needed": 20,
  "volunteers_assigned": 12,
  "assigned_volunteers": [
    {
      "volunteer_id": "507f1f77bcf86cd79943901c",
      "role": "Setup",
      "confirmed": true
    }
  ],
  "image_url": "https://example.com/events/adoption-fair.jpg",
  "organizer_id": "507f1f77bcf86cd799439011",
  "contact_email": "events@foundation.org",
  "contact_phone": "+1234567890",
  "tags": ["adoption", "community", "family-friendly"],
  "notes": "Bring weather-appropriate clothing",
  "cancellation_reason": "",
  "created_by": "507f1f77bcf86cd799439011",
  "created_at": "2025-10-01T00:00:00Z",
  "updated_at": "2025-11-08T15:30:00Z"
}
```

### Event Endpoints

#### GET /api/v1/events
**Description**: List events
**Authentication**: Required
**Permissions**: `PermissionViewEvents`

**Query Parameters:**
- `limit`, `offset`: Pagination
- `status` (string): Filter by status
- `event_type` (string): Filter by type
- `is_public` (bool): Filter public events
- `start_date`, `end_date`: Date range

**Response: 200 OK**

---

#### GET /api/v1/events/:id
**Description**: Get event by ID
**Authentication**: Required
**Permissions**: `PermissionViewEvents`

**Response: 200 OK**

---

#### POST /api/v1/events
**Description**: Create event
**Authentication**: Required
**Permissions**: `PermissionCreateEvents`

**Request Body:** (See Event Structure)

**Response: 201 Created**

---

#### PUT /api/v1/events/:id
**Description**: Update event
**Authentication**: Required
**Permissions**: `PermissionUpdateEvents`

**Response: 200 OK**

---

#### DELETE /api/v1/events/:id
**Description**: Delete event
**Authentication**: Required
**Permissions**: `PermissionDeleteEvents`

**Response: 200 OK**

---

#### GET /api/v1/events/upcoming
**Description**: Get upcoming events
**Authentication**: Required
**Permissions**: `PermissionViewEvents`

**Response: 200 OK**

---

#### GET /api/v1/events/active
**Description**: Get active events
**Authentication**: Required
**Permissions**: `PermissionViewEvents`

**Response: 200 OK**

---

#### GET /api/v1/events/public
**Description**: Get public events
**Authentication**: Required
**Permissions**: `PermissionViewEvents`

**Response: 200 OK**

---

#### GET /api/v1/events/featured
**Description**: Get featured events
**Authentication**: Required
**Permissions**: `PermissionViewEvents`

**Response: 200 OK**

---

#### GET /api/v1/events/needing-volunteers
**Description**: Get events needing volunteers
**Authentication**: Required
**Permissions**: `PermissionViewEvents`

**Response: 200 OK**

---

#### GET /api/v1/events/statistics
**Description**: Get event statistics
**Authentication**: Required
**Permissions**: `PermissionViewEvents`

**Response: 200 OK**

---

#### POST /api/v1/events/:id/activate
**Description**: Activate event
**Authentication**: Required
**Permissions**: `PermissionUpdateEvents`

**Response: 200 OK**

---

#### POST /api/v1/events/:id/complete
**Description**: Complete event
**Authentication**: Required
**Permissions**: `PermissionUpdateEvents`

**Response: 200 OK**

---

#### POST /api/v1/events/:id/cancel
**Description**: Cancel event
**Authentication**: Required
**Permissions**: `PermissionUpdateEvents`

**Request Body:**
```json
{
  "reason": "Weather conditions"
}
```

**Response: 200 OK**

---

#### POST /api/v1/events/:id/assign-volunteer
**Description**: Assign volunteer to event
**Authentication**: Required
**Permissions**: `PermissionUpdateEvents`

**Request Body:**
```json
{
  "volunteer_id": "507f1f77bcf86cd79943901c",
  "role": "Setup"
}
```

**Response: 200 OK**

---

#### POST /api/v1/events/:id/unassign-volunteer
**Description**: Unassign volunteer from event
**Authentication**: Required
**Permissions**: `PermissionUpdateEvents`

**Request Body:**
```json
{
  "volunteer_id": "507f1f77bcf86cd79943901c"
}
```

**Response: 200 OK**

---

## Volunteer Management

### Volunteer Structure

```json
{
  "id": "507f1f77bcf86cd79943901c",
  "user_id": "507f1f77bcf86cd799439012",
  "first_name": "Alice",
  "last_name": "Volunteer",
  "email": "alice@example.com",
  "phone": "+1234567890",
  "date_of_birth": "1990-05-15T00:00:00Z",
  "address": {
    "street": "123 Helper St",
    "city": "Boston",
    "state": "MA",
    "postal_code": "02101",
    "country": "USA"
  },
  "emergency_contact": {
    "name": "Bob Volunteer",
    "relationship": "Spouse",
    "phone": "+1234567891"
  },
  "application_date": "2024-01-15T00:00:00Z",
  "status": "active",
  "background_check_status": "completed",
  "background_check_date": "2024-02-01T00:00:00Z",
  "background_check_expiry": "2026-02-01T00:00:00Z",
  "orientation_completed": true,
  "orientation_date": "2024-02-15T00:00:00Z",
  "skills": ["animal_care", "event_planning", "fundraising"],
  "interests": ["dogs", "cats", "events"],
  "availability": {
    "monday": {"available": true, "times": ["morning", "afternoon"]},
    "tuesday": {"available": false, "times": []},
    "wednesday": {"available": true, "times": ["evening"]},
    "thursday": {"available": false, "times": []},
    "friday": {"available": true, "times": ["afternoon"]},
    "saturday": {"available": true, "times": ["morning", "afternoon"]},
    "sunday": {"available": false, "times": []}
  },
  "total_hours": 150.5,
  "hours_this_month": 12.0,
  "hours_this_year": 145.0,
  "assignments": [
    {
      "event_id": "507f1f77bcf86cd79943901b",
      "role": "Setup",
      "hours": 4.0,
      "date": "2025-12-15T00:00:00Z"
    }
  ],
  "certifications": [
    {
      "name": "Animal Handling",
      "issued_date": "2024-03-01T00:00:00Z",
      "expiry_date": "2026-03-01T00:00:00Z"
    }
  ],
  "commendations": [
    {
      "date": "2025-06-01T00:00:00Z",
      "reason": "Outstanding service",
      "given_by": "507f1f77bcf86cd799439011"
    }
  ],
  "warnings": [],
  "notes": "Very reliable and skilled with animals",
  "created_by": "507f1f77bcf86cd799439011",
  "created_at": "2024-01-15T00:00:00Z",
  "updated_at": "2025-11-08T15:30:00Z"
}
```

### Volunteer Endpoints

#### GET /api/v1/volunteers
**Description**: List volunteers
**Authentication**: Required
**Permissions**: `PermissionViewVolunteers`

**Query Parameters:**
- `limit`, `offset`: Pagination
- `status` (string): Filter by status
- `skill` (string): Filter by skill

**Response: 200 OK**

---

#### GET /api/v1/volunteers/:id
**Description**: Get volunteer by ID
**Authentication**: Required
**Permissions**: `PermissionViewVolunteers`

**Response: 200 OK**

---

#### POST /api/v1/volunteers
**Description**: Create volunteer record
**Authentication**: Required
**Permissions**: `PermissionCreateVolunteers`

**Request Body:** (See Volunteer Structure)

**Response: 201 Created**

---

#### PUT /api/v1/volunteers/:id
**Description**: Update volunteer
**Authentication**: Required
**Permissions**: `PermissionUpdateVolunteers`

**Response: 200 OK**

---

#### DELETE /api/v1/volunteers/:id
**Description**: Delete volunteer
**Authentication**: Required
**Permissions**: `PermissionDeleteVolunteers`

**Response: 200 OK**

---

#### GET /api/v1/volunteers/active
**Description**: Get active volunteers
**Authentication**: Required
**Permissions**: `PermissionViewVolunteers`

**Response: 200 OK**

---

#### GET /api/v1/volunteers/by-skill
**Description**: Get volunteers by skill
**Authentication**: Required
**Permissions**: `PermissionViewVolunteers`

**Query Parameters:**
- `skill` (string): Skill name

**Response: 200 OK**

---

#### GET /api/v1/volunteers/needing-background-check
**Description**: Get volunteers needing background check
**Authentication**: Required
**Permissions**: `PermissionViewVolunteers`

**Response: 200 OK**

---

#### GET /api/v1/volunteers/top
**Description**: Get top volunteers by hours
**Authentication**: Required
**Permissions**: `PermissionViewVolunteers`

**Query Parameters:**
- `limit` (int): Number of volunteers (default: 10)

**Response: 200 OK**

---

#### GET /api/v1/volunteers/statistics
**Description**: Get volunteer statistics
**Authentication**: Required
**Permissions**: `PermissionViewVolunteers`

**Response: 200 OK**

---

#### POST /api/v1/volunteers/:id/approve
**Description**: Approve volunteer application
**Authentication**: Required
**Permissions**: `PermissionUpdateVolunteers`

**Response: 200 OK**

---

#### POST /api/v1/volunteers/:id/suspend
**Description**: Suspend volunteer
**Authentication**: Required
**Permissions**: `PermissionUpdateVolunteers`

**Request Body:**
```json
{
  "reason": "Policy violation"
}
```

**Response: 200 OK**

---

#### POST /api/v1/volunteers/:id/add-hours
**Description**: Add volunteer hours
**Authentication**: Required
**Permissions**: `PermissionUpdateVolunteers`

**Request Body:**
```json
{
  "hours": 4.0,
  "date": "2025-11-08T00:00:00Z",
  "activity": "Animal care",
  "notes": "Walked dogs"
}
```

**Response: 200 OK**

---

#### POST /api/v1/volunteers/:id/commendation
**Description**: Add commendation
**Authentication**: Required
**Permissions**: `PermissionUpdateVolunteers`

**Request Body:**
```json
{
  "reason": "Exceptional dedication"
}
```

**Response: 200 OK**

---

#### POST /api/v1/volunteers/:id/warning
**Description**: Add warning
**Authentication**: Required
**Permissions**: `PermissionUpdateVolunteers`

**Request Body:**
```json
{
  "reason": "Late arrival",
  "severity": "minor"
}
```

**Response: 200 OK**

---

#### POST /api/v1/volunteers/:id/certification
**Description**: Add certification
**Authentication**: Required
**Permissions**: `PermissionUpdateVolunteers`

**Request Body:**
```json
{
  "name": "First Aid",
  "issued_date": "2025-11-08T00:00:00Z",
  "expiry_date": "2027-11-08T00:00:00Z"
}
```

**Response: 200 OK**

---

## Communication Management

### Communication Structure

```json
{
  "id": "507f1f77bcf86cd79943901d",
  "type": "email",
  "channel": "email",
  "recipient_type": "individual",
  "recipient_id": "507f1f77bcf86cd799439018",
  "recipient_email": "john.donor@example.com",
  "recipient_phone": "",
  "subject": "Thank you for your donation",
  "body": "Dear John, Thank you for your generous donation...",
  "template_id": "507f1f77bcf86cd79943901e",
  "campaign_id": "507f1f77bcf86cd79943901a",
  "batch_id": "",
  "status": "sent",
  "scheduled_date": null,
  "sent_date": "2025-11-08T10:00:00Z",
  "opened_date": "2025-11-08T11:30:00Z",
  "clicked_date": "2025-11-08T11:35:00Z",
  "bounced": false,
  "bounce_reason": "",
  "error_message": "",
  "retry_count": 0,
  "metadata": {
    "donation_id": "507f1f77bcf86cd799439019",
    "amount": "100.00"
  },
  "tracking": {
    "opens": 1,
    "clicks": 2,
    "last_opened": "2025-11-08T11:30:00Z",
    "last_clicked": "2025-11-08T11:35:00Z"
  },
  "created_by": "507f1f77bcf86cd799439011",
  "created_at": "2025-11-08T10:00:00Z",
  "updated_at": "2025-11-08T11:35:00Z"
}
```

### Email Template Structure

```json
{
  "id": "507f1f77bcf86cd79943901e",
  "name": "Donation Thank You",
  "slug": "donation-thank-you",
  "category": "donor_relations",
  "subject": "Thank you for your donation to {{foundation_name}}",
  "body_html": "<html><body>Dear {{first_name}},...</body></html>",
  "body_text": "Dear {{first_name}},...",
  "variables": ["first_name", "foundation_name", "donation_amount", "donation_date"],
  "is_active": true,
  "is_default": false,
  "language": "en",
  "tags": ["donation", "thank-you"],
  "created_by": "507f1f77bcf86cd799439011",
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-11-08T15:30:00Z"
}
```

### Communication Endpoints

#### GET /api/v1/communications
**Description**: List communications
**Authentication**: Required
**Permissions**: `PermissionViewCommunications`

**Query Parameters:**
- `limit`, `offset`: Pagination
- `type` (string): Filter by type
- `status` (string): Filter by status
- `recipient_id` (string): Filter by recipient

**Response: 200 OK**

---

#### GET /api/v1/communications/:id
**Description**: Get communication by ID
**Authentication**: Required
**Permissions**: `PermissionViewCommunications`

**Response: 200 OK**

---

#### POST /api/v1/communications
**Description**: Create and send communication
**Authentication**: Required
**Permissions**: `PermissionCreateCommunications`

**Request Body:**
```json
{
  "type": "email",
  "recipient_type": "individual",
  "recipient_id": "507f1f77bcf86cd799439018",
  "recipient_email": "john@example.com",
  "subject": "Subject",
  "body": "Message body",
  "scheduled_date": null
}
```

**Response: 201 Created**

---

#### POST /api/v1/communications/send-from-template
**Description**: Send communication from template
**Authentication**: Required
**Permissions**: `PermissionCreateCommunications`

**Request Body:**
```json
{
  "template_id": "507f1f77bcf86cd79943901e",
  "recipient_id": "507f1f77bcf86cd799439018",
  "recipient_email": "john@example.com",
  "variables": {
    "first_name": "John",
    "donation_amount": "100.00"
  }
}
```

**Response: 201 Created**

---

#### PUT /api/v1/communications/:id/status
**Description**: Update communication status
**Authentication**: Required
**Permissions**: `PermissionUpdateCommunications`

**Request Body:**
```json
{
  "status": "sent"
}
```

**Response: 200 OK**

---

#### DELETE /api/v1/communications/:id
**Description**: Delete communication
**Authentication**: Required
**Permissions**: `PermissionDeleteCommunications`

**Response: 200 OK**

---

#### GET /api/v1/communications/pending
**Description**: Get pending communications
**Authentication**: Required
**Permissions**: `PermissionViewCommunications`

**Response: 200 OK**

---

#### GET /api/v1/communications/retry
**Description**: Get communications for retry
**Authentication**: Required
**Permissions**: `PermissionViewCommunications`

**Response: 200 OK**

---

#### GET /api/v1/communications/by-recipient
**Description**: Get communications by recipient
**Authentication**: Required
**Permissions**: `PermissionViewCommunications`

**Query Parameters:**
- `recipient_id` (string): Recipient ID

**Response: 200 OK**

---

#### GET /api/v1/communications/statistics
**Description**: Get communication statistics
**Authentication**: Required
**Permissions**: `PermissionViewCommunications`

**Response: 200 OK**

---

#### POST /api/v1/communications/:id/track-open
**Description**: Track email open (webhook)
**Authentication**: Required
**Permissions**: None

**Response: 200 OK**

---

#### POST /api/v1/communications/:id/track-click
**Description**: Track email click (webhook)
**Authentication**: Required
**Permissions**: None

**Response: 200 OK**

---

#### GET /api/v1/campaigns/:id/communications
**Description**: Get communications for campaign
**Authentication**: Required
**Permissions**: `PermissionViewCommunications`

**Response: 200 OK**

---

#### GET /api/v1/batches/:id/communications
**Description**: Get communications for batch
**Authentication**: Required
**Permissions**: `PermissionViewCommunications`

**Response: 200 OK**

---

### Template Endpoints

#### GET /api/v1/templates
**Description**: List email templates
**Authentication**: Required
**Permissions**: `PermissionViewTemplates`

**Response: 200 OK**

---

#### GET /api/v1/templates/:id
**Description**: Get template by ID
**Authentication**: Required
**Permissions**: `PermissionViewTemplates`

**Response: 200 OK**

---

#### POST /api/v1/templates
**Description**: Create template
**Authentication**: Required
**Permissions**: `PermissionCreateTemplates`

**Request Body:** (See Email Template Structure)

**Response: 201 Created**

---

#### PUT /api/v1/templates/:id
**Description**: Update template
**Authentication**: Required
**Permissions**: `PermissionUpdateTemplates`

**Response: 200 OK**

---

#### DELETE /api/v1/templates/:id
**Description**: Delete template
**Authentication**: Required
**Permissions**: `PermissionDeleteTemplates`

**Response: 200 OK**

---

#### GET /api/v1/templates/active
**Description**: Get active templates
**Authentication**: Required
**Permissions**: `PermissionViewTemplates`

**Response: 200 OK**

---

#### GET /api/v1/templates/by-category
**Description**: Get templates by category
**Authentication**: Required
**Permissions**: `PermissionViewTemplates`

**Query Parameters:**
- `category` (string): Category name

**Response: 200 OK**

---

#### GET /api/v1/templates/default
**Description**: Get default template for category
**Authentication**: Required
**Permissions**: `PermissionViewTemplates`

**Query Parameters:**
- `category` (string): Category name

**Response: 200 OK**

---

## Notification System

### Notification Structure

```json
{
  "id": "507f1f77bcf86cd79943901f",
  "user_id": "507f1f77bcf86cd799439011",
  "type": "info",
  "category": "system",
  "title": "New Adoption Application",
  "message": "A new adoption application has been submitted for Buddy",
  "action_url": "/adoptions/applications/507f1f77bcf86cd799439016",
  "action_text": "View Application",
  "priority": "normal",
  "channels": ["in_app", "email"],
  "is_read": false,
  "read_at": null,
  "is_dismissed": false,
  "dismissed_at": null,
  "metadata": {
    "application_id": "507f1f77bcf86cd799439016",
    "animal_id": "507f1f77bcf86cd799439013"
  },
  "expires_at": "2025-12-08T00:00:00Z",
  "created_at": "2025-11-08T10:00:00Z",
  "updated_at": "2025-11-08T10:00:00Z"
}
```

### Notification Endpoints

#### GET /api/v1/notifications/me
**Description**: Get current user's notifications
**Authentication**: Required
**Permissions**: Authenticated user

**Response: 200 OK**

---

#### GET /api/v1/notifications/unread
**Description**: Get unread notifications
**Authentication**: Required
**Permissions**: Authenticated user

**Response: 200 OK**

---

#### GET /api/v1/notifications/unread/count
**Description**: Get unread notification count
**Authentication**: Required
**Permissions**: Authenticated user

**Response: 200 OK**
```json
{
  "count": 5
}
```

---

#### GET /api/v1/notifications
**Description**: List all notifications (admin)
**Authentication**: Required
**Permissions**: `PermissionViewNotifications`

**Response: 200 OK**

---

#### GET /api/v1/notifications/:id
**Description**: Get notification by ID
**Authentication**: Required
**Permissions**: `PermissionViewNotifications`

**Response: 200 OK**

---

#### POST /api/v1/notifications
**Description**: Create notification
**Authentication**: Required
**Permissions**: `PermissionCreateNotifications`

**Request Body:**
```json
{
  "user_id": "507f1f77bcf86cd799439011",
  "type": "info",
  "category": "system",
  "title": "Notification Title",
  "message": "Notification message",
  "action_url": "/path/to/action",
  "priority": "normal",
  "channels": ["in_app", "email"]
}
```

**Response: 201 Created**

---

#### POST /api/v1/notifications/:id/read
**Description**: Mark notification as read
**Authentication**: Required
**Permissions**: Authenticated user (own notifications)

**Response: 200 OK**

---

#### POST /api/v1/notifications/:id/unread
**Description**: Mark notification as unread
**Authentication**: Required
**Permissions**: Authenticated user (own notifications)

**Response: 200 OK**

---

#### POST /api/v1/notifications/read-all
**Description**: Mark all notifications as read
**Authentication**: Required
**Permissions**: Authenticated user

**Response: 200 OK**

---

#### POST /api/v1/notifications/:id/dismiss
**Description**: Dismiss notification
**Authentication**: Required
**Permissions**: Authenticated user (own notifications)

**Response: 200 OK**

---

#### DELETE /api/v1/notifications/:id
**Description**: Delete notification
**Authentication**: Required
**Permissions**: Authenticated user (own notifications)

**Response: 200 OK**

---

## Report Generation

### Report Structure

```json
{
  "id": "507f1f77bcf86cd799439020",
  "name": "Monthly Adoption Report",
  "slug": "monthly-adoption-report",
  "description": "Monthly summary of adoptions",
  "category": "adoptions",
  "report_type": "scheduled",
  "query": "SELECT * FROM adoptions WHERE...",
  "parameters": [
    {
      "name": "month",
      "type": "string",
      "required": true,
      "default_value": ""
    }
  ],
  "output_format": ["pdf", "excel", "csv"],
  "schedule": {
    "enabled": true,
    "frequency": "monthly",
    "day_of_month": 1,
    "time": "09:00",
    "timezone": "America/New_York"
  },
  "recipients": ["admin@example.com"],
  "is_active": true,
  "is_public": false,
  "created_by": "507f1f77bcf86cd799439011",
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-11-08T15:30:00Z"
}
```

### Report Endpoints

#### GET /api/v1/reports
**Description**: List reports
**Authentication**: Required
**Permissions**: `PermissionViewReports`

**Response: 200 OK**

---

#### GET /api/v1/reports/:id
**Description**: Get report by ID
**Authentication**: Required
**Permissions**: `PermissionViewReports`

**Response: 200 OK**

---

#### POST /api/v1/reports
**Description**: Create report
**Authentication**: Required
**Permissions**: `PermissionCreateReports`

**Request Body:** (See Report Structure)

**Response: 201 Created**

---

#### PUT /api/v1/reports/:id
**Description**: Update report
**Authentication**: Required
**Permissions**: `PermissionUpdateReports`

**Response: 200 OK**

---

#### DELETE /api/v1/reports/:id
**Description**: Delete report
**Authentication**: Required
**Permissions**: `PermissionDeleteReports`

**Response: 200 OK**

---

#### GET /api/v1/reports/active
**Description**: Get active reports
**Authentication**: Required
**Permissions**: `PermissionViewReports`

**Response: 200 OK**

---

#### GET /api/v1/reports/public
**Description**: Get public reports
**Authentication**: Required
**Permissions**: `PermissionViewReports`

**Response: 200 OK**

---

#### POST /api/v1/reports/:id/execute
**Description**: Execute report
**Authentication**: Required
**Permissions**: `PermissionViewReports`

**Request Body:**
```json
{
  "parameters": {
    "month": "2025-11"
  },
  "format": "pdf"
}
```

**Response: 200 OK**
```json
{
  "execution_id": "507f1f77bcf86cd799439021",
  "status": "completed",
  "file_url": "https://example.com/reports/monthly-adoption-2025-11.pdf"
}
```

---

#### GET /api/v1/reports/:id/executions
**Description**: Get report executions
**Authentication**: Required
**Permissions**: `PermissionViewReports`

**Response: 200 OK**

---

#### GET /api/v1/reports/executions/recent
**Description**: Get recent report executions
**Authentication**: Required
**Permissions**: `PermissionViewReports`

**Response: 200 OK**

---

## Dashboard & Analytics

### Dashboard Endpoints

#### GET /api/v1/dashboard
**Description**: Get full dashboard metrics
**Authentication**: Required
**Permissions**: `PermissionViewDashboard`

**Response: 200 OK**
```json
{
  "overview": {
    "total_animals": 150,
    "available_animals": 45,
    "animals_adopted_this_month": 12,
    "animals_in_treatment": 15,
    "total_donations_this_month": 5000.00,
    "total_volunteers": 50,
    "active_volunteers": 35,
    "upcoming_events": 3
  },
  "animals": {
    "by_status": {
      "available": 45,
      "adopted": 80,
      "under_treatment": 15,
      "fostered": 10
    },
    "by_species": {
      "dog": 90,
      "cat": 50,
      "other": 10
    },
    "recent_intakes": 8,
    "pending_adoptions": 15
  },
  "financial": {
    "total_donations_ytd": 35000.00,
    "recurring_revenue_monthly": 2500.00,
    "average_donation": 100.00,
    "top_campaigns": [
      {
        "id": "507f1f77bcf86cd79943901a",
        "name": "Annual Fundraiser",
        "raised": 35000.00,
        "goal": 50000.00
      }
    ]
  },
  "volunteers": {
    "total_hours_this_month": 250.0,
    "total_hours_ytd": 2500.0,
    "active_count": 35,
    "needing_background_check": 3
  },
  "events": {
    "upcoming_count": 3,
    "this_month_count": 5,
    "needing_volunteers_count": 2
  },
  "recent_activity": [
    {
      "type": "adoption",
      "message": "Buddy was adopted by Jane Smith",
      "timestamp": "2025-11-08T10:00:00Z"
    }
  ]
}
```

---

#### GET /api/v1/dashboard/overview
**Description**: Get overview metrics only
**Authentication**: Required
**Permissions**: `PermissionViewDashboard`

**Response: 200 OK**

---

## Settings Management

### Settings Structure

```json
{
  "id": "507f1f77bcf86cd799439022",
  "foundation_name": "Animal Foundation",
  "foundation_tagline": "Saving lives, one animal at a time",
  "contact_info": {
    "email": "info@foundation.org",
    "phone": "+1234567890",
    "address": {
      "street": "123 Foundation St",
      "city": "New York",
      "state": "NY",
      "postal_code": "10001",
      "country": "USA"
    },
    "social_media": {
      "facebook": "https://facebook.com/foundation",
      "twitter": "https://twitter.com/foundation",
      "instagram": "https://instagram.com/foundation"
    }
  },
  "operating_hours": {
    "monday": {"open": "09:00", "close": "17:00", "closed": false},
    "tuesday": {"open": "09:00", "close": "17:00", "closed": false},
    "wednesday": {"open": "09:00", "close": "17:00", "closed": false},
    "thursday": {"open": "09:00", "close": "17:00", "closed": false},
    "friday": {"open": "09:00", "close": "17:00", "closed": false},
    "saturday": {"open": "10:00", "close": "16:00", "closed": false},
    "sunday": {"open": "", "close": "", "closed": true}
  },
  "email_settings": {
    "smtp_host": "smtp.example.com",
    "smtp_port": 587,
    "smtp_username": "noreply@foundation.org",
    "smtp_from_name": "Animal Foundation",
    "smtp_from_email": "noreply@foundation.org"
  },
  "notification_settings": {
    "new_adoption_application": true,
    "donation_received": true,
    "event_reminder": true,
    "volunteer_signup": true
  },
  "feature_flags": {
    "enable_online_donations": true,
    "enable_event_registration": true,
    "enable_volunteer_portal": true,
    "enable_public_animal_profiles": true
  },
  "branding": {
    "logo_url": "https://example.com/logo.png",
    "primary_color": "#007bff",
    "secondary_color": "#6c757d",
    "favicon_url": "https://example.com/favicon.ico"
  },
  "updated_by": "507f1f77bcf86cd799439011",
  "updated_at": "2025-11-08T15:30:00Z"
}
```

### Settings Endpoints

#### GET /api/v1/settings
**Description**: Get all settings
**Authentication**: Required
**Permissions**: `PermissionViewSettings`

**Response: 200 OK**

---

#### PUT /api/v1/settings
**Description**: Update settings
**Authentication**: Required
**Permissions**: `PermissionUpdateSettings`

**Request Body:** (See Settings Structure)

**Response: 200 OK**

---

#### POST /api/v1/settings/initialize
**Description**: Initialize settings (first time setup)
**Authentication**: Required
**Permissions**: Super Admin

**Response: 201 Created**

---

#### PUT /api/v1/settings/email
**Description**: Update email settings
**Authentication**: Required
**Permissions**: `PermissionUpdateSettings`

**Request Body:**
```json
{
  "smtp_host": "smtp.example.com",
  "smtp_port": 587,
  "smtp_username": "noreply@foundation.org"
}
```

**Response: 200 OK**

---

#### PUT /api/v1/settings/notifications
**Description**: Update notification settings
**Authentication**: Required
**Permissions**: `PermissionUpdateSettings`

**Response: 200 OK**

---

#### PUT /api/v1/settings/features
**Description**: Update feature flags
**Authentication**: Required
**Permissions**: `PermissionUpdateSettings`

**Response: 200 OK**

---

#### PUT /api/v1/settings/branding
**Description**: Update branding
**Authentication**: Required
**Permissions**: `PermissionUpdateSettings`

**Response: 200 OK**

---

#### GET /api/v1/settings/contact
**Description**: Get contact info (public)
**Authentication**: None
**Permissions**: None

**Response: 200 OK**

---

#### GET /api/v1/settings/hours
**Description**: Get operating hours (public)
**Authentication**: None
**Permissions**: None

**Response: 200 OK**

---

## Task Management

### Task Structure

```json
{
  "id": "507f1f77bcf86cd799439023",
  "title": "Follow up with adopter",
  "description": "Call Jane Smith to check on Buddy's adaptation",
  "status": "pending",
  "priority": "high",
  "due_date": "2025-11-15T00:00:00Z",
  "assigned_to": "507f1f77bcf86cd799439011",
  "assigned_by": "507f1f77bcf86cd799439011",
  "category": "adoption",
  "related_entity_type": "adoption",
  "related_entity_id": "507f1f77bcf86cd799439017",
  "tags": ["follow-up", "adoption"],
  "checklist": [
    {
      "id": "check1",
      "item": "Prepare questions",
      "completed": true,
      "completed_at": "2025-11-08T10:00:00Z"
    },
    {
      "id": "check2",
      "item": "Make phone call",
      "completed": false,
      "completed_at": null
    }
  ],
  "estimated_hours": 0.5,
  "actual_hours": 0.0,
  "started_at": null,
  "completed_at": null,
  "notes": "",
  "created_by": "507f1f77bcf86cd799439011",
  "created_at": "2025-11-08T10:00:00Z",
  "updated_at": "2025-11-08T10:00:00Z"
}
```

### Task Endpoints

#### GET /api/v1/tasks
**Description**: List tasks
**Authentication**: Required
**Permissions**: `PermissionViewTasks`

**Query Parameters:**
- `limit`, `offset`: Pagination
- `status` (string): Filter by status
- `priority` (string): Filter by priority
- `assigned_to` (string): Filter by assignee

**Response: 200 OK**

---

#### GET /api/v1/tasks/:id
**Description**: Get task by ID
**Authentication**: Required
**Permissions**: `PermissionViewTasks`

**Response: 200 OK**

---

#### POST /api/v1/tasks
**Description**: Create task
**Authentication**: Required
**Permissions**: `PermissionCreateTasks`

**Request Body:** (See Task Structure)

**Response: 201 Created**

---

#### PUT /api/v1/tasks/:id
**Description**: Update task
**Authentication**: Required
**Permissions**: `PermissionUpdateTasks`

**Response: 200 OK**

---

#### DELETE /api/v1/tasks/:id
**Description**: Delete task
**Authentication**: Required
**Permissions**: `PermissionDeleteTasks`

**Response: 200 OK**

---

#### GET /api/v1/tasks/my
**Description**: Get current user's tasks
**Authentication**: Required
**Permissions**: `PermissionViewTasks`

**Response: 200 OK**

---

#### GET /api/v1/tasks/overdue
**Description**: Get overdue tasks
**Authentication**: Required
**Permissions**: `PermissionViewTasks`

**Response: 200 OK**

---

#### GET /api/v1/tasks/upcoming
**Description**: Get upcoming tasks
**Authentication**: Required
**Permissions**: `PermissionViewTasks`

**Query Parameters:**
- `days` (int): Number of days ahead (default: 7)

**Response: 200 OK**

---

#### GET /api/v1/tasks/statistics
**Description**: Get task statistics
**Authentication**: Required
**Permissions**: `PermissionViewTasks`

**Response: 200 OK**

---

#### GET /api/v1/tasks/assignee/:user_id
**Description**: Get tasks by assignee
**Authentication**: Required
**Permissions**: `PermissionViewTasks`

**Response: 200 OK**

---

#### POST /api/v1/tasks/:id/assign
**Description**: Assign task
**Authentication**: Required
**Permissions**: `PermissionUpdateTasks`

**Request Body:**
```json
{
  "assigned_to": "507f1f77bcf86cd799439011"
}
```

**Response: 200 OK**

---

#### POST /api/v1/tasks/:id/unassign
**Description**: Unassign task
**Authentication**: Required
**Permissions**: `PermissionUpdateTasks`

**Response: 200 OK**

---

#### POST /api/v1/tasks/:id/start
**Description**: Start task
**Authentication**: Required
**Permissions**: `PermissionUpdateTasks`

**Response: 200 OK**

---

#### POST /api/v1/tasks/:id/complete
**Description**: Complete task
**Authentication**: Required
**Permissions**: `PermissionUpdateTasks`

**Request Body:**
```json
{
  "actual_hours": 0.5,
  "notes": "Task completed successfully"
}
```

**Response: 200 OK**

---

#### POST /api/v1/tasks/:id/cancel
**Description**: Cancel task
**Authentication**: Required
**Permissions**: `PermissionUpdateTasks`

**Request Body:**
```json
{
  "reason": "No longer needed"
}
```

**Response: 200 OK**

---

#### POST /api/v1/tasks/:id/checklist
**Description**: Add checklist item
**Authentication**: Required
**Permissions**: `PermissionUpdateTasks`

**Request Body:**
```json
{
  "item": "New checklist item"
}
```

**Response: 200 OK**

---

#### POST /api/v1/tasks/:id/checklist/:item_id/complete
**Description**: Complete checklist item
**Authentication**: Required
**Permissions**: `PermissionUpdateTasks`

**Response: 200 OK**

---

#### DELETE /api/v1/tasks/:id/checklist/:item_id
**Description**: Remove checklist item
**Authentication**: Required
**Permissions**: `PermissionUpdateTasks`

**Response: 200 OK**

---

## Document Management

### Document Structure

```json
{
  "id": "507f1f77bcf86cd799439024",
  "name": "Adoption Contract - Buddy",
  "description": "Signed adoption contract",
  "file_name": "adoption-contract-buddy.pdf",
  "file_type": "application/pdf",
  "file_size": 102400,
  "file_url": "https://example.com/documents/adoption-contract-buddy.pdf",
  "document_type": "contract",
  "category": "adoption",
  "related_entity_type": "adoption",
  "related_entity_id": "507f1f77bcf86cd799439017",
  "version": 1,
  "is_latest_version": true,
  "previous_version_id": null,
  "tags": ["adoption", "contract", "signed"],
  "is_public": false,
  "access_level": "private",
  "allowed_users": ["507f1f77bcf86cd799439011"],
  "expiration_date": "2026-11-08T00:00:00Z",
  "is_signed": true,
  "signed_by": "Jane Smith",
  "signed_date": "2025-11-08T00:00:00Z",
  "checksum": "abc123def456",
  "metadata": {
    "adopter_name": "Jane Smith",
    "animal_name": "Buddy"
  },
  "uploaded_by": "507f1f77bcf86cd799439011",
  "created_at": "2025-11-08T10:00:00Z",
  "updated_at": "2025-11-08T10:00:00Z"
}
```

### Document Endpoints

#### GET /api/v1/documents
**Description**: List documents
**Authentication**: Required
**Permissions**: `PermissionViewDocuments`

**Query Parameters:**
- `limit`, `offset`: Pagination
- `document_type` (string): Filter by type
- `category` (string): Filter by category

**Response: 200 OK**

---

#### GET /api/v1/documents/:id
**Description**: Get document by ID
**Authentication**: Required
**Permissions**: `PermissionViewDocuments`

**Response: 200 OK**

---

#### POST /api/v1/documents
**Description**: Upload document
**Authentication**: Required
**Permissions**: `PermissionCreateDocuments`

**Content-Type**: `multipart/form-data`

**Form Data:**
- `file`: Document file
- `name` (string): Document name
- `description` (string): Description
- `document_type` (string): Type
- `category` (string): Category

**Response: 201 Created**

---

#### PUT /api/v1/documents/:id
**Description**: Update document metadata
**Authentication**: Required
**Permissions**: `PermissionUpdateDocuments`

**Response: 200 OK**

---

#### DELETE /api/v1/documents/:id
**Description**: Delete document
**Authentication**: Required
**Permissions**: `PermissionDeleteDocuments`

**Response: 200 OK**

---

#### GET /api/v1/documents/public
**Description**: Get public documents
**Authentication**: None
**Permissions**: None

**Response: 200 OK**

---

#### GET /api/v1/documents/my
**Description**: Get current user's documents
**Authentication**: Required
**Permissions**: `PermissionViewDocuments`

**Response: 200 OK**

---

#### GET /api/v1/documents/expired
**Description**: Get expired documents
**Authentication**: Required
**Permissions**: `PermissionViewDocuments`

**Response: 200 OK**

---

#### GET /api/v1/documents/expiring-soon
**Description**: Get documents expiring soon
**Authentication**: Required
**Permissions**: `PermissionViewDocuments`

**Query Parameters:**
- `days` (int): Days ahead (default: 30)

**Response: 200 OK**

---

#### GET /api/v1/documents/statistics
**Description**: Get document statistics
**Authentication**: Required
**Permissions**: `PermissionViewDocuments`

**Response: 200 OK**

---

#### GET /api/v1/documents/:id/versions
**Description**: Get document versions
**Authentication**: Required
**Permissions**: `PermissionViewDocuments`

**Response: 200 OK**

---

#### GET /api/v1/documents/:id/download
**Description**: Download document
**Authentication**: Required
**Permissions**: `PermissionViewDocuments`

**Response: 200 OK** (File download)

---

#### POST /api/v1/documents/:id/versions
**Description**: Create new version
**Authentication**: Required
**Permissions**: `PermissionUpdateDocuments`

**Content-Type**: `multipart/form-data`

**Response: 201 Created**

---

#### POST /api/v1/documents/:id/grant-access
**Description**: Grant user access
**Authentication**: Required
**Permissions**: `PermissionUpdateDocuments`

**Request Body:**
```json
{
  "user_id": "507f1f77bcf86cd799439011"
}
```

**Response: 200 OK**

---

#### POST /api/v1/documents/:id/revoke-access
**Description**: Revoke user access
**Authentication**: Required
**Permissions**: `PermissionUpdateDocuments`

**Request Body:**
```json
{
  "user_id": "507f1f77bcf86cd799439011"
}
```

**Response: 200 OK**

---

#### POST /api/v1/documents/:id/make-public
**Description**: Make document public
**Authentication**: Required
**Permissions**: `PermissionUpdateDocuments`

**Response: 200 OK**

---

#### POST /api/v1/documents/:id/make-private
**Description**: Make document private
**Authentication**: Required
**Permissions**: `PermissionUpdateDocuments`

**Response: 200 OK**

---

#### POST /api/v1/documents/:id/set-expiration
**Description**: Set expiration date
**Authentication**: Required
**Permissions**: `PermissionUpdateDocuments`

**Request Body:**
```json
{
  "expiration_date": "2026-11-08T00:00:00Z"
}
```

**Response: 200 OK**

---

## Partner Management

### Partner Structure

**✅ VERIFIED - Tested and working**

```json
{
  "id": "507f1f77bcf86cd799439025",
  "name": "Animal Rescue Network",
  "legal_name": "Animal Rescue Network Inc.",
  "type": "rescue",
  "status": "active",
  "contact_info": {
    "email": "contact@animalrescue.org",
    "phone": "+1555123456",
    "fax": "",
    "mobile": "",
    "website": "https://animalrescue.org"
  },
  "address": {
    "street": "789 Rescue Lane",
    "city": "Springfield",
    "state": "IL",
    "zip_code": "62701",
    "country": "USA"
  },
  "primary_contact": {
    "name": "Sarah Wilson",
    "title": "Operations Manager",
    "email": "sarah@animalrescue.org",
    "phone": "+1555123456",
    "mobile": ""
  },
  "partner_since": "2025-11-08T20:00:00Z",
  "agreement_number": "",
  "agreement_date": null,
  "agreement_expiry": null,
  "services_provided": ["rescue", "foster", "medical_care"],
  "specializations": ["dogs", "cats"],
  "max_capacity": 50,
  "current_capacity": 15,
  "accepts_intakes": true,
  "discount_percentage": 0.0,
  "standard_rate": 0.0,
  "rating": 0.0,
  "total_ratings": 0,
  "total_transfers_in": 0,
  "total_transfers_out": 0,
  "successful_placements": 0,
  "documents": [],
  "notes": "Reliable partner with excellent track record",
  "website": "",
  "social_media": {},
  "tags": [],
  "created_by": "507f1f77bcf86cd799439011",
  "created_at": "2025-11-08T20:00:00Z",
  "updated_at": "2025-11-08T20:00:00Z"
}
```

**Valid Partner Types:**
- `rescue` - Rescue organization
- `shelter` - Animal shelter
- `veterinary` - Veterinary clinic
- `foster_network` - Foster network
- `transport` - Transport service
- `sanctuary` - Animal sanctuary
- `government` - Government agency
- `corporate` - Corporate partner
- `other` - Other type

**Key Differences from Original Documentation:**
1. ✅ Added `legal_name` field
2. ✅ Nested `contact_info` object with email, phone, fax, mobile, website
3. ✅ Nested `address` object (postal_code → zip_code)
4. ✅ Nested `primary_contact` object with name, title, email, phone, mobile
5. ✅ `partnership_start_date` → `partner_since` (renamed)
6. ✅ `agreement_expiry_date` → `agreement_expiry` (renamed)
7. ✅ `capacity` → `max_capacity` (renamed), added `current_capacity`
8. ✅ Added `accepts_intakes`, `specializations`
9. ✅ `ratings` (array) → `rating` (float) + `total_ratings` (count)
10. ✅ `average_rating` → calculated from `rating`
11. ✅ Added statistics: `total_transfers_in`, `total_transfers_out`, `successful_placements`
12. ❌ Removed: `contact_person` (replaced by `primary_contact` object), `email`, `phone` (moved to nested objects), `current_usage`, `ratings` array

### Partner Endpoints

#### GET /api/v1/partners
**Description**: List partners
**Authentication**: Required
**Permissions**: `PermissionViewPartners`

**Response: 200 OK**

---

#### GET /api/v1/partners/:id
**Description**: Get partner by ID
**Authentication**: Required
**Permissions**: `PermissionViewPartners`

**Response: 200 OK**

---

#### POST /api/v1/partners
**Description**: Create partner
**Authentication**: Required
**Permissions**: `PermissionCreatePartners`

**Request Body:** (See Partner Structure)

**Response: 201 Created**

---

#### PUT /api/v1/partners/:id
**Description**: Update partner
**Authentication**: Required
**Permissions**: `PermissionUpdatePartners`

**Response: 200 OK**

---

#### DELETE /api/v1/partners/:id
**Description**: Delete partner
**Authentication**: Required
**Permissions**: `PermissionDeletePartners`

**Response: 200 OK**

---

#### GET /api/v1/partners/active
**Description**: Get active partners
**Authentication**: Required
**Permissions**: `PermissionViewPartners`

**Response: 200 OK**

---

#### GET /api/v1/partners/with-capacity
**Description**: Get partners with available capacity
**Authentication**: Required
**Permissions**: `PermissionViewPartners`

**Response: 200 OK**

---

#### GET /api/v1/partners/type/:type
**Description**: Get partners by type
**Authentication**: Required
**Permissions**: `PermissionViewPartners`

**Response: 200 OK**

---

#### GET /api/v1/partners/status/:status
**Description**: Get partners by status
**Authentication**: Required
**Permissions**: `PermissionViewPartners`

**Response: 200 OK**

---

#### GET /api/v1/partners/statistics
**Description**: Get partner statistics
**Authentication**: Required
**Permissions**: `PermissionViewPartners`

**Response: 200 OK**

---

#### POST /api/v1/partners/:id/activate
**Description**: Activate partner
**Authentication**: Required
**Permissions**: `PermissionUpdatePartners`

**Response: 200 OK**

---

#### POST /api/v1/partners/:id/suspend
**Description**: Suspend partner
**Authentication**: Required
**Permissions**: `PermissionUpdatePartners`

**Response: 200 OK**

---

#### POST /api/v1/partners/:id/deactivate
**Description**: Deactivate partner
**Authentication**: Required
**Permissions**: `PermissionUpdatePartners`

**Response: 200 OK**

---

#### POST /api/v1/partners/:id/add-rating
**Description**: Add rating
**Authentication**: Required
**Permissions**: `PermissionUpdatePartners`

**Request Body:**
```json
{
  "rating": 5,
  "comment": "Excellent service"
}
```

**Response: 200 OK**

---

#### POST /api/v1/partners/:id/update-capacity
**Description**: Update capacity
**Authentication**: Required
**Permissions**: `PermissionUpdatePartners`

**Request Body:**
```json
{
  "capacity": 100
}
```

**Response: 200 OK**

---

#### POST /api/v1/partners/:id/set-agreement-expiry
**Description**: Set agreement expiry
**Authentication**: Required
**Permissions**: `PermissionUpdatePartners`

**Request Body:**
```json
{
  "expiry_date": "2027-01-01T00:00:00Z"
}
```

**Response: 200 OK**

---

## Transfer Management

### Transfer Structure

```json
{
  "id": "507f1f77bcf86cd799439026",
  "animal_id": "507f1f77bcf86cd799439013",
  "from_location": "Main Shelter",
  "to_location": "Partner Facility",
  "partner_id": "507f1f77bcf86cd799439025",
  "transfer_type": "medical",
  "status": "pending",
  "reason": "Specialized medical treatment needed",
  "scheduled_date": "2025-11-10T10:00:00Z",
  "actual_transfer_date": null,
  "estimated_return_date": "2025-11-20T00:00:00Z",
  "actual_return_date": null,
  "transport_method": "vehicle",
  "transport_provider": "ABC Transport",
  "transport_cost": 50.00,
  "responsible_person": "507f1f77bcf86cd799439011",
  "approved_by": null,
  "approval_date": null,
  "documents": ["507f1f77bcf86cd799439024"],
  "notes": "Requires temperature-controlled transport",
  "follow_up_required": true,
  "follow_up_date": "2025-11-15T00:00:00Z",
  "created_by": "507f1f77bcf86cd799439011",
  "created_at": "2025-11-08T10:00:00Z",
  "updated_at": "2025-11-08T10:00:00Z"
}
```

### Transfer Endpoints

#### GET /api/v1/transfers
**Description**: List transfers
**Authentication**: Required
**Permissions**: `PermissionViewTransfers`

**Response: 200 OK**

---

#### GET /api/v1/transfers/:id
**Description**: Get transfer by ID
**Authentication**: Required
**Permissions**: `PermissionViewTransfers`

**Response: 200 OK**

---

#### POST /api/v1/transfers
**Description**: Create transfer
**Authentication**: Required
**Permissions**: `PermissionCreateTransfers`

**Request Body:** (See Transfer Structure)

**Response: 201 Created**

---

#### PUT /api/v1/transfers/:id
**Description**: Update transfer
**Authentication**: Required
**Permissions**: `PermissionUpdateTransfers`

**Response: 200 OK**

---

#### DELETE /api/v1/transfers/:id
**Description**: Delete transfer
**Authentication**: Required
**Permissions**: `PermissionDeleteTransfers`

**Response: 200 OK**

---

#### GET /api/v1/transfers/pending
**Description**: Get pending transfers
**Authentication**: Required
**Permissions**: `PermissionViewTransfers`

**Response: 200 OK**

---

#### GET /api/v1/transfers/upcoming
**Description**: Get upcoming transfers
**Authentication**: Required
**Permissions**: `PermissionViewTransfers`

**Response: 200 OK**

---

#### GET /api/v1/transfers/overdue
**Description**: Get overdue transfers
**Authentication**: Required
**Permissions**: `PermissionViewTransfers`

**Response: 200 OK**

---

#### GET /api/v1/transfers/follow-up
**Description**: Get transfers requiring follow-up
**Authentication**: Required
**Permissions**: `PermissionViewTransfers`

**Response: 200 OK**

---

#### GET /api/v1/transfers/animal/:animal_id
**Description**: Get transfers by animal
**Authentication**: Required
**Permissions**: `PermissionViewTransfers`

**Response: 200 OK**

---

#### GET /api/v1/transfers/partner/:partner_id
**Description**: Get transfers by partner
**Authentication**: Required
**Permissions**: `PermissionViewTransfers`

**Response: 200 OK**

---

#### GET /api/v1/transfers/status/:status
**Description**: Get transfers by status
**Authentication**: Required
**Permissions**: `PermissionViewTransfers`

**Response: 200 OK**

---

#### GET /api/v1/transfers/statistics
**Description**: Get transfer statistics
**Authentication**: Required
**Permissions**: `PermissionViewTransfers`

**Response: 200 OK**

---

#### POST /api/v1/transfers/:id/approve
**Description**: Approve transfer
**Authentication**: Required
**Permissions**: `PermissionUpdateTransfers`

**Response: 200 OK**

---

#### POST /api/v1/transfers/:id/reject
**Description**: Reject transfer
**Authentication**: Required
**Permissions**: `PermissionUpdateTransfers`

**Request Body:**
```json
{
  "reason": "Not approved by management"
}
```

**Response: 200 OK**

---

#### POST /api/v1/transfers/:id/start-transit
**Description**: Start transfer transit
**Authentication**: Required
**Permissions**: `PermissionUpdateTransfers`

**Response: 200 OK**

---

#### POST /api/v1/transfers/:id/complete
**Description**: Complete transfer
**Authentication**: Required
**Permissions**: `PermissionUpdateTransfers`

**Request Body:**
```json
{
  "actual_transfer_date": "2025-11-10T10:00:00Z"
}
```

**Response: 200 OK**

---

#### POST /api/v1/transfers/:id/cancel
**Description**: Cancel transfer
**Authentication**: Required
**Permissions**: `PermissionUpdateTransfers`

**Request Body:**
```json
{
  "reason": "No longer needed"
}
```

**Response: 200 OK**

---

#### POST /api/v1/transfers/:id/schedule
**Description**: Schedule transfer
**Authentication**: Required
**Permissions**: `PermissionUpdateTransfers`

**Request Body:**
```json
{
  "scheduled_date": "2025-11-10T10:00:00Z"
}
```

**Response: 200 OK**

---

## Inventory Management

### Inventory Item Structure

```json
{
  "id": "507f1f77bcf86cd799439027",
  "name": "Dog Food - Premium Adult",
  "sku": "DF-PREM-001",
  "category": "food",
  "description": "Premium adult dog food, 20kg bags",
  "unit_of_measure": "bag",
  "quantity_in_stock": 50,
  "reorder_level": 20,
  "reorder_quantity": 30,
  "unit_cost": 45.00,
  "supplier": "Pet Supply Co",
  "supplier_contact": "orders@petsupply.com",
  "location": "Storage Room A, Shelf 3",
  "expiration_date": "2026-06-01T00:00:00Z",
  "last_restock_date": "2025-10-01T00:00:00Z",
  "last_usage_date": "2025-11-08T00:00:00Z",
  "is_active": true,
  "tags": ["food", "dog", "premium"],
  "notes": "Keep in cool, dry place",
  "created_by": "507f1f77bcf86cd799439011",
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-11-08T15:30:00Z"
}
```

### Inventory Endpoints

#### GET /api/v1/inventory
**Description**: List inventory items
**Authentication**: Required
**Permissions**: `PermissionViewInventory`

**Query Parameters:**
- `limit`, `offset`: Pagination
- `category` (string): Filter by category
- `sku` (string): Search by SKU

**Response: 200 OK**

---

#### GET /api/v1/inventory/:id
**Description**: Get inventory item by ID
**Authentication**: Required
**Permissions**: `PermissionViewInventory`

**Response: 200 OK**

---

#### POST /api/v1/inventory
**Description**: Create inventory item
**Authentication**: Required
**Permissions**: `PermissionCreateInventory`

**Request Body:** (See Inventory Item Structure)

**Response: 201 Created**

---

#### PUT /api/v1/inventory/:id
**Description**: Update inventory item
**Authentication**: Required
**Permissions**: `PermissionUpdateInventory`

**Response: 200 OK**

---

#### DELETE /api/v1/inventory/:id
**Description**: Delete inventory item
**Authentication**: Required
**Permissions**: `PermissionDeleteInventory`

**Response: 200 OK**

---

#### GET /api/v1/inventory/low-stock
**Description**: Get low stock items
**Authentication**: Required
**Permissions**: `PermissionViewInventory`

**Response: 200 OK**

---

#### GET /api/v1/inventory/expired
**Description**: Get expired items
**Authentication**: Required
**Permissions**: `PermissionViewInventory`

**Response: 200 OK**

---

#### GET /api/v1/inventory/expiring-soon
**Description**: Get items expiring soon
**Authentication**: Required
**Permissions**: `PermissionViewInventory`

**Query Parameters:**
- `days` (int): Days ahead (default: 30)

**Response: 200 OK**

---

#### GET /api/v1/inventory/needing-reorder
**Description**: Get items needing reorder
**Authentication**: Required
**Permissions**: `PermissionViewInventory`

**Response: 200 OK**

---

#### GET /api/v1/inventory/category/:category
**Description**: Get items by category
**Authentication**: Required
**Permissions**: `PermissionViewInventory`

**Response: 200 OK**

---

#### GET /api/v1/inventory/statistics
**Description**: Get inventory statistics
**Authentication**: Required
**Permissions**: `PermissionViewInventory`

**Response: 200 OK**

---

#### GET /api/v1/inventory/sku/:sku
**Description**: Get item by SKU
**Authentication**: Required
**Permissions**: `PermissionViewInventory`

**Response: 200 OK**

---

#### POST /api/v1/inventory/:id/add-stock
**Description**: Add stock **✅ VERIFIED**
**Authentication**: Required
**Permissions**: `PermissionUpdateInventory`

**Request Body:**
```json
{
  "quantity": 20,
  "unit_cost": 25.00,
  "reference": "PO-2025-001",
  "notes": "Monthly restock from supplier"
}
```

**Required Fields:**
1. `quantity` (float, required, > 0) - Quantity to add
2. `unit_cost` (float, required, >= 0) - Cost per unit **✅ REQUIRED**
3. `reference` (string, optional) - Purchase order or reference number
4. `notes` (string, optional) - Additional notes

**Response: 200 OK**

---

#### POST /api/v1/inventory/:id/remove-stock
**Description**: Remove stock
**Authentication**: Required
**Permissions**: `PermissionUpdateInventory`

**Request Body:**
```json
{
  "quantity": 5,
  "reason": "Used for animals",
  "notes": "Fed to shelter dogs"
}
```

**Response: 200 OK**

---

#### POST /api/v1/inventory/:id/adjust-stock
**Description**: Adjust stock (inventory correction)
**Authentication**: Required
**Permissions**: `PermissionUpdateInventory`

**Request Body:**
```json
{
  "new_quantity": 48,
  "reason": "Physical count correction",
  "notes": "Annual inventory audit"
}
```

**Response: 200 OK**

---

#### POST /api/v1/inventory/:id/activate
**Description**: Activate item
**Authentication**: Required
**Permissions**: `PermissionUpdateInventory`

**Response: 200 OK**

---

#### POST /api/v1/inventory/:id/deactivate
**Description**: Deactivate item
**Authentication**: Required
**Permissions**: `PermissionUpdateInventory`

**Response: 200 OK**

---

## Stock Transactions

### Stock Transaction Structure

```json
{
  "id": "507f1f77bcf86cd799439028",
  "inventory_item_id": "507f1f77bcf86cd799439027",
  "transaction_type": "in",
  "quantity": 30,
  "unit_cost": 45.00,
  "total_cost": 1350.00,
  "previous_quantity": 20,
  "new_quantity": 50,
  "reason": "Restock",
  "reference_type": "",
  "reference_id": "",
  "supplier": "Pet Supply Co",
  "notes": "Monthly restock order #12345",
  "performed_by": "507f1f77bcf86cd799439011",
  "created_at": "2025-11-08T10:00:00Z"
}
```

### Stock Transaction Endpoints

#### GET /api/v1/stock-transactions
**Description**: List stock transactions
**Authentication**: Required
**Permissions**: `PermissionViewStockTransactions`

**Query Parameters:**
- `limit`, `offset`: Pagination
- `transaction_type` (string): in, out, adjustment
- `start_date`, `end_date`: Date range

**Response: 200 OK**

---

#### GET /api/v1/stock-transactions/:id
**Description**: Get transaction by ID
**Authentication**: Required
**Permissions**: `PermissionViewStockTransactions`

**Response: 200 OK**

---

#### GET /api/v1/stock-transactions/item/:item_id
**Description**: Get transactions by item
**Authentication**: Required
**Permissions**: `PermissionViewStockTransactions`

**Response: 200 OK**

---

#### GET /api/v1/stock-transactions/type/:type
**Description**: Get transactions by type
**Authentication**: Required
**Permissions**: `PermissionViewStockTransactions`

**Response: 200 OK**

---

#### GET /api/v1/stock-transactions/statistics
**Description**: Get transaction statistics
**Authentication**: Required
**Permissions**: `PermissionViewStockTransactions`

**Response: 200 OK**

---

## Audit Logs

### Audit Log Structure

```json
{
  "id": "507f1f77bcf86cd799439029",
  "user_id": "507f1f77bcf86cd799439011",
  "action": "create",
  "entity_type": "animal",
  "entity_id": "507f1f77bcf86cd799439013",
  "changes": {
    "name": {"old": null, "new": "Buddy"},
    "status": {"old": null, "new": "available"}
  },
  "ip_address": "192.168.1.1",
  "user_agent": "Mozilla/5.0...",
  "timestamp": "2025-11-08T10:00:00Z"
}
```

### Audit Log Endpoints

#### GET /api/v1/audit-logs
**Description**: List audit logs
**Authentication**: Required
**Permissions**: Admin

**Query Parameters:**
- `limit`, `offset`: Pagination
- `user_id` (string): Filter by user
- `action` (string): Filter by action
- `entity_type` (string): Filter by entity type
- `start_date`, `end_date`: Date range

**Response: 200 OK**

---

#### GET /api/v1/audit-logs/:id
**Description**: Get audit log by ID
**Authentication**: Required
**Permissions**: Admin

**Response: 200 OK**

---

#### GET /api/v1/audit-logs/recent
**Description**: Get recent activity
**Authentication**: Required
**Permissions**: Admin

**Query Parameters:**
- `limit` (int): Number of logs (default: 50)

**Response: 200 OK**

---

#### GET /api/v1/audit-logs/statistics
**Description**: Get audit statistics
**Authentication**: Required
**Permissions**: Admin

**Response: 200 OK**

---

#### GET /api/v1/audit-logs/user/:user_id
**Description**: Get user activity
**Authentication**: Required
**Permissions**: Admin

**Response: 200 OK**

---

#### GET /api/v1/audit-logs/action/:action
**Description**: Get logs by action
**Authentication**: Required
**Permissions**: Admin

**Response: 200 OK**

---

#### GET /api/v1/audit-logs/entity/:entity_type/:entity_id
**Description**: Get entity history
**Authentication**: Required
**Permissions**: Admin

**Response: 200 OK**

---

#### GET /api/v1/audit-logs/date-range
**Description**: Get logs for date range
**Authentication**: Required
**Permissions**: Admin

**Query Parameters:**
- `start_date` (string): ISO date
- `end_date` (string): ISO date

**Response: 200 OK**

---

#### DELETE /api/v1/audit-logs/cleanup
**Description**: Delete old logs
**Authentication**: Required
**Permissions**: Super Admin

**Query Parameters:**
- `days_old` (int): Delete logs older than X days

**Response: 200 OK**

---

## System Monitoring

### Monitoring Endpoints

#### GET /api/v1/monitoring/health
**Description**: Get system health status
**Authentication**: Required
**Permissions**: Admin

**Response: 200 OK**
```json
{
  "status": "healthy",
  "database": {
    "status": "connected",
    "response_time_ms": 5
  },
  "memory": {
    "used_mb": 512,
    "total_mb": 2048,
    "percentage": 25.0
  },
  "uptime_seconds": 86400,
  "timestamp": "2025-11-08T15:30:00Z"
}
```

---

#### GET /api/v1/monitoring/statistics
**Description**: Get usage statistics
**Authentication**: Required
**Permissions**: Admin

**Response: 200 OK**

---

#### GET /api/v1/monitoring/performance
**Description**: Get performance metrics
**Authentication**: Required
**Permissions**: Admin

**Response: 200 OK**

---

#### GET /api/v1/monitoring/database
**Description**: Get database statistics
**Authentication**: Required
**Permissions**: Admin

**Response: 200 OK**

---

#### GET /api/v1/monitoring/configuration
**Description**: Get system configuration
**Authentication**: Required
**Permissions**: Admin

**Response: 200 OK**

---

## Medical Records & Prescriptions

### Medical Condition Structure

```json
{
  "id": "507f1f77bcf86cd79943902a",
  "animal_id": "507f1f77bcf86cd799439013",
  "condition_name": "Hip Dysplasia",
  "diagnosis_date": "2025-10-15T00:00:00Z",
  "diagnosed_by": "507f1f77bcf86cd799439011",
  "severity": "moderate",
  "status": "active",
  "is_chronic": true,
  "symptoms": ["limping", "difficulty rising"],
  "treatment_plan_id": "507f1f77bcf86cd79943902b",
  "resolution_date": null,
  "resolution_notes": "",
  "notes": "Requires ongoing management",
  "created_by": "507f1f77bcf86cd799439011",
  "created_at": "2025-10-15T00:00:00Z",
  "updated_at": "2025-11-08T15:30:00Z"
}
```

### Medication Structure

```json
{
  "id": "507f1f77bcf86cd79943902c",
  "animal_id": "507f1f77bcf86cd799439013",
  "medication_name": "Carprofen",
  "dosage": "50mg",
  "frequency": "Twice daily",
  "route": "oral",
  "start_date": "2025-10-15T00:00:00Z",
  "end_date": "2025-11-15T00:00:00Z",
  "status": "active",
  "prescribed_by": "507f1f77bcf86cd799439011",
  "refills_remaining": 2,
  "next_refill_due": "2025-11-10T00:00:00Z",
  "side_effects": "Monitor for GI upset",
  "special_instructions": "Give with food",
  "administration_logs": [
    {
      "administered_at": "2025-11-08T08:00:00Z",
      "administered_by": "507f1f77bcf86cd799439011",
      "dosage_given": "50mg",
      "notes": "Given with breakfast"
    }
  ],
  "notes": "Pain management for hip dysplasia",
  "created_by": "507f1f77bcf86cd799439011",
  "created_at": "2025-10-15T00:00:00Z",
  "updated_at": "2025-11-08T15:30:00Z"
}
```

### Treatment Plan Structure

**✅ VERIFIED - Tested and working**

```json
{
  "id": "507f1f77bcf86cd79943902b",
  "animal_id": "507f1f77bcf86cd799439013",
  "condition_id": "507f1f77bcf86cd79943902a",
  "plan_name": "Dermatitis Treatment Plan",
  "description": "7-day comprehensive treatment for allergic dermatitis",
  "start_date": "2025-11-08T20:00:00Z",
  "end_date": null,
  "created_by": "507f1f77bcf86cd799439011",
  "status": "active",
  "goals": [
    "Reduce inflammation",
    "Eliminate itching",
    "Prevent recurrence"
  ],
  "medications": [],
  "procedures": [
    {
      "procedure_name": "Daily medication application",
      "description": "Apply topical medication twice daily",
      "scheduled_date": "2025-11-09T20:00:00Z",
      "completed_date": null,
      "performed_by": null,
      "status": "scheduled",
      "cost": 0.0,
      "notes": "Monitor for adverse reactions"
    }
  ],
  "dietary_plan": "Continue regular diet, monitor for food allergies",
  "exercise_plan": "Normal activity, avoid excessive running",
  "monitoring_plan": "Daily skin checks, weekly vet follow-up",
  "follow_up_schedule": [],
  "notes": "Initial treatment plan - may adjust based on response",
  "progress": [],
  "created_at": "2025-11-08T20:00:00Z",
  "updated_at": "2025-11-08T20:00:00Z"
}
```

**Key Differences from Original Documentation:**
1. ✅ Added REQUIRED field: `plan_name` (not "name")
2. ✅ Added `description` field
3. ✅ Added `created_by` field
4. ✅ Updated `procedures` structure to use PlannedProcedure:
   - `name` → `procedure_name` (field renamed)
   - Added: `description`, `scheduled_date`, `completed_date`, `performed_by`, `status`, `cost`, `notes`
   - Removed: `frequency`, `next_scheduled`
5. ✅ Added: `dietary_plan`, `exercise_plan`, `monitoring_plan`
6. ✅ `follow_up_schedule` simplified (uses TreatmentFollowUpSchedule structure)
7. ✅ `progress` uses ProgressNote structure with `date`, `recorded_by`, `note`, `improvement`

### Medical Condition Endpoints

#### GET /api/v1/medical-conditions
**Description**: List medical conditions
**Authentication**: Required
**Permissions**: `PermissionViewVeterinary`

**Query Parameters:**
- `limit`, `offset`: Pagination
- `animal_id` (string): Filter by animal
- `status` (string): Filter by status
- `severity` (string): Filter by severity

**Response: 200 OK**

---

#### GET /api/v1/medical-conditions/:id
**Description**: Get medical condition by ID
**Authentication**: Required
**Permissions**: `PermissionViewVeterinary`

**Response: 200 OK**

---

#### POST /api/v1/medical-conditions
**Description**: Create medical condition
**Authentication**: Required
**Permissions**: `PermissionCreateVeterinary`

**Request Body:** (See Medical Condition Structure)

**Response: 201 Created**

---

#### PUT /api/v1/medical-conditions/:id
**Description**: Update medical condition
**Authentication**: Required
**Permissions**: `PermissionUpdateVeterinary`

**Response: 200 OK**

---

#### DELETE /api/v1/medical-conditions/:id
**Description**: Delete medical condition
**Authentication**: Required
**Permissions**: `PermissionDeleteVeterinary`

**Response: 200 OK**

---

#### GET /api/v1/medical-conditions/chronic
**Description**: Get chronic conditions
**Authentication**: Required
**Permissions**: `PermissionViewVeterinary`

**Response: 200 OK**

---

#### POST /api/v1/medical-conditions/:id/resolve
**Description**: Resolve medical condition
**Authentication**: Required
**Permissions**: `PermissionUpdateVeterinary`

**Request Body:**
```json
{
  "resolution_date": "2025-11-08T00:00:00Z",
  "resolution_notes": "Fully recovered"
}
```

**Response: 200 OK**

---

### Medication Endpoints

#### GET /api/v1/medications
**Description**: List medications
**Authentication**: Required
**Permissions**: `PermissionViewVeterinary`

**Query Parameters:**
- `limit`, `offset`: Pagination
- `animal_id` (string): Filter by animal
- `status` (string): Filter by status

**Response: 200 OK**

---

#### GET /api/v1/medications/:id
**Description**: Get medication by ID
**Authentication**: Required
**Permissions**: `PermissionViewVeterinary`

**Response: 200 OK**

---

#### POST /api/v1/medications
**Description**: Create medication
**Authentication**: Required
**Permissions**: `PermissionCreateVeterinary`

**Request Body:** (See Medication Structure)

**Response: 201 Created**

---

#### PUT /api/v1/medications/:id
**Description**: Update medication
**Authentication**: Required
**Permissions**: `PermissionUpdateVeterinary`

**Response: 200 OK**

---

#### DELETE /api/v1/medications/:id
**Description**: Delete medication
**Authentication**: Required
**Permissions**: `PermissionDeleteVeterinary`

**Response: 200 OK**

---

#### GET /api/v1/medications/due-for-refill
**Description**: Get medications due for refill
**Authentication**: Required
**Permissions**: `PermissionViewVeterinary`

**Response: 200 OK**

---

#### GET /api/v1/medications/expiring-soon
**Description**: Get medications expiring soon
**Authentication**: Required
**Permissions**: `PermissionViewVeterinary`

**Query Parameters:**
- `days` (int): Days ahead (default: 30)

**Response: 200 OK**

---

#### POST /api/v1/medications/:id/administer
**Description**: Record medication administration
**Authentication**: Required
**Permissions**: `PermissionUpdateVeterinary`

**Request Body:**
```json
{
  "dosage_given": "50mg",
  "notes": "Given with food"
}
```

**Response: 200 OK**

---

#### POST /api/v1/medications/:id/refill
**Description**: Refill medication
**Authentication**: Required
**Permissions**: `PermissionUpdateVeterinary`

**Request Body:**
```json
{
  "refills_added": 1,
  "new_end_date": "2025-12-15T00:00:00Z"
}
```

**Response: 200 OK**

---

### Treatment Plan Endpoints

#### GET /api/v1/treatment-plans
**Description**: List treatment plans
**Authentication**: Required
**Permissions**: `PermissionViewVeterinary`

**Query Parameters:**
- `limit`, `offset`: Pagination
- `animal_id` (string): Filter by animal
- `status` (string): Filter by status

**Response: 200 OK**

---

#### GET /api/v1/treatment-plans/:id
**Description**: Get treatment plan by ID
**Authentication**: Required
**Permissions**: `PermissionViewVeterinary`

**Response: 200 OK**

---

#### POST /api/v1/treatment-plans
**Description**: Create treatment plan
**Authentication**: Required
**Permissions**: `PermissionCreateVeterinary`

**Request Body:** (See Treatment Plan Structure)

**Response: 201 Created**

---

#### PUT /api/v1/treatment-plans/:id
**Description**: Update treatment plan
**Authentication**: Required
**Permissions**: `PermissionUpdateVeterinary`

**Response: 200 OK**

---

#### DELETE /api/v1/treatment-plans/:id
**Description**: Delete treatment plan
**Authentication**: Required
**Permissions**: `PermissionDeleteVeterinary`

**Response: 200 OK**

---

#### POST /api/v1/treatment-plans/:id/progress
**Description**: Add progress note
**Authentication**: Required
**Permissions**: `PermissionUpdateVeterinary`

**Request Body:**
```json
{
  "note": "Showing significant improvement",
  "improvement": "good"
}
```

**Response: 200 OK**

---

#### POST /api/v1/treatment-plans/:id/activate
**Description**: Activate treatment plan
**Authentication**: Required
**Permissions**: `PermissionUpdateVeterinary`

**Response: 200 OK**

---

#### POST /api/v1/treatment-plans/:id/complete
**Description**: Complete treatment plan
**Authentication**: Required
**Permissions**: `PermissionUpdateVeterinary`

**Request Body:**
```json
{
  "completion_notes": "Goals achieved, condition resolved"
}
```

**Response: 200 OK**

---

## Health Check

#### GET /api/v1/ping
**Description**: Health check endpoint
**Authentication**: None
**Permissions**: None

**Response: 200 OK**
```json
{
  "message": "pong"
}
```

---

## Appendix

### Common Enums

#### User Roles
- `super_admin`
- `admin`
- `employee`
- `volunteer`
- `user`

#### User Status
- `active`
- `inactive`
- `suspended`

#### Animal Status
- `available`
- `adopted`
- `under_treatment`
- `quarantine`
- `fostered`
- `reserved`
- `deceased`
- `transferred`

#### Animal Category
- `mammal`
- `reptile`
- `bird`
- `amphibian`
- `fish`
- `invertebrate`
- `farm_animal`

#### Animal Sex
- `male`
- `female`
- `unknown`

#### Animal Size
- `small`
- `medium`
- `large`
- `xlarge`

#### Payment Method
- `credit_card`
- `debit_card`
- `bank_transfer`
- `cash`
- `check`
- `paypal`
- `other`

#### Payment Status
- `pending`
- `completed`
- `failed`
- `refunded`
- `cancelled`

#### Campaign Status
- `draft`
- `active`
- `paused`
- `completed`
- `cancelled`

#### Event Status
- `draft`
- `active`
- `completed`
- `cancelled`

#### Communication Type
- `email`
- `sms`
- `in_app`
- `push`

#### Communication Status
- `draft`
- `pending`
- `sent`
- `delivered`
- `failed`
- `bounced`

#### Notification Type
- `info`
- `success`
- `warning`
- `error`

#### Task Status
- `pending`
- `in_progress`
- `completed`
- `cancelled`

#### Task Priority
- `low`
- `normal`
- `high`
- `urgent`

#### Transfer Status
- `pending`
- `approved`
- `in_transit`
- `completed`
- `cancelled`
- `rejected`

#### Condition Severity
- `mild`
- `moderate`
- `severe`
- `critical`

#### Medication Route
- `oral`
- `topical`
- `injection`
- `intravenous`

---

**End of API Documentation**

*For support or questions, contact: api-support@foundation.org*
