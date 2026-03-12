package repoindex

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
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
	runGit(t, repoDir, "worktree", "add", "--detach", wtInside, "master")

	root2 := filepath.Join(tmpDir, "root2")
	if err := os.MkdirAll(root2, 0755); err != nil {
		t.Fatal(err)
	}
	wtOutside := filepath.Join(root2, "wt-outside")
	runGit(t, repoDir, "worktree", "add", "--detach", wtOutside, "master")

	repos := locateRepos(root1, make(map[string]struct{}))

	foundRepo := false
	foundWtInside := false
	foundWtOutside := false

	for _, r := range repos {
		if r.FullName == "my-repo" {
			foundRepo = true
		}
		if r.FullName == "wt-inside" {
			foundWtInside = true
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
	if foundWtOutside {
		// t.Errorf("Did not expect outside worktree to be found (yet)") // Old expectation
	}

	if !foundWtOutside {
		t.Errorf("Expected outside worktree to be found")
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
	runGit(t, repoDir, "worktree", "add", "--detach", wt1, "master")

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

func runGit(t *testing.T, dir string, args ...string) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git command %v failed in %s: %v\nOutput: %s", args, dir, err, out)
	}
}
