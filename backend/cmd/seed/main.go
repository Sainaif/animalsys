package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/pkg/security"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	foundationStartDate = "2022-11-09" // 3 years ago
	minStaff            = 20
	maxStaff            = 40
	minVolunteers       = 30
	maxVolunteers       = 60
)

var (
	db          *mongo.Database
	ctx         = context.Background()
	startDate   time.Time
	currentDate = time.Now()

	userIDs      []primitive.ObjectID
	animalIDs    []primitive.ObjectID
	donorIDs     []primitive.ObjectID
	campaignIDs  []primitive.ObjectID
	volunteerIDs []primitive.ObjectID
	adminIDs     []primitive.ObjectID
	employeeIDs  []primitive.ObjectID
	partnerIDs   []primitive.ObjectID
	speciesIDs   []primitive.ObjectID
)

func main() {
	log.Println("🌱 Starting Animal Foundation Seed Script")
	log.Println("📅 Foundation established:", foundationStartDate)

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	var err error
	startDate, err = time.Parse("2006-01-02", foundationStartDate)
	if err != nil {
		log.Fatal("Failed to parse start date:", err)
	}

	if err := connectDB(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer disconnectDB()

	log.Println("✅ Connected to MongoDB")

	if err := clearDatabase(); err != nil {
		log.Fatal("Failed to clear database:", err)
	}

	rand.Seed(time.Now().UnixNano())

	log.Println("\n📝 Starting data seeding...")

	seedSystemSettings()
	seedUsers()
	seedSpecies()
	seedAnimals()
	seedDonors()
	seedCampaigns()
	seedDonations()
	seedAdoptions()
	seedVolunteers()
	seedEvents()
	seedVeterinaryRecords()
	seedPartners()
	seedInventory()
	seedTasks()
	seedDocuments()

	log.Println("\n✅ Seeding completed!")
	printSummary()
}

func connectDB() error {
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://mongodb:27017"
	}

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	if err = client.Ping(ctx, nil); err != nil {
		return err
	}

	dbName := os.Getenv("MONGODB_DATABASE")
	if dbName == "" {
		dbName = "animalsys"
	}

	db = client.Database(dbName)
	return nil
}

func disconnectDB() {
	if db != nil {
		db.Client().Disconnect(ctx)
	}
}

func clearDatabase() error {
	log.Println("🗑️  Clearing existing data...")

	collections := []string{
		"users", "animals", "species", "adoptions", "adoption_applications",
		"donors", "donations", "campaigns", "events",
		"volunteers", "veterinary_visits", "vaccinations",
		"partners", "inventory_items", "tasks", "documents",
		"communications", "notifications", "audit_logs", "settings",
	}

	for _, coll := range collections {
		db.Collection(coll).Drop(ctx)
	}

	return nil
}

func seedSystemSettings() {
	log.Println("\n⚙️  Seeding system settings...")

	settings := map[string]interface{}{
		"_id": primitive.NewObjectID(),
		"organization": map[string]interface{}{
			"name":        "Happy Paws Animal Foundation",
			"description": "Dedicated to rescuing and rehoming animals in need",
			"address":     "123 Rescue Lane, Animal City, AC 12345",
			"phone":       "+1-555-ANIMALS",
			"email":       "info@happypaws.org",
			"website":     "https://happypaws.org",
		},
		"created_at": startDate,
		"updated_at": currentDate,
	}

	db.Collection("settings").InsertOne(ctx, settings)
	log.Println("✓ Created system settings")
}

