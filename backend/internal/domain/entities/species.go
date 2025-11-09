package entities

// SpeciesInfo holds information about a species
type SpeciesInfo struct {
	CommonName MultilingualName `json:"common_name"`
	Category   AnimalCategory   `json:"category"`
	Breeds     []string         `json:"breeds,omitempty"`
}

// PredefinedSpecies contains commonly encountered species organized by category
var PredefinedSpecies = map[AnimalCategory][]SpeciesInfo{
	CategoryMammal: {
		{
			CommonName: MultilingualName{English: "Dog", Polish: "Pies", Latin: "Canis lupus familiaris"},
			Category:   CategoryMammal,
			Breeds: []string{
				"Mixed Breed", "Labrador Retriever", "German Shepherd", "Golden Retriever",
				"Bulldog", "Beagle", "Poodle", "Rottweiler", "Yorkshire Terrier",
				"Boxer", "Dachshund", "Siberian Husky", "Great Dane", "Doberman Pinscher",
				"Shih Tzu", "Chihuahua", "Pomeranian", "Border Collie", "Australian Shepherd",
				"Cocker Spaniel", "Maltese", "Jack Russell Terrier",
			},
		},
		{
			CommonName: MultilingualName{English: "Cat", Polish: "Kot", Latin: "Felis catus"},
			Category:   CategoryMammal,
			Breeds: []string{
				"Mixed Breed", "Persian", "Maine Coon", "Siamese", "Ragdoll",
				"British Shorthair", "Abyssinian", "Birman", "Oriental Shorthair",
				"Sphynx", "Devon Rex", "American Shorthair", "Scottish Fold",
				"Bengal", "Russian Blue", "Norwegian Forest Cat",
			},
		},
		{
			CommonName: MultilingualName{English: "Rabbit", Polish: "Królik", Latin: "Oryctolagus cuniculus"},
			Category:   CategoryMammal,
			Breeds:     []string{"Mixed Breed", "Dutch", "Mini Lop", "Lionhead", "Flemish Giant", "Rex", "Angora"},
		},
		{
			CommonName: MultilingualName{English: "Guinea Pig", Polish: "Świnka morska", Latin: "Cavia porcellus"},
			Category:   CategoryMammal,
		},
		{
			CommonName: MultilingualName{English: "Hamster", Polish: "Chomik", Latin: "Cricetinae"},
			Category:   CategoryMammal,
		},
		{
			CommonName: MultilingualName{English: "Ferret", Polish: "Fretka", Latin: "Mustela putorius furo"},
			Category:   CategoryMammal,
		},
		{
			CommonName: MultilingualName{English: "Rat", Polish: "Szczur", Latin: "Rattus"},
			Category:   CategoryMammal,
		},
		{
			CommonName: MultilingualName{English: "Mouse", Polish: "Mysz", Latin: "Mus"},
			Category:   CategoryMammal,
		},
		{
			CommonName: MultilingualName{English: "Hedgehog", Polish: "Jeż", Latin: "Erinaceinae"},
			Category:   CategoryMammal,
		},
	},

	CategoryBird: {
		{
			CommonName: MultilingualName{English: "Parrot", Polish: "Papuga", Latin: "Psittaciformes"},
			Category:   CategoryBird,
			Breeds: []string{
				"African Grey", "Cockatiel", "Budgerigar", "Macaw", "Cockatoo",
				"Amazon Parrot", "Conure", "Lovebird", "Parakeet",
			},
		},
		{
			CommonName: MultilingualName{English: "Canary", Polish: "Kanarek", Latin: "Serinus canaria"},
			Category:   CategoryBird,
		},
		{
			CommonName: MultilingualName{English: "Finch", Polish: "Zięba", Latin: "Fringillidae"},
			Category:   CategoryBird,
		},
		{
			CommonName: MultilingualName{English: "Pigeon", Polish: "Gołąb", Latin: "Columba"},
			Category:   CategoryBird,
		},
		{
			CommonName: MultilingualName{English: "Dove", Polish: "Gołębica", Latin: "Columbidae"},
			Category:   CategoryBird,
		},
		{
			CommonName: MultilingualName{English: "Chicken", Polish: "Kurczak", Latin: "Gallus gallus domesticus"},
			Category:   CategoryBird,
		},
	},

	CategoryReptile: {
		{
			CommonName: MultilingualName{English: "Turtle", Polish: "Żółw", Latin: "Testudines"},
			Category:   CategoryReptile,
			Breeds:     []string{"Red-Eared Slider", "Box Turtle", "Russian Tortoise", "Sulcata Tortoise", "Greek Tortoise"},
		},
		{
			CommonName: MultilingualName{English: "Tortoise", Polish: "Żółw lądowy", Latin: "Testudinidae"},
			Category:   CategoryReptile,
		},
		{
			CommonName: MultilingualName{English: "Lizard", Polish: "Jaszczurka", Latin: "Lacertilia"},
			Category:   CategoryReptile,
			Breeds:     []string{"Bearded Dragon", "Leopard Gecko", "Iguana", "Chameleon", "Blue-Tongued Skink"},
		},
		{
			CommonName: MultilingualName{English: "Snake", Polish: "Wąż", Latin: "Serpentes"},
			Category:   CategoryReptile,
			Breeds:     []string{"Ball Python", "Corn Snake", "King Snake", "Garter Snake", "Boa Constrictor"},
		},
		{
			CommonName: MultilingualName{English: "Gecko", Polish: "Gekon", Latin: "Gekkota"},
			Category:   CategoryReptile,
		},
	},

	CategoryAmphibian: {
		{
			CommonName: MultilingualName{English: "Frog", Polish: "Żaba", Latin: "Anura"},
			Category:   CategoryAmphibian,
			Breeds:     []string{"African Dwarf Frog", "Pac-Man Frog", "Tree Frog", "Poison Dart Frog"},
		},
		{
			CommonName: MultilingualName{English: "Toad", Polish: "Ropucha", Latin: "Bufonidae"},
			Category:   CategoryAmphibian,
		},
		{
			CommonName: MultilingualName{English: "Salamander", Polish: "Salamandra", Latin: "Caudata"},
			Category:   CategoryAmphibian,
		},
		{
			CommonName: MultilingualName{English: "Newt", Polish: "Traszka", Latin: "Pleurodelinae"},
			Category:   CategoryAmphibian,
		},
	},

	CategoryFish: {
		{
			CommonName: MultilingualName{English: "Goldfish", Polish: "Złota rybka", Latin: "Carassius auratus"},
			Category:   CategoryFish,
		},
		{
			CommonName: MultilingualName{English: "Betta Fish", Polish: "Bojownik syjamski", Latin: "Betta splendens"},
			Category:   CategoryFish,
		},
		{
			CommonName: MultilingualName{English: "Guppy", Polish: "Gupik", Latin: "Poecilia reticulata"},
			Category:   CategoryFish,
		},
		{
			CommonName: MultilingualName{English: "Molly", Polish: "Molinezja", Latin: "Poecilia sphenops"},
			Category:   CategoryFish,
		},
		{
			CommonName: MultilingualName{English: "Tetra", Polish: "Tetra", Latin: "Characidae"},
			Category:   CategoryFish,
		},
		{
			CommonName: MultilingualName{English: "Koi", Polish: "Koi", Latin: "Cyprinus rubrofuscus"},
			Category:   CategoryFish,
		},
	},

	CategoryInvertebrate: {
		{
			CommonName: MultilingualName{English: "Tarantula", Polish: "Ptasznik", Latin: "Theraphosidae"},
			Category:   CategoryInvertebrate,
		},
		{
			CommonName: MultilingualName{English: "Hermit Crab", Polish: "Rak pustelnik", Latin: "Paguroidea"},
			Category:   CategoryInvertebrate,
		},
		{
			CommonName: MultilingualName{English: "Snail", Polish: "Ślimak", Latin: "Gastropoda"},
			Category:   CategoryInvertebrate,
		},
		{
			CommonName: MultilingualName{English: "Scorpion", Polish: "Skorpion", Latin: "Scorpiones"},
			Category:   CategoryInvertebrate,
		},
	},

	CategoryFarmAnimal: {
		{
			CommonName: MultilingualName{English: "Horse", Polish: "Koń", Latin: "Equus caballus"},
			Category:   CategoryFarmAnimal,
		},
		{
			CommonName: MultilingualName{English: "Pony", Polish: "Kucyk", Latin: "Equus ferus caballus"},
			Category:   CategoryFarmAnimal,
		},
		{
			CommonName: MultilingualName{English: "Donkey", Polish: "Osioł", Latin: "Equus asinus"},
			Category:   CategoryFarmAnimal,
		},
		{
			CommonName: MultilingualName{English: "Goat", Polish: "Koza", Latin: "Capra aegagrus hircus"},
			Category:   CategoryFarmAnimal,
		},
		{
			CommonName: MultilingualName{English: "Sheep", Polish: "Owca", Latin: "Ovis aries"},
			Category:   CategoryFarmAnimal,
		},
		{
			CommonName: MultilingualName{English: "Pig", Polish: "Świnia", Latin: "Sus scrofa domesticus"},
			Category:   CategoryFarmAnimal,
		},
		{
			CommonName: MultilingualName{English: "Cow", Polish: "Krowa", Latin: "Bos taurus"},
			Category:   CategoryFarmAnimal,
		},
		{
			CommonName: MultilingualName{English: "Duck", Polish: "Kaczka", Latin: "Anas platyrhynchos"},
			Category:   CategoryFarmAnimal,
		},
		{
			CommonName: MultilingualName{English: "Goose", Polish: "Gęś", Latin: "Anser anser"},
			Category:   CategoryFarmAnimal,
		},
		{
			CommonName: MultilingualName{English: "Turkey", Polish: "Indyk", Latin: "Meleagris gallopavo"},
			Category:   CategoryFarmAnimal,
		},
	},
}

// GetSpeciesByCategory returns all species for a given category
func GetSpeciesByCategory(category AnimalCategory) []SpeciesInfo {
	if species, ok := PredefinedSpecies[category]; ok {
		return species
	}
	return []SpeciesInfo{}
}

// FindSpecies searches for a species by name and category
func FindSpecies(name string, category AnimalCategory) *SpeciesInfo {
	species := GetSpeciesByCategory(category)
	for _, s := range species {
		if s.CommonName.English == name || s.CommonName.Polish == name {
			return &s
		}
	}
	return nil
}
