# LogDiet for Gemini

This directory contains Gemini instruction and hook templates for LogDiet.

## Status

- rules fallback: supported through `GEMINI.md`
- native adapter: hook rewrite template
- transparent rewrite: template / not verified
- trust required: yes, if your Gemini CLI setup supports command hooks

## Recommended setup

```sh
logdiet init --agent gemini --mode rules
logdiet doctor
```

Use `--mode native` only after reviewing the template and confirming a supported hook path.

## Default fallback behavior

When native rewrite is unavailable, use:

```sh
logdiet wrap -- <command>
```

## How it works

`GEMINI.md` provides rules fallback. The hook script is a template unless a tested Gemini hook config is added.

## Verification

Fixture tests verify central LogDiet rewrite decisions. Runtime Gemini hook behavior is not verified here.

## Limitations

- Gemini CLI hook behavior is environment-specific.
- Non-shell agent tools may bypass LogDiet.
- Raw logs stay local under `.logdiet/runs/`.
