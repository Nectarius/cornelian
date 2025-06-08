package repository

import (
	"fmt"
	"log"

	"github.com/nefarius/cornelian/underlying/app"

	"github.com/nefarius/cornelian/underlying/app/conf"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/net/context"
)

type SettingsRepository struct {
	Conf conf.MongoConf
}

func NewSettingsRepository(Conf conf.MongoConf) *SettingsRepository {
	return &SettingsRepository{Conf: Conf}
}

func (r *SettingsRepository) GetAll() []app.QuizSettings {

	var client = r.Conf.MongoClient
	collection := client.Database("taffeite").Collection("settings")

	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	var results []app.QuizSettings
	if err := cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}

	return results
}

func (r *SettingsRepository) GetCurrent() app.QuizSettings {

	var client = r.Conf.MongoClient
	collection := client.Database("taffeite").Collection("settings")
	var result app.QuizSettings
	var filter = bson.M{"current": true}

	// Define your filter (e.g., find all users)

	// Add options, such as sorting to define what "first" means
	//findOptions := options.Find().SetSort(filter).SetLimit(1) // Sort by creation time, ascending, and limit to 1

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		log.Fatalf("Failed to find users: %v", err)
	}
	defer cursor.Close(context.Background()) // Always close the cursor!

	if cursor.Next(context.Background()) {
		err = cursor.Decode(&result)
		if err != nil {
			log.Fatalf("Failed to decode user: %v", err)
		}
		fmt.Printf("Found first user via cursor: %+v\n", result)
	} else {
		// This means cursor.Next(ctx) returned false, indicating no documents or end of results
		fmt.Println("No user found via cursor.")
	}

	// Check for any errors that occurred during cursor iteration
	if err := cursor.Err(); err != nil {
		log.Fatalf("Cursor error: %v", err)
	}

	return result
}

func (r *SettingsRepository) InsertSettings(settings app.QuizSettings) error {
	var client = r.Conf.MongoClient
	collection := client.Database("taffeite").Collection("settings")

	var current = false
	var update = bson.M{
		"$set": bson.M{
			"current": current,
		},
	}

	_, err := collection.UpdateMany(context.Background(), bson.M{}, update)
	if err != nil {
		fmt.Println("failed to update settings: %w", err)
		return fmt.Errorf("failed to update settings: %w", err)
	}

	_, err2 := collection.InsertOne(context.Background(), settings)
	if err2 != nil {
		log.Fatal(err2)
		panic(err2)
	}

	return err2
}