func seedUsers() {
	log.Println("\n👥 Seeding users...")

	hasher := security.NewPasswordService()
	passwordHash, _ := hasher.HashPassword("password123")

	// Super admin
	superAdmin := createUser("Sarah", "Johnson", "sarah.johnson@happypaws.org", entities.RoleSuperAdmin, passwordHash, startDate)
	userIDs = append(userIDs, superAdmin)
	adminIDs = append(adminIDs, superAdmin)

	// Admins (2-3)
	adminNames := [][2]string{
		{"Michael", "Chen"},
		{"Emily", "Rodriguez"},
	}

	for _, name := range adminNames {
		admin := createUser(name[0], name[1], fmt.Sprintf("%s.%s@happypaws.org", strings.ToLower(name[0]), strings.ToLower(name[1])), entities.RoleAdmin, passwordHash, randomDate(startDate, startDate.AddDate(0, 6, 0)))
		userIDs = append(userIDs, admin)
		adminIDs = append(adminIDs, admin)
	}

	// Employees (15-35)
	employeeCount := randInt(minStaff-3, maxStaff-3)
	firstNames := []string{"John", "Jane", "Alex", "Maria", "Chris", "Lisa", "Tom", "Anna", "Mark", "Sophie",
		"Ryan", "Emma", "Nick", "Olivia", "Sam", "Grace", "Ben", "Lily", "Jake", "Mia"}
	lastNames := []string{"Smith", "Jones", "Brown", "Davis", "Miller", "Wilson", "Moore", "Taylor", "Anderson", "Thomas"}

	for i := 0; i < employeeCount; i++ {
		firstName := firstNames[rand.Intn(len(firstNames))]
		lastName := lastNames[rand.Intn(len(lastNames))]
		email := fmt.Sprintf("%s.%s%d@happypaws.org", strings.ToLower(firstName), strings.ToLower(lastName), randInt(1, 999))
		employee := createUser(firstName, lastName, email, entities.RoleEmployee, passwordHash, randomDate(startDate, currentDate.AddDate(0, -1, 0)))
		userIDs = append(userIDs, employee)
		employeeIDs = append(employeeIDs, employee)
	}

	log.Printf("✓ Created %d users (1 super admin, %d admins, %d employees)", 1+len(adminNames)+employeeCount, len(adminNames), employeeCount)
}

func createUser(firstName, lastName, email string, role entities.UserRole, passwordHash string, joinDate time.Time) primitive.ObjectID {
	id := primitive.NewObjectID()

	user := map[string]interface{}{
		"_id":           id,
		"email":         email,
		"password_hash": passwordHash,
		"first_name":    firstName,
		"last_name":     lastName,
		"role":          role,
		"status":        entities.StatusActive,
		"phone":         fmt.Sprintf("+1-555-%04d", rand.Intn(10000)),
		"language":      "en",
		"theme":         randomChoice([]string{"light", "dark"}),
		"created_at":    joinDate,
		"updated_at":    currentDate,
	}

	db.Collection("users").InsertOne(ctx, user)
	return id
}

func seedSpecies() {
	log.Println("\n🐾 Seeding species...")

	speciesList := []struct {
		name   string
		breeds []string
	}{
		{"Dog", []string{"Labrador Retriever", "German Shepherd", "Golden Retriever", "Bulldog", "Beagle", "Poodle", "Mixed Breed"}},
		{"Cat", []string{"Persian", "Maine Coon", "Siamese", "Ragdoll", "Bengal", "Mixed Breed", "Domestic Shorthair"}},
		{"Rabbit", []string{"Holland Lop", "Netherland Dwarf", "Flemish Giant", "Rex"}},
		{"Bird", []string{"Parakeet", "Cockatiel", "African Grey", "Macaw"}},
	}

	for _, sp := range speciesList {
		id := primitive.NewObjectID()
		species := map[string]interface{}{
			"_id":         id,
			"name":        sp.name,
			"breeds":      sp.breeds,
			"description": fmt.Sprintf("%s species", sp.name),
			"created_at":  startDate,
			"updated_at":  currentDate,
		}
		db.Collection("species").InsertOne(ctx, species)
		speciesIDs = append(speciesIDs, id)
	}

	log.Printf("✓ Created %d species", len(speciesList))
}

