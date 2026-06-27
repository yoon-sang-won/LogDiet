# Demo

This demo shows LogDiet's core loop:

1. capture full command output locally;
2. show compact evidence to the agent;
3. expand raw evidence only when needed.

## Manual wrapper demo

```sh
logdiet wrap -- sh -c 'echo "running tests"; echo "panic: synthetic failure" >&2; exit 7'
```

Expected:

- exit code is preserved;
- compact output is printed;
- raw logs are stored under `.logdiet/runs/`;
- `raw`, `grep`, and `show` can expand evidence.

```sh
logdiet raw latest
logdiet grep latest "panic"
```

## Setup demo

```sh
logdiet install
eval "$(logdiet env)"
logdiet doctor
```

## Agent setup demos

```sh
logdiet setup codex
logdiet setup claude
logdiet setup cursor
logdiet setup antigravity
```

## Real project demo

In a Go repo:

```sh
logdiet wrap -- go test ./...
```

In a Python repo:

```sh
logdiet wrap -- pytest -q
```

In a Node repo:

```sh
logdiet wrap -- npm test
```

## Notes

This demo uses synthetic failures. LogDiet does not discard raw output. Use `logdiet raw latest` whenever compact evidence is insufficient.
