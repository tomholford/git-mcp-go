package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/geropl/git-mcp-go/pkg"
	"github.com/geropl/git-mcp-go/pkg/gitops"
	"github.com/geropl/git-mcp-go/pkg/gitops/gogit"
	"github.com/geropl/git-mcp-go/pkg/gitops/shell"
)

func main() {
	// Parse command line arguments
	var repoPath string
	var verbose bool
	var mode string
	var writeAccess bool

	// Update flags to use double dashes for non-shorthand flags
	flag.StringVar(&repoPath, "repository", "", "Git repository path")
	flag.StringVar(&repoPath, "r", "", "Git repository path (shorthand)")
	flag.StringVar(&mode, "mode", "shell", "Git operation mode: 'shell' or 'go-git'")
	flag.BoolVar(&verbose, "v", false, "Enable verbose logging")
	flag.BoolVar(&writeAccess, "write-access", false, "Enable write access for remote operations (push)")
	flag.Parse()

	// Create the appropriate GitOperations implementation
	var gitOps gitops.GitOperations
	switch strings.ToLower(mode) {
	case "go-git":
		if verbose {
			fmt.Println("Using go-git implementation")
		}
		gitOps = gogit.NewGoGitOperations()
	case "shell":
		if verbose {
			fmt.Println("Using shell implementation")
		}
		gitOps = shell.NewShellGitOperations()
	default:
		if verbose {
			fmt.Println("Using shell implementation")
		}
		gitOps = shell.NewShellGitOperations()
	}

	// Create and configure the Git MCP server
	gitServer := pkg.NewGitServer(repoPath, gitOps, writeAccess)

	// Register all Git tools
	gitServer.RegisterTools()

	// Start the server
	if verbose {
		fmt.Println("Starting Git MCP Server...")
	}
	if err := gitServer.Serve(); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
