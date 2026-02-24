# Cerberus — Architecture Contract

## What
Go + Bubbletea TUI password manager. Local-only, single encrypted vault file.

## Stack
- Go 1.25, Bubbletea v1.3, Bubbles, Lipgloss
- XChaCha20-Poly1305 (cipher), Argon2id (KDF), memguard (secure memory)
- Envelope encryption: master_password → KDF → wrapping_key → decrypts random vault_key
- TOML config (`~/.cerberus/config.toml`), single vault file (`~/.cerberus/vault.enc`)

## Layout
```
cmd/cerberus/         main entrypoint
internal/
  config/             TOML config (pelletier/go-toml/v2)
  vault/              domain types (Entry, Vault, Category) — NO crypto/IO imports
  crypto/             Argon2id KDF, XChaCha20 encrypt/decrypt, envelope key wrapping
  storage/            vault file read/write (atomic temp+rename)
  clipboard/          copy + timed clear (Wayland/X11 detection)
  generator/          password (random chars) + passphrase (diceware, dash separator)
  tui/                Bubbletea model, screens, components
  style/              lipgloss theme — red #E53935, gold #FDD835
```

## Conventions
- `internal/` layout, layered arch: domain types never import persistence/UI
- vim keybindings (j/k/g/G)
- Functional style where possible
- Minimal comments — only where logic isn't self-evident
- Tests alongside source files (`_test.go`)
- crypto/rand for ALL randomness (never math/rand)
- Never store secrets as Go `string` — use `[]byte` and zero explicitly
- memguard LockedBuffer for active key material, Enclave for at-rest

## Config Defaults
- Argon2id: memory=64MiB, iterations=3, parallelism=4, salt=16B, output=32B
- Clipboard timeout: 30s
- Password length: 20
- Diceware separator: dash
- Diceware words: 6

## Theme
- Primary: #E53935 (red)
- Accent: #FDD835 (gold)
- Dark terminal background assumed
