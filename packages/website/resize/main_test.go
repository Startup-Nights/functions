package main

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestResize(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
	}{
		{
			name:   "resize a png file",
			input:  "dreiseenstafette.png",
			output: "dreiseenstafette_scaled.png",
		},
		{
			name:   "resize a jpg file",
			input:  "dreiseenstafette.jpg",
			output: "dreiseenstafette_jpg_scaled.png",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := os.ReadFile(tt.input)
			if err != nil {
				t.Fatal(err)
			}

			result, err := resize(bytes.NewBuffer(data), typeFromFilename(tt.input), 600, 300)
			if err != nil {
				t.Fatal(err)
			}

			data, err = io.ReadAll(result)
			if err != nil {
				t.Fatal(err)
			}

			if err := os.WriteFile(tt.output, data, 0644); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestAdjustCanvas(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
	}{
		{
			name:   "adjust a png file",
			input:  "dreiseenstafette_scaled.png",
			output: "dreiseenstafette_adjusted.png",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := os.ReadFile(tt.input)
			if err != nil {
				t.Fatal(err)
			}

			result, err := adjustCanvas(bytes.NewBuffer(data), 600, 300)
			if err != nil {
				t.Fatal(err)
			}

			data, err = io.ReadAll(result)
			if err != nil {
				t.Fatal(err)
			}

			if err := os.WriteFile(tt.output, data, 0644); err != nil {
				t.Fatal(err)
			}
		})
	}
}
