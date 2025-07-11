package service

import (
	"fmt"
	"log"
	"math"
	"sort"
	"time"

	"github.com/nefarius/cornelian/underlying/app"
	"github.com/nefarius/cornelian/underlying/app/conf"
	"github.com/nefarius/cornelian/underlying/app/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// QuestionService provides methods to manage questions and quizzes.
type QuestionService struct {
	CacheConf          conf.CacheConf
	PersonRepository   repository.PersonRepository
	QuestionRepository repository.QuestionRepository
	QuizRepository     repository.QuizRepository
	QuizInfoRepository repository.QuizInfoRepository
	SettingsRepository repository.SettingsRepository
}

func (r *QuestionService) ResetCurrentQuiz(personId string) error {
	currentQuiz := r.GetQuiz()
	var person = r.PersonRepository.GetPersonById(personId)

	if person.Email == "" {
		log.Printf("Critical error happened")
		return fmt.Errorf("person not found")
	}
	// Reset quiz answers for the person

	var err = r.QuizRepository.ResetQuizAnswersByEmail(currentQuiz.Id, person.Email)

	// Reset quiz info for the person
	r.QuizInfoRepository.ResetQuizInfo(person.Email, currentQuiz.Id)

	r.CacheConf.Cache.Del(conf.CURRENT_TAG)

	return err

}

func (r *QuestionService) GetParticipants() []app.ParticipantView {
	var persons = r.PersonRepository.GetAll()
	var timeLimit int = 17
	participants := make([]app.ParticipantView, 0, len(persons))
	currentQuiz := r.GetQuiz()
	questions := currentQuiz.Questions
	for _, person := range persons {
		var quizInfo = r.QuizInfoRepository.GetQuizByIdAndEmail(currentQuiz.Id, person.Email)

		var correctAnswers int = 0
		var declinedDuetoTime int = 0
		var quizDuration float64 = 0.0
		userAnswersMap := make(map[string]app.AnswerInfo)
		for _, personAnswer := range quizInfo.Answers {
			userAnswersMap[personAnswer.QuestionId] = personAnswer
		}

		correctAnswersByQuestionMap := make(map[string]string)
		for _, q := range questions {
			for _, ac := range q.AnswerChoices {
				if ac.CorrectResponse {
					correctAnswersByQuestionMap[q.ID] = ac.Text
					break // Found the correct answer for this question, move to the next question
				}
			}
		}

		for questionID, correctText := range correctAnswersByQuestionMap {
			if userAnswer, ok := userAnswersMap[questionID]; ok { // Check if user answered this question
				var duration = userAnswer.Completed.Sub(userAnswer.Started).Seconds()
				quizDuration += duration
				if userAnswer.Text == correctText {
					if int(math.Floor(duration)) <= timeLimit {
						correctAnswers++
					} else {
						declinedDuetoTime++
					}
				}
			}
		}

		participants = append(participants, app.ParticipantView{
			Person:    person,
			Questions: questions,
			Answers:   quizInfo.Answers,
			SummaryView: app.SummaryView{
				CorrectResponses:  correctAnswers,
				DeclinedDuetoTime: declinedDuetoTime,
				QuizDuration:      quizDuration,
			},
		})
	}

	return participants
}

func (r *QuestionService) StartAnswering(quizId primitive.ObjectID, email string, questionId string) {
	var quizInfo = r.QuizInfoRepository.GetQuizByIdAndEmail(quizId, email)
	if quizInfo.Email == "" {
		log.Printf("Critical error happened")
		return
	}
	r.QuizInfoRepository.InsertNewAnswer(quizId, email, questionId, time.Now())

}

func (r *QuestionService) HandleAnswer(id primitive.ObjectID, email string, questionID string, answerText string) {
	quizInfo := r.QuizInfoRepository.GetQuizByIdAndEmail(id, email)
	if quizInfo.Email == "" {
		log.Printf("HandleAnswer: quiz info not found for quizId=%v, email=%s", id.Hex(), email)
		return
	}

	// Check if the question exists in the quiz info's answers
	found := false
	for _, ans := range quizInfo.Answers {
		if ans.QuestionId == questionID {
			found = true
			break
		}
	}
	if !found {
		log.Printf("HandleAnswer: questionID %s not found in answers for quizId=%v, email=%s", questionID, id.Hex(), email)
		return
	}

	err := r.QuizInfoRepository.UpdateAnswer(id, email, questionID, answerText)
	if err != nil {
		log.Printf("HandleAnswer: failed to update answer for quizId=%v, email=%s, questionID=%s: %v", id.Hex(), email, questionID, err)
	}
}

func (r *QuestionService) HandleNewAnswer(id primitive.ObjectID, email string, questionID string, started time.Time) {
	var quizInfo = r.QuizInfoRepository.GetQuizByIdAndEmail(id, email)
	if quizInfo.Email == "" {
		log.Printf("Critical error happened")
		return
	}
	r.QuizInfoRepository.InsertNewAnswer(id, email, questionID, started)
}

// NewQuestionService creates a new instance of QuestionService.
func NewQuestionService(cache *conf.CacheConf, personRepository *repository.PersonRepository, questionRepository *repository.QuestionRepository, quizRepository *repository.QuizRepository, quizInfoRepository *repository.QuizInfoRepository, settingsRepository *repository.SettingsRepository) *QuestionService {
	return &QuestionService{CacheConf: *cache, PersonRepository: *personRepository, QuestionRepository: *questionRepository, QuizRepository: *quizRepository, QuizInfoRepository: *quizInfoRepository, SettingsRepository: *settingsRepository}
}

// InsertQuizAndMakeCurrent inserts a new quiz and makes it the current quiz.
func (r *QuestionService) InsertQuizAndMakeCurrent(quizHeader string, quizText string, email string) error {
	quiz := app.Quiz{
		Id:          primitive.NewObjectID(),
		Header:      quizHeader,
		Description: quizText,
		Active:      true,
		Current:     true,
		Tag:         conf.CURRENT_TAG,
		Questions:   []app.Question{},
		AssignedTo:  []string{},
	}

	if err := r.QuizRepository.InsertQuizAndMakeCurrent(quiz); err != nil {
		log.Printf("Failed to insert quiz: %v", err)
		return err
	}
	return nil
}

func (r *QuestionService) InsertQuizWithQuestionsAndMakeCurrent(quizHeader string, quizText string, email string, questions []app.Question) error {

	var available = r.PersonRepository.GetAll()

	var participantList = make([]string, 0, len(available)) // Pre-allocate capacity for efficiency
	for _, p := range available {
		participantList = append(participantList, p.Email)
	}

	quiz := app.Quiz{
		Id:          primitive.NewObjectID(),
		Header:      quizHeader,
		Description: quizText,
		Active:      true,
		Current:     true,
		Tag:         conf.CURRENT_TAG,
		Questions:   questions,
		AssignedTo:  participantList,
	}

	if err := r.QuizRepository.InsertQuizAndMakeCurrent(quiz); err != nil {
		log.Printf("Failed to insert quiz: %v", err)
		return err
	}
	return nil
}

// AssignQuizIfApplicable assigns the current quiz to the given email if applicable.
func (r *QuestionService) AssignQuizIfApplicable(email string) error {
	r.CacheConf.Cache.Del(conf.CURRENT_TAG)
	if err := r.QuestionRepository.AssignToCurrentQuiz(email); err != nil {
		log.Printf("Failed to assign quiz: %v", err)
		return err
	}
	return nil
}

func (r *QuestionService) CreateIfNotExist(id primitive.ObjectID, email string) {
	var quizInfo = r.QuizInfoRepository.GetQuizByIdAndEmail(id, email)
	if quizInfo.Email == "" {
		var newQuizInfo = app.QuizInfo{
			QuizId:  id,
			Email:   email,
			Started: time.Now(),
			Answers: []app.AnswerInfo{},
		}
		r.QuizInfoRepository.InsertQuizInfo(newQuizInfo)
	}
}

// GetQuiz retrieves the current quiz from the cache or database.
func (r *QuestionService) GetQuiz() app.Quiz {
	key := conf.CURRENT_TAG
	cachedData, found := r.CacheConf.Cache.Get(key)
	if found {
		return cachedData.(app.Quiz)
	}

	var currentSettings = r.SettingsRepository.GetCurrent()
	quiz := r.QuizRepository.GetQuizById(currentSettings.QuizChoice)

	// Put question into cache
	// e.g., "q1"
	value := quiz
	cost := int64(1) // Cost can be 1 for simplicity or based on size (e.g., len(question.Text))

	// Set with no TTL (stays until evicted)
	if r.CacheConf.Cache.Set(conf.CURRENT_TAG, value, cost) {
		fmt.Printf("Successfully cached question %s\n", key)
	} else {
		fmt.Printf("Failed to cache question %s\n", key)
	}

	return quiz
}

// UpdateQuestion updates an existing question.
func (r *QuestionService) UpdateQuestion(questionId string, text string, answeredBy string) error {
	r.CacheConf.Cache.Del(conf.CURRENT_TAG)
	if err := r.QuestionRepository.UpdateQuestion(questionId, text, answeredBy); err != nil {
		log.Printf("Failed to update question: %v", err)
		return err
	}
	return nil
}

// AddQuestion adds a new question.
func (r *QuestionService) AddQuestion(text string, answeredBy string) error {
	r.CacheConf.Cache.Del(conf.CURRENT_TAG)
	if err := r.QuestionRepository.AddQuestion(text, answeredBy); err != nil {
		log.Printf("Failed to add question: %v", err)
		return err
	}
	return nil
}

// SaveAnswer saves an answer to a question.
func (r *QuestionService) SaveAnswer(id string, text string, answeredBy string) error {
	r.CacheConf.Cache.Del(conf.CURRENT_TAG)
	if err := r.QuestionRepository.SaveAnswer(id, text, answeredBy); err != nil {
		log.Printf("Failed to save answer: %v", err)
		return err
	}
	return nil
}

// AllForAuthorInOpenStatus returns all open questions for the given author.
func (r *QuestionService) AllForAuthorInOpenStatus(email string) []app.Question {
	questions := r.AllForAssignedTo(email)
	filtered := make([]app.Question, 0)
	for _, q := range questions {
		if q.Status == app.StatusOpen || len(q.Answers) == 0 {
			filtered = append(filtered, q)
		}
	}
	return filtered
}

// GetQuestion retrieves a question by its ID.
func (r *QuestionService) GetQuestion(id string) (app.Question, error) {
	return r.QuestionRepository.GetQuestion(id)
}

// GetQuizzes retrieves all active quizzes.
func (r *QuestionService) GetQuizzes() []app.Quiz {
	quizzes := r.QuizRepository.GetQuizzes()
	return filterBy(quizzes, func(p app.Quiz) bool {
		return p.Active
	})
}

// UpdateQuiz updates an existing quiz.
func (r *QuestionService) UpdateQuiz(id string, header string, description string, current bool) error {
	identifier, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to update quiz: %w", err)
	}
	r.CacheConf.Cache.Del(conf.CURRENT_TAG)
	if err := r.QuizRepository.UpdateQuiz(identifier, header, description, current); err != nil {
		log.Printf("Failed to update quiz: %v", err)
		return err
	}
	return nil
}

