package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_GetShortFQDN(t *testing.T) {
	var tests = []struct {
		in  string
		out string
	}{
		{
			"foo.qwerty", "foo.qwerty",
		},
		{
			"foo.bar", "foo",
		},
		{
			"foo.com.ua", "foo",
		},
		{
			"foo42.domain.net", "foo42",
		},
		{
			"foo.domain.com.ua", "foo",
		},
		{
			"foo.bar.domain.com.ua", "foo.bar",
		},
	}

	r := require.New(t)

	for i, test := range tests {
		got := GetShortFQDN(test.in)

		r.Equalf(
			got, test.out,
			"GetShortFQDN(\"%s\"): tc %d, Expected \"%s\", Got \"%s\"",
			test.in, i, test.out, got,
		)
	}
}

func Test_GetFQDNWithOutPublicSuffix(t *testing.T) {
	var tests = []struct {
		in  string
		out string
	}{
		{
			"foo.qwerty", "foo.qwerty",
		},
		{
			"foo.bar", "foo",
		},
		{
			"foo.com.ua", "foo",
		},
		{
			"foo42.domain.net", "foo42.domain",
		},
		{
			"foo.domain.com.ua", "foo.domain",
		},
		{
			"foo.bar.domain.com.ua", "foo.bar.domain",
		},
	}

	r := require.New(t)

	for i, test := range tests {
		got := GetFQDNWithOutPublicSuffix(test.in)

		r.Equalf(
			got, test.out,
			"GetFQDNWithOutPublicSuffix(\"%s\"): tc %d, Expected \"%s\", Got \"%s\"",
			test.in, i, test.out, got,
		)
	}
}
