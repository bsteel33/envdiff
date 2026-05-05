# envdiff

> Utility to diff `.env` files across environments and flag missing or mismatched keys.

---

## Installation

```bash
go install github.com/yourusername/envdiff@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envdiff.git
cd envdiff
go build -o envdiff .
```

---

## Usage

Compare two `.env` files and see what's missing or mismatched:

```bash
envdiff .env.development .env.production
```

**Example output:**

```
MISSING in .env.production:
  - DEBUG_MODE
  - LOCAL_DB_URL

MISMATCHED keys (present in both, different values):
  ~ API_BASE_URL
  ~ LOG_LEVEL

OK: 12 keys match across both files.
```

### Flags

| Flag | Description |
|------|-------------|
| `--keys-only` | Compare only key names, ignore values |
| `--quiet` | Exit with non-zero status if diff found, no output |
| `--format json` | Output results as JSON |

---

## Why envdiff?

Keeping `.env` files in sync across environments is error-prone. `envdiff` makes it easy to catch missing secrets or misconfigured variables before they cause issues in staging or production.

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

---

## License

[MIT](LICENSE)