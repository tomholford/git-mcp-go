# Test Refactoring Plan

## Current Structure

Currently, our test cases are structured as follows:

```go
testCases := []struct {
    name           string
    toolParam      string
    repoPath       string
    writeAccess    bool
    autoApprove    string
    expectedFiles  map[string]string // map of tool to expected file
    expectedErrors []string
    expectedExit   int
}{
    // Test cases...
}
```

And we verify the expectations with separate code blocks for each type of expectation.

## Proposed Refactoring

### 1. Create an Expectations Sub-struct

We'll refactor the test case structure to include an expectations sub-struct:

```go
type expectations struct {
    files       map[string]fileExpectation // map of tool to file expectation
    errors      []string
    exitCode    int
    outputLines []string // optional: additional output lines to check
}

type fileExpectation struct {
    path     string
    content  map[string]interface{} // expected JSON content
    mustExist bool
}

testCases := []struct {
    name        string
    toolParam   string
    repoPath    string
    writeAccess bool
    autoApprove string
    expect      expectations
}{
    // Test cases...
}
```

### 2. Enhance File Content Verification

Instead of just checking if files exist, we'll verify their content against expected values:

```go
// Helper function to verify file expectations
func verifyFileExpectations(t *testing.T, fileExpects map[string]fileExpectation, clineConfigDir, rooCodeConfigDir string) {
    for tool, expect := range fileExpects {
        var filePath string
        if tool == "cline" {
            filePath = filepath.Join(clineConfigDir, expect.path)
        } else if tool == "roo-code" {
            filePath = filepath.Join(rooCodeConfigDir, expect.path)
        }

        // Check if file exists
        fileInfo, err := os.Stat(filePath)
        if os.IsNotExist(err) {
            if expect.mustExist {
                t.Errorf("Expected file %s was not created for %s", filePath, tool)
            }
            continue
        }

        // File exists, verify content if expected
        if expect.content != nil {
            data, err := os.ReadFile(filePath)
            if err != nil {
                t.Fatalf("Failed to read configuration file %s: %v", filePath, err)
            }

            var actualContent map[string]interface{}
            if err := json.Unmarshal(data, &actualContent); err != nil {
                t.Fatalf("Failed to parse JSON in file %s: %v", filePath, err)
            }

            // Verify specific fields in the content
            verifyJSONContent(t, expect.content, actualContent, filePath)
        }
    }
}

// Helper function to verify JSON content
func verifyJSONContent(t *testing.T, expected, actual map[string]interface{}, filePath string) {
    // Verify mcpServers exists
    mcpServers, ok := actual["mcpServers"].(map[string]interface{})
    if !ok {
        t.Errorf("mcpServers not found in settings in file %s", filePath)
        return
    }

    // Verify git exists in mcpServers
    git, ok := mcpServers["git"].(map[string]interface{})
    if !ok {
        t.Errorf("git not found in mcpServers in file %s", filePath)
        return
    }

    // Verify expected fields in git
    expectedGit, ok := expected["mcpServers"].(map[string]interface{})["git"].(map[string]interface{})
    if ok {
        // Verify command
        if expectedCmd, hasCmd := expectedGit["command"]; hasCmd {
            actualCmd, ok := git["command"].(string)
            if !ok || !strings.Contains(actualCmd, expectedCmd.(string)) {
                t.Errorf("Expected command to contain '%s', got: %s in file %s", 
                    expectedCmd, actualCmd, filePath)
            }
        }

        // Verify args
        if expectedArgs, hasArgs := expectedGit["args"]; hasArgs {
            actualArgs, ok := git["args"].([]interface{})
            if !ok {
                t.Errorf("args not found in git in file %s", filePath)
            } else {
                expectedArgsList := expectedArgs.([]interface{})
                if len(actualArgs) < len(expectedArgsList) {
                    t.Errorf("Expected at least %d args, got %d in file %s", 
                        len(expectedArgsList), len(actualArgs), filePath)
                } else {
                    for i, arg := range expectedArgsList {
                        if actualArgs[i] != arg {
                            t.Errorf("Expected arg[%d] to be '%s', got '%s' in file %s", 
                                i, arg, actualArgs[i], filePath)
                        }
                    }
                }
            }
        }

        // Verify disabled
        if expectedDisabled, hasDisabled := expectedGit["disabled"]; hasDisabled {
            actualDisabled, ok := git["disabled"].(bool)
            if !ok || actualDisabled != expectedDisabled.(bool) {
                t.Errorf("Expected disabled to be %v, got: %v in file %s", 
                    expectedDisabled, actualDisabled, filePath)
            }
        }
    }
}
```

### 3. Update Test Cases

We'll update the test cases to use the new structure:

```go
testCases := []struct {
    name        string
    toolParam   string
    repoPath    string
    writeAccess bool
    autoApprove string
    expect      expectations
}{
    {
        name:        "Cline Only",
        toolParam:   "cline",
        repoPath:    "/mock/repo",
        writeAccess: true,
        autoApprove: "allow-read-only",
        expect: expectations{
            files: map[string]fileExpectation{
                "cline": {
                    path:     "cline_mcp_settings.json",
                    mustExist: true,
                    content: map[string]interface{}{
                        "mcpServers": map[string]interface{}{
                            "git": map[string]interface{}{
                                "command":  "git-mcp-go",
                                "args":     []interface{}{"serve", "--repository=/mock/repo", "--write-access=true"},
                                "disabled": false,
                            },
                        },
                    },
                },
            },
            exitCode: 0,
        },
    },
    // More test cases...
}
```

### 4. Update Verification Code

Finally, we'll update the verification code in the test function:

```go
// Verify exit code
if exitCode != tc.expect.exitCode {
    t.Errorf("Expected exit code %d, got %d", tc.expect.exitCode, exitCode)
}

// Verify expected files
verifyFileExpectations(t, tc.expect.files, clineConfigDir, rooCodeConfigDir)

// Verify expected errors in output
output := stdout.String() + stderr.String()
for _, expectedError := range tc.expect.errors {
    if !strings.Contains(output, expectedError) {
        t.Errorf("Expected output to contain '%s', got: %s", expectedError, output)
    }
}

// Verify expected output lines
for _, expectedLine := range tc.expect.outputLines {
    if !strings.Contains(output, expectedLine) {
        t.Errorf("Expected output to contain '%s', got: %s", expectedLine, output)
    }
}
```

## Benefits of This Refactoring

1. **Better organization**: All expectations are grouped together in a dedicated struct
2. **More thorough testing**: We verify the actual content of the files, not just their existence
3. **More flexible**: The new structure makes it easier to add new types of expectations in the future
4. **More maintainable**: The verification logic is encapsulated in dedicated helper functions
5. **More readable**: The test cases are more concise and focused on the test scenario

## Implementation Steps

1. Define the new structs for expectations
2. Implement the helper functions for verification
3. Update the test cases to use the new structure
4. Update the verification code in the test function
5. Run the tests to ensure they still pass