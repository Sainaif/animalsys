# Quick Start: Seed Database

## Run the Seed Script

```bash
# Drop the DB and reseed in one command
make reseed

# Using Docker (Recommended)
docker-compose exec backend ./seed

# The script will:
# 1. Clear ALL existing data
# 2. Generate 3 years of realistic data
# 3. Create 20-40 staff members
# 4. Create 30-60 volunteers
# 5. Generate 200-400 animals (60% adopted, 40% available)
# 6. Create 150-300 donors and 500-1,500 donations
# 7. Generate complete veterinary, adoption, and event records
```

## Login Credentials

**Super Admin (configurable via `ORG_*` in `docker-compose.yml`):**
- Email: `sarah.johnson@happypaws.org`
- Password: `password123`

**Other Admins (generated with your domain):**
- `michael.chen@happypaws.org`
- `emily.rodriguez@happypaws.org`
- Password: `password123` (all users)

## What You'll Get

| Category | Count | Details |
|----------|-------|---------|
| **Users** | 20-40 | 1 Super Admin, 2 Admins, 17-37 Employees |
| **Volunteers** | 30-60 | Active volunteers with logged hours |
| **Animals** | 200-400 | Dogs, Cats, Rabbits, Birds |
| **Adoptions** | ~120-240 | 60% of animals adopted |
| **Donors** | 150-300 | Individual donors |
| **Donations** | 500-1,500 | $10-$5,000 range |
| **Campaigns** | 5 | Fundraisers 2023-2025 |
| **Events** | 10-30 | Distributed over 3 years |
| **Vet Visits** | 200-1,200 | 1-3 per animal |
| **Vaccinations** | 160-320 | 80% of animals |
| **Partners** | 4 | Vet clinics, suppliers, foster networks |
| **Tasks** | 100-200 | 70% completed |
| **Documents** | 50-150 | Contracts, medical, legal |

## Foundation Timeline

- **Start Date**: November 9, 2022 (3 years ago)
- **Data Range**: Nov 2022 → Today
- **Operations**: Fully simulated 3-year history

## Quick Verification

```bash
# Check user count
docker-compose exec mongodb mongosh animalsys --eval "db.users.countDocuments()"

# Check animals
docker-compose exec mongodb mongosh animalsys --eval "db.animals.countDocuments()"

# Check donations total
docker-compose exec mongodb mongosh animalsys --eval "db.donations.countDocuments()"
```

## ⚠️ Warning

**THIS DELETES ALL EXISTING DATA!**

Only run on development/demo environments.

---

For detailed information, see [backend/cmd/seed/README.md](backend/cmd/seed/README.md)
