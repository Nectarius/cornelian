package repository

import (
	"fmt"
	"log"

	"github.com/nefarius/cornelian/underlying/app"

	"github.com/nefarius/cornelian/underlying/app/conf"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
)

type QuizRepository struct {
	Conf conf.MongoConf
}

func (r QuizRepository) ResetQuizAnswersByEmail(id primitive.ObjectID, email string) error {

	var client = r.Conf.MongoClient
	collection := client.Database("taffeite").Collection("quiz-data")
	filter := bson.M{"id": id}
	update := bson.M{"$pull": bson.M{"questions.$[].answers": bson.M{"answeredby": email}}}
	opts := options.Update()
	fmt.Printf("ResetQuizAnswersByEmail: %s", email)
	fmt.Printf("id: %s", id.Hex())
	result, err := collection.UpdateMany(context.Background(), filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to update answers: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("no document matched the filter with id: %s", id)
	}

	return nil
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
			"tag":     " ",
		},
	}

	_, err := collection.UpdateMany(context.Background(), bson.M{}, update)
	if err != nil {
		fmt.Println("failed to update question: %w", err)
		return fmt.Errorf("failed to update question: %w", err)
	}

	quiz.Current = true
	quiz.Tag = "TheOne"

	_, err2 := collection.InsertOne(context.Background(), quiz)
	if err2 != nil {
		log.Fatal(err2)
		panic(err2)
	}
	return err2
}

func (r *QuizRepository) UpdateQuiz(id primitive.ObjectID, header string, description string, current bool) error {

	var client = r.Conf.MongoClient
	collection := client.Database("taffeite").Collection("quiz-data")
	var filter = bson.M{"id": id}
	var tag string
	if current {
		var update = bson.M{
			"$set": bson.M{
				"current": false,
				"tag":     " ",
			},
		}
		tag = "TheOne"
		_, err := collection.UpdateMany(context.Background(), bson.M{}, update)
		if err != nil {
			fmt.Println("failed to update question: %w", err)
			return fmt.Errorf("failed to update question: %w", err)
		}
	} else {
		tag = " "
	}

	var update = bson.M{
		"$set": bson.M{
			"description": description,
			"header":      header,
			"current":     current,
			"tag":         tag,
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
