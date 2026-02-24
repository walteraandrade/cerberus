package config

import (
	"os"
	"path/filepath"

	toml "github.com/pelletier/go-toml/v2"
)

type Argon2Params struct {
	Memory      uint32 `toml:"memory_kib"`
	Iterations  uint32 `toml:"iterations"`
	Parallelism uint8  `toml:"parallelism"`
	SaltLen     int    `toml:"salt_len"`
	KeyLen      uint32 `toml:"key_len"`
}

type GeneratorParams struct {
	PasswordLen    int    `toml:"password_len"`
	DicewareWords  int    `toml:"diceware_words"`
	DicewareSep    string `toml:"diceware_sep"`
	IncludeUpper   bool   `toml:"include_upper"`
	IncludeLower   bool   `toml:"include_lower"`
	IncludeDigits  bool   `toml:"include_digits"`
	IncludeSymbols bool   `toml:"include_symbols"`
}

type Theme struct {
	Primary string `toml:"primary"`
	Accent  string `toml:"accent"`
}

type Config struct {
	DataDir          string          `toml:"-"`
	ClipboardTimeout int             `toml:"clipboard_timeout_sec"`
	LockTimeout      int             `toml:"lock_timeout_sec"`
	Argon2           Argon2Params    `toml:"argon2"`
	Generator        GeneratorParams `toml:"generator"`
	Theme            Theme           `toml:"theme"`
}

func Default() *Config {
	return &Config{
		ClipboardTimeout: 30,
		LockTimeout:      300,
		Argon2: Argon2Params{
			Memory:      64 * 1024,
			Iterations:  3,
			Parallelism: 4,
			SaltLen:     16,
			KeyLen:      32,
		},
		Generator: GeneratorParams{
			PasswordLen:    20,
			DicewareWords:  6,
			DicewareSep:    "-",
			IncludeUpper:   true,
			IncludeLower:   true,
			IncludeDigits:  true,
			IncludeSymbols: true,
		},
		Theme: Theme{
			Primary: "#E53935",
			Accent:  "#FDD835",
		},
	}
}

func DataDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".cerberus"), nil
}

func Load() (*Config, error) {
	cfg := Default()

	dir, err := DataDir()
	if err != nil {
		return cfg, err
	}
	cfg.DataDir = dir

	path := filepath.Join(dir, "config.toml")
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return cfg, nil
	}
	if err != nil {
		return cfg, err
	}

	if err := toml.Unmarshal(data, cfg); err != nil {
		return cfg, err
	}
	cfg.DataDir = dir
	return cfg, nil
}

func (c *Config) Save() error {
	if err := os.MkdirAll(c.DataDir, 0700); err != nil {
		return err
	}
	data, err := toml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(c.DataDir, "config.toml"), data, 0600)
}

func (c *Config) VaultPath() string {
	return filepath.Join(c.DataDir, "vault.enc")
}
