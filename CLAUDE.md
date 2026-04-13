@AI_CONTEXT.md

## Adding a new command

1. Create `src/app/commands/my_cmd.py` with `handle(args: str) -> None`
2. Register in `src/app/core.py`: `router.register("my")(my_cmd.handle)`
3. Add tests in `tests/test_commands.py`

## Architecture rules

- `workflow/scripts/entry.py` is the **only** file Alfred executes. No business logic here.
- `src/alfred/` contains **only** Alfred SDK helpers — no application logic.
- Commands call services. Services call clients. Never skip layers.
- All `output()` calls go through `alfred.response.output()`.
- Always wrap `main()` in `safe_run()` — unhandled exceptions = blank Alfred.

## Testing rules

- Test `src/app/` (commands, services, clients) — these are pure Python.
- Mock external API calls in `ApiClient`. Never make real HTTP calls in tests.
- `conftest.py` sets Alfred env vars to tmp dirs automatically.
- `alfred/` SDK helpers are tested in `tests/test_alfred.py`.

## Performance target

Script Filter response < 100ms.
Use `alfred.cache.Cache` for any network calls.
Cache TTL default: 300s (5 min).

## Dependency management

Runtime dependencies → `requirements.txt` → vendored into `workflow/vendor/`
Dev dependencies → `pyproject.toml [project.optional-dependencies.dev]`

Keep runtime deps minimal. Every package adds to workflow size.

## Pre-coding checklist

Before starting work, confirm if any of these are unclear:

- Goal / completion criteria
- Language / framework / version constraints
- New code vs existing code modification
- Whether tests are required
- Scope of impact

Do **not** ask about code style, file placement, or minor implementation
details — follow existing patterns in the codebase.

## Error handling stance

When an error occurs: **diagnose root cause → explain fix plan → implement**.
Never retry the same failing command. Never skip hooks (`--no-verify`).
