# CHANGELOG

## [v0.6.0](https://github.com/tab/cmt/releases/tag/v0.6.0)

### Features
- **feat:** Use gpt-4.1-nano as default model for faster, cost-effective GPT-4.1

### Chore
- **chore:** Change GPT API error logging from error to debug mode

### Dependencies
- **chore(deps):** Bump go.uber.org/fx from 1.23.0 to 1.24.0
- **chore(deps):** Bump go.uber.org/mock from 0.5.1 to 0.5.2
- **chore(deps):** Bump codecov/codecov-action from 5.4.2 to 5.4.3

### Documentation
- **docs:** Add configuration section with example YAML in README

## [v0.5.0](https://github.com/tab/cmt/releases/tag/v0.5.0)

### Features
- **feat:** Refactor application architecture with FX framework for dependency injection and lifecycle management

### Chore
- **chore:** Update dependencies (httpmock from v1.3.1 to v1.4.0, mock from v0.5.0 to v0.5.1)
- **chore:** Bump goreleaser/goreleaser-action from 6.1.0 to 6.3.0

## [v0.4.2](https://github.com/tab/cmt/releases/tag/v0.4.2)

### Documentation
- **docs:** Added godoc comments

### Chore
- **chore:** Bumped dependency github.com/go-resty/resty/v2 from 2.16.2 to 2.16.5
- **chore:** Updated GitHub workflows for linter, staticcheck, tests, and coverage
- **chore:** Added staticcheck action to CI workflow
- **chore:** Added CODEOWNERS files for repository management

## [v0.4.1](https://github.com/tab/cmt/releases/tag/v0.4.1)

### Chore
- **chore:** Add Makefile to run common Go development tasks
- **chore:** Add golangci-lint
- **chore:** Update GitHub Actions workflow
- **chore:** Bump github.com/stretchr/testify from 1.9.0 to 1.10.0
- **chore:** Bump github.com/go-resty/resty/v2 from 2.15.3 to 2.16.2
- **chore:** Bump go.uber.org/mock from 0.4.0 to 0.5.0

## [v0.4.0](https://github.com/tab/cmt/releases/tag/v0.4.0)

### Features
- **feat:** Interactive commit message editing

### Refactor
- **refactor:** Replace hardcoded timeout with configurable timeout context

### Chore
- **chore:** Update codecov configuration
- **chore:** Update Go file formatting rules

## [v0.3.0](https://github.com/tab/cmt/releases/tag/v0.3.0)

### Features
- **feat:** Add codecov configuration file (`codecov.yaml`) for coverage reporting.
- **feat:** Refactor and enhance command structure for better functionality, including a new command approach and improved error handling.
- **feat:** Add changelog command.

### Refactor
- **refactor:** Simplify input reading process in command handling.

### Documentation
- **docs:** Update README to enhance clarity, add new features, and improve formatting.
- **docs:** Revise commit message examples for better clarity and add changelog generation instructions.

### Chore
- **chore:** Add CHANGELOG.md for tracking releases.

## [v0.2.0](https://github.com/tab/cmt/releases/tag/v0.2.0)

### Features
- **feat(loader):** Add loading indicator functionality

### Bug Fixes
- **fix(cmd/main):** Handle error when reading user response.

### Refactor
- **refactor(loader):** Modified Start and Stop methods for the Loader to ensure safe concurrent access.

### Chore
- **chore(github/workflows):** Add Golang CI linter to checks workflow.%
- **chore(github/workflows):** Add Dependabot configuration file for version updates.
- **chore(deps):** Bump codecov/codecov-action from 3 to 4.

## [v0.1.0](https://github.com/tab/cmt/releases/tag/v0.1.0)

### Features
- **feat(core):** Add optional prefix for commit messages with '--prefix' flag
- **feat(github/workflows):** Add GitHub Actions workflow for master branch
- **feat(github/workflows):** Add release workflow for building and distributing binaries

### Refactor
- **refactor(flags):** Update PrintVersion and PrintHelp methods to accept io.Writer for output

### Chore
- **chore(github/workflows):** Add release workflow for building and distributing binaries
