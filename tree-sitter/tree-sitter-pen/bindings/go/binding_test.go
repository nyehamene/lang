package tree_sitter_penny_test

import (
	"testing"

	tree_sitter "github.com/smacker/go-tree-sitter"
	"github.com/tree-sitter/tree-sitter-penny"
)

func TestCanLoadGrammar(t *testing.T) {
	language := tree_sitter.NewLanguage(tree_sitter_penny.Language())
	if language == nil {
		t.Errorf("Error loading Penny grammar")
	}
}
