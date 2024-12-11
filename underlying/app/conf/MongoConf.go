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

func NewMongoConf() *MongoConf {

	var mongoUrlPath = "mongodb://admin:8BlanchE8@80.190.84.21:27017/?directConnection=true&serverSelectionTimeoutMS=2000"

	clientOptions := options.Client().ApplyURI(mongoUrlPath)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	return &MongoConf{
		MongoUri:    mongoUrlPath,
		MongoClient: client,
	}
}

func (r *MongoConf) Clear() {
	err := r.MongoClient.Disconnect(context.Background())
	log.Println("MongoConf Disconnected ")
	if err != nil {
		log.Fatal(err)
	}

}
