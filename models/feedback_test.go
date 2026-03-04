package models

import "testing"

func TestFeedback_Validate(t *testing.T) {
	tests := []struct {
		name    string
		fb      Feedback
		wantErr bool
	}{
		{
			name:    "valid feedback",
			fb:      Feedback{Rating: 3, Page: "/dashboard"},
			wantErr: false,
		},
		{
			name:    "valid feedback with comment",
			fb:      Feedback{Rating: 4, Comment: "Great app!", Page: "/dashboard"},
			wantErr: false,
		},
		{
			name:    "rating too low",
			fb:      Feedback{Rating: 0, Page: "/dashboard"},
			wantErr: true,
		},
		{
			name:    "rating too high",
			fb:      Feedback{Rating: 5, Page: "/dashboard"},
			wantErr: true,
		},
		{
			name:    "missing page",
			fb:      Feedback{Rating: 3},
			wantErr: true,
		},
		{
			name:    "min rating",
			fb:      Feedback{Rating: 1, Page: "/home"},
			wantErr: false,
		},
		{
			name:    "max rating",
			fb:      Feedback{Rating: 4, Page: "/home"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fb.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
