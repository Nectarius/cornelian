package repository

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/nefarius/cornelian/underlying/app"

	"github.com/nefarius/cornelian/underlying/app/conf"
	"github.com/nefarius/cornelian/underlying/app/store"
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
	var filter = bson.M{"tag": conf.CURRENT_TAG}

	count, err2 := collection.CountDocuments(context.Background(), filter)
	if err2 != nil {
		panic(err2)
	}

	if count == 0 {
		return store.GetDefaultQuizData2()
	}

	result := collection.FindOne(context.Background(), filter)

	var panelView = app.Quiz{}
	var err = result.Decode(&panelView)
	if err != nil {
		log.Fatal(err)
	}

	return panelView
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

func (r *QuestionRepository) AssignToCurrentQuiz(email string) error {
	client := r.Conf.MongoClient
	collection := client.Database("taffeite").Collection("quiz-data")
	filter := bson.M{"tag": conf.CURRENT_TAG}

	fmt.Println("push email :" + email)
	assign := bson.M{"$push": bson.M{"assignedto": email}}

	_, err := collection.UpdateOne(context.Background(), filter, assign)
	if err != nil {
		return fmt.Errorf("failed to update question: %w", err)
	}

	return nil
}

func (r *QuestionRepository) UpdateQuestion(questionId string, text string, answeredBy string) error {
	client := r.Conf.MongoClient
	collection := client.Database("taffeite").Collection("quiz-data")
	filter := bson.M{"tag": conf.CURRENT_TAG, "questions.id": questionId}

	update := bson.M{
		"$set": bson.M{
			"questions.$.text": text,
			"questions.$.from": answeredBy,
		},
	}

	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to update question: %w", err)
	}

	return nil
}

//result := collection.FindOne(context.Background(), filter)

//var quiz = app.Quiz{}
//var err = result.Decode(&quiz)
//if err != nil {
//	log.Fatal(err)
//}

//var question = app.Question{
//	ID:        uuid.NewString(),
//	Talk:      conf.CURRENT_TALK,
//	Text:      text,
//	From:      answeredBy,
//	CreatedAt: time.Now(),
//	Status:    app.StatusOpen,
//}
// update := bson.M{"$set": bson.M{"inner.$[elem]": bson.M{"field1": "new_value"}}}
//
//	var arrayFilter = bson.M{"elem.ID": questionId}
// Update the embedded document with ID "embedded_doc_id_1"
//arrayFilters :=

//var updateOptions = options.Update().SetArrayFilters(arrayFilters)
//var updateQuestion = bson.M{"$set": bson.M{"questions.$[elem].text": text}}
//var filter = bson.M{
//	{"tag": conf.CURRENT_TAG},
//{"questions.id": questionId},
//	}
//var filter = bson.M{"tag": conf.CURRENT_TAG, "questions.id": questionId}

//   var filters= append(
// make([]interface{}, 0),
// bson.D{{"elem", bson.D{
//                   {"questions._id", questionId},
//                 }
// })
//  	options.Update().SetArrayFilters(
// 	options.ArrayFilters{
// 	Filters:

// )
//addQuestion := bson.M{"$push": bson.M{"questions": question}}
// arrayFilters := []interface{}{bson.M{"elem._id": questionId}}
//_, err2 := collection.UpdateOne(context.Background(), filter, updateQuestion,
//            options.Update().SetArrayFilters(
//            Filters : []interface{}{bson.M{"elem.id": questionId}},
//             )
//)
//	Filters: []interface{}{bson.M{"elem.id": questionId}},
//	fmt.Println(err2)

//	return error
//}

func (r *QuestionRepository) AddQuestion(text string, answeredBy string) error {
	var client = r.Conf.MongoClient
	collection := client.Database("taffeite").Collection("quiz-data")
	var filter = bson.M{"tag": conf.CURRENT_TAG}
	result := collection.FindOne(context.Background(), filter)

	var quiz = app.Quiz{}
	var err = result.Decode(&quiz)
	if err != nil {
		log.Fatal(err)
	}

	var question = app.Question{
		ID:        uuid.NewString(),
		Text:      text,
		From:      answeredBy,
		CreatedAt: time.Now(),
		Status:    app.StatusOpen,
		Answers:   []app.Answer{},
	}

	addQuestion := bson.M{"$push": bson.M{"questions": question}}
	_, err2 := collection.UpdateOne(context.Background(), filter, addQuestion)

	return err2
}

func (r *QuestionRepository) InsertQuiz(quiz app.Quiz) primitive.ObjectID {
	var client = r.Conf.MongoClient
	collection := client.Database("taffeite").Collection("quiz-data")
	one, err := collection.InsertOne(context.Background(), quiz)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	return one.InsertedID.(primitive.ObjectID)
}

func (r *QuestionRepository) SaveAnswer(id string, text string, answeredBy string) error {
	client := r.Conf.MongoClient
	collection := client.Database("taffeite").Collection("quiz-data")
	filter := bson.M{"tag": conf.CURRENT_TAG, "questions.id": id}

	update := bson.M{
		"$push": bson.M{
			"questions.$.answers": app.Answer{
				ID:         uuid.NewString(),
				AnsweredBy: answeredBy,
				Text:       text,
				AnsweredAt: time.Now(),
			},
		},
		"$set": bson.M{
			"questions.$.status": app.StatusAnswered,
		},
	}

	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to update answer: %w", err)
	}

	return nil
}
