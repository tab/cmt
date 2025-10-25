# Changelog Refactoring: Remove TUI, Keep Console Output

## Overview
Refactor the changelog command to use simple console output (like help/version commands) instead of TUI, while keeping the commit command in TUI.

## Analysis Summary

### Current Architecture
- **TUI State Machine**: 7 states (Init, Fetch, ViewCommit, EditCommit, ViewChangelog, ViewLogs, Exit)
- **Flow Types**: FlowCommit and FlowChangelog
- **Changelog TUI**: Split-panel view with commit list (left) and changelog content (right)
- **Help/Version**: Direct console output handled in cmd/main.go before FX initialization

### Components Used by Changelog TUI
1. **Model fields**:
   - `UserFlow` (FlowChangelog)
   - `ChangelogViewport` (viewport for changelog content)
   - `CommitList` (*components.CommitList)
   - `Between` (string for range parameter)

2. **State**: StateViewChangelog

3. **Controller**: state_controller_changelog.go

4. **Components**:
   - `components/commits.go` - CommitEntry, CommitList, ParseCommitLog, RenderCommits, ConvertCommitLogForGPT
   - `components/ui.go` - RenderChangelogHints (key bindings)

5. **Workflow**: workflow.Service.GenerateChangelog (keep for console command)

6. **View**: viewChangelog() in model/view.go

### New Design
- Handle `changelog` command in cmd/main.go before FX initialization
- Create `helpers/changelog.go` with IsChangelogCmd and RenderChangelog
- RenderChangelog calls workflow.GenerateChangelog directly
- Output plain text changelog to console
- No TUI, no split panels, no interactive features

## Implementation Plan

### Phase 1: Add Console Changelog Command
- [ ] Create `internal/app/cli/helpers/changelog.go`
  - IsChangelogCmd(cmd string) bool
  - RenderChangelog(args []string) - orchestrate changelog generation and output
  - Need direct access to git and gpt clients without FX
  - Handle range parsing from args

- [ ] Update `cmd/main.go`
  - Add changelog case to handleCommands switch
  - Call helpers.IsChangelogCmd and helpers.RenderChangelog

### Phase 2: Remove Changelog TUI Components
- [ ] Remove `internal/app/cli/model/state_controller_changelog.go`
- [ ] Remove `FlowChangelog` from `internal/app/cli/model/types.go`
- [ ] Remove `StateViewChangelog` from `internal/app/cli/model/types.go`
- [ ] Remove `ChangelogViewport` field from Model struct
- [ ] Remove `CommitList` field from Model struct
- [ ] Remove `Between` field from Model struct

- [ ] Update `internal/app/cli/model/init.go`
  - Remove FlowChangelog from parseArgs logic
  - Remove changelog argument parsing
  - Simplify to only handle commit flow with optional prefix

- [ ] Update `internal/app/cli/model/view.go`
  - Remove StateViewChangelog case
  - Remove viewChangelog(), renderCommitsList(), renderChangelogContent()

- [ ] Update `internal/app/cli/model/state_controller_fetch.go`
  - Remove FlowChangelog case handling
  - Remove CommitList handling in Result

- [ ] Update `internal/app/cli/ui.go`
  - Remove FlowChangelog flow name logging
  - Simplify Run() - only handle commit flow

### Phase 3: Clean Up Unused Components
- [ ] Update `internal/app/cli/components/commits.go`
  - Remove CommitEntry, CommitList structs
  - Remove ParseCommitLog, RenderCommits, RenderCommitEntry functions
  - Keep ConvertCommitLogForGPT (used by workflow/changelog generation)

- [ ] Update `internal/app/cli/components/ui.go`
  - Remove RenderChangelogHints function

- [ ] Check `internal/app/cli/components/style.go`
  - Remove commit-related styles if unused: CommitHashStyle, CommitSubjectStyle, CommitDateStyle, CommitAuthorStyle, CommitEntryStyle

- [ ] Update `internal/app/cli/model/types.go`
  - Remove Result.CommitList field (if no longer used)
  - Keep Result.Content, Result.FileTree, Result.Err for commit flow

### Phase 4: Update Tests
- [ ] Remove `internal/app/cli/model/state_controller_changelog_test.go` (if exists)
- [ ] Update `internal/app/cli/model/init_test.go` - remove changelog flow tests
- [ ] Update `internal/app/cli/model/state_controller_fetch_test.go` - remove changelog cases
- [ ] Update `internal/app/cli/model/view_test.go` (if exists) - remove changelog view tests
- [ ] Update `internal/app/cli/ui_test.go` - remove changelog flow tests
- [ ] Add tests for `helpers/changelog.go`
- [ ] Update `cmd/main_test.go` - add changelog command test case
- [ ] Update `internal/app/cli/components/commits_test.go` - remove tests for deleted functions

### Phase 5: Update Documentation
- [ ] Update `README.md`
  - Remove TUI description for changelog
  - Remove split-panel changelog description
  - Remove changelog key bindings (c, r, l, q, arrows)
  - Update to show console output example
  - Simplify changelog usage section

- [ ] Update `CLAUDE.md`
  - Remove FlowChangelog from architecture overview
  - Remove StateViewChangelog from state machine description
  - Remove changelog TUI workflow description
  - Update component list (remove CommitList, etc.)
  - Update CLI package structure
  - Update execution flow (remove changelog TUI flow)
  - Update available commands section

- [ ] Update `internal/app/cli/helpers/help.go`
  - Keep changelog command in help text
  - Simplify description (remove TUI key bindings)

### Phase 6: Final Cleanup
- [ ] Run `go generate ./...` to regenerate mocks
- [ ] Run `go fmt ./...`
- [ ] Run `make lint`
- [ ] Run `make vet`
- [ ] Run `make test`
- [ ] Verify changelog command works: `./cmt changelog`
- [ ] Verify changelog with range works: `./cmt changelog v0.6.0..HEAD`
- [ ] Verify commit command still works: `./cmt`

## Files to Modify

### New Files
- `internal/app/cli/helpers/changelog.go`

### Modified Files
- `cmd/main.go`
- `internal/app/cli/model/types.go`
- `internal/app/cli/model/init.go`
- `internal/app/cli/model/view.go`
- `internal/app/cli/model/state_controller_fetch.go`
- `internal/app/cli/ui.go`
- `internal/app/cli/components/commits.go`
- `internal/app/cli/components/ui.go`
- `internal/app/cli/components/style.go`
- `internal/app/cli/helpers/help.go`
- `README.md`
- `CLAUDE.md`

### Deleted Files
- `internal/app/cli/model/state_controller_changelog.go`
- Tests for deleted functionality

### Test Files to Update
- `cmd/main_test.go`
- `internal/app/cli/model/init_test.go`
- `internal/app/cli/model/state_controller_fetch_test.go`
- `internal/app/cli/ui_test.go`
- `internal/app/cli/components/commits_test.go`

## Risk Assessment

### Low Risk
- Adding new console changelog command (isolated change)
- Removing unused TUI states and controllers (dead code removal)

### Medium Risk
- Modifying parseArgs logic (affects command routing)
- Removing fields from Model struct (could affect serialization, if any)

### Testing Strategy
- Unit tests for new helpers/changelog.go
- Integration test: run changelog command and verify output
- Regression test: ensure commit flow still works
- Manual testing of all command variations

## Dependencies
- Keep `workflow.Service.GenerateChangelog` - reuse for console command
- Keep `components.ConvertCommitLogForGPT` - needed by workflow
- Remove all TUI-specific changelog components
