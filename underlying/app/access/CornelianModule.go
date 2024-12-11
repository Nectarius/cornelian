package access

import (
	"github.com/nefarius/cornelian/underlying/app/conf"

	"github.com/nefarius/cornelian/underlying/app/repository"
)

type CornelianModule struct {
	QuestionRepository repository.QuestionRepository
	MongoConf          *conf.MongoConf
}

func NewCornelianModule() *CornelianModule {
	var mongoConf = conf.NewMongoConf()
	var panelViewRepository = repository.QuestionRepository{Conf: *mongoConf}

	//var defaultData = store.GetDefaultQuizData()
	//panelViewRepository.InsertQuiz(defaultData)
	return &CornelianModule{QuestionRepository: panelViewRepository, MongoConf: mongoConf}
}

func (tf *CornelianModule) Clear() {
	tf.MongoConf.Clear()
}
