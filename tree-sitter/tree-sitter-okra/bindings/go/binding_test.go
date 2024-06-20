package tree_sitter_okra_test

import (
	"testing"

	tree_sitter "github.com/smacker/go-tree-sitter"
	"github.com/tree-sitter/tree-sitter-okra"
)

func TestCanLoadGrammar(t *testing.T) {
	language := tree_sitter.NewLanguage(tree_sitter_okra.Language())
	if language == nil {
		t.Errorf("Error loading Okra grammar")
	}
}
