# Git MCP Server (Go)

A Model Context Protocol (MCP) server for Git repository interaction and automation, written in Go. This server provides tools to read, search, and manipulate Git repositories via Large Language Models.

## Features

This MCP server provides the following Git operations as tools:

- **git_status**: Shows the working tree status
- **git_diff_unstaged**: Shows changes in the working directory that are not yet staged
- **git_diff_staged**: Shows changes that are staged for commit
- **git_diff**: Shows differences between branches or commits
- **git_commit**: Records changes to the repository
- **git_add**: Adds file contents to the staging area
- **git_reset**: Unstages all staged changes
- **git_log**: Shows the commit logs
- **git_create_branch**: Creates a new branch from an optional base branch
- **git_checkout**: Switches branches
- **git_show**: Shows the contents of a commit
- **git_init**: Initialize a new Git repository
- **git_push**: Pushes local commits to a remote repository

## Installation

### Prerequisites

- Go 1.18 or higher
- Git installed on your system

### Download Prebuilt Binaries

You can download prebuilt binaries for your platform from the [GitHub Releases](https://github.com/geropl/git-mcp-go/releases) page.

### Building from Source

```bash
# Clone the repository
git clone https://github.com/geropl/git-mcp-go.git
cd git-mcp-go

# Build the server
go build -o git-mcp-go .
```

## Usage

### Command Line Options

```
Usage of git-mcp-go:
  -r string
        Git repository path (shorthand)
  -repository string
        Git repository path
  -mode string
        Git operation mode: 'shell' or 'go-git' (default "shell")
  -v    Enable verbose logging
```

The `-mode` flag allows you to choose between two different implementations:

- **shell**: Uses the Git CLI commands via shell execution (default)
- **go-git**: Uses the go-git library for Git operations where possible

### Running the Server

```bash
# Run with a specific repository
./git-mcp-go -r /path/to/git/repository

# Run with verbose logging
./git-mcp-go -v -r /path/to/git/repository

# Run with go-git implementation
./git-mcp-go -mode go-git -r /path/to/git/repository
```

### Integration with Claude Desktop

Add this to your `claude_desktop_config.json`:

```json
"mcpServers": {
  "git": {
    "command": "/path/to/git-mcp-go",
    "args": ["-mode", "shell", "-r", "/path/to/git/repository"]
  }
}
```

Or if you prefer the go-git implementation:

```json
"mcpServers": {
  "git": {
    "command": "/path/to/git-mcp-go",
    "args": ["-mode", "go-git", "-r", "/path/to/git/repository"]
  }
}
```

## Implementation Details

This server is implemented using:

- [mcp-go](https://github.com/mark3labs/mcp-go): Go SDK for the Model Context Protocol
- [go-git](https://github.com/go-git/go-git): Pure Go implementation of Git

For operations not supported by go-git, the server falls back to using the Git CLI.

## Testing

The server includes comprehensive tests for all Git operations. The tests are designed to run against both implementation modes:

```bash
# Run all tests
go test ./pkg -v

# Run specific tests
go test ./pkg -v -run TestGitOperations/push
```

The test suite creates temporary repositories for each test case and verifies that the operations work correctly in both modes.

## Continuous Integration

This project uses GitHub Actions for continuous integration and deployment:

- Automated tests run on every pull request to the main branch
- Releases are created when a tag with the format `v*` is pushed
- Each release includes binaries for multiple platforms:
  - Linux (amd64, arm64)
  - macOS (amd64, arm64)
  - Windows (amd64)

To create a new release:
```bash
# Tag the current commit
git tag v1.0.0

# Push the tag to GitHub
git push origin v1.0.0
```

## License

This project is licensed under the MIT License.
