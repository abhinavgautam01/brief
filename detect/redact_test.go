package detect

import (
	"fmt"
	"strings"
	"testing"
)

func TestRedactURL(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{
			"https token as username",
			"https://github_pat_11ABCDEF0abcdefghijklmnop@github.com/owner/repo.git",
			"https://REDACTED@github.com/owner/repo.git",
		},
		{
			"https user and password",
			"https://x-access-token:ghp_abcdefghijklmnopqrstuvwxyz012345@github.com/owner/repo.git",
			"https://REDACTED@github.com/owner/repo.git",
		},
		{
			"https basic auth",
			"https://deploy:hunter2@gitlab.example.com/group/proj.git",
			"https://REDACTED@gitlab.example.com/group/proj.git",
		},
		{
			"http long opaque token",
			"http://0123456789abcdef0123456789abcdef01234567@bitbucket.org/team/repo.git",
			"http://REDACTED@bitbucket.org/team/repo.git",
		},
		{
			"gitlab pat prefix",
			"https://glpat-xxx@gitlab.com/group/proj.git",
			"https://REDACTED@gitlab.com/group/proj.git",
		},
		{
			"ssh url with git user left alone",
			"ssh://git@github.com/owner/repo.git",
			"ssh://git@github.com/owner/repo.git",
		},
		{
			"scp syntax left alone",
			"git@github.com:owner/repo.git",
			"git@github.com:owner/repo.git",
		},
		{
			"scp syntax with password",
			"deploy:secret@host.example.com:path/repo.git",
			"REDACTED@host.example.com:path/repo.git",
		},
		{
			"plain https no userinfo",
			"https://github.com/owner/repo.git",
			"https://github.com/owner/repo.git",
		},
		{
			"short real username preserved",
			"https://andrew@git.example.com/repo.git",
			"https://andrew@git.example.com/repo.git",
		},
		{
			"local path",
			"/srv/git/repo.git",
			"/srv/git/repo.git",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := redactURL(tc.in)
			if got != tc.want {
				t.Errorf("redactURL(%q)\n  got  %q\n  want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestRedactURLNeverContainsToken(t *testing.T) {
	// Fixtures deliberately don't match real provider token regexes so
	// secret scanners don't block pushes of this file.
	tokens := []string{
		"github_pat_" + strings.Repeat("X", 80),
		"ghp_" + strings.Repeat("X", 36),
		"glpat-" + strings.Repeat("X", 20),
	}
	wrappers := []string{
		"https://%s@github.com/o/r.git",
		"https://x:%s@github.com/o/r.git",
		"https://%s:x-oauth-basic@github.com/o/r.git",
	}
	for _, tok := range tokens {
		for _, w := range wrappers {
			in := fmt.Sprintf(w, tok)
			out := redactURL(in)
			if strings.Contains(out, tok) {
				t.Fatalf("token leaked: in=%q out=%q", in, out)
			}
		}
	}
}
