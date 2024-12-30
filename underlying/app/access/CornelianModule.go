package access

import (
	"github.com/nefarius/cornelian/underlying/app/conf"
	"github.com/nefarius/cornelian/underlying/app/repository"
	"github.com/nefarius/cornelian/underlying/app/service"
)

type CornelianModule struct {
	QuestionRepository *repository.QuestionRepository
	QuestionService    *service.QuestionService
	MongoConf          *conf.MongoConf
	CacheConf          *conf.CacheConf
}

func NewCornelianModule() *CornelianModule {
	var mongoConf = conf.NewMongoConf()
	var cacheConf = conf.NewCacheConf()
	var repository = &repository.QuestionRepository{Conf: *mongoConf}
	var questionService = service.NewQuestionService(repository)
	// var defaultData = store.GetDefaultQuizData()
	// panelViewRepository.InsertQuiz(defaultData)
	return &CornelianModule{QuestionRepository: repository, QuestionService: questionService, MongoConf: mongoConf, CacheConf: cacheConf}
}

func (tf *CornelianModule) Clear() {
	tf.MongoConf.Clear()
}
