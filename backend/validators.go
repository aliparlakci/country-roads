package main

import (
	"time"

	"github.com/go-playground/validator/v10"
)

var futureDate validator.Func = func(fl validator.FieldLevel) bool {
	date, ok := fl.Field().Interface().(time.Time)
	if ok {
		today := time.Now()
		if today.Year() < date.Year() {
			return true
		} else if today.Year() == date.Year() {
			if today.Month() < date.Month() {
				return true
			} else if today.Month() == date.Month() {
				if today.Day() <= date.Day() {
					return true
				}
			}
		}
	}
	return false
}

var validRideType validator.Func = func(fl validator.FieldLevel) bool {
	if field, ok := fl.Field().Interface().(string); ok {
		switch field {
		case "offer":
			return true
		case "request":
			return true
		case "taxi":
			return true
		default:
			return false
		}
	}

	return false
}

var validDirection validator.Func = func(fl validator.FieldLevel) bool {
	if field, ok := fl.Field().Interface().(string); ok {
		switch field {
		case "to_campus":
			return true
		case "from_campus":
			return true
		default:
			return false
		}
	}

	return false
}
