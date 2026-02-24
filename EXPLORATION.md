# Cerberus — Exploration Notes

## Series Context

| Project | Stack | Purpose |
|---|---|---|
| mr-argus | Bun/Ink | Claude Code monitor |
| arachne | Rust/ratatui | Git graph viewer |
| homer | Go/Bubbletea | Commit chronicle via LLM |
| apollo | Go/Bubbletea | Commit review triage |

### Shared Go Conventions (homer + apollo)

- Go 1.25 + Bubbletea v1.3 + Bubbles + Lipgloss
- SQLite (modernc, pure Go, WAL mode)
- `internal/` layout: `config/`, `db/`, `tui/`, `style/`, domain packages
- Layered arch: domain types never import persistence/UI
- TOML config, `~/.projectname/` data dir
- vim keybindings (j/k/g/G), fsnotify watcher w/ 300ms debounce
- CLAUDE.md as architecture contract

---

## Go Password Manager Landscape

No mature Go + Bubbletea password manager exists. This is an open niche.

| Project | Stars | Cipher | KDF | Storage | Key Insight |
|---|---|---|---|---|---|
| gopass | 6.7k | GPG or age | scrypt | one file/secret | leaks metadata (filenames) |
| gokey (Cloudflare) | 2.4k | AES-256-GCM | custom | vaultless | nothing to steal, inflexible |
| goldwarden | 966 | Bitwarden protocol | PBKDF2/Argon2id | Bitwarden sync | memguard, SSH agent |
| passwall-server | 764 | AES | — | PostgreSQL/SQLite | server-side, less relevant |
| masterkey | 279 | XChaCha20-Poly1305 | Argon2id | single .db | all metadata encrypted |
| kure | 167 | AES-256-GCM | Argon2id per-record | bbolt | memguard, per-entry isolation |
| go-hash | 102 | AES-256 | Argon2 | custom binary | academic-informed, separate enc+MAC keys |
| passgo | — | XSalsa20-Poly1305 (NaCl) | scrypt | directory tree | asymmetric design, HMAC integrity |
| go2fa | — | RSA | N/A | JSON | only Go+Bubbletea secret tool found |

---

## Crypto Stack (Recommended)

| Layer | Choice | Why |
|---|---|---|
| KDF | Argon2id | OWASP #1, memory-hard, resists GPU/ASIC |
| Cipher | XChaCha20-Poly1305 | 192-bit nonce (safe random), no AES-NI dependency, AEAD |
| Memory | memguard | allocates outside Go heap, mlock, guard pages |
| Clipboard | atotto/clipboard + timed clear | checksum-verify before clearing |
| Randomness | crypto/rand | CSPRNG for all salts/nonces/keys |

### Argon2id Params (local vault unlock)

| Parameter | Value | Notes |
|---|---|---|
| Variant | Argon2id | resists side-channel and GPU |
| Memory | 64 MiB | tune to ~300ms on target hardware |
| Iterations | 3 | |
| Parallelism | 4 | |
| Salt | 16 bytes | random, stored with ciphertext |
| Output | 32 bytes | key material |

Go: `golang.org/x/crypto/argon2.IDKey(password, salt, 3, 64*1024, 4, 32)`

### Cipher Notes

**XChaCha20-Poly1305** (preferred):
- 192-bit nonce — collision-safe with random generation
- Constant-time, no timing side-channels
- Go: `golang.org/x/crypto/chacha20poly1305.NewX(key)`

**AES-256-GCM** (alternative):
- 96-bit nonce — risky for random generation at scale
- Hardware-accelerated (AES-NI)
- Go stdlib: `crypto/cipher`

### KDF Hierarchy (OWASP)

1. Argon2id (recommended)
2. scrypt (if Argon2id unavailable)
3. PBKDF2 (FIPS compliance only)
4. bcrypt (legacy)

---

## Key Architecture Decisions

### Vault Format Options

| Option | Pros | Cons |
|---|---|---|
| Single encrypted file | simple, all metadata hidden, easy backup | full decrypt per read, no partial update |
| BBolt + app-level encryption | ACID, partial reads, pure Go | key names may leak metadata |
| SQLite + encrypted fields | familiar (series pattern), SQL queries | schema/table names leak structure |

### Key Hierarchy Patterns

**Direct KDF** (simple):
`master_password + salt → Argon2id → vault_key`
Changing password = re-encrypt everything.

**Envelope encryption** (Bitwarden/KDBX model, recommended):
`master_password + salt → Argon2id → wrapping_key → decrypts random vault_key`
Changing password = only re-wrap the vault_key.

