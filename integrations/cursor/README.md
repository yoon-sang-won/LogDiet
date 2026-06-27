# LogDiet for Cursor

This directory contains a Cursor rules file and a hook rewrite template.

## Status

- rules fallback: supported through `logdiet.mdc`
- native adapter: hook rewrite template
- transparent rewrite: template / not verified
- trust required: yes, if your Cursor setup supports command hooks

## Recommended setup

```sh
logdiet init --agent cursor --mode rules
logdiet doctor
```

Use `--mode native` only when you have a trusted Cursor hook path to attach the template to.

## Default fallback behavior

When hook support is unavailable, Cursor should run noisy commands as:

```sh
logdiet wrap -- <command>
```

## How it works

Install `logdiet.mdc` as a rules file where appropriate for your workspace. The hook script is a template unless a tested Cursor hook config is added.

## Verification

Fixture tests verify LogDiet rewrite decisions. Runtime Cursor hook behavior is not claimed.

## Limitations

- Built-in Cursor file/search/editor tools may bypass shell hooks.
- Native hook behavior depends on Cursor configuration outside this repository.
- Raw logs stay local under `.logdiet/runs/`.
