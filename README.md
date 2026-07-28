# reclaim-cli
A CLI for reclaim.ai

Set `RECLAIM_API_KEY` to a personal Reclaim API key before use.

## Interactive or not

Every command still prompts for anything you do not pass as a flag, so running
`reclaim create --title "..."` behaves exactly as it always has.

Pass the flags instead and nothing is prompted for, which is what makes the CLI
usable from a script or an agent. Prompting is skipped automatically when stdin is
not a terminal, and `--non-interactive` forces the same behaviour when it is. In
either case a missing flag is an immediate error rather than a prompt nobody can
answer.

Data goes to stdout and diagnostics to stderr, so `--json` output is safe to pipe.

## `create`

```bash
reclaim create --title "Review the auth migration" \
  --mins 30 --min-chunk 30 --priority P3 \
  --start "in 1 week" --due 2026-08-14 --json
```

| Flag | |
|---|---|
| `--title` | required |
| `--notes` | free text; a better home for links than the title |
| `--mins` | total length: 15, 30, 45, 60, 90, 120, 180, 240 |
| `--min-chunk` | shortest block it may be split into: 15, 30, 45, 60 |
| `--priority` | P1, P2, P3, P4 |
| `--start` | when to start working on it; snoozes the task until then |
| `--due` | when it is due; omit for no due date |
| `--json` | write the created task to stdout as one line of JSON |
| `--non-interactive` | never prompt |

`--mins` and `--min-chunk` are validated against the lists above, so a value the
interactive picker would not have offered is rejected rather than turned into an
odd chunk count.

## `snooze`

```bash
reclaim snooze --id 13354745 --until "in 3 days"
reclaim snooze --title "auth migration" --until 2026-09-01
```

`--id` and `--title` both skip the picker; `--title` matches on a case-insensitive
substring of an open task's title and fails if more than one task matches.
`--until` defaults to `in 1 day`.

## Time values

`--start`, `--due` and `--until` all accept:

| Form | Example |
|---|---|
| `now` | `--start now` |
| `in <n> hours\|days\|weeks` | `--until "in 3 days"` |
| a date | `--due 2026-08-14` |
| an RFC3339 timestamp | `--due 2026-08-14T17:30:00Z` |

Bare dates are read in the local timezone. A bare `--start` or `--until` date
begins at midnight; a bare `--due` date is due at 23:59, since a task due at
midnight would be overdue for the whole of its due date.

## Other commands

| | |
|---|---|
| `reclaim dedupe` | merge tasks sharing a title, and delete tasks for GitLab MRs that are closed or merged. Needs `GITLAB_URL` and `GITLAB_TOKEN`. |
| `reclaim archive` | archive completed tasks. `--auto-archive-age` archives anything finished more than that many days ago without asking; `--max-count` bounds how many you are asked about. `--non-interactive` archives on age alone. |
| `reclaim meeting` | book a meeting through a scheduling link |
| `reclaim version` | print the version |
