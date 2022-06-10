package main

import (
	"reflect"
	"testing"
)

func TestShuffle(t *testing.T) {
	type args struct {
		groups []Group[string, int]
		num    int
	}
	tests := []struct {
		name string
		args args
		want []ShuffleGroup[string, int]
	}{
		{
			"happy",
			args{
				groups: []Group[string, int]{
					&group[string, int]{"a", make([]int, 5)},
					&group[string, int]{"b", make([]int, 2)},
					&group[string, int]{"c", make([]int, 3)},
					&group[string, int]{"d", make([]int, 3)},
					&group[string, int]{"e", make([]int, 9)},
				},
				num: 3,
			},
			[]ShuffleGroup[string, int]{
				{&group[string, int]{"a", make([]int, 5)}, 0},
				{&group[string, int]{"b", make([]int, 2)}, 0},
				{&group[string, int]{"c", make([]int, 3)}, 1},
				{&group[string, int]{"d", make([]int, 3)}, 1},
				{&group[string, int]{"e", make([]int, 9)}, 2},
			},
		},
		{
			"happy2",
			args{
				groups: []Group[string, int]{
					&group[string, int]{"a", make([]int, 5)},
					&group[string, int]{"b", make([]int, 2)},
					&group[string, int]{"c", make([]int, 3)},
					&group[string, int]{"d", make([]int, 3)},
					&group[string, int]{"e", make([]int, 9)},
				},
				num: 2,
			},
			[]ShuffleGroup[string, int]{
				{&group[string, int]{"a", make([]int, 5)}, 0},
				{&group[string, int]{"b", make([]int, 2)}, 0},
				{&group[string, int]{"c", make([]int, 3)}, 0},
				{&group[string, int]{"d", make([]int, 3)}, 1},
				{&group[string, int]{"e", make([]int, 9)}, 1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Shuffle(tt.args.groups, tt.args.num); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\nExpected: %v,\ngot:      %v", tt.want, got)
			}
		})
	}
}
