package date

import (
	"testing"
)

func TestFormatDateString(t *testing.T) {
	tests := []struct {
		name       string
		dateString string
		want       string
	}{
		{
			name:       "Format successfully",
			dateString: "2025-03-27T20:57:58+01:00",
			want:       "27.03.2025 20:57:58",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatDateString(tt.dateString)
			if got != tt.want {
				t.Errorf("FormatDateString() = %v, want %v", got, tt.want)
			}
		})
	}
}
