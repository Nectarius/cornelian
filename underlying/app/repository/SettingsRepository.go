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
