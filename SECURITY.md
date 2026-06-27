# Security

LogDiet does not make network calls or send telemetry. It stores raw command output locally under `.logdiet/runs/`.

Raw logs may contain secrets, tokens, credentials, private file paths, or proprietary code. Do not commit `.logdiet/runs/`.

Do not commit `.logdiet/backup/` either. Review logs before sharing them outside your machine.

If you report a security issue, avoid including sensitive raw logs unless necessary.
