package pkg

// Models for Git operations

// GitStatus represents the input for git status operation
type GitStatus struct {
	RepoPath string `json:"repo_path"`
}

// GitDiffUnstaged represents the input for git diff (unstaged) operation
type GitDiffUnstaged struct {
	RepoPath string `json:"repo_path"`
}

// GitDiffStaged represents the input for git diff (staged) operation
type GitDiffStaged struct {
	RepoPath string `json:"repo_path"`
}

// GitDiff represents the input for git diff with target operation
type GitDiff struct {
	RepoPath string `json:"repo_path"`
	Target   string `json:"target"`
}

// GitCommit represents the input for git commit operation
type GitCommit struct {
	RepoPath string `json:"repo_path"`
	Message  string `json:"message"`
}

// GitAdd represents the input for git add operation
type GitAdd struct {
	RepoPath string   `json:"repo_path"`
	Files    []string `json:"files"`
}

// GitReset represents the input for git reset operation
type GitReset struct {
	RepoPath string `json:"repo_path"`
}

// GitLog represents the input for git log operation
type GitLog struct {
	RepoPath string `json:"repo_path"`
	MaxCount int    `json:"max_count,omitempty"`
}

// GitCreateBranch represents the input for git branch creation operation
type GitCreateBranch struct {
	RepoPath   string `json:"repo_path"`
	BranchName string `json:"branch_name"`
	BaseBranch string `json:"base_branch,omitempty"`
}

// GitCheckout represents the input for git checkout operation
type GitCheckout struct {
	RepoPath   string `json:"repo_path"`
	BranchName string `json:"branch_name"`
}

// GitShow represents the input for git show operation
type GitShow struct {
	RepoPath string `json:"repo_path"`
	Revision string `json:"revision"`
}

// GitInit represents the input for git init operation
type GitInit struct {
	RepoPath string `json:"repo_path"`
}
