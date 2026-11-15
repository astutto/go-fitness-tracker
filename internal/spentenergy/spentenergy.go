// Package spentenergy provides calculations of distance, mean speed, and energy expenditure.
package spentenergy

import (
	"errors"
	"time"
)

// Основные константы, необходимые для расчётов.
const (
	mInKm                      = 1000 // количество метров в километре
	minInH                     = 60   // количество минут в часе
	stepLengthCoefficient      = 0.45 // коэффициент длины шага относительно роста
	walkingCaloriesCoefficient = 0.5  // коэффициент калорий при ходьбе
)

// Distance возвращает пройденную дистанцию в километрах.
func Distance(steps int, height float64) float64 {
	if steps <= 0 || height <= 0 {
		return 0
	}
	return float64(steps) * height * stepLengthCoefficient / mInKm
}

// MeanSpeed возвращает среднюю скорость в километрах в час.
func MeanSpeed(steps int, height float64, duration time.Duration) float64 {
	if steps <= 0 || height <= 0 || duration <= 0 {
		return 0
	}
	return Distance(steps, height) / duration.Hours()
}

// RunningSpentCalories вычисляет количество калорий, потраченных при беге.
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if err := validate(steps, weight, height, duration); err != nil {
		return 0, err
	}
	calories := weight * MeanSpeed(steps, height, duration) * duration.Minutes() / minInH
	return calories, nil
}

// WalkingSpentCalories вычисляет количество калорий, потраченных при ходьбе.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	calories, err := RunningSpentCalories(steps, weight, height, duration)
	if err != nil {
		return 0, err
	}
	return calories * walkingCaloriesCoefficient, nil
}

// validate проверяет данные, необходимые для расчёта калорий.
func validate(steps int, weight, height float64, duration time.Duration) error {
	switch {
	case steps <= 0:
		return errors.New("количество шагов должно быть больше нуля")
	case weight <= 0:
		return errors.New("вес должен быть больше нуля")
	case height <= 0:
		return errors.New("рост должен быть больше нуля")
	case duration <= 0:
		return errors.New("продолжительность должна быть больше нуля")
	default:
		return nil
	}
}
