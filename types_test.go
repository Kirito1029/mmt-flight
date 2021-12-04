package main

import (
	"testing"
	"time"
)

func TestTimeBetween(t *testing.T) {
	type args struct {
		start string
		end   string
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		// TODO: Add test cases.
		{
			name: "T",
			args: args{
				start: "2300",
				end:   "0100",
			},
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TimeBetween(tt.args.start, tt.args.end); got != tt.want {
				t.Errorf("TimeBetween() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddTime(t *testing.T) {
	type args struct {
		currentTime string
		addTime     time.Duration
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "T",
			args: args{
				// currentTime: "",
				addTime: time.Minute * (120 + time.Duration(65)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AddTime(tt.args.currentTime, tt.args.addTime); got != tt.want {
				t.Errorf("AddTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
