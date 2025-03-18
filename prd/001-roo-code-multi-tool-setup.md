# PRD: Extending Git MCP Setup Mechanism

## Overview

This document outlines the requirements for extending the current setup mechanism of the Git MCP server to support:
1. A new tool called "roo-code" with a different configuration location
2. Multiple tools specified in the --tool parameter

## Background

The Git MCP server currently supports setting up for the "cline" tool, which involves:
- Taking a tool name as input (--tool parameter)
- Looking up the tool's configuration in ~/.config/cline/tools/{tool-name}.json
- Using that configuration to register the tool with Cline

We need to extend this functionality to support a new tool called "roo-code" and to allow multiple tools to be specified in the --tool parameter.

## Goals

1. Support the "roo-code" tool with configuration stored in ~/.config/roo/tools/{tool-name}.json
2. Allow multiple tools to be specified in the --tool parameter (comma-separated)
3. Maintain backward compatibility with existing functionality
4. Ensure a smooth user experience when setting up multiple tools

## Non-Goals

1. Changing the overall architecture of the Git MCP server
2. Modifying the functionality of existing tools
3. Adding support for tools other than "cline" and "roo-code"

## Requirements

### 1. Support for "roo-code" Tool

#### 1.1 Configuration Location
- The "roo-code" tool configuration should be read from ~/.config/roo/tools/{tool-name}.json
- The MCP server config should be written to $HOME/.vscode-server/data/User/globalStorage/rooveterinaryinc.roo-cline/settings/cline_mcp_settings.json

#### 1.2 Tool Setup
- When the --tool parameter is set to "roo-code", the setup command should:
  - Read the tool configuration from the appropriate location
  - Configure the Git MCP server for use with "roo-code"
  - Update the appropriate settings file

### 2. Support for Multiple Tools

#### 2.1 Parameter Format
- The --tool parameter should accept a comma-separated list of tool names
- Example: --tool=cline,roo-code

#### 2.2 Setup Process
- When multiple tools are specified, the setup command should:
  - Process each tool in the order specified
  - Read each tool's configuration from its appropriate location
  - Configure the Git MCP server for each tool
  - Update each tool's settings file

#### 2.3 Error Handling
- If an error occurs during the setup of one tool, the command should:
  - Log the error
  - Continue with the setup of the remaining tools
  - Return a non-zero exit code at the end

### 3. User Experience

#### 3.1 Command Output
- The setup command should provide clear feedback about:
  - Which tools are being set up
  - The progress of the setup process
  - Any errors that occur
  - Successful completion of the setup

#### 3.2 Help Text
- The help text for the --tool parameter should be updated to indicate:
  - The supported tools (cline, roo-code)
  - The ability to specify multiple tools (comma-separated)
  - Example usage

## Technical Design

### Architecture

The implementation follows a modular approach with the following key components:

1. **Command Entry Point**: The `setupCmd` in setup.go handles parameter parsing and orchestration
2. **Tool-specific Setup Functions**: Separate functions for each supported tool
3. **Common Setup Logic**: A shared `setupTool` function that handles the common configuration tasks

### Key Components

#### 1. Multi-tool Processing

```go
// Process each tool
tools := strings.Split(tool, ",")
hasErrors := false

for _, t := range tools {
    t = strings.TrimSpace(t)
    if t == "" {
        continue
    }

    fmt.Printf("Setting up tool: %s\n", t)

    // Set up the tool-specific configuration
    var err error
    switch strings.ToLower(t) {
    case "cline":
        err = setupCline(binaryPath, repoPath, writeAccess, autoApprove)
    case "roo-code":
        err = setupRooCode(binaryPath, repoPath, writeAccess, autoApprove)
    default:
        fmt.Printf("Unsupported tool: %s\n", t)
        fmt.Println("Currently supported tools: cline, roo-code")
        hasErrors = true
        continue
    }

    // Error handling
    if err != nil {
        fmt.Printf("Error setting up %s: %v\n", t, err)
        hasErrors = true
    } else {
        fmt.Printf("git-mcp-go binary successfully set up for %s\n", t)
    }
}

// Return non-zero exit code if any errors occurred
if hasErrors {
    os.Exit(1)
}
```

#### 2. Common Setup Logic

The `setupTool` function handles the common configuration tasks for all tools:

```go
func setupTool(toolName string, binaryPath string, repoPath string, writeAccess bool, autoApprove string, configDir string) error {
    // Create the config directory if it doesn't exist
    if err := os.MkdirAll(configDir, 0755); err != nil {
        return fmt.Errorf("failed to create config directory: %w", err)
    }

    // Build server arguments
    serverArgs := []string{"serve"}
    if repoPath != "" {
        serverArgs = append(serverArgs, "--repository="+repoPath)
    }
    if writeAccess {
        serverArgs = append(serverArgs, "--write-access=true")
    }

    // Process auto-approve flag
    autoApproveTools := []string{}
    if autoApprove != "" {
        // Handle special values and comma-separated lists
        // ...
    }

    // Create or update the MCP settings file
    // ...
}
```

