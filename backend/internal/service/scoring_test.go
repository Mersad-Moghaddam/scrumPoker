package service

import (
	"testing"

	"scrum-poker/backend/internal/models"
)

func TestCalculateAverageLabel(t *testing.T) {
	half := 0.5
	three := 3.0
	label, avg := CalculateAverageLabel([]models.Vote{
		{Value: "0.5", NumericValue: &half},
		{Value: "3", NumericValue: &three},
		{Value: "Coffee"},
		{Value: "?"},
	})

	if label == nil || *label != "1.75" {
		t.Fatalf("expected average label 1.75, got %#v", label)
	}
	if avg == nil || *avg != 1.75 {
		t.Fatalf("expected numeric average 1.75, got %#v", avg)
	}
}

func TestCalculateAverageLabelWithoutNumericVotes(t *testing.T) {
	label, avg := CalculateAverageLabel([]models.Vote{
		{Value: "Coffee"},
		{Value: "?"},
		{Value: "∞"},
	})

	if label == nil || *label != "بدون برآورد عددی" {
		t.Fatalf("expected no numeric estimate label, got %#v", label)
	}
	if avg != nil {
		t.Fatalf("expected nil numeric average, got %#v", avg)
	}
}

func TestIsValidVoteValue(t *testing.T) {
	if !IsValidVoteValue("Coffee") {
		t.Fatal("Coffee should be valid")
	}
	if IsValidVoteValue("55") {
		t.Fatal("55 should be invalid")
	}
}

func TestShouldReveal(t *testing.T) {
	votes := []models.Vote{
		{MemberID: 1, Value: "3"},
		{MemberID: 2, Value: "5"},
	}
	if ShouldReveal([]int64{1, 2, 3}, votes) {
		t.Fatal("should not reveal before all active members vote")
	}
	if !ShouldReveal([]int64{1, 2}, votes) {
		t.Fatal("should reveal when all active members voted")
	}
}
