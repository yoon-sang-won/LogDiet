# Security

LogDiet does not make network calls or send telemetry. It stores raw command output locally under `.logdiet/runs/`.

Raw logs may contain secrets, tokens, credentials, private file paths, or proprietary code. Do not commit `.logdiet/runs/`.

If you report a security issue, avoid including sensitive raw logs unless necessary.
