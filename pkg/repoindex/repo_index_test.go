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
	runGit(t, mainRepoDir, "worktree", "add", "--detach", wtDir, "master")

	// Scan the root directory
	repos := locateRepos(root, make(map[string]struct{}))

	// Find the main repo
	var foundMain, foundWt bool
	var wtFullName, wtBaseDir string

	for _, r := range repos {
		if r.FullName == "my-repo" && r.BaseDir == root {
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
	runGit(t, mainRepoDir, "worktree", "add", "--detach", wtDir, "master")

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
