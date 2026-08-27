package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/git-pkgs/brief"
)

const scanHelperRootEnv = "BRIEF_SCAN_HELPER_ROOT"
const diffHelperEnv = "BRIEF_DIFF_HELPER"

func TestScanDefaultsBoundRecursiveDetection(t *testing.T) {
	if root := os.Getenv(scanHelperRootEnv); root != "" {
		cmdScan([]string{"-json", root})
		return
	}

	root := t.TempDir()
	writeScanFixture(t, root, "pyproject.toml", "[project]\nname = \"example\"\nversion = \"1.0.0\"\n")
	writeScanFixture(t, root, "pipelines/example/Snakefile", "rule all:\n")
	deep := filepath.Join(
		"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "asv.conf.json",
	)
	writeScanFixture(t, root, deep, "{}\n")

	cmd := exec.Command(os.Args[0], "-test.run=^TestScanDefaultsBoundRecursiveDetection$")
	cmd.Env = append(os.Environ(), scanHelperRootEnv+"="+root, "PATH=")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("scan command failed: %v", err)
	}

	var report brief.Report
	if err := json.Unmarshal(out, &report); err != nil {
		t.Fatalf("parsing scan output: %v\n%s", err, out)
	}
	if !reportHasTool(&report, "build", "Snakemake") {
		t.Fatal("expected nested Snakemake project within the default depth")
	}
	if reportHasTool(&report, "test", "ASV") {
		t.Fatal("did not expect ASV beyond the default depth")
	}
	if !report.Stats.ScanTruncated {
		t.Fatal("expected scan output to report the depth truncation")
	}
}

func TestDiffAppliesScanOverrides(t *testing.T) {
	if os.Getenv(diffHelperEnv) != "" {
		cmdDiff([]string{"-scan-limit=-1", "HEAD"})
		return
	}

	dir := t.TempDir()
	writeScanFixture(t, dir, "go.mod", "module example.com/project\n\ngo 1.22\n")
	writeScanFixture(t, dir, "main.go", "package main\n")
	runGitFixture(t, dir, "init", "-q")
	runGitFixture(t, dir, "add", "go.mod", "main.go")
	runGitFixture(t, dir, "-c", "user.name=Test", "-c", "user.email=test@example.com", "commit", "-q", "-m", "initial")
	writeScanFixture(t, dir, "main.go", "package main\n\nfunc main() {}\n")

	cmd := exec.Command(os.Args[0], "-test.run=^TestDiffAppliesScanOverrides$")
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), diffHelperEnv+"=1", "PATH="+os.Getenv("PATH"))
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("diff command succeeded with a negative scan limit\n%s", out)
	}
	if !strings.Contains(string(out), "scan limit must not be negative") {
		t.Fatalf("diff command output = %q, want scan limit validation error", out)
	}
}

func runGitFixture(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

func writeScanFixture(t *testing.T, root, name, content string) {
	t.Helper()
	full := filepath.Join(root, name)
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func reportHasTool(report *brief.Report, category, name string) bool {
	for _, tool := range report.Tools[category] {
		if tool.Name == name {
			return true
		}
	}
	return false
}
