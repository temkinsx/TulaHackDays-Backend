package utils

import (
	"math"
	"strconv"
	"strings"
)

// ParseInt parses string to int with default value
func ParseInt(s string) int {
	if s == "" {
		return 0
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

// ParseFloat parses string to float64 with default value
func ParseFloat(s string) float64 {
	if s == "" {
		return 0.0
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.0
	}
	return f
}

// Contains checks if slice contains string
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// SanitizeString removes potentially harmful characters
func SanitizeString(s string) string {
	// Remove HTML tags and trim spaces
	s = strings.TrimSpace(s)

	// You can add more sanitization logic here
	// For example, remove scripts, special characters, etc.

	return s
}

// ValidateEmail performs basic email validation
func ValidateEmail(email string) bool {
	// Simple email validation
	// In production, you might want to use a more robust validation library
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

// TruncateString truncates string to specified length
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// CalculateDistance calculates distance between two points using Haversine formula
func CalculateDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const earthRadius = 6371 // km

	lat1Rad := lat1 * math.Pi / 180.0
	lat2Rad := lat2 * math.Pi / 180.0
	dLat := (lat2 - lat1) * math.Pi / 180.0
	dLng := (lng2 - lng1) * math.Pi / 180.0

	sinLat := math.Sin(dLat / 2)
	sinLng := math.Sin(dLng / 2)

	a := sinLat*sinLat + math.Cos(lat1Rad)*math.Cos(lat2Rad)*sinLng*sinLng
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}
