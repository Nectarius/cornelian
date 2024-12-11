package views

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/nefarius/cornelian/internal/app"
	"github.com/stretchr/testify/assert"
)

func TestAnswersTemplateRendersNameCorrectly(t *testing.T) {
	answers := []app.Answer{
		{"1", "Answer text 1", "kalle.karlsson@domain.com", time.Now()},
		{"2", "Answer text 2", "hanna.hansson@domain.com", time.Now()},
	}
	buffer := new(strings.Builder)
	assert.NoError(t, Answers(answers).Render(context.Background(), buffer))

	assert.Contains(t, buffer.String(), "<div class=\"col-10\">Kalle Karlsson</div>")
	assert.Contains(t, buffer.String(), "<div class=\"col-10\">Hanna Hansson</div>")
}
