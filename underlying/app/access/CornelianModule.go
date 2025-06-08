package access

import (
	"log"

	"github.com/nefarius/cornelian/underlying/app/conf"
	"github.com/nefarius/cornelian/underlying/app/repository"
	"github.com/nefarius/cornelian/underlying/app/service"
)

// CornelianModule holds the repositories, services, and configurations for the application.
type CornelianModule struct {
	PersonRepository   *repository.PersonRepository
	SettingsRepository *repository.SettingsRepository
	QuestionRepository *repository.QuestionRepository
	QuestionService    *service.QuestionService
	MongoConf          *conf.MongoConf
	CacheConf          *conf.CacheConf
}

// NewCornelianModule creates and initializes a new instance of CornelianModule.
func NewCornelianModule() *CornelianModule {
	mongoConf, err := conf.NewMongoConf()
	if err != nil {
		log.Fatalf("Failed to initialize Mongo configuration: %v", err)
	}

	cacheConf, err := conf.NewCacheConf()
	if err != nil {
		log.Fatalf("Failed to initialize Cache configuration: %v", err)
	}

	questionRepository := &repository.QuestionRepository{Conf: *mongoConf}

	personRepository := &repository.PersonRepository{Conf: *mongoConf}
	settingsRepository := &repository.SettingsRepository{Conf: *mongoConf}
	quizRepository := &repository.QuizRepository{Conf: *mongoConf}
	quizInfoRepository := &repository.QuizInfoRepository{Conf: *mongoConf}
	questionService := service.NewQuestionService(cacheConf, personRepository, questionRepository, quizRepository, quizInfoRepository)

	//var defaultData = store.GetDefaultQuizData2()
	//quizRepository.InsertQuizAndMakeCurrent(defaultData)

	return &CornelianModule{
		SettingsRepository: settingsRepository,
		PersonRepository:   personRepository,
		QuestionRepository: questionRepository,
		QuestionService:    questionService,
		MongoConf:          mongoConf,
		CacheConf:          cacheConf,
	}
}

// Clear clears the MongoDB configuration.
func (tf *CornelianModule) Clear() {
	tf.MongoConf.Clear()
}
