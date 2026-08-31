package detect

import (
	"slices"
	"testing"

	"github.com/git-pkgs/brief"
)

func TestMissingDevelopmentToolSignals(t *testing.T) {
	tests := []struct {
		name           string
		files          map[string]string
		category       string
		tool           string
		packageManager bool
		notTool        string
	}{
		{
			name:     "Babashka project",
			files:    map[string]string{"bb.edn": "{:tasks {test {:task (println \"test\")}}}\n"},
			category: "build",
			tool:     "Babashka",
		},
		{
			name:     "nested OpenTofu configuration",
			files:    map[string]string{"infra/main.tofu": "terraform { required_version = \">= 1.8\" }\n"},
			category: "infrastructure",
			tool:     "OpenTofu",
		},
		{
			name:           "Stack project",
			files:          map[string]string{"stack.yaml": "resolver: lts-23.18\npackages: [.]\n"},
			tool:           "Stack",
			packageManager: true,
		},
		{
			name:     "Chef Ruby metadata",
			files:    map[string]string{"metadata.rb": "name 'example'\nversion '1.0.0'\n"},
			category: "infrastructure",
			tool:     "Chef",
		},
		{
			name:     "Chef JSON metadata",
			files:    map[string]string{"metadata.json": "{\"name\":\"example\",\"version\":\"1.0.0\"}\n"},
			category: "infrastructure",
			tool:     "Chef",
		},
		{
			name:     "Chef Berksfile",
			files:    map[string]string{"Berksfile": "source 'https://supermarket.chef.io'\nmetadata\n"},
			category: "infrastructure",
			tool:     "Chef",
		},
		{
			name:     "BitBake recipe",
			files:    map[string]string{"recipes/example/example_1.0.bb": "SUMMARY = \"Example recipe\"\n"},
			category: "build",
			tool:     "BitBake",
		},
		{
			name:     "BitBake recipe append",
			files:    map[string]string{"recipes/example/example_1.0.bbappend": "FILESEXTRAPATHS:prepend := \"${THISDIR}/files:\"\n"},
			category: "build",
			tool:     "BitBake",
		},
		{
			name:     "Buck build file",
			files:    map[string]string{"app/BUCK": "cxx_binary(name = \"app\", srcs = [\"main.cc\"])\n"},
			category: "monorepo",
			tool:     "Buck",
		},
		{
			name:     "Buck metadata file",
			files:    map[string]string{"prelude/METADATA.bzl": "METADATA = {}\n"},
			category: "monorepo",
			tool:     "Buck",
		},
		{
			name:     "Devbox config",
			files:    map[string]string{"devbox.json": "{\"packages\": [\"ripgrep@latest\"]}\n"},
			category: "environment",
			tool:     "Devbox",
		},
		{
			name:     "Devbox lockfile",
			files:    map[string]string{"devbox.lock": "{\"lockfile_version\": \"1\"}\n"},
			category: "environment",
			tool:     "Devbox",
		},
		{
			name:           "Jsonnet Bundler manifest",
			files:          map[string]string{"jsonnetfile.json": "{\"version\": 1, \"dependencies\": []}\n"},
			tool:           "Jsonnet Bundler",
			packageManager: true,
		},
		{
			name:           "Jsonnet Bundler lockfile",
			files:          map[string]string{"jsonnetfile.lock.json": "{\"version\": 1, \"dependencies\": []}\n"},
			tool:           "Jsonnet Bundler",
			packageManager: true,
		},
		{
			name:     "Puppet control repository",
			files:    map[string]string{"Puppetfile": "forge 'https://forge.puppet.com'\nmod 'puppetlabs/stdlib'\n"},
			category: "infrastructure",
			tool:     "Puppet",
		},
		{
			name:           "Mint package list",
			files:          map[string]string{"Mintfile": "realm/SwiftLint@0.57.0\n"},
			tool:           "Mint",
			packageManager: true,
		},
		{
			name:     "Meteor package definition",
			files:    map[string]string{"package.js": "Package.describe({ name: 'example:package', version: '1.0.0' });\n"},
			category: "build",
			tool:     "Meteor",
		},
		{
			name:     "Meteor versions file",
			files:    map[string]string{".meteor/versions": "meteor-base@1.5.2\n"},
			category: "build",
			tool:     "Meteor",
		},
		{
			name:     "Meteor JSON versions file",
			files:    map[string]string{"versions.json": "{\"meteor\": \"3.0.0\"}\n"},
			category: "build",
			tool:     "Meteor",
		},
		{
			name:     "Helmfile config",
			files:    map[string]string{"helmfile.yaml": "releases: []\n"},
			category: "infrastructure",
			tool:     "Helmfile",
		},
		{
			name:     "templated Helmfile config",
			files:    map[string]string{"helmfile.yaml.gotmpl": "releases: []\n"},
			category: "infrastructure",
			tool:     "Helmfile",
		},
		{
			name:     "Helmfile directory config",
			files:    map[string]string{"helmfile.d/production.yaml": "releases: []\n"},
			category: "infrastructure",
			tool:     "Helmfile",
		},
		{
			name:     "templated Helmfile directory config",
			files:    map[string]string{"helmfile.d/production.yaml.gotmpl": "releases: []\n"},
			category: "infrastructure",
			tool:     "Helmfile",
		},
		{
			name:     "Argo CD YAML manifest",
			files:    map[string]string{"deploy/argocd.yaml": "apiVersion: argoproj.io/v1alpha1\nkind: Application\n"},
			category: "infrastructure",
			tool:     "Argo CD",
			notTool:  "Flux",
		},
		{
			name:     "Argo CD YML manifest",
			files:    map[string]string{"deploy/argocd.yml": "apiVersion: argoproj.io/v1alpha1\nkind: AppProject\n"},
			category: "infrastructure",
			tool:     "Argo CD",
			notTool:  "Flux",
		},
		{
			name:     "Argo CD ApplicationSet manifest",
			files:    map[string]string{"deploy/applicationset.yaml": "apiVersion: argoproj.io/v1alpha1\nkind: ApplicationSet\n"},
			category: "infrastructure",
			tool:     "Argo CD",
			notTool:  "Flux",
		},
		{
			name:     "Flux YAML manifest",
			files:    map[string]string{"clusters/flux.yaml": "apiVersion: source.toolkit.fluxcd.io/v1\nkind: GitRepository\n"},
			category: "infrastructure",
			tool:     "Flux",
			notTool:  "Argo CD",
		},
		{
			name:     "Flux YML manifest",
			files:    map[string]string{"clusters/flux.yml": "apiVersion: kustomize.toolkit.fluxcd.io/v1\nkind: Kustomization\n"},
			category: "infrastructure",
			tool:     "Flux",
			notTool:  "Argo CD",
		},
	}

	knowledgeBase := loadKB(t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dir := t.TempDir()
			for path, content := range test.files {
				writeProjectFile(t, dir, path, content)
			}

			report, err := New(knowledgeBase, dir).Run()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if test.packageManager {
				if !slices.ContainsFunc(report.PackageManagers, func(d brief.Detection) bool {
					return d.Name == test.tool
				}) {
					t.Errorf("expected %s package manager, got %v", test.tool, packageManagerNames(report))
				}
				return
			}
			assertToolDetected(t, report, test.category, test.tool)
			if test.notTool != "" {
				assertToolNotDetected(t, report, test.category, test.notTool)
			}
		})
	}
}