func seedAnimals() {
	log.Println("\n🐕 Seeding animals...")

	type speciesProfile struct {
		name         string
		polishName   string
		category     entities.AnimalCategory
		englishNames []string
		polishNames  []string
		breeds       []string
		colors       []string
		sizes        []entities.AnimalSize
		minWeight    float64
		maxWeight    float64
	}

	speciesProfiles := []speciesProfile{
		{
			name:         "Dog",
			polishName:   "Pies",
			category:     entities.CategoryMammal,
			englishNames: []string{"Buddy", "Max", "Charlie", "Lucy", "Bella", "Cooper", "Rocky", "Sadie"},
			polishNames:  []string{"Borys", "Maks", "Figa", "Luna", "Bella", "Kora", "Rocky"},
			breeds:       []string{"Labrador Retriever", "German Shepherd", "Golden Retriever", "Border Collie", "Beagle", "Mixed Breed"},
			colors:       []string{"Black", "Brown", "Golden", "White", "Brindle"},
			sizes:        []entities.AnimalSize{entities.SizeSmall, entities.SizeMedium, entities.SizeLarge},
			minWeight:    8,
			maxWeight:    45,
		},
		{
			name:         "Cat",
			polishName:   "Kot",
			category:     entities.CategoryMammal,
			englishNames: []string{"Whiskers", "Shadow", "Mittens", "Simba", "Cleo", "Felix", "Luna"},
			polishNames:  []string{"Mruczek", "Ciapek", "Luna", "Kleo", "Feliks"},
			breeds:       []string{"Persian", "Maine Coon", "Siamese", "Ragdoll", "Bengal", "Domestic Shorthair"},
			colors:       []string{"Black", "White", "Orange", "Gray", "Calico", "Tabby"},
			sizes:        []entities.AnimalSize{entities.SizeSmall, entities.SizeMedium},
			minWeight:    3,
			maxWeight:    8,
		},
		{
			name:         "Rabbit",
			polishName:   "Królik",
			category:     entities.CategoryMammal,
			englishNames: []string{"Fluffy", "Cotton", "Nibbles", "Thumper", "Snowball"},
			polishNames:  []string{"Puszek", "Bąbel", "Śnieżek", "Uszek"},
			breeds:       []string{"Holland Lop", "Netherland Dwarf", "Flemish Giant", "Rex"},
			colors:       []string{"White", "Gray", "Brown", "Spotted"},
			sizes:        []entities.AnimalSize{entities.SizeSmall},
			minWeight:    1,
			maxWeight:    6,
		},
		{
			name:         "Bird",
			polishName:   "Ptak",
			category:     entities.CategoryBird,
			englishNames: []string{"Sunny", "Skye", "Kiwi", "Rio", "Echo"},
			polishNames:  []string{"Słoneczko", "Niebieska", "Kiwi", "Rio"},
			breeds:       []string{"Parakeet", "Cockatiel", "African Grey", "Macaw"},
			colors:       []string{"Green", "Blue", "Yellow", "Grey"},
			sizes:        []entities.AnimalSize{entities.SizeSmall},
			minWeight:    0.1,
			maxWeight:    1.5,
		},
	}

	colorTranslations := map[string]string{
		"black":   "czarny",
		"brown":   "brązowy",
		"golden":  "złoty",
		"white":   "biały",
		"brindle": "pręgowany",
		"orange":  "pomarańczowy",
		"gray":    "szary",
		"grey":    "szary",
		"calico":  "szylkretowy",
		"tabby":   "pręgowany",
		"spotted": "cętkowany",
		"green":   "zielony",
		"blue":    "niebieski",
		"yellow":  "żółty",
	}

	totalAnimals := randInt(200, 400)
	adoptedCount := 0
	availableCount := 0

	for i := 0; i < totalAnimals; i++ {
		profile := speciesProfiles[rand.Intn(len(speciesProfiles))]

		englishName := fmt.Sprintf("%s %d", profile.englishNames[rand.Intn(len(profile.englishNames))], randInt(1, 999))
		polishName := fmt.Sprintf("%s %d", profile.polishNames[rand.Intn(len(profile.polishNames))], randInt(1, 999))

		birthDate := randomDate(startDate.AddDate(-8, 0, 0), currentDate.AddDate(0, -1, 0))
		intakeDate := randomDate(birthDate, currentDate)

		sexes := []entities.AnimalSex{entities.SexMale, entities.SexFemale}
		status := entities.AnimalStatusAvailable
		if rand.Float64() < 0.45 {
			status = entities.AnimalStatusAdopted
			adoptedCount++
		} else {
			availableCount++
		}

		caretakerID := employeeIDs[rand.Intn(len(employeeIDs))]
		assignedCaretaker := caretakerID
		adoptionInfo := entities.AdoptionInfo{
			AdoptionFee: float64(randInt(50, 300)),
			Requirements: []string{
				"Completed adoption form",
				"Meet & greet",
			},
		}
		if status == entities.AnimalStatusAdopted {
			adoptionDate := randomDate(intakeDate, currentDate)
			adoptionInfo.AdoptionDate = &adoptionDate
		}

		color := randomChoice(profile.colors)
		colorPL := strings.ToLower(color)
		if translated, ok := colorTranslations[strings.ToLower(color)]; ok {
			colorPL = translated
		}

		descriptionEN := fmt.Sprintf("%s is a %s %s who enjoys playtime and cozy naps.", englishName, strings.ToLower(color), strings.ToLower(profile.name))
		descriptionPL := fmt.Sprintf("%s to %s %s, który lubi zabawę i drzemki.", polishName, colorPL, strings.ToLower(profile.polishName))

		animal := entities.Animal{
			ID:           primitive.NewObjectID(),
			Name:         entities.MultilingualName{English: englishName, Polish: polishName},
			Category:     profile.category,
			Species:      profile.name,
			Breed:        randomChoice(profile.breeds),
			Sex:          sexes[rand.Intn(len(sexes))],
			Status:       status,
			DateOfBirth:  &birthDate,
			AgeEstimated: rand.Intn(2) == 0,
			Color:        color,
			Size:         profile.sizes[rand.Intn(len(profile.sizes))],
			Weight:       randomFloat(profile.minWeight, profile.maxWeight),
			Description: entities.MultilingualName{
				English: descriptionEN,
				Polish:  descriptionPL,
			},
			Images: entities.AnimalImages{
				Primary: fmt.Sprintf("https://placehold.co/600x400?text=%s", strings.ReplaceAll(profile.name, " ", "+")),
			},
			Medical: entities.MedicalInfo{
				Vaccinated:      true,
				Sterilized:      rand.Intn(2) == 0,
				Microchipped:    true,
				MicrochipNumber: fmt.Sprintf("MC%012d", rand.Intn(1000000000000)),
				HealthStatus:    randomChoice([]string{"healthy", "recovering"}),
			},
			Behavior: entities.BehaviorInfo{
				Temperament: []entities.Temperament{
					entities.TemperamentFriendly,
					entities.TemperamentPlayful,
				},
				GoodWithKids: rand.Intn(2) == 0,
				GoodWithDogs: rand.Intn(2) == 0,
				GoodWithCats: rand.Intn(2) == 0,
				HouseTrained: rand.Intn(2) == 0,
			},
			Shelter: entities.ShelterInfo{
				IntakeDate:        intakeDate,
				IntakeReason:      randomChoice([]string{"rescue", "surrender", "transfer"}),
				Location:          fmt.Sprintf("%s %d", randomChoice([]string{"Kennel", "Suite", "Room"}), randInt(1, 40)),
				AssignedCaretaker: &assignedCaretaker,
			},
			Adoption:  adoptionInfo,
			CreatedBy: caretakerID,
			UpdatedBy: caretakerID,
			CreatedAt: intakeDate,
			UpdatedAt: currentDate,
		}

		if _, err := db.Collection("animals").InsertOne(ctx, animal); err == nil {
			animalIDs = append(animalIDs, animal.ID)
		}
	}

	log.Printf("✓ Created %d animals (%d adopted, %d available)", totalAnimals, adoptedCount, availableCount)
}

