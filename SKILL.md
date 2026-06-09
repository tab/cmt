---
name: cmt
description: Generate a Conventional Commits message (title + optional body) for the current changes. Use when the user invokes /cmt, or asks to "write a commit message", "generate commit", "make a commit message". Generates the message only — never runs git commit.
allowed-tools: Bash(git status:*), Bash(git diff:*), Bash(git log:*), Bash(git branch:*), Bash(git rev-parse:*), Read
---

# cmt

Generate a commit message that follows the [Conventional Commits](https://www.conventionalcommits.org/) spec, matching the style already established in the repository's `git log`.

**Output the message only. NEVER run `git commit`** — the user commits themselves. If they want to commit, hand them the command to run via `!`.

## Invocation

```
/cmt [hint]
```

`[hint]` — optional free-text steer (e.g. a type, scope, or "this is a fix not a feature"). Use it as written; otherwise infer everything from the diff.

## Step 1 — Read the changes

Run in one batch:

```bash
git status --short
git diff --cached --stat
git diff --cached
git log --pretty=format:"%s" -30
```

- If **anything is staged**, the staged diff is the source of truth — describe only what's staged.
- If **nothing is staged**, fall back to unstaged + untracked: `git diff` and `git status --short` (read new files with the Read tool if their purpose isn't obvious from the name).
- If **nothing to commit at all** — no staged changes, no unstaged changes, no untracked files — say so in one line and stop. Don't invent a message.
- State which set you used in one line so the user knows what the message covers.
- Use the recent `git log` titles to learn the repository's scope vocabulary — keep terms consistent across commits.

If the diff is large, read the `--stat` first to find the dominant change, then the full diff for the relevant files.

## Step 2 — Compose the message

### Title — `<type>(<scope>): <Description>`

- **type** — one of: `feat`, `fix`, `chore`, `refactor`, `docs`, `test`, `perf`, `build`, `ci`, `style`, `revert`.
  Pick by intent of the change, not the file touched (e.g. a bug fix in a test file is still `fix` only if it fixes behavior; adding tests is `test`; use `revert` when the change undoes a prior commit).
- **scope** — the affected component, lowercase. **Always include a scope** — never emit a scopeless `type: description`. **Prefer a scope already used in recent `git log`** — keep the vocabulary consistent. Typical patterns:
  - a sub-component or package name when the change is localized (e.g. `ui`, `log`, `worker`, `config`)
  - a top-level area when the change is broader (e.g. `docs`, `deps`, `ci`)
  - if the change genuinely spans many unrelated areas, that's a signal to split the commit — propose splitting instead of dropping the scope
- **description** — imperative mood, **capitalized first letter**, **no trailing period**, concise. Aim ≤ 50 chars for the title line; hard cap 72 (Tim Pope's 50/72 convention).
  - Good: `fix(log): Fix lifecycle logs`
  - Good: `feat(ui): Add service aside panel`
  - Good: `refactor(ui): Hard-wrap env values in aside panel`
  - Bad: `fix(log): fix lifecycle logs` (lowercase description)
  - Bad: `fix(log): Fixed lifecycle logs.` (past tense, trailing period)
- **Breaking change** — append `!` before the colon (`feat(api)!: ...`) and add a `BREAKING CHANGE:` footer.

### Body — optional, only when it adds non-obvious information

- **Skip the body** for trivial, self-evident single-purpose changes — most commits ship with just the title.
- **Add a body** when there are decisions, constraints, or non-obvious behavior worth recording. Two styles are commonly used — match whichever fits the change:
  - **Prose lines** — each line a complete thought, capitalized, no trailing period:
    ```
    Fix API and relay shutdown messages to debug so a normal info-level run no longer prints 'API server shutting down', 'API server stopped' and 'Server stopped' on every exit
    ```
  - **Bullet list** — `- ` prefix, capitalized, no trailing period:
    ```
    - Move close(a.done) to after flushing sentry
    - Ensures proper shutdown signaling
    ```
- Separate body from title with one blank line.
- Wrap at ~72 chars (prose) — bullets can flow long if natural.
- **No ticket footer, no "Test plan" section, no mention of Claude / Claude Code / Codex / any AI assistant** in the message.

### House rules

- One logical change per commit — if the diff mixes unrelated concerns, say so and propose splitting it into multiple messages rather than one bloated message.
- Clear, short, meaningful. Provide context, not noise.

## Step 3 — Present the result

Output the message in a single copy-pasteable code block:

```
fix(log): Fix lifecycle logs

Fix API and relay shutdown messages to debug so a normal info-level run no longer prints 'API server shutting down', 'API server stopped' and 'Server stopped' on every exit
```

Then, offer the ready-to-run command for the user (do not run it). Pick the right form for the message you generated:

- **Title only** — single `-m`:
  > To commit: `!git commit -m "fix(log): Fix lifecycle logs"`
- **Title + body** — heredoc, so multi-line bodies and bullet lists survive intact:
  ```
  !git commit -F- <<'EOF'
  fix(log): Fix lifecycle logs

  Fix API and relay shutdown messages to debug so a normal info-level run no longer prints 'API server shutting down', 'API server stopped' and 'Server stopped' on every exit
  EOF
  ```

Never use stacked `-m "title" -m "body"` — each `-m` becomes a separate paragraph and bullet lists get mangled.

If you proposed splitting into multiple commits, present each message in its own block with its own command.