func TestMissingDevelopmentToolFixtures(t *testing.T) {
	tests := []struct {
		fixture        string
		category       string
		tool           string
		packageManager bool
	}{
		{fixture: "babashka-project", category: "build", tool: "Babashka"},
		{fixture: "opentofu-project", category: "infrastructure", tool: "OpenTofu"},
		{fixture: "stack-project", tool: "Stack", packageManager: true},
		{fixture: "chef-project", category: "infrastructure", tool: "Chef"},
		{fixture: "bitbake-project", category: "build", tool: "BitBake"},
		{fixture: "buck-project", category: "monorepo", tool: "Buck"},
		{fixture: "devbox-project", category: "environment", tool: "Devbox"},
		{fixture: "jsonnet-bundler-project", tool: "Jsonnet Bundler", packageManager: true},
		{fixture: "puppet-project", category: "infrastructure", tool: "Puppet"},
		{fixture: "mint-project", tool: "Mint", packageManager: true},
		{fixture: "meteor-project", category: "build", tool: "Meteor"},
		{fixture: "helmfile-project", category: "infrastructure", tool: "Helmfile"},
		{fixture: "argocd-project", category: "infrastructure", tool: "Argo CD"},
		{fixture: "flux-project", category: "infrastructure", tool: "Flux"},
	}

	for _, test := range tests {
		t.Run(test.tool, func(t *testing.T) {
			report := runOn(t, "../testdata/"+test.fixture)
			if test.packageManager {
				if !slices.ContainsFunc(report.PackageManagers, func(d brief.Detection) bool {
					return d.Name == test.tool
				}) {
					t.Errorf("expected %s package manager, got %v", test.tool, packageManagerNames(report))
				}
				return
			}
			assertToolDetected(t, report, test.category, test.tool)
		})
	}
}

func TestGitOpsContentDetectionIgnoresGenericKubernetesYAML(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{
			name:    "generic Kubernetes manifest",
			content: "apiVersion: apps/v1\nkind: Deployment\n",
		},
		{
			name:    "Argo Workflows manifest",
			content: "apiVersion: argoproj.io/v1alpha1\nkind: Workflow\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dir := t.TempDir()
			writeProjectFile(t, dir, "deploy/manifest.yaml", test.content)

			report := runOn(t, dir)
			assertToolNotDetected(t, report, "infrastructure", "Argo CD")
			assertToolNotDetected(t, report, "infrastructure", "Flux")
		})
	}
}
