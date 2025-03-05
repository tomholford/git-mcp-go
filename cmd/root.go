package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "git-mcp-go",
	Short: "Git MCP Server",
	Long: `A Model Context Protocol (MCP) server for Git.

This server provides tools for interacting with Git repositories through the MCP protocol.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
