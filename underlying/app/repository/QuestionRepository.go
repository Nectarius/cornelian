package repository

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/nefarius/cornelian/underlying/app"
	"github.com/nefarius/cornelian/underlying/app/conf"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
)

var CURRENT_TAG = "Test8"

type QuestionRepository struct {
	Conf conf.MongoConf
}

func NewPanelViewRepository(Conf conf.MongoConf) *QuestionRepository {
	return &QuestionRepository{Conf: Conf}
}

func (r *QuestionRepository) GetQuiz() app.Quiz {
	var client = r.Conf.MongoClient
	collection := client.Database("taffeite").Collection("quiz-data")
	var filter = bson.M{"tag": CURRENT_TAG}
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

func (r *QuestionRepository) SaveAnswer(id string, text string, answeredBy string) error {
	var client = r.Conf.MongoClient
	collection := client.Database("taffeite").Collection("quiz-data")
	var filter = bson.M{"tag": CURRENT_TAG}
	result := collection.FindOne(context.Background(), filter)

	var panelView = app.Quiz{}
	var err = result.Decode(&panelView)
	if err != nil {
		log.Fatal(err)
	}

	for _, q := range panelView.Questions {

		if q.ID == id {
			q.Answers = append(q.Answers, app.Answer{ID: uuid.NewString(), AnsweredBy: answeredBy, Text: text, AnsweredAt: time.Now()})
			q.Status = app.StatusAnswered
			removePreviousVersion := bson.M{"$pull": bson.M{"questions": bson.M{"id": id}}}

			_, err1 := collection.UpdateOne(context.Background(), filter, removePreviousVersion)
			if err1 != nil {
				log.Fatal(err)
			}
			addWithAnswer := bson.M{"$push": bson.M{"questions": q}}
			_, err2 := collection.UpdateOne(context.Background(), filter, addWithAnswer)

			if err2 != nil {
				log.Fatal(err)
			}

		}
	}

	return nil
}
