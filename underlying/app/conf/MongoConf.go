package conf

import (
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
)

type MongoConf struct {
	MongoUri    string
	MongoClient *mongo.Client
}

func NewMongoConf() (*MongoConf, error) {

	var mongoUrlPath = "mongodb://admin:8BlanchE8@80.190.84.21:27017/?directConnection=true&serverSelectionTimeoutMS=2000"
	//var mongoUrlPath = "mongodb://nefarius:8BlanchE8@156.67.30.28:27017/?directConnection=true&serverSelectionTimeoutMS=2000"
	//
	clientOptions := options.Client().ApplyURI(mongoUrlPath)
	client, err := mongo.Connect(context.Background(), clientOptions)

	return &MongoConf{
		MongoUri:    mongoUrlPath,
		MongoClient: client,
	}, err
}

func (r *MongoConf) Clear() {
	err := r.MongoClient.Disconnect(context.Background())
	log.Println("MongoConf Disconnected ")
	if err != nil {
		log.Fatal(err)
	}

}
