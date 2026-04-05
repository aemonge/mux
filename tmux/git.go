package tmux

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const gitInfoCacheTTL = 5 * time.Second

// GitInfo holds git metadata for a directory.
type GitInfo struct {
	Branch     string
	IsWorktree bool
}

type cachedGitInfo struct {
	info      GitInfo
	expiresAt time.Time
}

var (
	gitCache   = make(map[string]cachedGitInfo)
	gitCacheMu sync.Mutex
)

// LookupGitInfo returns git branch and worktree info for the given directory.
// Results are cached with a TTL.
func LookupGitInfo(dir string) GitInfo {
	if dir == "" {
		return GitInfo{}
	}

	gitCacheMu.Lock()
	if cached, ok := gitCache[dir]; ok && time.Now().Before(cached.expiresAt) {
		gitCacheMu.Unlock()
		return cached.info
	}
	gitCacheMu.Unlock()

	info := detectGitInfo(dir)

	gitCacheMu.Lock()
	gitCache[dir] = cachedGitInfo{
		info:      info,
		expiresAt: time.Now().Add(gitInfoCacheTTL),
	}
	gitCacheMu.Unlock()

	return info
}

func detectGitInfo(dir string) GitInfo {
	gitDir := filepath.Join(dir, ".git")

	fi, err := os.Lstat(gitDir)
	if err != nil {
		return GitInfo{}
	}

	var info GitInfo

	if fi.IsDir() {
		// Normal repo: .git is a directory
		info.Branch = readBranch(filepath.Join(gitDir, "HEAD"))
	} else {
		// Worktree: .git is a file containing "gitdir: /path/to/main/.git/worktrees/name"
		info.IsWorktree = true
		info.Branch = readBranchFromWorktreeGitFile(gitDir)
	}

	return info
}

// readBranch reads the branch name from a git HEAD file.
// HEAD contains "ref: refs/heads/<branch>" or a commit hash.
func readBranch(headPath string) string {
	data, err := os.ReadFile(headPath)
	if err != nil {
		return ""
	}
	line := strings.TrimSpace(string(data))
	if strings.HasPrefix(line, "ref: refs/heads/") {
		return strings.TrimPrefix(line, "ref: refs/heads/")
	}
	// Detached HEAD — return short hash
	if len(line) >= 8 {
		return line[:8]
	}
	return ""
}

// readBranchFromWorktreeGitFile reads branch from a worktree's .git file.
// The file contains "gitdir: /path/to/.git/worktrees/<name>".
// The actual HEAD is at that path/HEAD.
func readBranchFromWorktreeGitFile(gitFilePath string) string {
	data, err := os.ReadFile(gitFilePath)
	if err != nil {
		return ""
	}
	line := strings.TrimSpace(string(data))
	if !strings.HasPrefix(line, "gitdir: ") {
		return ""
	}
	worktreeGitDir := strings.TrimPrefix(line, "gitdir: ")
	return readBranch(filepath.Join(worktreeGitDir, "HEAD"))
}
