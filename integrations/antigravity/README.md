# LogDiet for Antigravity

Use `.agents/rules/logdiet.md` as the install target.

## Status

- rules fallback: supported
- native adapter: not verified; see `native-template.md` for conservative notes
- transparent rewrite: unknown
- trust required: unknown

## Recommended setup

```sh
logdiet init --agent antigravity --mode rules
logdiet doctor
```

## Default fallback behavior

Antigravity should run noisy commands as:

```sh
logdiet wrap -- <command>
```

## How it works

This package is conservative: rules fallback first, with no claim of automatic command rewrite unless a tested hook path is added later.

## Verification

Rules installation is covered by LogDiet tests. Runtime native adapter behavior is not verified.

## Limitations

- Automatic command rewrite is not guaranteed.
- Native hook support is unknown.
- Raw logs stay local under `.logdiet/runs/`.
