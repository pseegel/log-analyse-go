package main

import (
	"testing"
)

func TestParseLine(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Result
		wantErr bool
		// erwartete Werte
	}{
		// Testfälle hier
		{
			name:    "valid line",
			input:   "2024-01-15 10:23:45 GET /api/users 200 145",
			want:    Result{Line: "2024-01-15 10:23:45 GET /api/users 200 145", Status: 200},
			wantErr: false,
		},
		{
			name:    "malformed: more than 6 fields",
			input:   "2024-01-15 10:23:45 GET /api/users 200 145 extra_field",
			want:    Result{Line: "2024-01-15 10:23:45 GET /api/users 200 145 extra_field", Status: 0},
			wantErr: true,
		},
		{
			name:    "malformed: less than 6 fields",
			input:   "2024-01-15 10:23:45 GET /api/users 200",
			want:    Result{Line: "2024-01-15 10:23:45 GET /api/users 200", Status: 0},
			wantErr: true,
		},
		{
			name:    "malformed: status not an integer",
			input:   "2024-01-15 10:23:45 GET /api/users OK 145",
			want:    Result{Line: "2024-01-15 10:23:45 GET /api/users OK 145", Status: 0},
			wantErr: true,
		},
		{
			name:    "malformed: empty line",
			input:   "",
			want:    Result{Line: "", Status: 0},
			wantErr: true,
		},
		{
			name:    "valid: extra whitespace between fields",
			input:   "2024-01-15    10:23:45  GET /api/users    200 145",
			want:    Result{Line: "2024-01-15    10:23:45  GET /api/users    200 145", Status: 200},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Aufruf der zu testenden Funktion + Assertions
			got, err := parseLine(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseLine() = %v, want %v", got, tt.want)
			}
		})
	}
}
