package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_FilterStringSlice(t *testing.T) {
	var tests = []struct {
		in  []string
		out []string
	}{
		{
			[]string{"1", "3", "2", "5"}, []string{"1", "2", "3", "5"},
		},
		{
			[]string{"1", "1", "2", "2"}, []string{"1", "2"},
		},
		{
			[]string{"foo", "bar", "", "ozz", "aaa"}, []string{"aaa", "bar", "foo", "ozz"},
		},
		{
			[]string{"foo", "foo", "", "", "foo"}, []string{"foo"},
		},
	}

	r := require.New(t)

	for i, test := range tests {
		got := FilterStringSlice(test.in)

		r.Equalf(
			got, test.out,
			"FilterStringSlice(\"%v\"): tc %d, Expected \"%v\", Got \"%v\"",
			test.in, i, test.out, got,
		)
	}
}
