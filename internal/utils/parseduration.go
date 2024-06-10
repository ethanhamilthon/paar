package utils

import (
	"regexp"
	"strconv"
	"time"
)

func ParseDuration(durationStr string) (time.Duration, error) {
	var totalDuration time.Duration
	// Регулярное выражение для извлечения чисел и их единиц измерения (s, m, h, d)
	re := regexp.MustCompile(`(\d+)([smhd])`)
	matches := re.FindAllStringSubmatch(durationStr, -1)

	// Парсинг найденных подстрок и суммирование их в общую продолжительность
	for _, match := range matches {
		value, err := strconv.Atoi(match[1])
		if err != nil {
			return 0, err
		}

		switch match[2] {
		case "s":
			totalDuration += time.Duration(value) * time.Second
		case "m":
			totalDuration += time.Duration(value) * time.Minute
		case "h":
			totalDuration += time.Duration(value) * time.Hour
		case "d":
			totalDuration += time.Duration(value*24) * time.Hour
		}
	}

	return totalDuration, nil
}