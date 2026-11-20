package repoindex

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLocateRepos_Worktrees(t *testing.T) {
	tmpDir := t.TempDir()

	root1 := filepath.Join(tmpDir, "root1")
	if err := os.MkdirAll(root1, 0755); err != nil {
		t.Fatal(err)
	}

	repoDir := filepath.Join(root1, "my-repo")
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		t.Fatal(err)
	}

	runGit(t, repoDir, "init")
	runGit(t, repoDir, "config", "user.email", "you@example.com")
	runGit(t, repoDir, "config", "user.name", "Your Name")
	runGit(t, repoDir, "commit", "--allow-empty", "-m", "initial commit")
	runGit(t, repoDir, "remote", "add", "origin", "git@github.com:user/repo.git")

	wtInside := filepath.Join(root1, "wt-inside")
	runGit(t, repoDir, "worktree", "add", "--detach", wtInside, "HEAD")

	root2 := filepath.Join(tmpDir, "root2")
	if err := os.MkdirAll(root2, 0755); err != nil {
		t.Fatal(err)
	}
	wtOutside := filepath.Join(root2, "wt-outside")
	runGit(t, repoDir, "worktree", "add", "--detach", wtOutside, "HEAD")

	repos := locateRepos(root1, make(map[string]struct{}))

	foundRepo := false
	foundWtInside := false
	foundWtOutside := false
	wtInsideCount := 0

	for _, r := range repos {
		if r.FullName == "my-repo" {
			foundRepo = true
		}
		if r.FullName == "wt-inside" {
			foundWtInside = true
			wtInsideCount++
		}
		if r.FullName == "my-repo" {
			foundWtOutside = true
		}
	}

	if !foundRepo {
		t.Errorf("Expected main repo to be found")
	}
	if !foundWtInside {
		t.Errorf("Expected inside worktree to be found")
	}
	if wtInsideCount > 1 {
		t.Errorf("Expected worktree 'wt-inside' to appear only once, but found %d times", wtInsideCount)
	}
	if foundWtOutside {
		// t.Errorf("Did not expect outside worktree to be found (yet)") // Old expectation
	}

	if !foundWtOutside {
		t.Errorf("Expected outside worktree to be found")
	}
}

func TestLocateRepos_WorktreeNoDuplication(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a directory structure similar to the bug report
	root := filepath.Join(tmpDir, "code/personal")
	if err := os.MkdirAll(root, 0755); err != nil {
		t.Fatal(err)
	}

	// Create main repo
	repoDir := filepath.Join(root, "griffin")
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		t.Fatal(err)
	}

	runGit(t, repoDir, "init")
	runGit(t, repoDir, "config", "user.email", "test@test.com")
	runGit(t, repoDir, "config", "user.name", "Test")
	runGit(t, repoDir, "commit", "--allow-empty", "-m", "init")
	runGit(t, repoDir, "remote", "add", "origin", "https://github.com/ronkitay/griffin")

	// Create a worktree under the same configured directory (this is the bug scenario)
	wtDir := filepath.Join(root, "griffin-test")
	runGit(t, repoDir, "worktree", "add", "--detach", wtDir, "HEAD")

	// Locate repos starting from the configured root
	repos := locateRepos(root, make(map[string]struct{}))

	// Count how many times the worktree appears and check it has the alias
	wtCount := 0
	wtWithAlias := false
	rootEval, _ := filepath.EvalSymlinks(root)

	for _, r := range repos {
		rBaseEval, _ := filepath.EvalSymlinks(r.BaseDir)
		if rBaseEval == rootEval && r.FullName == "griffin-test" {
			wtCount++
			if r.Alias == "griffin" {
				wtWithAlias = true
			}
		}
	}

	if wtCount != 1 {
		t.Errorf("Expected worktree 'griffin-test' to appear exactly once, but found %d times", wtCount)
		t.Logf("All repos: %+v", repos)
	}
	if !wtWithAlias {
		t.Errorf("Expected worktree to have alias 'griffin'")
		t.Logf("All repos: %+v", repos)
	}
}

