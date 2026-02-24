package generator

import (
	"crypto/rand"
	"math/big"
	"strings"
)

const (
	lowerChars  = "abcdefghijklmnopqrstuvwxyz"
	upperChars  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digitChars  = "0123456789"
	symbolChars = "!@#$%^&*()-_=+[]{}|;:,.<>?"
)

type PasswordOpts struct {
	Length         int
	IncludeUpper   bool
	IncludeLower   bool
	IncludeDigits  bool
	IncludeSymbols bool
}

func DefaultPasswordOpts() PasswordOpts {
	return PasswordOpts{
		Length:         20,
		IncludeUpper:   true,
		IncludeLower:   true,
		IncludeDigits:  true,
		IncludeSymbols: true,
	}
}

func Password(opts PasswordOpts) (string, error) {
	var charset string
	if opts.IncludeLower {
		charset += lowerChars
	}
	if opts.IncludeUpper {
		charset += upperChars
	}
	if opts.IncludeDigits {
		charset += digitChars
	}
	if opts.IncludeSymbols {
		charset += symbolChars
	}
	if charset == "" {
		charset = lowerChars + upperChars + digitChars
	}

	result := make([]byte, opts.Length)
	for i := range result {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[idx.Int64()]
	}
	return string(result), nil
}

type PassphraseOpts struct {
	Words     int
	Separator string
}

func DefaultPassphraseOpts() PassphraseOpts {
	return PassphraseOpts{
		Words:     6,
		Separator: "-",
	}
}

func Passphrase(opts PassphraseOpts) (string, error) {
	words := make([]string, opts.Words)
	for i := range words {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(effWordlist))))
		if err != nil {
			return "", err
		}
		words[i] = effWordlist[idx.Int64()]
	}
	return strings.Join(words, opts.Separator), nil
}
