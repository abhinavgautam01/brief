package main

import (
	"strings"
	"testing"

	"github.com/git-pkgs/brief/kb"
)

func TestWriteToolsReadmeOmitsKnowledgeBaseTotals(t *testing.T) {
	var out strings.Builder
	writeToolsReadme(&out, &kb.KnowledgeBase{})

	const introduction = "Language ecosystems and development tools across multiple categories."
	if !strings.Contains(out.String(), introduction) {
		t.Errorf("README output does not contain count-free introduction %q", introduction)
	}
	if strings.Contains(out.String(), "language ecosystems with") {
		t.Errorf("README output contains numeric knowledge-base totals: %q", out.String())
	}
}
