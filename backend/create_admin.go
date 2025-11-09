package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("animalsys")
	usersCollection := db.Collection("users")

	// Hash the password
	password := "AdminPass123!"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	// Check if admin user already exists
	email := "admin@animalsys.com"
	var existingUser bson.M
	err = usersCollection.FindOne(ctx, bson.M{"email": email}).Decode(&existingUser)

	if err == mongo.ErrNoDocuments {
		// Create new admin user
		user := bson.M{
			"_id":           primitive.NewObjectID(),
			"email":         email,
			"password_hash": string(hashedPassword),
			"first_name":    "System",
			"last_name":     "Administrator",
			"role":          "super_admin",
			"status":        "active",
			"created_at":    time.Now(),
			"updated_at":    time.Now(),
		}

		_, err = usersCollection.InsertOne(ctx, user)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("✓ Created admin user: %s\n", email)
		fmt.Printf("  Password: %s\n", password)
	} else if err != nil {
		log.Fatal(err)
	} else {
		// Update existing user's password
		update := bson.M{
			"$set": bson.M{
				"password_hash": string(hashedPassword),
				"updated_at":    time.Now(),
			},
		}
		_, err = usersCollection.UpdateOne(ctx, bson.M{"email": email}, update)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("✓ Updated password for existing admin user: %s\n", email)
		fmt.Printf("  New Password: %s\n", password)
	}

	fmt.Println("\nYou can now login with these credentials.")
}
