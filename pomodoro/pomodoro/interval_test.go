package pomodoro_test

import (
	"go-cmd-book/pomodoro/pomodoro"
	"go-cmd-book/pomodoro/repository"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		var repo pomodoro.Repository
		config := pomodoro.NewConfig(
			repo,
			0*time.Minute,
			0*time.Minute,
			0*time.Minute,
		)

		assert.Equal(t, 25*time.Minute, config.PomodoroDuration)
		assert.Equal(t, 5*time.Minute, config.ShortBreakDuration)
		assert.Equal(t, 15*time.Minute, config.LongBreakDuration)
	})
}

func getRepo(t *testing.T) (pomodoro.Repository, func()) {
	t.Helper()
	return repository.NewInMemoryRepo(), func() {}
}
func TestGetInterval(t *testing.T) {
	repo, cleanup := getRepo(t)
	defer cleanup()

	const duration = 1 * time.Millisecond
	config := pomodoro.NewConfig(repo, 3*duration, duration, 2*duration)

	i, err := pomodoro.GetInterval(config)
	assert.NoError(t, err)
	assert.Equal(t, pomodoro.CategoryPomodoro, i.Category)

}