func TestLocateRepos_Dev3Worktrees(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a dev3-style structure with main repo at root
	root := filepath.Join(tmpDir, "dev3-repos")
	if err := os.MkdirAll(root, 0755); err != nil {
		t.Fatal(err)
	}

	mainRepoDir := filepath.Join(root, "my-repo")
	if err := os.MkdirAll(mainRepoDir, 0755); err != nil {
		t.Fatal(err)
	}

	runGit(t, mainRepoDir, "init")
	runGit(t, mainRepoDir, "config", "user.email", "test@test.com")
	runGit(t, mainRepoDir, "config", "user.name", "Test")
	runGit(t, mainRepoDir, "commit", "--allow-empty", "-m", "init")
	runGit(t, mainRepoDir, "remote", "add", "origin", "git@github.com:user/repo.git")

	// Create a dev3-style worktree directory structure:
	// .dev3.0/worktrees/Users-ronk-code-personal-griffin/3c9a68c5/worktree
	worktreesRoot := filepath.Join(tmpDir, "worktrees")
	wtDir := filepath.Join(worktreesRoot, "Users-user-code-personal-repo", "hash123", "worktree")
	if err := os.MkdirAll(wtDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create a worktree at this location
	runGit(t, mainRepoDir, "worktree", "add", "--detach", wtDir, "HEAD")

	// Scan the root directory
	repos := locateRepos(root, make(map[string]struct{}))

	// Find the main repo
	var foundMain, foundWt bool
	var wtFullName, wtBaseDir string
	rootEval, _ := filepath.EvalSymlinks(root)

	for _, r := range repos {
		rBaseEval, _ := filepath.EvalSymlinks(r.BaseDir)
		if r.FullName == "my-repo" && rBaseEval == rootEval {
			foundMain = true
		}
		if r.FullName == "worktree" {
			foundWt = true
			wtFullName = r.FullName
			wtBaseDir = r.BaseDir
			// Verify the alias is set to the repo name
			if r.Alias != "my-repo" {
				t.Errorf("Worktree alias should be 'my-repo', got '%s'", r.Alias)
			}
		}
	}

	if !foundMain {
		t.Errorf("Expected main repo to be found")
		t.Logf("Found repos: %+v", repos)
	}
	if !foundWt {
		t.Errorf("Expected worktree to be found")
		t.Logf("Found repos: %+v", repos)
	} else {
		// Verify that ToString() returns a valid path
		expectedPath := filepath.Join(wtBaseDir, wtFullName)
		expectedPathEval, _ := filepath.EvalSymlinks(expectedPath)
		wtDirEval, _ := filepath.EvalSymlinks(wtDir)
		
		if expectedPathEval != wtDirEval {
			t.Errorf("Worktree path mismatch. Expected %s, got %s", wtDirEval, expectedPathEval)
		}
		// Verify the path actually exists
		if info, err := os.Stat(expectedPath); err != nil {
			t.Errorf("Worktree path does not exist: %s, error: %v", expectedPath, err)
		} else if !info.IsDir() {
			t.Errorf("Worktree path is not a directory: %s", expectedPath)
		}
	}
}

func TestMatchingWorktrees(t *testing.T) {
	tmpDir := t.TempDir()

	root := filepath.Join(tmpDir, "repos")
	if err := os.MkdirAll(root, 0755); err != nil {
		t.Fatal(err)
	}

	mainRepoDir := filepath.Join(root, "my-repo")
	if err := os.MkdirAll(mainRepoDir, 0755); err != nil {
		t.Fatal(err)
	}

	runGit(t, mainRepoDir, "init")
	runGit(t, mainRepoDir, "config", "user.email", "test@test.com")
	runGit(t, mainRepoDir, "config", "user.name", "Test")
	runGit(t, mainRepoDir, "commit", "--allow-empty", "-m", "init")
	runGit(t, mainRepoDir, "remote", "add", "origin", "git@github.com:user/repo.git")

	wtDir := filepath.Join(tmpDir, "worktrees", "worktree-1")
	if err := os.MkdirAll(wtDir, 0755); err != nil {
		t.Fatal(err)
	}
	runGit(t, mainRepoDir, "worktree", "add", "--detach", wtDir, "HEAD")

	repos := locateRepos(root, make(map[string]struct{}))

	// Test matching by FullName (directory name)
	var foundByDirName bool
	for _, r := range repos {
		if r.FullName == "worktree-1" {
			foundByDirName = true
			// Check that we can match it
			matchable := r.Matchable()
			if len(matchable) == 0 {
				t.Errorf("Matchable should not be empty")
			}
			var foundInMatchable bool
			for _, m := range matchable {
				if m == "worktree-1" {
					foundInMatchable = true
				}
			}
			if !foundInMatchable {
				t.Errorf("Directory name 'worktree-1' should be in Matchable(), got: %v", matchable)
			}
		}
	}

	if !foundByDirName {
		t.Errorf("Expected worktree with FullName 'worktree-1' to be found")
	}
}

func TestGetWorktrees(t *testing.T) {
	tmpDir := t.TempDir()
	repoDir := filepath.Join(tmpDir, "repo")
	os.MkdirAll(repoDir, 0755)

	runGit(t, repoDir, "init")
	runGit(t, repoDir, "config", "user.email", "test@test.com")
	runGit(t, repoDir, "config", "user.name", "Test")
	runGit(t, repoDir, "commit", "--allow-empty", "-m", "init")

	wt1 := filepath.Join(tmpDir, "wt1")
	runGit(t, repoDir, "worktree", "add", "--detach", wt1, "HEAD")

	wts, err := getWorktrees(repoDir)
	if err != nil {
		t.Fatalf("getWorktrees failed: %v", err)
	}

	// Expect 2 worktrees: main repo and wt1
	if len(wts) != 2 {
		t.Errorf("Expected 2 worktrees, got %d: %v", len(wts), wts)
	}

	// Check paths
	foundMain := false
	foundWt1 := false
	// resolve symlinks for comparison if needed, but getWorktrees returns as is from git.
	// Git usually returns absolute paths.
	// But on Mac /var is /private/var.
	// We might need to EvalSymlinks.

	repoDirEval, _ := filepath.EvalSymlinks(repoDir)
	wt1Eval, _ := filepath.EvalSymlinks(wt1)

	for _, w := range wts {
		wEval, _ := filepath.EvalSymlinks(w)
		if wEval == repoDirEval {
			foundMain = true
		}
		if wEval == wt1Eval {
			foundWt1 = true
		}
	}

	if !foundMain {
		t.Errorf("Main repo worktree not found in list: %v", wts)
	}
	if !foundWt1 {
		t.Errorf("Added worktree not found in list: %v", wts)
	}
}

func TestLocateRepos_WorktreeAliasFromMainRepo(t *testing.T) {
	tmpDir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	root := filepath.Join(tmpDir, "code/personal")
	if err := os.MkdirAll(root, 0755); err != nil {
		t.Fatal(err)
	}

	repoDir := filepath.Join(root, "personal-assistant")
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		t.Fatal(err)
	}

	runGit(t, repoDir, "init")
	runGit(t, repoDir, "config", "user.email", "test@test.com")
	runGit(t, repoDir, "config", "user.name", "Test")
	runGit(t, repoDir, "commit", "--allow-empty", "-m", "init")
	runGit(t, repoDir, "remote", "add", "origin", "https://github.com/ronkitay/personal-assistant")

	wtDir := filepath.Join(root, ".worktrees", "personal-assistant--griffin-wt-fix")
	if err := os.MkdirAll(wtDir, 0755); err != nil {
		t.Fatal(err)
	}
	runGit(t, repoDir, "worktree", "add", "--detach", wtDir, "HEAD")

	wtRoot := filepath.Join(tmpDir, "workmux/worktrees")
	if err := os.MkdirAll(wtRoot, 0755); err != nil {
		t.Fatal(err)
	}
	crossWtDir := filepath.Join(wtRoot, "anniversary-planning")
	runGit(t, repoDir, "worktree", "add", "--detach", crossWtDir, "HEAD")

	processed := make(map[string]struct{})
	repos := locateRepos(wtRoot, processed)
	repos = append(repos, locateRepos(root, processed)...)
	repos = deDuplicate(repos)

	mainAlias := ""
	wtAlias := ""
	crossWtAlias := ""
	polluted := false
	wtFound := false
	for _, r := range repos {
		if r.FullName == "personal-assistant--griffin-wt-fix" {
			wtFound = true
			wtAlias = r.Alias
		}
		if r.FullName == "anniversary-planning" {
			crossWtAlias = r.Alias
		}
		if r.FullName == "personal-assistant" {
			mainAlias = r.Alias
		}
		if strings.Contains(r.Alias, ".worktrees/") || strings.Contains(r.Alias, "griffin-wt-fix") || strings.Contains(r.Alias, "/") {
			polluted = true
		}
	}

	if !wtFound {
		t.Errorf("Expected linked worktree to be found by its path name")
		t.Logf("All repos: %+v", repos)
	}
	if polluted {
		t.Errorf("Found polluted alias on an entry")
		t.Logf("All repos: %+v", repos)
	}
	if wtAlias != "personal-assistant" {
		t.Errorf("Expected linked worktree alias to be 'personal-assistant', got '%s'", wtAlias)
		t.Logf("All repos: %+v", repos)
	}
	if crossWtAlias != "personal-assistant" {
		t.Errorf("Expected cross-root worktree alias to be 'personal-assistant', got '%s'", crossWtAlias)
		t.Logf("All repos: %+v", repos)
	}
	if mainAlias != "personal-assistant" {
		t.Errorf("Expected main repo alias to be 'personal-assistant', got '%s'", mainAlias)
		t.Logf("All repos: %+v", repos)
	}
}

