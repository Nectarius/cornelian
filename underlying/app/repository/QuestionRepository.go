package repository

import (
	"fmt"
	"log"
	"slices"
	"sort"
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

func (r *QuestionRepository) AllForAuthorInStatus(email string, status app.Status) []app.Question {
	var quiz = r.GetQuiz()
	out := make([]app.Question, 0)
	for _, q := range quiz.Questions {
		if slices.Contains(q.Talk.AssignedTo, email) && q.Status == status {
			out = append(out, q)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out
}

func (r *QuestionRepository) GetQuestion(id string) (app.Question, error) {
	var quiz = r.GetQuiz()
	for _, q := range quiz.Questions {
		if q.ID == id {
			return q, nil
		}
	}
	return app.Question{}, fmt.Errorf("question identified by %v not found", id)
}

func (r *QuestionRepository) AllQuestions() []app.Question {
	out := r.GetQuiz().Questions
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
		//return out[i].CreatedAt.After(out[j].CreatedAt) && out[i].Status == app.StatusOpen
	})
	return out
}

func (r *QuestionRepository) AllInStatus(status app.Status) []app.Question {
	var quiz = r.GetQuiz()
	out := make([]app.Question, 0)
	for _, q := range quiz.Questions {
		if q.Status == status {
			out = append(out, q)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		//return (out[j].Status == app.StatusOpen && out[i].Status != app.StatusOpen) && out[i].CreatedAt.After(out[j].CreatedAt)
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out
}

func (r *QuestionRepository) AllForAssignedTo(email string) []app.Question {
	var quiz = r.GetQuiz()
	out := make([]app.Question, 0)
	for _, q := range quiz.Questions {
		if slices.Contains(q.Talk.AssignedTo, email) {
			out = append(out, q)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out
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
