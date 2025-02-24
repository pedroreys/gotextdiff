package gotextdiff_test

import (
	_ "embed"
	"fmt"
	"github.com/pedroreys/gotextdiff/myers"
	"testing"

	diff "github.com/pedroreys/gotextdiff"
	"github.com/pedroreys/gotextdiff/difftest"
	"github.com/pedroreys/gotextdiff/span"
)

func TestApplyEdits(t *testing.T) {
	for _, tc := range difftest.TestCases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Helper()
			if got := diff.ApplyEdits(tc.In, tc.Edits); got != tc.Out {
				t.Errorf("ApplyEdits edits got %q, want %q", got, tc.Out)
			}
			if tc.LineEdits != nil {
				if got := diff.ApplyEdits(tc.In, tc.LineEdits); got != tc.Out {
					t.Errorf("ApplyEdits lineEdits got %q, want %q", got, tc.Out)
				}
			}
		})
	}
}

func TestLineEdits(t *testing.T) {
	for _, tc := range difftest.TestCases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Helper()
			// if line edits not specified, it is the same as edits
			edits := tc.LineEdits
			if edits == nil {
				edits = tc.Edits
			}
			if got := diff.LineEdits(tc.In, tc.Edits); diffEdits(got, edits) {
				t.Errorf("LineEdits got %q, want %q", got, edits)
			}
		})
	}
}

func TestUnified(t *testing.T) {
	for _, tc := range difftest.TestCases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Helper()
			unified := fmt.Sprint(diff.ToUnified(difftest.FileA, difftest.FileB, tc.In, tc.Edits))
			if unified != tc.Unified {
				t.Errorf("edits got diff:\n%v\nexpected:\n%v", unified, tc.Unified)
			}
			if tc.LineEdits != nil {
				unified := fmt.Sprint(diff.ToUnified(difftest.FileA, difftest.FileB, tc.In, tc.LineEdits))
				if unified != tc.Unified {
					t.Errorf("lineEdits got diff:\n%v\nexpected:\n%v", unified, tc.Unified)
				}
			}
		})
	}
}

//go:embed testdata/petstore.json
var FileA string

//go:embed testdata/petstore2.json
var FileB string

//go:embed testdata/expected_default.diff
var expectedDefault string

//go:embed testdata/expected_fullfile.diff
var expectedFullFile string

func TestUnified_WithContextLines(t *testing.T) {
	edits := myers.ComputeEdits(span.URIFromPath("petstore.json"), FileA, FileB)

	t.Run("should default to 3 line of context", func(t *testing.T) {
		actual := fmt.Sprint(diff.ToUnified("petstore.json", "petstore2.json", FileA, edits))
		if expectedDefault != actual {
			t.Errorf("unified:\n%q\nexpected:\n%q", actual, expectedDefault)
		}
	})

	t.Run("should expand context to the configured number of lines", func(t *testing.T) {
		actual := fmt.Sprint(diff.ToUnified("petstore.json", "petstore2.json", FileA, edits, diff.WithContextLines(180)))
		if expectedFullFile != actual {
			t.Errorf("unified:\n%q\nexpected:\n%q", actual, expectedFullFile)
		}
	})
}

func diffEdits(got, want []diff.TextEdit) bool {
	if len(got) != len(want) {
		return true
	}
	for i, w := range want {
		g := got[i]
		if span.Compare(w.Span, g.Span) != 0 {
			return true
		}
		if w.NewText != g.NewText {
			return true
		}
	}
	return false
}