#### 3. Tool-specific Setup Functions

Each tool has its own setup function that determines the appropriate configuration directory:

```go
func setupCline(binaryPath string, repoPath string, writeAccess bool, autoApprove string) error {
    // Determine the Cline config directory based on OS
    // ...
    return setupTool("Cline", binaryPath, repoPath, writeAccess, autoApprove, configDir)
}

func setupRooCode(binaryPath string, repoPath string, writeAccess bool, autoApprove string) error {
    // Determine the Roo Code config directory based on OS
    // ...
    return setupTool("Roo Code", binaryPath, repoPath, writeAccess, autoApprove, configDir)
}
```

### Configuration Structure

The configuration file for each tool follows this structure:

```json
{
  "mcpServers": {
    "git": {
      "command": "/path/to/git-mcp-go",
      "args": ["serve", "--repository=/path/to/repo", "--write-access=true"],
      "disabled": false,
      "autoApprove": ["git_status", "git_diff_unstaged", "git_diff_staged", "git_diff", "git_log", "git_show"]
    }
  }
}
```

### Configuration Paths

#### Cline
- MCP Settings: `$HOME/.vscode-server/data/User/globalStorage/saoudrizwan.claude-dev/settings/cline_mcp_settings.json`

#### Roo Code
- MCP Settings: `$HOME/.vscode-server/data/User/globalStorage/rooveterinaryinc.roo-cline/settings/cline_mcp_settings.json`

## Testing Plan

### Integration Tests

We use a table-driven approach for integration testing to verify the functionality of the extended setup mechanism. This approach allows us to test multiple scenarios with different combinations of parameters in a structured way.

#### Test Structure

Our tests follow this structure:

1. **Define Expectations**: We use a dedicated struct to define what we expect from each test case:
   ```go
   type fileExpectation struct {
       path      string
       content   string
       mustExist bool
   }

   type expectations struct {
       files     map[string]fileExpectation
       errors    []string
       exitCode  int
   }
   ```

2. **Test Cases**: Each test case defines:
   - Input parameters (tool, repository path, etc.)
   - Expected outcomes (files created, content, errors, exit code)

3. **Test Environment**: For each test case:
   - Build the binary from the project root
   - Create a temporary directory structure
   - Set the HOME environment variable to the temp directory
   - Execute the binary with the appropriate parameters
   - Verify the results against expectations
   - Clean up by removing the temp directory

4. **Content Verification**: We verify the JSON content of configuration files by:
   - Parsing both expected and actual content as JSON
   - Normalizing the content (sorting arrays, handling dynamic fields)
   - Comparing the normalized structures

#### Key Test Cases

1. **Cline Only**: Verify setup works for just the Cline tool
2. **Roo Code Only**: Verify setup works for just the Roo Code tool
3. **Multiple Tools**: Verify setup works for both tools specified together
4. **Invalid Tool**: Verify error handling when an invalid tool is specified

### Unit Tests

In addition to integration tests, we implement unit tests for individual components:

1. **Parameter Parsing**: Test parsing of the --tool parameter
2. **Tool Setup Functions**: Test the setupCline and setupRooCode functions
3. **Error Handling**: Test behavior when various errors occur

### Manual Testing

Manual testing steps to verify functionality:

1. **Single Tool Setup**: Test each tool individually
2. **Multiple Tools Setup**: Test setting up multiple tools at once
3. **Error Handling**: Test behavior when errors occur

### Test Execution

```bash
# Run all tests
go test -v ./cmd

# Run specific tests
go test -v ./cmd -run TestSetupCommand
```

## Implementation Notes

### Key Implementation Details

1. **Tool Configuration Paths**:
   - Cline: `$HOME/.vscode-server/data/User/globalStorage/saoudrizwan.claude-dev/settings/cline_mcp_settings.json`
   - Roo Code: `$HOME/.vscode-server/data/User/globalStorage/rooveterinaryinc.roo-cline/settings/cline_mcp_settings.json`

2. **Configuration Content**:
   - Both tools use the same configuration format
   - The `autoApprove` field contains a list of allowed Git operations
   - The command path is dynamically determined based on the binary location

3. **Error Handling**:
   - When one tool setup fails, the command continues with other tools
   - A non-zero exit code is returned if any tool setup fails
   - Clear error messages are provided for each failure

### Example Usage

```bash
# Set up for cline only
./git-mcp-go setup --tool=cline -r /path/to/git/repository

# Set up for roo-code only
./git-mcp-go setup --tool=roo-code -r /path/to/git/repository

# Set up for both cline and roo-code
./git-mcp-go setup --tool=cline,roo-code -r /path/to/git/repository --write-access
```

### Future Considerations

1. **Extensibility**: The current design allows for easy addition of new tools by:
   - Adding a new setup function for the tool
   - Adding the tool to the switch statement in the Run function

2. **Potential Improvements**:
   - Add a way to unregister tools
   - Add validation for tool configurations
   - Create a more pluggable architecture for adding new tools