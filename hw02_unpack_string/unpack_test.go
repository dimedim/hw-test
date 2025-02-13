package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},
		{input: "🙃0", expected: ""},
		{input: "aaф0b", expected: "aab"},
		{input: "abcd", expected: "abcd"},

		// uncomment if task with asterisk completed
		{input: `qwe\4\5`, expected: `qwe45`},
		{input: `qwe\45`, expected: `qwe44444`},
		{input: `qwe\\5`, expected: `qwe\\\\\`},
		{input: `qwe\\\3`, expected: `qwe\3`},
		// more tests
		{input: "d\n5abc", expected: "d\n\n\n\n\nabc"},
		{input: "a1b1c1", expected: "abc"},
		{input: "a0b0c0", expected: ""},
		{input: "🤖3💥2", expected: "🤖🤖🤖💥💥"},
		{input: "😃😃2", expected: "😃😃😃"},
		{input: "😒4👌3", expected: "😒😒😒😒👌👌👌"},
		{input: "🔥0👇5", expected: "👇👇👇👇👇"},

		{input: `qwe\\\4`, expected: `qwe\4`},
		{input: `qwe\\\\4`, expected: `qwe\\\\\`},
		{input: `qwe\\0`, expected: `qwe`},
		{input: `\\`, expected: `\`},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{
		"3abc",
		"45",
		"aaa10b",
		// more tests
		`qw\ne`,
		`1`,
		`10`,
		`a2b3c10`,
		`🙃🙃10`,
		`a\`,
		`a\z`,
		`qwe\\\`,
		`qwe\ `,
	}
	for _, tc := range invalidStrings {
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}