// GetQuizById retrieves a quiz by its ID.
func (r *QuestionService) GetQuizById(id string) (app.Quiz, error) {
	identifier, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return app.Quiz{}, fmt.Errorf("invalid quiz ID: %w", err)
	}
	return r.QuizRepository.GetQuizById(identifier), nil
}

// AllQuestions retrieves all questions sorted by creation date.
func (r *QuestionService) AllQuestions() []app.Question {
	quiz := r.GetQuiz()
	out := quiz.Questions
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out
}

func (r *QuestionService) CountQuestions() []app.Question {
	quiz := r.GetQuiz()
	out := make([]app.Question, 0)
	for _, q := range quiz.Questions {
		out = append(out, q)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out
}

// AllInStatus retrieves all questions with the given status.
func (r *QuestionService) AllInStatus(status app.Status) []app.Question {
	quiz := r.GetQuiz()
	out := make([]app.Question, 0)
	for _, q := range quiz.Questions {
		if q.Status == status {
			out = append(out, q)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out
}

// AllForAssignedTo retrieves all questions assigned to the given email.
func (r *QuestionService) AllForAssignedTo(email string) []app.Question {
	quiz := r.GetQuiz()
	if !contains(quiz.AssignedTo, email) {
		return []app.Question{}
	}

	questions := quiz.Questions
	filtered := make([]app.Question, 0)
	for _, question := range questions {
		if question.Answers != nil {
			question.Answers = filterByAssignedTo(question.Answers, email)
		}
		filtered = append(filtered, question)
	}
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].CreatedAt.After(filtered[j].CreatedAt)
	})
	return filtered
}

// Utility function to filter quizzes by a condition.
func filterBy(quizzes []app.Quiz, filter func(app.Quiz) bool) []app.Quiz {
	result := make([]app.Quiz, 0)
	for _, p := range quizzes {
		if filter(p) {
			result = append(result, p)
		}
	}
	return result
}

// Utility function to filter answers by the assigned email.
func filterByAssignedTo(answers []app.Answer, email string) []app.Answer {
	filtered := make([]app.Answer, 0)
	for _, answer := range answers {
		if answer.AnsweredBy == email {
			filtered = append(filtered, answer)
		}
	}
	return filtered
}

// Utility function to check if a string is in a slice.
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
