package access

import (
	"github.com/nefarius/cornelian/underlying/app/conf"
	"github.com/nefarius/cornelian/underlying/app/repository"
	"github.com/nefarius/cornelian/underlying/app/service"
	"github.com/nefarius/cornelian/underlying/app/store"
)

type CornelianModule struct {
	PersonRepository   *repository.PersonRepository
	QuestionRepository *repository.QuestionRepository
	QuestionService    *service.QuestionService
	MongoConf          *conf.MongoConf
	CacheConf          *conf.CacheConf
}

func NewCornelianModule() *CornelianModule {
	var mongoConf = conf.NewMongoConf()
	var cacheConf = conf.NewCacheConf()
	var questionRepository = &repository.QuestionRepository{Conf: *mongoConf}
	var personRepository = &repository.PersonRepository{Conf: *mongoConf}
	var quizRepository = &repository.QuizRepository{Conf: *mongoConf}
	var questionService = service.NewQuestionService(cacheConf, questionRepository, quizRepository)
	var defaultData = store.GetDefaultQuizData()
	questionRepository.InsertQuiz(defaultData)
	return &CornelianModule{PersonRepository: personRepository, QuestionRepository: questionRepository, QuestionService: questionService, MongoConf: mongoConf, CacheConf: cacheConf}
}

func (tf *CornelianModule) Clear() {
	tf.MongoConf.Clear()
}
