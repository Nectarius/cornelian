package service

import (
	"fmt"
	"slices"
	"sort"

	"github.com/nefarius/cornelian/underlying/app"
	"github.com/nefarius/cornelian/underlying/app/repository"
)

type QuestionService struct {
	QuestionRepository repository.QuestionRepository
}

func NewQuestionService(questionRepository *repository.QuestionRepository) *QuestionService {
	return &QuestionService{QuestionRepository: *questionRepository}
}

func (r *QuestionService) GetQuiz() app.Quiz {
	var questionRepository = r.QuestionRepository

	return questionRepository.GetQuiz()
}

func (r *QuestionService) AllForAuthorInStatus(email string, status app.Status) []app.Question {
	var quiz = r.GetQuiz()
	out := make([]app.Question, 0)
	for _, q := range quiz.Questions {
		if slices.Contains(q.Talk.AssignedTo, email) && q.Status == status {
			out = append(out, q)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out
}

func (r *QuestionService) GetQuestion(id string) (app.Question, error) {
	var quiz = r.GetQuiz()
	for _, q := range quiz.Questions {
		if q.ID == id {
			return q, nil
		}
	}
	return app.Question{}, fmt.Errorf("question identified by %v not found", id)
}

func (r *QuestionService) AllQuestions() []app.Question {
	out := r.GetQuiz().Questions
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
		//return out[i].CreatedAt.After(out[j].CreatedAt) && out[i].Status == app.StatusOpen
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

func (r *QuestionService) AllForAssignedTo(email string) []app.Question {
	var quiz = r.GetQuiz()
	out := make([]app.Question, 0)
	for _, q := range quiz.Questions {
		if slices.Contains(q.Talk.AssignedTo, email) {
			out = append(out, q)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out
}