**Asymmetric per-entry** (passgo/NaCl):
`master_password → KDF → symmetric key → encrypts master private key`
Each entry encrypted to master public key. Allows public-key operations without master password.

### Memory Security in Go

Go's GC copies memory freely — `mlock()` on a Go slice is unreliable.

**memguard** (`github.com/awnumar/memguard`):
- Allocates via `mmap()` outside Go heap
- `mlock()` prevents swap
- Guard pages + canaries detect overflow
- Encrypts data at rest in RAM (XSalsa20Poly1305)
- Two types: `Enclave` (encrypted/sealed) and `LockedBuffer` (active use)

Key practices:
- Never store passwords as Go `string` (immutable, can't zero)
- Use `[]byte` and zero explicitly
- Avoid `fmt.Sprintf` with secrets

### Clipboard Security

- Copy + auto-clear after 30-45s (configurable)
- Checksum-verify clipboard still holds the secret before clearing
- atotto/clipboard has limited Wayland support — may need wl-clipboard fallback
- Clipboard managers (copyq, etc.) persist history; no perfect OS-level fix

### How Major Vaults Work

**KeePass KDBX4**: Outer header (KDF params, cipher ID) → HMAC integrity check → encrypted payload (GZip + ChaCha20). Key: composite hash → Argon2d → SHA-256(master_seed || transformed_key). All metadata encrypted.

**Bitwarden**: master_password + email → PBKDF2/Argon2id → master_key → HKDF → stretched_key → decrypts protected_symmetric_key → that key encrypts vault items (AES-256-CBC + HMAC-SHA-256). Zero-knowledge cloud sync.

**1Password**: Three-level key hierarchy (derived key → overview key + master key → per-item keys). All metadata encrypted since OPVault format. Secret Key + master password prevents server-side brute-force.

**age**: Header with recipient stanzas + HKDF → ChaCha20-Poly1305 payload. Passphrase mode uses scrypt (not Argon2id).

---

## Go Libraries Reference

### Standard + Extended

| Package | Use |
|---|---|
| `crypto/rand` | CSPRNG for nonces/keys/salts |
| `crypto/subtle` | Constant-time comparison |
| `golang.org/x/crypto/argon2` | Argon2id KDF |
| `golang.org/x/crypto/chacha20poly1305` | XChaCha20-Poly1305 (.NewX) |
| `golang.org/x/crypto/nacl/secretbox` | XSalsa20-Poly1305 |
| `golang.org/x/crypto/scrypt` | scrypt KDF |
| `golang.org/x/crypto/pbkdf2` | PBKDF2 (Bitwarden compat) |

### Third-Party

| Library | Role |
|---|---|
| `github.com/awnumar/memguard` | Secure memory (mlock, guard pages) |
| `filippo.io/age` | File encryption (production-ready) |
| `github.com/alexedwards/argon2id` | Argon2id wrapper with sane defaults |
| `github.com/atotto/clipboard` | Cross-platform clipboard |
| `go.etcd.io/bbolt` | Embedded ACID KV store |
| `github.com/charmbracelet/bubbletea` | TUI framework |
| `github.com/charmbracelet/bubbles` | TUI components |
| `github.com/charmbracelet/lipgloss` | TUI styling |
| `github.com/pelletier/go-toml/v2` | TOML config |

---

## Proposed Structure

```
~/.cerberus/
  config.toml          # plaintext (Argon2 params, clipboard timeout, theme)
  vault.enc            # single encrypted blob

internal/
  config/              # TOML config
  vault/               # domain types (Entry, Vault, Category)
  crypto/              # Argon2id KDF, XChaCha20 encrypt/decrypt, key wrapping
  storage/             # vault file read/write
  clipboard/           # copy + timed clear
  generator/           # password/passphrase generation
  tui/                 # Bubbletea model, screens (unlock, list, detail, edit)
  style/               # lipgloss theme
```

---

## Open Questions

1. Single encrypted file vs bbolt? (single file hides all metadata; bbolt better for growth)
2. Sync/backup strategy? (git-backed like gopass? plain file copy? none?)
3. Password generation — random chars only, or also diceware passphrases?
4. TOTP/2FA support — store TOTP seeds alongside passwords?
5. Import/export — support KeePass KDBX or CSV import?
6. Agent/daemon for caching — cache decrypted vault key in background process, or re-enter master password each session?
7. Wayland clipboard — atotto/clipboard or wl-clipboard directly?
