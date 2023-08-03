package librelinkup

import (
	"testing"
	"time"
)

func TestToTime(t *testing.T) {
	type TestCase struct {
		ts       string
		expected time.Time
	}
	tests := []TestCase{
		{
			ts:       "6/20/2023 10:01:57 PM",
			expected: time.Date(2023, time.June, 20, 22, 1, 57, 0, time.UTC),
		},
		{
			ts:       "12/20/2023 10:01:57 PM",
			expected: time.Date(2023, time.December, 20, 22, 1, 57, 0, time.UTC),
		},
		{
			ts:       "12/20/2023 10:01:57 AM",
			expected: time.Date(2023, time.December, 20, 10, 1, 57, 0, time.UTC),
		},
		{
			ts:       "8/2/2023 6:51:00 AM",
			expected: time.Date(2023, time.August, 2, 6, 51, 0, 0, time.UTC),
		},
	}

	for _, test := range tests {
		result, err := ToTime(test.ts)
		if err != nil {
			t.Error(err)
		}
		if result != test.expected {
			t.Errorf("expected %v, got %v", test.expected, result)
		}
	}
}