func seedDonors() {
	log.Println("\n💰 Seeding donors...")

	donorCount := randInt(150, 300)

	for i := 0; i < donorCount; i++ {
		firstName := randomChoice([]string{"John", "Jane", "Michael", "Sarah", "David"})
		lastName := randomChoice([]string{"Smith", "Johnson", "Williams", "Brown", "Jones"})

		id := primitive.NewObjectID()
		donor := map[string]interface{}{
			"_id":        id,
			"type":       "individual",
			"first_name": firstName,
			"last_name":  lastName,
			"email":      fmt.Sprintf("%s.%s%d@email.com", strings.ToLower(firstName), strings.ToLower(lastName), randInt(1, 9999)),
			"phone":      fmt.Sprintf("+1-555-%04d", rand.Intn(10000)),
			"status":     "active",
			"created_at": randomDate(startDate, currentDate),
			"updated_at": currentDate,
		}

		db.Collection("donors").InsertOne(ctx, donor)
		donorIDs = append(donorIDs, id)
	}

	log.Printf("✓ Created %d donors", donorCount)
}

func seedCampaigns() {
	log.Println("\n📢 Seeding campaigns...")

	campaigns := []struct {
		name string
		goal float64
		year int
	}{
		{"Annual Fundraiser 2023", 50000, 0},
		{"Medical Fund 2023", 30000, 0},
		{"Annual Fundraiser 2024", 60000, 1},
		{"Medical Fund 2024", 35000, 1},
		{"Annual Fundraiser 2025", 70000, 2},
	}

	for _, c := range campaigns {
		campaignDate := startDate.AddDate(c.year, 0, 0)
		endDate := campaignDate.AddDate(0, 3, 0)

		var status string
		if endDate.Before(currentDate) {
			status = "completed"
		} else {
			status = "active"
		}

		id := primitive.NewObjectID()
		campaign := map[string]interface{}{
			"_id":         id,
			"name":        c.name,
			"description": fmt.Sprintf("Campaign: %s", c.name),
			"goal":        c.goal,
			"raised":      float64(randInt(0, int(c.goal))),
			"start_date":  campaignDate,
			"end_date":    endDate,
			"status":      status,
			"created_by":  adminIDs[rand.Intn(len(adminIDs))],
			"created_at":  campaignDate.AddDate(0, 0, -7),
			"updated_at":  currentDate,
		}

		db.Collection("campaigns").InsertOne(ctx, campaign)
		campaignIDs = append(campaignIDs, id)
	}

	log.Printf("✓ Created %d campaigns", len(campaigns))
}