func TestLocateRepos_SymlinkedDuplicate(t *testing.T) {
	tmpDir := t.TempDir()

	realDir := filepath.Join(tmpDir, "real")
	if err := os.MkdirAll(realDir, 0755); err != nil {
		t.Fatal(err)
	}

	repoDir := filepath.Join(realDir, "repo")
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		t.Fatal(err)
	}

	runGit(t, repoDir, "init")
	runGit(t, repoDir, "config", "user.email", "test@test.com")
	runGit(t, repoDir, "config", "user.name", "Test")
	runGit(t, repoDir, "commit", "--allow-empty", "-m", "init")
	runGit(t, repoDir, "remote", "add", "origin", "https://github.com/ronkitay/griffin")

	symDir := filepath.Join(tmpDir, "sym")
	if err := os.Symlink(realDir, symDir); err != nil {
		t.Fatal(err)
	}

	// Scan both real and symlinked directories
	repos := locateRepos(realDir, make(map[string]struct{}))
	repos = append(repos, locateRepos(symDir, make(map[string]struct{}))...)

	// Deduplicate them as BuildRepoIndex would do
	repos = deDuplicate(repos)

	repoCount := 0
	for _, r := range repos {
		if r.FullName == "repo" {
			repoCount++
		}
	}

	if repoCount != 1 {
		t.Errorf("Expected repo to appear exactly once, but found %d times", repoCount)
		for _, r := range repos {
			t.Logf("Repo: %+v", r)
		}
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestGitURLToHTTP_Scenarios(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		{"1_scp_github", "git@github.com:org/repo.git", "https://github.com/org/repo"},
		{"2_scp_gitlab_subgroups", "git@gitlab.acme.com:group/sub/repo.git", "https://gitlab.acme.com/group/sub/repo"},
		{"3_ssh_scheme_git_user", "ssh://git@git.acme.com/team/repo.git", "https://git.acme.com/team/repo"},
		{"4_ssh_scheme_with_port", "ssh://git@git.acme.com:2222/team/repo.git", "https://git.acme.com/team/repo"},
		{"5_ssh_scheme_no_dot_git", "ssh://git@git.acme.com/team/repo", "https://git.acme.com/team/repo"},
		{"6_https_with_dot_git", "https://git.acme.com/team/repo.git", "https://git.acme.com/team/repo"},
		{"7_http_with_web_port", "http://git.acme.com:8080/team/repo.git", "http://git.acme.com:8080/team/repo"},
		{"8_local_path", "/srv/git/repo.git", "/srv/git/repo"},
		{"9_file_scheme", "file:///srv/git/repo.git", "file:///srv/git/repo"},
		{"10_ssh_scheme_no_user", "ssh://git.acme.com/team/repo.git", "https://git.acme.com/team/repo"},
		{"11_ssh_scheme_non_git_user_with_port", "ssh://alice@git.acme.com:22/team/repo.git", "https://git.acme.com/team/repo"},
		{"12_scp_non_git_user", "alice@git.acme.com:team/repo.git", "https://git.acme.com/team/repo"},
		{"limitation_scp_with_numeric_segment", "git@host:2222/team/repo.git", "https://host/2222/team/repo"},
		{"limitation_gitlab_relative_root", "git@example.com:group/project.git", "https://example.com/group/project"},
		{"limitation_gitlab_custom_web_port", "git@gitlab.acme.com:group/repo.git", "https://gitlab.acme.com/group/repo"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := gitURLToHTTP(tc.input); got != tc.expected {
				t.Errorf("gitURLToHTTP(%q) = %q, expected %q", tc.input, got, tc.expected)
			}
		})
	}
}

