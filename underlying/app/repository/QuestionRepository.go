package repository

import (
	"log"

	"github.com/nefarius/cornelian/underlying/app"
	"github.com/nefarius/cornelian/underlying/app/conf"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
)

type QuestionRepository struct {
	Conf conf.MongoConf
}

func NewPanelViewRepository(Conf conf.MongoConf) *QuestionRepository {
	return &QuestionRepository{Conf: Conf}
}

func (r *QuestionRepository) GetQuiz() app.Quiz {
	var client = r.Conf.MongoClient
	collection := client.Database("taffeite").Collection("quiz-data")
	var filter = bson.M{"tag": "Test"}
	result := collection.FindOne(context.Background(), filter)

	var panelView = app.Quiz{}
	var err = result.Decode(&panelView)
	if err != nil {
		log.Fatal(err)
	}

	return panelView
}

func (r *QuestionRepository) InsertQuiz(panelView app.Quiz) primitive.ObjectID {
	var client = r.Conf.MongoClient
	collection := client.Database("taffeite").Collection("quiz-data")
	one, err := collection.InsertOne(context.Background(), panelView)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	return one.InsertedID.(primitive.ObjectID)
}
