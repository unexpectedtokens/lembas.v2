package types

import "time"

type MealplanEntry struct {
	Date time.Time
	Breakfast,
	Lunch,
	Dinner,
	Snacks []Recipe
}

const dateFormat = "02 Jan"

func FormatMealplanDate(date time.Time) string {
	return date.Format(dateFormat)
}

func ParseMealplanDate(date string) (time.Time, error) {
	return time.Parse(dateFormat, date)
}
