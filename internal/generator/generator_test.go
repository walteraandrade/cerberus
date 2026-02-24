package generator

import (
	"strings"
	"testing"
	"unicode"
)

func TestPassword_DefaultLength(t *testing.T) {
	pw, err := Password(DefaultPasswordOpts())
	if err != nil {
		t.Fatal(err)
	}
	if len(pw) != 20 {
		t.Fatalf("expected 20 chars, got %d", len(pw))
	}
}

func TestPassword_CustomLength(t *testing.T) {
	opts := DefaultPasswordOpts()
	opts.Length = 32
	pw, err := Password(opts)
	if err != nil {
		t.Fatal(err)
	}
	if len(pw) != 32 {
		t.Fatalf("expected 32 chars, got %d", len(pw))
	}
}

func TestPassword_OnlyLower(t *testing.T) {
	opts := PasswordOpts{Length: 50, IncludeLower: true}
	pw, err := Password(opts)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range pw {
		if !unicode.IsLower(c) {
			t.Fatalf("expected only lowercase, got %c", c)
		}
	}
}

func TestPassword_Uniqueness(t *testing.T) {
	pw1, _ := Password(DefaultPasswordOpts())
	pw2, _ := Password(DefaultPasswordOpts())
	if pw1 == pw2 {
		t.Fatal("two passwords should not be equal")
	}
}

func TestPassphrase_DefaultWords(t *testing.T) {
	pp, err := Passphrase(DefaultPassphraseOpts())
	if err != nil {
		t.Fatal(err)
	}
	words := strings.Split(pp, "-")
	if len(words) != 6 {
		t.Fatalf("expected 6 words, got %d", len(words))
	}
}

func TestPassphrase_CustomSeparator(t *testing.T) {
	opts := PassphraseOpts{Words: 4, Separator: "."}
	pp, err := Passphrase(opts)
	if err != nil {
		t.Fatal(err)
	}
	words := strings.Split(pp, ".")
	if len(words) != 4 {
		t.Fatalf("expected 4 words, got %d", len(words))
	}
}

func TestPassphrase_Uniqueness(t *testing.T) {
	pp1, _ := Passphrase(DefaultPassphraseOpts())
	pp2, _ := Passphrase(DefaultPassphraseOpts())
	if pp1 == pp2 {
		t.Fatal("two passphrases should not be equal")
	}
}
