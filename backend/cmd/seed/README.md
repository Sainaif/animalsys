# Database Seed Script

This seed script generates **3 years of realistic data** for an animal foundation (defaults to Happy Paws) and can now be fully customized via environment variables.

## What Gets Seeded

### Staff & Personnel (20-40 workers)
- **1 Super Admin**: sarah.johnson@happypaws.org
- **2-3 Admins**: Joined within first 6 months
- **17-37 Employees**: Various join dates over 3 years

### Volunteers (30-60 people)
- Active volunteers with realistic hours logged (10-500 hours each)
- Joined at different times throughout the 3-year period

### Animals (200-400)
- **~60% Adopted** - Animals that found their forever homes
- **~40% Available** - Currently seeking adoption
- Species: Dogs, Cats, Rabbits, Birds
- Complete intake records with microchip IDs

### Donors (150-300)
- Individual donors with contact information
- Donation history linked to specific donors

### Donations (500-1,500)
- Distributed over 3 years
- Various amounts: small ($10-50), medium ($51-500), large ($501-5,000)
- 40% linked to fundraising campaigns
- Multiple payment methods

### Campaigns (5)
- Annual Fundraisers for 2023, 2024, 2025
- Medical Funds for 2023, 2024
- Increasing goals year over year ($50k → $70k for annual fundraisers)
- Realistic raised amounts

### Adoptions
- Complete adoption records for all adopted animals
- Applications → Approvals → Finalized adoptions
- Realistic timeline (7-60 days from intake to adoption)
- Adoption fees ($50-300)

### Veterinary Records
- 1-3 vet visits per animal
- 80% of animals vaccinated (Rabies, DHPP, Bordetella)
- Visit types: Checkups, Vaccinations, Surgery

### Events (Variable)
- Adoption Fairs
- Volunteer Training sessions  
- Fundraiser Galas
- Community Outreach
- Distributed throughout 3 years

### Partners (4)
- City Vet Clinic
- Pet Supply Co
- Foster Network
- Transport Service

### Inventory (6 items)
- Dog Food
- Cat Food
- Cat Litter
- Toys
- Cleaning Supplies

### Tasks (100-200)
- 70% completed, 30% pending
- Various priorities (high/medium/low)
- Assigned to employees

### Documents (50-150)
- Contracts
- Medical records
- Legal documents
- Uploaded over 3-year period

## How to Run

### Using Docker (Recommended)

```bash
# From the project root
docker-compose exec backend ./seed

# Or rebuild and run
docker-compose down
docker-compose up -d
docker-compose exec backend ./seed
```

### Reset & Seed in one step

```bash
make reseed
```

This command runs `scripts/reseed.sh`, drops the MongoDB database, and executes the seed script from inside the backend container.

### Locally

```bash
# Make sure MongoDB is running and environment variables are set
cd backend
go run ./cmd/seed/main.go
```

## Environment Variables

The script uses these environment variables (or defaults):

```bash
MONGODB_URI=mongodb://mongodb:27017  # Default for Docker
MONGODB_DATABASE=animalsys           # Default database name
# Organization overrides come from docker-compose.yml (backend service environment)
```

## Default Credentials

After seeding, you can log in with:

**Super Admin:**
- Email: `sarah.johnson@happypaws.org`
- Password: `password123`

**All other users** have the same password: `password123`

## Sample User Accounts

- **sarah.johnson@happypaws.org** - Super Admin
- **michael.chen@happypaws.org** - Admin
- **emily.rodriguez@happypaws.org** - Admin
- Various employees: john.smith###@happypaws.org (### = random number)

## What Happens When You Run It

1. ✅ Connects to MongoDB
2. 🗑️ **Clears all existing data** (destructive!)
3. 🌱 Seeds data in order:
   - System settings
   - Users (staff)
   - Species & Animals
   - Donors & Campaigns
   - Donations
   - Adoptions
   - Volunteers
   - Events
   - Veterinary records
   - Partners
   - Inventory
   - Tasks
   - Documents
4. 📊 Prints summary of created records

## Timeline

- **Foundation Start Date**: November 9, 2022 (exactly 3 years ago)
- **Data Range**: November 2022 → Present day
- **Realistic Progression**: Data distributed naturally over time

## Warning

⚠️ **THIS SCRIPT DELETES ALL EXISTING DATA**

The script drops all collections before seeding. Only run this on:
- Development environments
- Fresh installations
- When you want to reset to demo data

## Customization

- **Organization branding**: update the `ORG_*` variables under the `backend` service in `docker-compose.yml` to change the name, contact info, domain, center name, and default password used throughout the seed data.
- **Super admin profile**: override `ORG_SUPER_ADMIN_*` values to control the initial account.
- **Data volumes**: adjust the constants near the top of `main.go` (`minStaff`, `maxStaff`, etc.) if you need larger or smaller datasets.

## Verification

After seeding, you can verify the data:

```bash
# Connect to MongoDB
docker-compose exec mongodb mongosh animalsys

# Count documents
db.users.countDocuments()
db.animals.countDocuments()
db.donations.countDocuments()
```

## Use Cases

This seeded database is perfect for:
- **Development**: Work with realistic data volumes
- **Testing**: Test features with comprehensive data
- **Demos**: Show off the system with believable data
- **Performance Testing**: Test with substantial data sets
- **Training**: Onboard new developers with realistic scenarios

## Next Steps

After seeding:
1. Log in as super admin
2. Explore the dashboard to see statistics
3. Browse animals, adoptions, donations
4. Test workflows with the pre-populated data
5. Create new test data as needed

---

**Happy Paws Animal Foundation** - Seed Script v1.0
Generated: 3 years of operational data (Nov 2022 - Nov 2025)
