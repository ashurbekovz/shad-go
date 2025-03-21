package varfmt

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormat(t *testing.T) {
	s := make([]any, 1002)
	s[10] = 1
	s[100] = 2
	s[1000] = 3

	for _, tc := range []struct {
		format string
		args   []any
		result string
	}{
		{
			format: "{}",
			args:   []any{0},
			result: "0",
		},
		{
			format: "{0} {0}",
			args:   []any{1},
			result: "1 1",
		},
		{
			format: "{1} {5}",
			args:   []any{0, 1, 2, 3, 4, 5, 6},
			result: "1 5",
		},
		{
			format: "{} {} {} {} {}",
			args:   []any{0, 1, 2, 3, 4},
			result: "0 1 2 3 4",
		},
		{
			format: "{} {0} {0} {0} {}",
			args:   []any{0, 1, 2, 3, 4},
			result: "0 0 0 0 4",
		},
		{
			format: "Hello, {2}",
			args:   []any{0, 1, "World"},
			result: "Hello, World",
		},
		{
			format: "He{2}o",
			args:   []any{0, 1, "ll"},
			result: "Hello",
		},
		{
			format: "{10}",
			args:   s[:11],
			result: "1",
		},
		{
			format: "{100}",
			args:   s[:101],
			result: "2",
		},
		{
			format: "{1000}",
			args:   s[:1001],
			result: "3",
		},
	} {
		t.Run(tc.result, func(t *testing.T) {
			require.Equal(t, tc.result, Sprintf(tc.format, tc.args...))
		})
	}
}

func BenchmarkFormat(b *testing.B) {
	for _, tc := range []struct {
		name   string
		format string
		args   []any
	}{
		{
			name:   "small int",
			format: "{}",
			args:   []any{42},
		},
		{
			name:   "small string",
			format: "{} {}",
			args:   []any{"Hello", "World"},
		},
		{
			name:   "big",
			format: strings.Repeat("{0}{1}", 1000),
			args:   []any{42, 43},
		},
	} {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = Sprintf(tc.format, tc.args...)
			}
		})
	}
}

func BenchmarkSprintf(b *testing.B) {
	for _, tc := range []struct {
		name   string
		format string
		args   []any
	}{
		{
			name:   "small",
			format: "%d",
			args:   []any{42},
		},
		{
			name:   "small string",
			format: "%v %v",
			args:   []any{"Hello", "World"},
		}, {
			name:   "big",
			format: strings.Repeat("%[0]v%[1]v", 1000),
			args:   []any{42, 43},
		},
	} {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()

			for b.Loop() {
				_ = fmt.Sprintf(tc.format, tc.args...)
			}
		})
	}
}
