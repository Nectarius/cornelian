package repository

import (
	"fmt"
	"log"

	"github.com/nefarius/cornelian/underlying/app"

	"github.com/nefarius/cornelian/underlying/app/conf"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
)

type QuizRepository struct {
	Conf conf.MongoConf
}

func NewQuizRepository(Conf conf.MongoConf) *QuizRepository {
	return &QuizRepository{Conf: Conf}
}

func (r *QuizRepository) GetQuizById(id primitive.ObjectID) app.Quiz {

	var client = r.Conf.MongoClient
	collection := client.Database("taffeite").Collection("quiz-data")
	var filter = bson.M{"id": id}
	result := collection.FindOne(context.Background(), filter)

	var quiz = app.Quiz{}
	var err = result.Decode(&quiz)
	if err != nil {
		log.Fatal(err)
	}

	return quiz
}

func (r *QuizRepository) InsertQuizAndMakeCurrent(quiz app.Quiz) error {
	var client = r.Conf.MongoClient
	collection := client.Database("taffeite").Collection("quiz-data")
	var current = false
	var update = bson.M{
		"$set": bson.M{
			"current": current,
		},
	}

	_, err := collection.UpdateMany(context.Background(), bson.M{}, update)
	if err != nil {
		fmt.Println("failed to update question: %w", err)
		return fmt.Errorf("failed to update question: %w", err)
	}

	quiz.Current = true

	_, err = collection.InsertOne(context.Background(), quiz)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	return err
}

func (r *QuizRepository) UpdateQuiz(id primitive.ObjectID, header string, description string) error {

	var client = r.Conf.MongoClient
	collection := client.Database("taffeite").Collection("quiz-data")
	var filter = bson.M{"id": id}

	var update = bson.M{
		"$set": bson.M{
			"description": description,
			"header":      header,
		},
	}

	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		fmt.Println("failed to update question: %w", err)
		return fmt.Errorf("failed to update question: %w", err)
	}

	return nil
}

func (r *QuizRepository) GetQuizzes() []app.Quiz {
	var client = r.Conf.MongoClient
	collection := client.Database("taffeite").Collection("quiz-data")
	cursor, err := collection.Find(context.TODO(), bson.M{}) // Empty filter selects all documents
	if err != nil {
		panic(err)
	}

	defer cursor.Close(context.TODO())

	var results []app.Quiz
	if err := cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	return results
}
