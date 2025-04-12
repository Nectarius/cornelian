package repository

import (
	"log"

	"github.com/nefarius/cornelian/underlying/app"

	"github.com/nefarius/cornelian/underlying/app/conf"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
)

type PersonRepository struct {
	Conf conf.MongoConf
}

func (r PersonRepository) GetPersonById(personId string) app.Person {
	var client = r.Conf.MongoClient
	collection := client.Database("taffeite").Collection("person")
	objectID, err2 := primitive.ObjectIDFromHex(personId)
	if err2 != nil {
		log.Printf("error happened: %v", err2)
		return app.Person{}
	}
	var filter = bson.M{"id": objectID}
	result := collection.FindOne(context.Background(), filter)

	var person = app.Person{}
	var err = result.Decode(&person)
	if err != nil {
		log.Printf("error happened: %v", err)
		return app.Person{}
	}

	return person
}

func NewPersonRepository(Conf conf.MongoConf) *PersonRepository {
	return &PersonRepository{Conf: Conf}
}

func (r *PersonRepository) GetAll() []app.Person {

	var client = r.Conf.MongoClient
	collection := client.Database("taffeite").Collection("person")

	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	//var person = app.Person{}
	//var err = result.Decode(&person)
	var results []app.Person
	if err := cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}

	return results
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

	filter := bson.M{"email": person.Email}
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
