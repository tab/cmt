# Development Guide

## Project Overview

**cmt** is a command-line utility that generates [Conventional Commit](https://www.conventionalcommits.org/) messages using OpenAI's GPT models based on staged Git changes. It automates the process of writing clear and structured commit messages, enhancing Git workflow and ensuring consistency across projects.

### Features
- Automated commit message generation following Conventional Commits specification
- Interactive TUI for commit approval and editing
- Custom prefix support for task IDs and issue numbers
- Console-based changelog generation from commit history
- Integration with OpenAI GPT models

## Architecture Overview

### Core Components

1. **Entry Point** (`cmd/`)
  - `main.go` - Application bootstrap with FX dependency injection and configuration loading
  - Early command routing for help, version, and changelog (handled before FX)

2. **Core Packages** (`internal/`)
  - **app/** - Main application container and lifecycle management
    - **cli/** - Bubble Tea TUI implementation for commit flow and command routing
      - **helpers/** - Help, version, and changelog console commands
      - **workflow/** - GPT + git orchestration layer (used by both TUI and console commands)
    - **git/** - Git operations (diff, log, commit) with command executor abstraction
    - **gpt/** - OpenAI GPT client (using sashabaranov/go-openai)
    - **errors/** - Application-specific error definitions
  - **config/** - Configuration loading, parsing, and data structures
    - **logger/** - Structured logging with zerolog

### Key Interfaces and Abstractions

1. **cli.UI** – Interface for terminal user interface operations:
   ```go
   type UI interface {
       Run(args []string) error
   }
   ```

2. **cli.gitClient** – Narrow interface used by the UI after the TUI exits:
   ```go
   type gitClient interface {
       Commit(ctx context.Context, msg string) (string, error)
   }
   ```

3. **workflow.Service** – High-level orchestration for generating commits and changelogs:
   ```go
   type Service interface {
       GenerateCommit(ctx context.Context, prefix string) (CommitResult, error)
       GenerateChangelog(ctx context.Context, between string) (ChangelogResult, error)
   }
   ```

4. **git.Client** – Git operations (provider interface used by the workflow service):
   ```go
   type Client interface {
       Diff(ctx context.Context) (string, error)
       Status(ctx context.Context) (string, error)
       Log(ctx context.Context, opts []string) (string, error)
       Commit(ctx context.Context, message string) (string, error)
   }
   ```

5. **git.Executor** – Command execution abstraction for git operations:
   ```go
   type Executor interface {
       Run(ctx context.Context, name string, args ...string) *exec.Cmd
   }
   ```

6. **gpt.Client** – GPT model client (provider interface consumed by the workflow service):
   ```go
   type Client interface {
       FetchCommitMessage(ctx context.Context, diff string) (string, error)
       FetchChangelog(ctx context.Context, commits string) (string, error)
   }
   ```

7. **logger.Logger** – Structured logging interface using zerolog

### CLI Package Structure

- **interface.go** – UI interface and consumer-defined interfaces (gitClient)
- **ui.go** – UI implementation, Bubble Tea program initialization, commit execution after TUI exit
- **module.go** – FX module definition for dependency injection (wires workflow service + UI)
- **workflow/** – Commit/changelog orchestration independent of Bubble Tea
- **helpers/** – Help and version rendering utilities
  - **help.go** - Help text generation and rendering
  - **version.go** - Version display and command detection
- **components/** - Reusable TUI components (file tree, commit list, styles, rendering)
- **model/** – Bubble Tea model implementation
  - **types.go** – Model struct, state definitions, flow/action enums, edit modes
  - **init.go** – Model initialization and argument parsing
  - **update.go** – Delegates to state controllers and handles app-wide events (window size, Ctrl+C)
  - **state_controller_*.go** – Focused handlers for fetch/review/edit/changelog/log states
  - **view.go** – Bubble Tea View() function for rendering UI states
  - **cmd.go** – Async command execution (fetchCommit, fetchChangelog)
  - **constants.go** – Command string constants (changelog, prefix flags, etc.)
  - **logs.go** – Log entry formatting for log viewer

### TUI State Machine

The TUI uses a 5-state finite state machine for commit workflow. Each state is backed by a dedicated controller, keeping transitions small and testable:

- **StateInit** - Initial state (not actively used, transitions immediately to StateFetch)
- **StateFetch** - Fetching data (git diff + GPT generation)
- **StateViewCommit** - Reviewing generated commit message with file tree
- **StateEditCommit** - Editing commit message with vim-like modal editing (normal/insert mode)
- **StateViewLogs** - Showing application logs in viewport (accessible via 'l' key)
- **StateExit** - Exit state (errors stored in model.Err)

**Model Structure**: The Bubble Tea model contains:
- State tracking (State, UserFlow, UserAction, Mode, PreviousState)
- Dependencies (Ctx, Cfg, Workflow, Log)
- Bubble Tea components (Spinner, Textarea, Viewport)
- Data (Content, FileTree, Err, Prefix)
- UI dimensions (Width, Height)
- Log buffering (LogBuffer) for displaying application logs
- Controller registry (lazily initialised map of state → controller)

**Flow Types**:
- `FlowCommit` - Generate commit message (only flow in TUI)

**Actions**:
- `ActionNone` - No action taken yet
- `ActionAccept` - User accepted commit (triggers actual git commit after TUI exits)
- `ActionEdit` - User wants to edit
- `ActionRefresh` - User wants to regenerate
- `ActionCancel` - User cancelled

**Edit Modes**:
- `EditModeNormal` - Vim normal mode (navigation and commands)
- `EditModeInsert` - Vim insert mode (text entry)

Key bindings follow simple single-letter patterns:
- Review (StateViewCommit): `a` (accept), `e` (edit), `r` (refresh), `l` (logs), `q` (quit), arrow/vim keys (scroll file tree)
- Edit Normal (StateEditCommit): `i` (insert), `o` (new line), `h/j/k/l/arrows` (move), `c` (delete word), `x` (delete char), `d` (delete line), `s` (save), `q` (quit without save)
- Edit Insert (StateEditCommit): `ESC` (normal mode), `Ctrl+S` (save)
- Logs (StateViewLogs): `q` or `l` (return to previous state), arrow/vim keys (scroll)

### Execution Flow

1. **Application Bootstrap** (`cmd/main.go` → `internal/app/app.go`)
   - Check for early commands first (handled without FX)
   - `help/--help/-h` → Render help text directly via `helpers.RenderHelp()`
   - `version/--version/-v` → Render version info directly via `helpers.RenderVersion()`
   - `changelog [range]` → Generate and output changelog via `helpers.RenderChangelog()`
     - Loads config and creates dependencies manually
     - Uses workflow service to generate changelog
     - Outputs directly to console
   - All other commands (commit flow): Load configuration from `cmt.yaml` (or use defaults)
   - Validate OPENAI_API_KEY environment variable during config load
   - Initialize FX dependency injection container
   - Supply configuration to all modules
   - Start FX lifecycle and run application

2. **Dependency Injection** (FX modules)
   - `logger.Module` - Logger creation
   - `git.Module` - Git client and executor
   - `gpt.Module` - GPT client (using sashabaranov/go-openai)
   - `cli.Module` - UI implementation
   - `app.Module` - Application container with lifecycle management

3. **Application Lifecycle** (`internal/app/app.go`)
   - FX OnStart: Launch UI.Run() in goroutine
   - UI.Run() parses args and starts Bubble Tea program
   - After TUI exits, if commit was accepted, execute git commit
   - Shutdown FX application

4. **TUI Workflow** (`internal/app/cli/`)
   - Parse args to determine flow (commit/changelog) and extract prefix/between parameters
   - Initialize Bubble Tea model with dependencies via `model.New()` (wires workflow.Service + controller registry)
   - Start state machine in StateFetch (StateInit immediately transitions on first spinner tick)
   - Transition through states based on controller decisions and async results
   - Handle all user interactions through `Model.Update` delegating to controllers
   - Render current state through View() function
   - Support live log viewing via 'l' key (when format: ui in config)

5. **Commit Message Generation Flow** (in TUI)
   - StateFetch: Get git diff + git status, send diff to GPT, parse file tree, transition to StateViewCommit
   - StateViewCommit: Display split view (file tree left, commit message right), user can accept (a), edit (e), refresh (r), view logs (l), or quit (q)
   - If edit: StateViewCommit → StateEditCommit → StateViewCommit (on save)
   - If accept: StateViewCommit → StateExit, commit happens after TUI exits in ui.Run()
   - If refresh: StateViewCommit → StateFetch → StateViewCommit (regenerate message)
   - If logs: StateViewCommit → StateViewLogs (toggle log viewer)

6. **Changelog Generation Flow** (console-only, no TUI)
   - Handled in `cmd/main.go` before FX initialization
   - Loads config and creates dependencies (git, gpt, workflow, logger)
   - Calls workflow.GenerateChangelog() with optional range parameter
   - Outputs markdown-formatted changelog directly to console
   - No interactive features - one-shot generation and output

### Configuration Capabilities

Configuration is loaded from `cmt.yaml` in the current directory. If not found, default values are used.

**Configuration Structure:**
```yaml
api:
  retry_count: 3      # Number of retry attempts for API requests
  timeout: 60s        # Timeout duration for API requests

model:
  name: gpt-4.1-nano  # OpenAI model to use
  max_tokens: 500     # Maximum tokens for model response
  temperature: 0.7    # Controls randomness (0.0-2.0)

logging:
  format: console     # Output format: console, json, or ui (for log buffer viewing)
  level: info         # Log level: debug, info, warn, error
```

**Default Values** (defined in `internal/config/config.go`):
- Model: `gpt-4.1-nano`
- Max Tokens: `500`
- Temperature: `0.7`
- Retry Count: `3`
- Timeout: `60s`
- Log Level: `info`
- Log Format: `console`
- Version: `0.7.0` (current version)

**Required Environment Variables:**
- `OPENAI_API_KEY` - OpenAI API key for GPT access

### Testing Patterns

1. **Mock Generation**
  - Uses `go.uber.org/mock` for interface mocking
  - Generated mocks with `//go:generate` directives in source files
  - Mock files named with `_mock.go` suffix (e.g., `types_mock.go`, `interface_mock.go`)

2. **Test Structure**
  - Table-driven tests with subtests using testify
  - Comprehensive error case coverage
  - Output capturing for CLI command testing
  - Mock expectation setup and verification
  - Entry point testing with extracted testable functions
  - Integration test skipping for complex application lifecycle scenarios

3. **Table Tests with Mocks Pattern**
  - Mocks are created once at the test function level
  - Each test case has a `before func()` that sets up mock expectations
  - Test data and mock expectations are co-located in the same test case
  - `tt.before()` is called just before executing the test logic
  - Example structure:
    ```go
    func Test_Example(t *testing.T) {
        ctrl := gomock.NewController(t)
        defer ctrl.Finish()

        mockDep := NewMockDependency(ctrl)
        subject := &Implementation{dep: mockDep}

        tests := []struct {
            name   string
            before func()
            input  string
            expect bool
        }{
            {
                name: "success case",
                input: "test-input",
                before: func() {
                    mockDep.EXPECT().Method("test-input").Return(nil)
                },
                expect: true,
            },
        }

        for _, tt := range tests {
            t.Run(tt.name, func(t *testing.T) {
                tt.before()
                result := subject.TestMethod(tt.input)
                assert.Equal(t, tt.expect, result)
            })
        }
    }
    ```

## Primary Guidelines

- provide brutally honest and realistic assessments of requests, feasibility, and potential issues. no sugar-coating. no vague possibilities where concrete answers are needed.
- always operate under the assumption that the user might be incorrect, misunderstanding concepts, or providing incomplete/flawed information. critically evaluate statements and ask clarifying questions when needed.
- don't be flattering or overly positive. be honest and direct.
- we work as equal partners and treat each other with respect as two senior developers with equal expertise and experience.
- prefer simple and focused solutions that are easy to understand, maintain and test.
- use table-driven tests ONLY when testing multiple scenarios with different inputs/outputs; for single test cases, use plain test functions instead of table tests with one entry
- table tests are appropriate when you have 2+ test cases with meaningful variations in input/output/behavior


## Build, Lint and Test Commands

```bash
# Build binary
go build -o cmd/cmt ./cmd

# Run all tests (always run from the top level)
make test

# Lint code (always run from the top level)
make lint

# Coverage report (always run from the top level)
make coverage

# Format code
go fmt ./...

# Run completion sequence (formatting, linting and testing)
go fmt ./... && make lint && make vet && make test
```

**IMPORTANT:** NEVER commit without running tests, formatter and linters for the entire codebase!

## Important Workflow Notes

- always run tests, linter BEFORE committing anything
- run formatting, code generation, linting and testing on completion
- never commit without running completion sequence
- run tests and linter after making significant changes to verify functionality
- IMPORTANT: never put into commit message any mention of Claude or Claude Code
- IMPORTANT: never include "Test plan" sections in PR descriptions
- do not add comments that describe changes, progress, or historical modifications
- comments should only describe the current state and purpose of the code, not its history or evolution
- use `go:generate` for generating mocks, never modify generated files manually
- mocks are generated with `go.uber.org/mock` and stored alongside source files
- after important functionality added, update README.md accordingly
- when merging master changes to an active branch, make sure both branches are pulled and up to date first
- don't leave commented out code in place
- if working with github repos use `gh`
- avoid multi-level nesting
- avoid deeply nested conditionals (more than 3 levels)
- never use goto
- prefer early returns to reduce nesting, but else/else if are acceptable when they improve readability
- write tests in compact form by fitting struct fields to a single line (up to 130 characters)
- before any significant refactoring, ensure all tests pass and consider creating a new branch
- when refactoring, editing, or fixing failed tests:
  - do not redesign fundamental parts of the code architecture
  - if unable to fix an issue with the current approach, report the problem and ask for guidance
  - focus on minimal changes to address the specific issue at hand
  - preserve the existing patterns and conventions of the codebase

## Code Style Guidelines

### Import Organization
- Organize imports in the following order:
  1. Standard library packages first (e.g., "fmt", "context")
  2. A blank line separator
  3. Third-party packages
  4. A blank line separator
  5. Project imports (e.g., "cmt/internal/*")
- Example:
  ```go
  import (
      "context"
      "fmt"
      "os"

      "github.com/rs/zerolog"
      "go.uber.org/fx"

      "cmt/internal/config"
  )
  ```

### Error Handling
- return errors to the caller rather than using panics
- use descriptive error messages that help with debugging
- use error wrapping: `fmt.Errorf("failed to process request: %w", err)`
- check errors immediately after function calls
- return early when possible to avoid deep nesting
- for functions that return multiple values including errors, handle both the primary result and the error appropriately
- when logging errors, include contextual information: `c.log.Error().Err(err).Msgf("Failed to run profile '%s'", profile)`

### Variable Naming
- use descriptive camelCase names for variables and functions
  - good: `serviceProcess`, `dependencyGraph`, `profileConfig`
  - bad: `sp`, `x`, `temp1`
- be consistent with abbreviations
- local scope variables can be short (e.g., "cfg" instead of "configuration")

### Function Parameters
- group related parameters together logically
- use descriptive parameter names that indicate their purpose
- consider using parameter structs for functions with many (4+) parameters
- if function returns 3 or more results, consider wrapping in result/response struct
- if function accepts 3 or more input parameters, consider wrapping in request/input struct (but never add context to struct)

### Documentation
- all exported functions, types, and methods must have clear godoc comments
- begin comments with the name of the element being documented
- include usage examples for complex functions
- document any non-obvious behavior or edge cases
- follow standard Go comment conventions: complete sentences should start with capital letters and end with periods
- godoc comments for exported functions should start with the function name
- keep internal comments concise and clear

### Code Structure
- keep code modular with focused responsibilities
- limit file sizes to 300-500 lines when possible
- group related functionality in the same package
- use interfaces to define behavior and enable mocking for tests
- keep code minimal and avoid unnecessary complexity
- don't keep old functions for imaginary compatibility
- interfaces should be defined on the consumer side (idiomatic Go)
- aim to pass interfaces but return concrete types when possible
- consider nested functions when they simplify complex functions
- always place `ctx context.Context` as the first field/parameter in struct/method definitions for consistency
- always place `cfg *config.Config` as the second field/parameter in struct/method definitions for consistency (when `ctx` is present) or first (when `ctx` is absent)
- always place `log logger.Logger` as the last field/parameter in struct/method definitions for consistency

### Code Layout
- keep cyclomatic complexity under 30
- function size preferences:
  - keep functions focused on a single responsibility
  - break down large functions (100+ lines) into smaller, logical pieces
  - avoid functions that are too small if they reduce readability
- keep lines readable; while gofmt doesn't enforce line length, consider breaking very long lines for clarity
- manage conditional complexity:
  - for many discrete values, prefer switch statements over long if-else-if chains
  - use early returns to reduce nesting depth when appropriate
  - extract complex conditions into well-named boolean functions or variables
  - use context structs or functional options instead of multiple boolean flags
- for CLI command processing, use switch statements with multiple conditions per case (e.g., `case cmd == "help" || cmd == "--help" || cmd == "-h":`)
- when handling default values, check for empty strings and provide sensible defaults (e.g., `if profile == "" { profile = config.DefaultProfile }`)
- for functions that need to be testable, separate return values from system calls: return exit codes and errors instead of calling os.Exit() directly

### Testing
- write thorough tests with descriptive names (e.g., `Test_Runner_ResolvesComplexDependencies`)
- prefer subtests or table-based tests, using testify
- use table-driven tests ONLY when testing multiple scenarios (2+ test cases) with different inputs/outputs; for single test cases, use plain test functions instead of table tests with one entry
- table-driven tests for testing multiple cases with the same logic
- test both success and error scenarios
- mock external dependencies to ensure unit tests are isolated and fast
- aim for at least 80% code coverage
- keep tests compact but readable
- if test has too many subtests, consider splitting it to multiple tests
- never disable tests without a good reason and approval
- important: never update code with special conditions to just pass tests
- don't create new test files if one already exists matching the source file name
- add new tests to existing test files following the same naming and structuring conventions
- don't add comments before subtests, t.Run("description") already communicates what test case is doing
- never use godoc-style comments for test functions
- for main package testing, extract testable functions from main() and runApp() to enable unit testing
- skip integration tests that would cause hanging or require subprocess execution (e.g., os.Exit(), long-running FX apps)
- when testing CLI applications, use simple skip statements for complex integration scenarios to maintain test suite stability
- for mocking external dependencies:
  - create a local interface in the package that needs the mock
  - generate mocks using `go.uber.org/mock` with: `//go:generate mockgen -source=file.go -destination=file_mock.go`
  - the mock should be located alongside the source file
  - always use mockgen-generated mocks, not testify mock
- for testing functions that can fail due to external dependencies (like config loading), use `t.Skip()` with descriptive messages rather than failing the test
- use descriptive test names that explain the scenario being tested (e.g., "No arguments - default profile", "Run command with --run= (empty profile defaults to default profile)")
- when testing CLI return values, test both exit codes and error conditions separately in table test fields like `expectedExit` and `expectedError`
- for testable code extraction: separate business logic from system calls (os.Exit, os.Args) by creating internal functions that can be unit tested
- always test multiple command variations (e.g., "help", "--help", "-h") to ensure all aliases work correctly

## Git Workflow

### After merging a PR
```bash
# switch back to the master branch
git checkout master

# pull latest changes including the merged PR
git pull

# delete the temporary branch (might need -D for force delete if squash merged)
git branch -D feature-branch-name
```

## Commonly Used Libraries
- dependency injection: `go.uber.org/fx`
- configuration: `github.com/spf13/viper`
- logging: `github.com/rs/zerolog`
- testing: `github.com/stretchr/testify`
- mock generation: `go.uber.org/mock`
- TUI framework: `github.com/charmbracelet/bubbletea`
- TUI components: `github.com/charmbracelet/bubbles` (spinner, textarea, viewport)
- styling: `github.com/charmbracelet/lipgloss`
- OpenAI client: `github.com/sashabaranov/go-openai`

## Formatting Guidelines
- always use `go fmt` for code formatting
- run `go generate` for mock generation

## Logging Guidelines
- use structured logging with zerolog
- never use fmt.Printf for logging, only log methods

## Available Commands

### Commit Message Generation
```bash
cmt                           # Generate commit message for staged changes
cmt prefix <PREFIX>           # Generate commit with custom prefix
cmt --prefix <PREFIX>         # Alternative prefix syntax
cmt -p <PREFIX>               # Short form prefix syntax
```

### Changelog Generation
```bash
cmt changelog                 # Generate changelog from initial commit to HEAD (console output)
cmt changelog v1.0.0..v1.1.0  # Generate changelog between tags (console output)
cmt changelog SHA1..SHA2      # Generate changelog between commits (console output)
cmt changelog SHA..HEAD       # Generate changelog from commit to HEAD (console output)
```

### Help and Version
```bash
cmt help                      # Display help (interactive TUI)
cmt --help                    # Display help
cmt -h                        # Display help
cmt version                   # Display version information
cmt --version                 # Display version
cmt -v                        # Display version
```
