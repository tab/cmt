# cmt

**cmt** is a [Claude Code](https://claude.com/claude-code) skill that generates [Conventional Commit](https://www.conventionalcommits.org/) messages from your staged Git changes.

It reads your diff, matches the commit style already established in the repository's `git log`, and hands you a ready-to-paste message — so your history stays clear and consistent without the busywork.

## What it does

- **Conventional Commits** — generates a `<type>(<scope>): <description>` title, plus an optional body when there's something non-obvious worth recording.
- **Matches your repo's style** — learns scope vocabulary and conventions from recent `git log` history.
- **Staged-first** — describes staged changes when present, otherwise falls back to unstaged and untracked files.
- **Never commits for you** — outputs the message only and offers a ready-to-run `git commit` command. You stay in control.
- **Optional hint** — steer the type, scope, or framing with a free-text argument.

## Requirements

- [Claude Code](https://claude.com/claude-code)
- Git

No API key, build step, or configuration file required — the skill runs inside Claude Code.

## Installation

**cmt** runs as a Claude Code skill — install it by cloning this repository into a skills directory.

**Personal** (recommended) — available in every project on your machine:

```sh
git clone https://github.com/tab/cmt.git ~/.claude/skills/cmt
```

**Project** — scoped to the current repository only:

```sh
git clone https://github.com/tab/cmt.git .claude/skills/cmt
```

The skill lives in `SKILL.md` at the repository root, so cloning into a folder named `cmt` inside a skills directory is all Claude Code needs to discover it. Start a new Claude Code session and `/cmt` will be available.

### Updating

Run `git pull` in the directory you cloned into:

```sh
git -C ~/.claude/skills/cmt pull
```

## Usage

Stage the changes you want to commit:

```sh
git add .
```

Then invoke the skill in Claude Code:

```
/cmt
```

You'll get a Conventional Commits message in a copy-pasteable block, followed by a ready-to-run commit command. Review it, then run the command yourself if you're happy.

### Hint

Pass an optional free-text hint to steer the result — a type, a scope, or a framing:

```
/cmt this is a fix, not a feature
/cmt scope: api
```

The hint is used as written; everything else is inferred from the diff.

## Conventions

The generated message follows the [Conventional Commits](https://www.conventionalcommits.org/) specification and Tim Pope's 50/72 rule for the subject line. Type is chosen by the intent of the change (`feat`, `fix`, `chore`, `refactor`, `docs`, `test`, `perf`, `build`, `ci`, `style`, `revert`), and a scope is always included. See [`SKILL.md`](SKILL.md) for the full rule set.

## License

Distributed under the MIT License. See [`LICENSE`](LICENSE) for more information.

## Acknowledgements

- [Conventional Commits](https://www.conventionalcommits.org/) for the commit message specification.