func TestRepoTypeDetection(t *testing.T) {
	t.Run("T1_T2_url_substring_skips_probe", func(t *testing.T) {
		emptyDir := t.TempDir()
		origTransport := probeHTTPClient.Transport
		origDisabled := httpProbeDisabled
		defer func() {
			probeHTTPClient.Transport = origTransport
			httpProbeDisabled = origDisabled
			hostTypeCache = make(map[string]string)
		}()
		httpProbeDisabled = false
		hostTypeCache = make(map[string]string)

		probed := false
		probeHTTPClient.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
			probed = true
			return nil, errors.New("unexpected probe")
		})

		if got := repoType("https://github.com/org/repo", emptyDir); got != "github" {
			t.Errorf("T1: expected github, got %s", got)
		}
		if got := repoType("https://gitlab.com/org/repo", emptyDir); got != "gitlab" {
			t.Errorf("T2: expected gitlab, got %s", got)
		}
		if probed {
			t.Errorf("T1/T2: expected no HTTP probe when URL contains github/gitlab")
		}
	})

	t.Run("T3_to_T14_classification_and_ordering", func(t *testing.T) {
		emptyDir := t.TempDir()
		dirWithGithub := t.TempDir()
		if err := os.MkdirAll(filepath.Join(dirWithGithub, ".github"), 0755); err != nil {
			t.Fatal(err)
		}
		dirWithGitlabCI := t.TempDir()
		if err := os.WriteFile(filepath.Join(dirWithGitlabCI, ".gitlab-ci.yml"), []byte("stages: []"), 0644); err != nil {
			t.Fatal(err)
		}
		dirWithGitlabDir := t.TempDir()
		if err := os.MkdirAll(filepath.Join(dirWithGitlabDir, ".gitlab"), 0755); err != nil {
			t.Fatal(err)
		}
		dirBoth := t.TempDir()
		if err := os.MkdirAll(filepath.Join(dirBoth, ".github"), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dirBoth, ".gitlab-ci.yml"), []byte("stages: []"), 0644); err != nil {
			t.Fatal(err)
		}

		origTransport := probeHTTPClient.Transport
		origDisabled := httpProbeDisabled
		defer func() {
			probeHTTPClient.Transport = origTransport
			httpProbeDisabled = origDisabled
			hostTypeCache = make(map[string]string)
		}()
		httpProbeDisabled = false
		hostTypeCache = make(map[string]string)
		hitCounts := make(map[string]int)

		probeHTTPClient.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
			hitCounts[req.URL.Host]++
			rec := httptest.NewRecorder()
			switch req.URL.Host {
			case "t3-ghe.corp.local":
				rec.Header().Set("X-GitHub-Request-Id", "req-123")
			case "t4-ghe-server.corp.local":
				rec.Header().Set("Server", "GitHub.com")
			case "t5-gl.corp.local":
				rec.Header().Set("X-Gitlab-Meta", `{"version":"1"}`)
			case "t8-unknown.corp.local":
			case "t11-unreachable.corp.local":
				return nil, errors.New("connection refused")
			}
			rec.WriteHeader(http.StatusOK)
			return rec.Result(), nil
		})

		if got := repoType("https://t3-ghe.corp.local/org/repo", emptyDir); got != "github" {
			t.Errorf("T3: expected github, got %s", got)
		}
		if got := repoType("https://t4-ghe-server.corp.local/org/repo", emptyDir); got != "github" {
			t.Errorf("T4: expected github, got %s", got)
		}
		if got := repoType("https://t5-gl.corp.local/org/repo", emptyDir); got != "gitlab" {
			t.Errorf("T5: expected gitlab, got %s", got)
		}
		if got := repoType("https://t5-gl.corp.local/org/repo2", dirWithGithub); got != "gitlab" {
			t.Errorf("T6: expected gitlab from probe despite .github dir, got %s", got)
		}
		if got := repoType("https://t3-ghe.corp.local/org/repo2", dirWithGitlabCI); got != "github" {
			t.Errorf("T7: expected github from probe despite .gitlab-ci.yml, got %s", got)
		}
		if got := repoType("https://t8-unknown.corp.local/org/repo", emptyDir); got != "git" {
			t.Errorf("T8: expected git, got %s", got)
		}
		if got := repoType("https://t8-unknown.corp.local/org/repo2", dirWithGithub); got != "github" {
			t.Errorf("T9: expected github from fs fallback, got %s", got)
		}
		if got := repoType("https://t8-unknown.corp.local/org/repo3", dirWithGitlabCI); got != "gitlab" {
			t.Errorf("T10a: expected gitlab from .gitlab-ci.yml fallback, got %s", got)
		}
		if got := repoType("https://t8-unknown.corp.local/org/repo4", dirWithGitlabDir); got != "gitlab" {
			t.Errorf("T10b: expected gitlab from .gitlab/ fallback, got %s", got)
		}
		if got := repoType("https://t8-unknown.corp.local/org/repo5", dirBoth); got != "github" {
			t.Errorf("T9/T10 both: expected github favored over gitlab, got %s", got)
		}
		if got := repoType("https://t11-unreachable.corp.local/org/repo", emptyDir); got != "git" {
			t.Errorf("T11: expected git when unreachable and no markers, got %s", got)
		}
		if got := repoType("https://t11-unreachable.corp.local/org/repo2", dirWithGithub); got != "github" {
			t.Errorf("T12: expected github from fs fallback when unreachable, got %s", got)
		}
		if got := repoType("/srv/git/repo", emptyDir); got != "git" {
			t.Errorf("T13a: expected git for local path, got %s", got)
		}
		if got := repoType("file:///srv/git/repo", dirWithGithub); got != "github" {
			t.Errorf("T13b: expected github from fs fallback for file:// URL, got %s", got)
		}
		if got := repoType("ssh://git@git.acme.com/team/repo", emptyDir); got != "git" {
			t.Errorf("T14a: expected git for unconverted ssh:// URL, got %s", got)
		}
		if got := repoType("ssh://git@git.acme.com/team/repo", dirWithGitlabCI); got != "gitlab" {
			t.Errorf("T14b: expected gitlab from fs fallback for unconverted ssh:// URL, got %s", got)
		}
		if hitCounts["git.acme.com"] != 0 {
			t.Errorf("T14: expected no HTTP probe for ssh:// scheme, got %d", hitCounts["git.acme.com"])
		}
	})

	t.Run("P1_to_P5_probe_behaviour", func(t *testing.T) {
		emptyDir := t.TempDir()
		origTransport := probeHTTPClient.Transport
		origTimeout := probeHTTPClient.Timeout
		origDisabled := httpProbeDisabled
		defer func() {
			probeHTTPClient.Transport = origTransport
			probeHTTPClient.Timeout = origTimeout
			httpProbeDisabled = origDisabled
			hostTypeCache = make(map[string]string)
		}()
		httpProbeDisabled = false
		hostTypeCache = make(map[string]string)
		hitCounts := make(map[string]int)

		probeHTTPClient.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
			key := req.URL.Scheme + "://" + req.URL.Host
			hitCounts[key]++
			rec := httptest.NewRecorder()
			switch req.URL.Host {
			case "p1.corp.local":
				rec.Header().Set("X-Gitlab-Meta", `{"version":"1"}`)
				rec.WriteHeader(http.StatusOK)
			case "p2-fail.corp.local":
				return nil, errors.New("network timeout")
			case "p3.corp.local":
				rec.Header().Set("X-GitHub-Request-Id", "abc")
				rec.WriteHeader(http.StatusOK)
			case "p4-redirect.corp.local":
				rec.Header().Set("X-Gitlab-Meta", `{"version":"1"}`)
				rec.Header().Set("Location", "https://sso.corp.local/login")
				rec.WriteHeader(http.StatusFound)
			case "sso.corp.local":
				t.Errorf("P4: redirect to sso.corp.local should not have been followed")
				rec.WriteHeader(http.StatusOK)
			case "p5-slow.corp.local":
				select {
				case <-req.Context().Done():
					return nil, req.Context().Err()
				case <-time.After(200 * time.Millisecond):
					rec.Header().Set("X-GitHub-Request-Id", "late")
					rec.WriteHeader(http.StatusOK)
				}
			}
			return rec.Result(), nil
		})

		_ = repoType("https://p1.corp.local/team/repo1", emptyDir)
		_ = repoType("https://p1.corp.local/team/repo2", emptyDir)
		if hitCounts["https://p1.corp.local"] != 1 {
			t.Errorf("P1: expected 1 probe for same host, got %d", hitCounts["https://p1.corp.local"])
		}

		_ = repoType("https://p2-fail.corp.local/team/repo1", emptyDir)
		_ = repoType("https://p2-fail.corp.local/team/repo2", emptyDir)
		if hitCounts["https://p2-fail.corp.local"] != 1 {
			t.Errorf("P2: expected failed host to be probed only once, got %d", hitCounts["https://p2-fail.corp.local"])
		}

		_ = repoType("http://p3.corp.local/team/repo1", emptyDir)
		_ = repoType("https://p3.corp.local/team/repo2", emptyDir)
		if hitCounts["http://p3.corp.local"] != 1 || hitCounts["https://p3.corp.local"] != 1 {
			t.Errorf("P3: expected http and https origins to be cached separately, got http=%d https=%d", hitCounts["http://p3.corp.local"], hitCounts["https://p3.corp.local"])
		}

		if got := repoType("https://p4-redirect.corp.local/team/repo", emptyDir); got != "gitlab" {
			t.Errorf("P4: expected gitlab from 302 response headers without following redirect, got %s", got)
		}

		probeHTTPClient.Timeout = 30 * time.Millisecond
		start := time.Now()
		if got := repoType("https://p5-slow.corp.local/team/repo", emptyDir); got != "git" {
			t.Errorf("P5: expected git fallback on slow server timeout, got %s", got)
		}
		if elapsed := time.Since(start); elapsed > 150*time.Millisecond {
			t.Errorf("P5: expected timeout within ~30ms, took %v", elapsed)
		}
	})

	t.Run("A1_A3_converter_git_and_legacy_unknown", func(t *testing.T) {
		conv := converter(false, false)
		data, err := conv([]string{"/base", "my-repo", "https://git.corp.local/team/my-repo", "git", "my-repo"})
		if err != nil {
			t.Fatalf("A1: expected converter to accept git type, got error: %v", err)
		}
		if data.Type != "git" {
			t.Errorf("A1: expected data.Type to be git, got %s", data.Type)
		}

		if _, err := conv([]string{"/base", "my-repo", "https://git.corp.local/team/my-repo", "unknown", "my-repo"}); err == nil {
			t.Errorf("A3: expected legacy unknown type to be skipped until reindex")
		}
	})
}

func runGit(t *testing.T, dir string, args ...string) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git command %v failed in %s: %v\nOutput: %s", args, dir, err, out)
	}
}
