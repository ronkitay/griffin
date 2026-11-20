package alfred

import (
	"strings"
	"testing"

	repo "ronkitay.com/griffin/pkg/repoindex"
)

func TestReposAsAlfred_GitType(t *testing.T) {
	output := ReposAsAlfred([]repo.RepoData{
		{
			BaseDir:  "/Users/test/code",
			FullName: "my-generic-repo",
			Url:      "https://git.acme.com/team/my-generic-repo",
			Type:     "git",
		},
	})

	if !strings.Contains(output, `"path": "icons/git.jpg"`) {
		t.Errorf("expected icons/git.jpg in alfred output, got: %s", output)
	}
	if !strings.Contains(output, `"arg": "https://git.acme.com/team/my-generic-repo"`) {
		t.Errorf("expected web URL in alfred output, got: %s", output)
	}
}
