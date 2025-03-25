package repository

import (
	"log"

	"github.com/nefarius/cornelian/underlying/app"

	"github.com/nefarius/cornelian/underlying/app/conf"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/net/context"
)

type PersonRepository struct {
	Conf conf.MongoConf
}

func NewPersonRepository(Conf conf.MongoConf) *PersonRepository {
	return &PersonRepository{Conf: Conf}
}

func (r *PersonRepository) GetPersonByEmail(email string) app.Person {

	var client = r.Conf.MongoClient
	collection := client.Database("taffeite").Collection("person")
	var filter = bson.M{"email": email}
	result := collection.FindOne(context.Background(), filter)

	var person = app.Person{}
	var err = result.Decode(&person)
	if err != nil {
		log.Fatal(err)
	}

	return person
}

func (r *PersonRepository) InsertPersonIfNotPresent(person app.Person) error {
	var client = r.Conf.MongoClient
	collection := client.Database("taffeite").Collection("person")

	filter := bson.M{"email": person.Email, "admin": false}
	count, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		panic(err)
	}
	if count == 0 {
		_, err := collection.InsertOne(context.Background(), person)
		if err != nil {
			log.Fatal(err)
			panic(err)
		}
	}
	return err
}
