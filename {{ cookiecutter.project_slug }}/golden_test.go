package main

import (
	"go/format"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Verify the generator generates a proper code.
func TestGolden(t *testing.T) {
	for _, tc := range []struct {
		title string
		want  string
	}{
		{
			title: "example",
			want:  `func Generated() {}`,
		},
	} {
		tc := tc
		t.Run(tc.title, func(t *testing.T) {
			g := NewGenerator()
			assert.Nil(t, g.Generate())
			got, err := format.Source(g.Bytes())
			assert.Nil(t, err)
			want, err := format.Source([]byte(tc.want))
			assert.Nil(t, err, "want")
			assert.Equal(t, strings.TrimRight(string(want), "\n"), strings.TrimRight(string(got), "\n"))
		})
	}
}
