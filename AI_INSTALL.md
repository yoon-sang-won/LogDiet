# AI Install Guide for LogDiet

You are an AI coding agent. The user asked you to install LogDiet and use it to reduce noisy terminal output.

LogDiet is agent-first and CLI-powered.

Your goal is to install the local LogDiet engine, add rules for your agent, and use LogDiet for noisy test/build/git/search output.

## Default installation flow

1. Install the binary:

```sh
go install github.com/yoon-sang-won/LogDiet/cmd/logdiet@latest
```

2. Bootstrap LogDiet for this agent:

```sh
logdiet bootstrap --agent auto
```

3. Verify setup:

```sh
logdiet doctor
```

4. Print session rules:

```sh
logdiet agent-instructions --agent auto
```

5. For the rest of this session, use LogDiet for noisy commands:

```sh
logdiet wrap -- go test ./...
logdiet wrap -- pytest -q
logdiet wrap -- npm test
logdiet wrap -- git diff
logdiet wrap -- rg "TODO"
```

6. Expand compact evidence only when needed:

```sh
logdiet show latest:F1 --around 40
logdiet grep latest "panic"
logdiet raw latest
```

## Agent behavior rules

After installing LogDiet:

- use `logdiet wrap -- <command>` for noisy test, build, typecheck, git, and search commands;
- read compact evidence first;
- use `show` to expand one handle;
- use `grep` to search exact raw output;
- use `raw` only when compact evidence is insufficient;
- do not ask the user to paste full terminal logs unless LogDiet expansion commands are insufficient;
- keep final responses concise:
  - success: changed files, tests run, risks;
  - failure: exact error, likely cause, next command.

## Native hooks

Hooks are optional advanced mode.

Native hooks are optional advanced mode.

If the current agent supports trusted command hooks, you may offer:

```sh
logdiet setup <agent> --mode native
```

But do not claim hooks are active until the agent's trust/review flow confirms them.

## If agent detection fails

Use:

```sh
logdiet bootstrap --agent generic
```

Then follow the printed instructions.
