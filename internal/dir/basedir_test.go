package dir

import (
	"fmt"
	"os"
	"testing"

	"github.com/tmknom/cross/internal/testlib"
)

func TestBaseDir_Abs(t *testing.T) {
	currentDir, _ := os.Getwd()

	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    ".",
			expected: currentDir,
		},
		{
			input:    "testdata/terraform/../",
			expected: fmt.Sprintf("%s/%s", currentDir, "testdata"),
		},
		{
			input:    "/path/to/dir",
			expected: "/path/to/dir",
		},
		{
			input:    "../../internal/dir/foo/bar/baz/../",
			expected: fmt.Sprintf("%s/%s", currentDir, "foo/bar"),
		},
	}

	for _, tc := range cases {
		sut := NewBaseDir(tc.input)
		actual := sut.Abs()
		if actual != tc.expected {
			t.Errorf(testlib.Format(sut, tc.expected, actual, tc.input))
		}
	}
}