func seedDonations() {
	log.Println("\n💵 Seeding donations...")

	donationCount := randInt(500, 1500)

	for i := 0; i < donationCount; i++ {
		donor := donorIDs[rand.Intn(len(donorIDs))]
		donationDate := randomDate(startDate, currentDate)

		var amount float64
		if rand.Float64() < 0.6 {
			amount = float64(randInt(10, 50))
		} else if rand.Float64() < 0.9 {
			amount = float64(randInt(51, 500))
		} else {
			amount = float64(randInt(501, 5000))
		}

		var campaignID *primitive.ObjectID
		if rand.Float64() < 0.4 && len(campaignIDs) > 0 {
			cid := campaignIDs[rand.Intn(len(campaignIDs))]
			campaignID = &cid
		}

		donation := map[string]interface{}{
			"_id":        primitive.NewObjectID(),
			"donor_id":   donor,
			"amount":     amount,
			"currency":   "USD",
			"type":       randomChoice([]string{"one_time", "recurring"}),
			"method":     randomChoice([]string{"credit_card", "paypal", "check"}),
			"status":     "completed",
			"date":       donationDate,
			"created_at": donationDate,
			"updated_at": currentDate,
		}

		if campaignID != nil {
			donation["campaign_id"] = *campaignID
		}

		db.Collection("donations").InsertOne(ctx, donation)
	}

	log.Printf("✓ Created %d donations", donationCount)
}

func seedAdoptions() {
	log.Println("\n🏠 Seeding adoptions...")

	cursor, _ := db.Collection("animals").Find(ctx, primitive.M{"status": "adopted"})
	var adoptedAnimals []struct {
		ID         primitive.ObjectID `bson:"_id"`
		IntakeDate time.Time          `bson:"intake_date"`
	}
	cursor.All(ctx, &adoptedAnimals)

	for _, animal := range adoptedAnimals {
		applicationDate := randomDate(animal.IntakeDate.AddDate(0, 0, 7), animal.IntakeDate.AddDate(0, 2, 0))

		appID := primitive.NewObjectID()
		application := map[string]interface{}{
			"_id":        appID,
			"animal_id":  animal.ID,
			"first_name": randomChoice([]string{"John", "Jane", "Michael", "Sarah"}),
			"last_name":  randomChoice([]string{"Smith", "Johnson", "Williams"}),
			"email":      fmt.Sprintf("adopter%d@email.com", randInt(1, 9999)),
			"phone":      fmt.Sprintf("+1-555-%04d", rand.Intn(10000)),
			"status":     "approved",
			"created_at": applicationDate,
			"updated_at": currentDate,
		}

		db.Collection("adoption_applications").InsertOne(ctx, application)

		adoptionDate := applicationDate.AddDate(0, 0, randInt(7, 21))

		adoption := map[string]interface{}{
			"_id":            primitive.NewObjectID(),
			"animal_id":      animal.ID,
			"application_id": appID,
			"adopter_name":   application["first_name"].(string) + " " + application["last_name"].(string),
			"adopter_email":  application["email"],
			"adoption_date":  adoptionDate,
			"adoption_fee":   float64(randInt(50, 300)),
			"status":         "finalized",
			"processed_by":   employeeIDs[rand.Intn(len(employeeIDs))],
			"created_at":     adoptionDate,
			"updated_at":     currentDate,
		}

		db.Collection("adoptions").InsertOne(ctx, adoption)
	}

	log.Printf("✓ Created %d adoptions", len(adoptedAnimals))
}

