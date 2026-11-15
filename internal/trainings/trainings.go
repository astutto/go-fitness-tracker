// Package trainings parses training records and calculates their summary.
package trainings

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/personaldata"
	"github.com/Yandex-Practicum/tracker/internal/spentenergy"
)

// Training содержит сведения о тренировке и персональные данные.
type Training struct {
	Steps        int
	TrainingType string
	Duration     time.Duration
	personaldata.Personal
}

// Parse разбирает запись вида "3456,Ходьба,3h00m".
func (t *Training) Parse(datastring string) (err error) {
	parts := strings.Split(datastring, ",")
	if len(parts) != 3 {
		return errors.New("неверный формат данных тренировки")
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("неверное количество шагов: %w", err)
	}
	if steps <= 0 {
		return errors.New("количество шагов должно быть больше нуля")
	}

	// Отрицательное слагаемое недопустимо даже при положительной сумме.
	if strings.Contains(parts[2], "-") {
		return errors.New("продолжительность не должна содержать отрицательных значений")
	}
	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return fmt.Errorf("неверная продолжительность тренировки: %w", err)
	}
	if duration <= 0 {
		return errors.New("продолжительность тренировки должна быть больше нуля")
	}

	t.Steps = steps
	t.TrainingType = parts[1]
	t.Duration = duration
	return nil
}

// ActionInfo возвращает информацию о тренировке или ошибку расчёта.
func (t Training) ActionInfo() (string, error) {
	distance := spentenergy.Distance(t.Steps, t.Height)
	speed := spentenergy.MeanSpeed(t.Steps, t.Height, t.Duration)

	var calories float64
	var err error
	switch t.TrainingType {
	case "Бег":
		calories, err = spentenergy.RunningSpentCalories(t.Steps, t.Weight, t.Height, t.Duration)
	case "Ходьба":
		calories, err = spentenergy.WalkingSpentCalories(t.Steps, t.Weight, t.Height, t.Duration)
	default:
		return "", errors.New("неизвестный тип тренировки")
	}
	if err != nil {
		return "", err
	}

	info := fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		t.TrainingType, t.Duration.Hours(), distance, speed, calories,
	)
	return info, nil
}
