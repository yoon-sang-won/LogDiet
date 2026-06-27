# Antigravity Native Adapter Notes

LogDiet does not currently claim a verified Antigravity native hook adapter.

Use rules fallback first:

```sh
logdiet wrap -- COMMAND
```

If Antigravity exposes a trusted shell command hook in your environment, adapt a thin local hook that delegates to:

```sh
logdiet hook rewrite --command "$COMMAND"
```

Do not enable automatic rewrite until you have reviewed the hook and verified it inside Antigravity.
