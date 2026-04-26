package service

import (
	"fmt"
	"slices"
	"strconv"

	"scrum-poker/backend/internal/models"
)

func IsValidVoteValue(value string) bool {
	return slices.Contains(models.ValidVoteValues, value)
}

func NumericVoteValue(value string) (*float64, bool) {
	switch value {
	case "?", "∞", "Coffee":
		return nil, false
	default:
		parsed, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, false
		}
		return &parsed, true
	}
}

func CalculateAverageLabel(votes []models.Vote) (*string, *float64) {
	var total float64
	var count float64
	for _, vote := range votes {
		if vote.NumericValue != nil {
			total += *vote.NumericValue
			count++
		}
	}
	if count == 0 {
		label := "بدون برآورد عددی"
		return &label, nil
	}
	avg := total / count
	label := trimAverage(avg)
	return &label, &avg
}

func trimAverage(value float64) string {
	if value == float64(int64(value)) {
		return fmt.Sprintf("%.0f", value)
	}
	if value*10 == float64(int64(value*10)) {
		return fmt.Sprintf("%.1f", value)
	}
	return fmt.Sprintf("%.2f", value)
}

func ShouldReveal(activeMemberIDs []int64, votes []models.Vote) bool {
	if len(activeMemberIDs) == 0 {
		return len(votes) > 0
	}
	voted := make(map[int64]struct{}, len(votes))
	for _, vote := range votes {
		voted[vote.MemberID] = struct{}{}
	}
	for _, memberID := range activeMemberIDs {
		if _, ok := voted[memberID]; !ok {
			return false
		}
	}
	return true
}
