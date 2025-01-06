package service

import (
	"fmt"
	"sort"

	"github.com/nefarius/cornelian/underlying/app"
	"github.com/nefarius/cornelian/underlying/app/conf"
	"github.com/nefarius/cornelian/underlying/app/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QuestionService struct {
	CacheConf          conf.CacheConf
	QuestionRepository repository.QuestionRepository
	QuizRepository     repository.QuizRepository
}

func NewQuestionService(cache *conf.CacheConf, questionRepository *repository.QuestionRepository, quizRepository *repository.QuizRepository) *QuestionService {
	return &QuestionService{CacheConf: *cache, QuestionRepository: *questionRepository, QuizRepository: *quizRepository}
}

func (r *QuestionService) InsertQuizAndMakeCurrent(quizHeader string, quizText string, email string) error {
	var quizRepository = r.QuizRepository

	//	talk := app.Talk{
	//		ID:         primitive.NewObjectID().Hex(),
	//		Title:      quizHeader,
	//		AssignedTo: []string{"redvelvet@gmail.com"}, // change to assigned to
	//	}

	var quiz = app.Quiz{
		Id:          primitive.NewObjectID(),
		Header:      quizHeader,
		Description: quizText,
		Active:      true,
		Current:     true,
		Tag:         conf.CURRENT_TAG,
		Questions:   []app.Question{},
		AssignedTo:  []string{},
	}

	return quizRepository.InsertQuizAndMakeCurrent(quiz)
}

func (r *QuestionService) AssignQuizIfApplicable(email string) error {
	var questionRepository = r.QuestionRepository

	var key = conf.CURRENT_TAG
	r.CacheConf.Cache.Del(key)

	return questionRepository.AssignToCurrentQuiz(email)
}

func (r *QuestionService) GetQuiz() app.Quiz {
	var questionRepository = r.QuestionRepository

	var key = conf.CURRENT_TAG
	var cachedData, found = r.CacheConf.Cache.Get(key)

	if found {
		//		fmt.Println("cached data present")
		return cachedData.(app.Quiz)
	}

	var quiz = questionRepository.GetQuiz()
	//	fmt.Println("cached data set")
	r.CacheConf.Cache.Set(key, quiz, 0)

	return questionRepository.GetQuiz()
}

func (r *QuestionService) UpdateQuestion(questionId string, text string, answeredBy string) error {
	var questionRepository = r.QuestionRepository
	r.CacheConf.Cache.Del(conf.CURRENT_TAG)

	return questionRepository.UpdateQuestion(questionId, text, answeredBy)
}

func (r *QuestionService) AddQuestion(text string, answeredBy string) error {
	var questionRepository = r.QuestionRepository
	r.CacheConf.Cache.Del(conf.CURRENT_TAG)

	return questionRepository.AddQuestion(text, answeredBy)
}

func (r *QuestionService) SaveAnswer(id string, text string, answeredBy string) error {
	var questionRepository = r.QuestionRepository
	r.CacheConf.Cache.Del(conf.CURRENT_TAG)

	return questionRepository.SaveAnswer(id, text, answeredBy)
}

func (r *QuestionService) AllForAuthorInStatus(email string, status app.Status) []app.Question {
	var questions = r.AllForAssignedTo(email)
	filtered := make([]app.Question, 0)
	for _, q := range questions {
		if q.Status == status {
			filtered = append(filtered, q)
		}
	}

	return filtered
}

func (r *QuestionService) GetQuestion(id string) (app.Question, error) {
	var questionRepository = r.QuestionRepository
	return questionRepository.GetQuestion(id)
}

func (r *QuestionService) GetQuizzes() []app.Quiz {
	var quizRepository = r.QuizRepository
	var quizzes = quizRepository.GetQuizzes()

	filteredQuizzes := filterBy(quizzes, func(p app.Quiz) bool {
		return p.Active == true
	})

	return filteredQuizzes
}

func filterBy(quizzes []app.Quiz, filter func(app.Quiz) bool) []app.Quiz {
	var result []app.Quiz
	for _, p := range quizzes {
		if filter(p) {
			result = append(result, p)
		}
	}
	return result
}

func (r *QuestionService) UpdateQuiz(id string, header string, description string, current bool) error {
	var quizRepository = r.QuizRepository
	var identifier, err = primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to update question: %w", err)
	}
	r.CacheConf.Cache.Del(conf.CURRENT_TAG)
	return quizRepository.UpdateQuiz(identifier, header, description, current)
}

func (r *QuestionService) GetQuizById(id string) (app.Quiz, error) {
	var quizRepository = r.QuizRepository
	var identifier, err = primitive.ObjectIDFromHex(id)

	return quizRepository.GetQuizById(identifier), err
}

func (r *QuestionService) AllQuestions() []app.Question {
	var quiz = r.GetQuiz()

	out := quiz.Questions
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out
}

func (r *QuestionService) AllInStatus(status app.Status) []app.Question {
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

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func (r *QuestionService) AllForAssignedTo(email string) []app.Question {
	var quiz = r.GetQuiz()
	fmt.Println("email : " + email)
	var contains = contains(quiz.AssignedTo, email)
	if contains {
		fmt.Println("contains : " + email)
	}

	if !contains {
		return []app.Question{}
	}

	questions := quiz.Questions

	for _, question := range questions {
		if question.Answers != nil {
			question.Answers = filterByAssignedTo(question.Answers, email)
			break
		}
	}

	sort.Slice(questions, func(i, j int) bool {
		return questions[i].CreatedAt.After(questions[j].CreatedAt)
	})
	return questions
}

func filterByAssignedTo(answers []app.Answer, email string) []app.Answer {
	var filtered []app.Answer
	for _, answer := range answers {
		if answer.AnsweredBy == email {
			filtered = append(filtered, answer)
		}
	}
	return filtered
}
