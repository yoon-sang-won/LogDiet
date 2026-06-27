# LogDiet Adapter Contract

Agent-specific adapters are thin delegates.

An adapter may:

1. receive an agent-specific hook payload;
2. extract the shell command;
3. call:

```sh
logdiet hook rewrite --command "COMMAND"
```

4. return the rewritten command in the agent-specific response format.

Adapters must not:

- implement independent rewrite policy;
- summarize command output;
- store raw logs;
- send command contents over the network;
- hide user trust requirements;
- claim unsupported automatic behavior.

## Required fallback

Every agent integration must include rules that tell the agent to use:

```sh
logdiet wrap -- COMMAND
```

when native hooks are unavailable.

## Verification

Each adapter should include:

- README;
- example payload fixtures where possible;
- a smoke test or script;
- clear trust/review notes.