func seedVolunteers() {
	log.Println("\n🙋 Seeding volunteers...")

	volunteerCount := randInt(minVolunteers, maxVolunteers)

	for i := 0; i < volunteerCount; i++ {
		firstName := randomChoice([]string{"Alex", "Sam", "Jordan", "Taylor", "Morgan"})
		lastName := randomChoice([]string{"Anderson", "Baker", "Clark", "Davis"})

		joinDate := randomDate(startDate, currentDate.AddDate(0, -3, 0))

		id := primitive.NewObjectID()
		volunteer := map[string]interface{}{
			"_id":          id,
			"first_name":   firstName,
			"last_name":    lastName,
			"email":        fmt.Sprintf("%s.%s%d@email.com", strings.ToLower(firstName), strings.ToLower(lastName), randInt(1, 999)),
			"phone":        fmt.Sprintf("+1-555-%04d", rand.Intn(10000)),
			"status":       "active",
			"hours_logged": float64(randInt(10, 500)),
			"start_date":   joinDate,
			"created_at":   joinDate,
			"updated_at":   currentDate,
		}

		db.Collection("volunteers").InsertOne(ctx, volunteer)
		volunteerIDs = append(volunteerIDs, id)
	}

	log.Printf("✓ Created %d volunteers", volunteerCount)
}

func seedEvents() {
	log.Println("\n📅 Seeding events...")

	eventTemplates := []string{"Adoption Fair", "Volunteer Training", "Fundraiser Gala", "Community Outreach"}
	eventCount := 0

	currentEventDate := startDate
	for currentEventDate.Before(currentDate) {
		if rand.Float64() < 0.3 {
			eventDate := currentEventDate.AddDate(0, 0, randInt(0, 90))
			if eventDate.After(currentDate) {
				break
			}

			event := map[string]interface{}{
				"_id":         primitive.NewObjectID(),
				"name":        randomChoice(eventTemplates),
				"description": "Community event",
				"start_date":  eventDate,
				"end_date":    eventDate.Add(time.Hour * 4),
				"location":    "Happy Paws Center",
				"status":      "completed",
				"created_by":  adminIDs[rand.Intn(len(adminIDs))],
				"created_at":  eventDate.AddDate(0, 0, -14),
				"updated_at":  currentDate,
			}

			db.Collection("events").InsertOne(ctx, event)
			eventCount++
		}

		currentEventDate = currentEventDate.AddDate(0, 3, 0)
	}

	log.Printf("✓ Created %d events", eventCount)
}

func seedVeterinaryRecords() {
	log.Println("\n🏥 Seeding veterinary records...")

	visitCount := 0
	vaccinationCount := 0

	for _, animalID := range animalIDs {
		numVisits := randInt(1, 3)

		for i := 0; i < numVisits; i++ {
			visitDate := randomDate(startDate, currentDate)

			visit := map[string]interface{}{
				"_id":          primitive.NewObjectID(),
				"animal_id":    animalID,
				"visit_date":   visitDate,
				"visit_type":   randomChoice([]string{"Checkup", "Vaccination", "Surgery"}),
				"veterinarian": "Dr. " + randomChoice([]string{"Smith", "Johnson", "Williams"}),
				"diagnosis":    "Healthy",
				"created_at":   visitDate,
				"updated_at":   currentDate,
			}

			db.Collection("veterinary_visits").InsertOne(ctx, visit)
			visitCount++
		}

		if rand.Float64() < 0.8 {
			vaccineDate := randomDate(startDate, currentDate)
			vaccination := map[string]interface{}{
				"_id":               primitive.NewObjectID(),
				"animal_id":         animalID,
				"vaccine_name":      randomChoice([]string{"Rabies", "DHPP", "Bordetella"}),
				"administered_date": vaccineDate,
				"veterinarian":      "Dr. " + randomChoice([]string{"Smith", "Johnson"}),
				"created_at":        vaccineDate,
				"updated_at":        currentDate,
			}

			db.Collection("vaccinations").InsertOne(ctx, vaccination)
			vaccinationCount++
		}
	}

	log.Printf("✓ Created %d vet visits and %d vaccinations", visitCount, vaccinationCount)
}

