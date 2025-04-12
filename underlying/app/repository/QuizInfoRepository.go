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

type QuizInfoRepository struct {
	Conf conf.MongoConf
}

func (r QuizInfoRepository) ResetQuizInfo(email string, quizId primitive.ObjectID) error {

	var client = r.Conf.MongoClient
	collection := client.Database("taffeite").Collection("quiz-info-data")
	filter := bson.M{"quizid": quizId, "email": email}

	update := bson.M{
		"$set": bson.M{
			"answers": []string{},
		},
	}

	_, err := collection.UpdateOne(context.Background(), filter, update)

	return err

}

func NewQuizInfoRepository(Conf conf.MongoConf) *QuizInfoRepository {
	return &QuizInfoRepository{Conf: Conf}
}

func (r *QuizInfoRepository) GetQuizByIdAndEmail(id primitive.ObjectID, email string) app.QuizInfo {

	var client = r.Conf.MongoClient
	collection := client.Database("taffeite").Collection("quiz-info-data")
	var filter = bson.M{"quizid": id}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}
	var results []app.QuizInfo
	for cursor.Next(context.TODO()) {
		var data app.QuizInfo
		if err := cursor.Decode(&data); err != nil {
			log.Fatal(err)
		}
		results = append(results, data)
	}

	if len(results) < 0 {
		return app.QuizInfo{Email: ""}
	}

	var filteredData []app.QuizInfo
	for _, item := range results {
		if item.Email == email && item.Answers != nil {
			filteredData = append(filteredData, item)
		}
	}

	if len(filteredData) > 0 {
		return filteredData[0]
	} else {
		return app.QuizInfo{Email: ""}
	}

}

func (r *QuizInfoRepository) InsertNewAnswer(id primitive.ObjectID, email string, questionID string, started time.Time) error {
	var client = r.Conf.MongoClient
	collection := client.Database("taffeite").Collection("quiz-info-data")
	var filter = bson.M{"quizid": id, "email": email}
	result := collection.FindOne(context.Background(), filter)

	var quiz = app.QuizInfo{}
	var err = result.Decode(&quiz)
	if err != nil {
		print("quiz info was not handled" + id.Hex() + " email " + email)
		log.Fatal(err)
	}

	for _, answer := range quiz.Answers {
		if answer.Current == true {
			return nil
		}
	}

	var answer = app.AnswerInfo{
		ID:         uuid.NewString(),
		QuestionId: questionID,
		Text:       "",
		Current:    true,
		Started:    started,
	}

	addAnswer := bson.M{"$push": bson.M{"answers": answer}}
	_, err2 := collection.UpdateOne(context.Background(), filter, addAnswer)

	return err2
}

func (r *QuizInfoRepository) UpdateAnswer(id primitive.ObjectID, email string, questionID string, answerText string) error {
	var client = r.Conf.MongoClient
	collection := client.Database("taffeite").Collection("quiz-info-data")
	filter := bson.M{"quizid": id, "email": email, "answers.questionid": questionID}

	println("quizid" + id.Hex() + " email " + email + " question id " + questionID)
	update := bson.M{
		"$set": bson.M{
			"answers.$.text":      answerText,
			"answers.$.completed": time.Now(),
			"answers.$.current":   false,
		},
	}

	_, err := collection.UpdateOne(context.Background(), filter, update)

	return err
}

func (r *QuizInfoRepository) InsertQuizInfo(quiz app.QuizInfo) error {
	var client = r.Conf.MongoClient
	collection := client.Database("taffeite").Collection("quiz-info-data")

	_, err2 := collection.InsertOne(context.Background(), quiz)
	if err2 != nil {
		log.Fatal(err2)
		panic(err2)
	}
	return err2
}
