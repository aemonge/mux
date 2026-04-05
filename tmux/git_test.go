package tmux

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectGitInfoNormalRepo(t *testing.T) {
	dir := t.TempDir()
	gitDir := filepath.Join(dir, ".git")
	os.Mkdir(gitDir, 0755)
	os.WriteFile(filepath.Join(gitDir, "HEAD"), []byte("ref: refs/heads/main\n"), 0644)

	info := detectGitInfo(dir)
	if info.Branch != "main" {
		t.Errorf("Branch = %q, want %q", info.Branch, "main")
	}
	if info.IsWorktree {
		t.Error("IsWorktree should be false for normal repo")
	}
}

func TestDetectGitInfoWorktree(t *testing.T) {
	// Set up a fake main repo
	mainDir := t.TempDir()
	mainGitDir := filepath.Join(mainDir, ".git")
	os.Mkdir(mainGitDir, 0755)

	// Worktree git dir
	worktreeGitDir := filepath.Join(mainGitDir, "worktrees", "feat")
	os.MkdirAll(worktreeGitDir, 0755)
	os.WriteFile(filepath.Join(worktreeGitDir, "HEAD"), []byte("ref: refs/heads/feat-branch\n"), 0644)

	// Set up worktree directory with .git file
	wtDir := t.TempDir()
	gitFile := filepath.Join(wtDir, ".git")
	os.WriteFile(gitFile, []byte("gitdir: "+worktreeGitDir+"\n"), 0644)

	info := detectGitInfo(wtDir)
	if info.Branch != "feat-branch" {
		t.Errorf("Branch = %q, want %q", info.Branch, "feat-branch")
	}
	if !info.IsWorktree {
		t.Error("IsWorktree should be true for worktree")
	}
}

func TestDetectGitInfoDetachedHead(t *testing.T) {
	dir := t.TempDir()
	gitDir := filepath.Join(dir, ".git")
	os.Mkdir(gitDir, 0755)
	os.WriteFile(filepath.Join(gitDir, "HEAD"), []byte("abc123def456789\n"), 0644)

	info := detectGitInfo(dir)
	if info.Branch != "abc123de" {
		t.Errorf("Branch = %q, want short hash %q", info.Branch, "abc123de")
	}
}

func TestDetectGitInfoNotARepo(t *testing.T) {
	dir := t.TempDir()
	info := detectGitInfo(dir)
	if info.Branch != "" || info.IsWorktree {
		t.Errorf("expected empty GitInfo for non-repo, got %+v", info)
	}
}

func TestLookupGitInfoCache(t *testing.T) {
	dir := t.TempDir()
	gitDir := filepath.Join(dir, ".git")
	os.Mkdir(gitDir, 0755)
	os.WriteFile(filepath.Join(gitDir, "HEAD"), []byte("ref: refs/heads/cached\n"), 0644)

	// First call populates cache
	info1 := LookupGitInfo(dir)
	if info1.Branch != "cached" {
		t.Fatalf("Branch = %q, want %q", info1.Branch, "cached")
	}

	// Change HEAD — should still return cached value
	os.WriteFile(filepath.Join(gitDir, "HEAD"), []byte("ref: refs/heads/changed\n"), 0644)
	info2 := LookupGitInfo(dir)
	if info2.Branch != "cached" {
		t.Errorf("Cache miss: Branch = %q, want cached value %q", info2.Branch, "cached")
	}
}

func TestLookupGitInfoEmptyDir(t *testing.T) {
	info := LookupGitInfo("")
	if info.Branch != "" {
		t.Errorf("expected empty branch for empty dir, got %q", info.Branch)
	}
}
