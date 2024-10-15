# CHANGELOG

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