func seedPartners() {
	log.Println("\n🤝 Seeding partners...")

	partners := []string{"City Vet Clinic", "Pet Supply Co", "Foster Network", "Transport Service"}

	for _, name := range partners {
		id := primitive.NewObjectID()
		partner := map[string]interface{}{
			"_id":          id,
			"name":         name,
			"type":         randomChoice([]string{"veterinary", "supplier", "foster", "transport"}),
			"contact_name": randomChoice([]string{"John Smith", "Jane Doe"}),
			"email":        fmt.Sprintf("contact@%s.org", strings.ReplaceAll(strings.ToLower(name), " ", "")),
			"phone":        fmt.Sprintf("+1-555-%04d", rand.Intn(10000)),
			"status":       "active",
			"created_at":   randomDate(startDate, currentDate),
			"updated_at":   currentDate,
		}

		db.Collection("partners").InsertOne(ctx, partner)
		partnerIDs = append(partnerIDs, id)
	}

	log.Printf("✓ Created %d partners", len(partners))
}

func seedInventory() {
	log.Println("\n📦 Seeding inventory...")

	items := []string{"Dog Food", "Cat Food", "Cat Litter", "Dog Toys", "Cat Toys", "Cleaning Supplies"}

	for _, name := range items {
		item := map[string]interface{}{
			"_id":          primitive.NewObjectID(),
			"name":         name,
			"category":     randomChoice([]string{"food", "supplies", "cleaning"}),
			"quantity":     randInt(50, 500),
			"unit":         "units",
			"min_quantity": randInt(10, 30),
			"created_at":   startDate,
			"updated_at":   currentDate,
		}

		db.Collection("inventory_items").InsertOne(ctx, item)
	}

	log.Printf("✓ Created %d inventory items", len(items))
}

func seedTasks() {
	log.Println("\n✅ Seeding tasks...")

	taskCount := randInt(100, 200)

	for i := 0; i < taskCount; i++ {
		createdDate := randomDate(startDate, currentDate)

		var status string
		if rand.Float64() < 0.7 {
			status = "completed"
		} else {
			status = "pending"
		}

		task := map[string]interface{}{
			"_id":         primitive.NewObjectID(),
			"title":       randomChoice([]string{"Clean kennels", "Feed animals", "Process applications", "Update website"}),
			"priority":    randomChoice([]string{"high", "medium", "low"}),
			"status":      status,
			"assigned_to": employeeIDs[rand.Intn(len(employeeIDs))],
			"created_by":  adminIDs[rand.Intn(len(adminIDs))],
			"created_at":  createdDate,
			"updated_at":  currentDate,
		}

		db.Collection("tasks").InsertOne(ctx, task)
	}

	log.Printf("✓ Created %d tasks", taskCount)
}

func seedDocuments() {
	log.Println("\n📄 Seeding documents...")

	docCount := randInt(50, 150)

	for i := 0; i < docCount; i++ {
		uploadDate := randomDate(startDate, currentDate)

		document := map[string]interface{}{
			"_id":         primitive.NewObjectID(),
			"name":        fmt.Sprintf("Document-%d.pdf", randInt(1000, 9999)),
			"type":        randomChoice([]string{"contract", "medical", "legal"}),
			"uploaded_by": employeeIDs[rand.Intn(len(employeeIDs))],
			"created_at":  uploadDate,
			"updated_at":  currentDate,
		}

		db.Collection("documents").InsertOne(ctx, document)
	}

	log.Printf("✓ Created %d documents", docCount)
}

func printSummary() {
	log.Println("\n" + strings.Repeat("=", 60))
	log.Println("📊 SEEDING SUMMARY")
	log.Println(strings.Repeat("=", 60))

	collections := []string{
		"users", "animals", "adoptions", "donors", "donations",
		"campaigns", "events", "volunteers", "veterinary_visits",
		"vaccinations", "partners", "inventory_items", "tasks", "documents",
	}

	for _, coll := range collections {
		count, _ := db.Collection(coll).CountDocuments(ctx, primitive.M{})
		log.Printf("%-25s: %d", coll, count)
	}

	log.Println(strings.Repeat("=", 60))
	log.Println("\n✨ Database seeded with 3 years of data!")
	log.Println("📧 Default password for all users: password123")
	log.Println("🔐 Super Admin: sarah.johnson@happypaws.org")
	log.Println(strings.Repeat("=", 60))
}

// Helper functions
func randInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}

func randomDate(min, max time.Time) time.Time {
	delta := max.Unix() - min.Unix()
	sec := rand.Int63n(delta)
	return min.Add(time.Duration(sec) * time.Second)
}

func randomFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func randomChoice(choices []string) string {
	return choices[rand.Intn(len(choices))]
}
