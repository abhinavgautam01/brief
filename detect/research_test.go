package detect

import (
	"slices"
	"testing"

	"github.com/git-pkgs/brief"
)

func TestResearchToolDetectors(t *testing.T) {
	dir := t.TempDir()
	files := map[string]string{
		"analysis/notebook.ipynb":  "{}\n",
		"report.qmd":               "---\ntitle: Analysis\n---\n",
		"workflow/Snakefile":       "rule all:\n    input: []\n",
		"nextflow.config":          "manifest { name = 'research-pipeline' }\n",
		"main.nf":                  "workflow { }\n",
		"pyproject.toml":           "[project]\nname = \"research-tools\"\nversion = \"0.1.0\"\n",
		"tox.ini":                  "[tox]\nenv_list = py312\n",
		"benchmarks/asv.conf.json": `{"repo": "."}`,
		"DESCRIPTION": `Package: researchtools
Version: 0.1.0
Suggests:
    covr,
    knitr,
    pkgdown,
    rmarkdown,
    testthat
Config/roxygen2/version: 7.3.3
Config/testthat/edition: 3
VignetteBuilder: knitr
`,
		"README.Rmd":                   "---\ntitle: Research tools\n---\n",
		"_pkgdown.yml":                 "url: https://example.com\n",
		"renv.lock":                    `{"R":{"Version":"4.5.0"}}`,
		"Project.toml":                 "name = \"ResearchTools\"\nuuid = \"11111111-1111-1111-1111-111111111111\"\n",
		"src/ResearchTools.jl":         "module ResearchTools\nend\n",
		"docs/Project.toml":            "[deps]\nDocumenter = \"e30172f5-a6a5-5a46-863b-614d45cd2de4\"\n",
		"docs/make.jl":                 "using Documenter\n",
		".JuliaFormatter.toml":         "margin = 92\n",
		".github/workflows/format.yml": "name: format\non: push\njobs:\n  format:\n    runs-on: ubuntu-latest\n    steps:\n      - uses: fredrikekre/runic-action@v1\n",
		"model.f90":                    "program model\nend program model\n",
		"fortitude.toml":               "line-length = 100\n",
	}
	for path, content := range files {
		writeProjectFile(t, dir, path, content)
	}

	report := runOn(t, dir)
	want := []struct {
		category string
		name     string
	}{
		{category: "environment", name: "Jupyter"},
		{category: "docs", name: "Quarto"},
		{category: "build", name: "Snakemake"},
		{category: "build", name: "Nextflow"},
		{category: "test", name: "tox"},
		{category: "test", name: "ASV"},
		{category: "test", name: "testthat"},
		{category: "coverage", name: "covr"},
		{category: "docs", name: "pkgdown"},
		{category: "docs", name: "roxygen2"},
		{category: "docs", name: "knitr"},
		{category: "docs", name: "R Markdown"},
		{category: "docs", name: "Documenter.jl"},
		{category: "format", name: "JuliaFormatter"},
		{category: "format", name: "Runic"},
		{category: "lint", name: "Fortitude"},
	}
	for _, item := range want {
		assertToolDetected(t, report, item.category, item.name)
	}
	if !slices.ContainsFunc(report.PackageManagers, func(d brief.Detection) bool {
		return d.Name == "renv"
	}) {
		t.Error("expected renv package manager")
	}
}
