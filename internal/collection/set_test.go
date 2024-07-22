package collection

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tmknom/cross/internal/testlib"
)

func TestTreeSet_AddAndSlice(t *testing.T) {
	cases := []struct {
		items    []string
		expected []string
	}{
		{
			items:    []string{"foo", "bar", "baz"},
			expected: []string{"bar", "baz", "foo"},
		},
		{
			items:    []string{"foo", "bar", "foo"},
			expected: []string{"bar", "foo"},
		},
	}

	for _, tc := range cases {
		sut := NewTreeSet()
		for _, item := range tc.items {
			sut.Add(item)
		}

		actual := sut.Slice()
		if diff := cmp.Diff(tc.expected, actual); diff != "" {
			t.Errorf(testlib.FormatWithoutInput(sut, tc.expected, actual))
		}
	}
}
