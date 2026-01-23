package ui

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ValidateAmount parses a string input and ensures it is a positive number.
func ValidateAmount(input string) (float64, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return 0, errors.New("amount is required")
	}
	val, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return 0, errors.New("invalid amount format")
	}
	if val < 0 {
		return 0, errors.New("amount cannot be negative")
	}
	return val, nil
}

// ValidateDate checks if the input string matches the YYYY-MM-DD format.
func ValidateDate(input string) (time.Time, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return time.Time{}, errors.New("date is required")
	}
	date, err := time.Parse("2006-01-02", input)
	if err != nil {
		return time.Time{}, errors.New("date must be in YYYY-MM-DD format")
	}
	return date, nil
}

// ValidateRequired checks if a string is empty.
func ValidateRequired(input, fieldName string) error {
	if strings.TrimSpace(input) == "" {
		return fmt.Errorf("%s is required", fieldName)
	}
	return nil
}
