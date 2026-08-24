package detect

import "testing"

func TestResearchToolsProject(t *testing.T) {
	report := runOn(t, "../testdata/research-tools-project")

	assertToolDetected(t, report, "test", "nf-test")
	assertToolDetected(t, report, "build", "nf-core")
	assertToolDetected(t, report, "docs", "MultiQC")
	assertToolDetected(t, report, "infrastructure", "Dockstore")
	assertToolDetected(t, report, "infrastructure", "DVC")
	assertToolDetected(t, report, "build", "cibuildwheel")
	assertToolDetected(t, report, "docs", "MyST-Parser")
	assertToolDetected(t, report, "test", "BenchmarkTools.jl")
	assertToolDetected(t, report, "lint", "lintr")
	assertToolDetected(t, report, "format", "styler")
	assertToolDetected(t, report, "build", "targets")
	assertToolDetected(t, report, "test", "vdiffr")
	assertToolDetected(t, report, "test", "tinytest")
}

func TestResearchToolSignals(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string
		category string
		tool     string
	}{
		{
			name: "nf-test test file",
			files: map[string]string{
				"nextflow.config":                   "nextflow.enable.dsl=2\n",
				"tests/modules/fastqc/main.nf.test": "nextflow_process { script \"main.nf\" }\n",
			},
			category: "test",
			tool:     "nf-test",
		},
		{
			name: "nf-core metadata",
			files: map[string]string{
				".nf-core.yml": "repository_type: pipeline\n",
			},
			category: "build",
			tool:     "nf-core",
		},
		{
			name: "nested MultiQC config",
			files: map[string]string{
				"assets/multiqc_config.yml": "title: Example analysis\n",
			},
			category: "docs",
			tool:     "MultiQC",
		},
		{
			name: "Dockstore GitHub descriptor",
			files: map[string]string{
				".github/.dockstore.yml": "version: 1.2\nworkflows:\n  - subclass: NFL\n    primaryDescriptorPath: /nextflow.config\n",
			},
			category: "infrastructure",
			tool:     "Dockstore",
		},
		{
			name: "DVC project config",
			files: map[string]string{
				".dvc/config": "[core]\n    remote = storage\n",
			},
			category: "infrastructure",
			tool:     "DVC",
		},
		{
			name: "cibuildwheel pyproject config",
			files: map[string]string{
				"pyproject.toml": "[project]\nname = \"native-example\"\nversion = \"1.0.0\"\n\n[tool.cibuildwheel]\ntest-command = \"pytest {project}/tests\"\n",
			},
			category: "build",
			tool:     "cibuildwheel",
		},
		{
			name: "MyST-Parser dependency",
			files: map[string]string{
				"pyproject.toml": "[project]\nname = \"research-docs\"\nversion = \"1.0.0\"\ndependencies = [\"myst-parser>=4\"]\n",
			},
			category: "docs",
			tool:     "MyST-Parser",
		},
		{
			name: "BenchmarkTools Julia dependency",
			files: map[string]string{
				"Project.toml": "name = \"NumericalModel\"\nuuid = \"12345678-1234-1234-1234-123456789abc\"\nversion = \"0.1.0\"\n\n[deps]\nBenchmarkTools = \"6e4b80f9-dd63-53aa-95a3-0cdb28fa8baf\"\n",
			},
			category: "test",
			tool:     "BenchmarkTools.jl",
		},
		{
			name: "lintr R dependency",
			files: map[string]string{
				"DESCRIPTION": "Package: modeltools\nVersion: 1.0.0\nSuggests: lintr\n",
			},
			category: "lint",
			tool:     "lintr",
		},
		{
			name: "styler R dependency",
			files: map[string]string{
				"DESCRIPTION": "Package: modeltools\nVersion: 1.0.0\nSuggests: styler\n",
			},
			category: "format",
			tool:     "styler",
		},
		{
			name: "targets pipeline",
			files: map[string]string{
				"_targets.R": "library(targets)\nlist(tar_target(result, run_model()))\n",
			},
			category: "build",
			tool:     "targets",
		},
		{
			name: "vdiffr R dependency",
			files: map[string]string{
				"DESCRIPTION": "Package: plottools\nVersion: 1.0.0\nSuggests: vdiffr\n",
			},
			category: "test",
			tool:     "vdiffr",
		},
		{
			name: "tinytest package runner",
			files: map[string]string{
				"DESCRIPTION":      "Package: modeltools\nVersion: 1.0.0\nSuggests: tinytest\n",
				"tests/tinytest.R": "if (requireNamespace(\"tinytest\")) tinytest::test_package(\"modeltools\")\n",
			},
			category: "test",
			tool:     "tinytest",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			for path, content := range tt.files {
				writeProjectFile(t, dir, path, content)
			}

			report := runOn(t, dir)
			assertToolDetected(t, report, tt.category, tt.tool)
		})
	}
}
