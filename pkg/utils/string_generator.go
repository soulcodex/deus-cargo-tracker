package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

type StringGenerator interface {
	Generate(size int) (string, error)
}

type RandomStringGenerator struct{}

func NewRandomStringGenerator() *RandomStringGenerator {
	return &RandomStringGenerator{}
}

// Generate returns a cryptographically secure random string
func (rsg *RandomStringGenerator) Generate(size int) (string, error) {
	return randomString(size)
}

// MustGenerate is a convenience method that panics
// if the generation fails.
func (rsg *RandomStringGenerator) MustGenerate(size int) string {
	str, err := rsg.Generate(size)
	if err != nil {
		panic(err)
	}
	return str
}

type FixedStringGenerator struct {
	value string
}

func NewFixedStringGenerator(value string) *FixedStringGenerator {
	return &FixedStringGenerator{
		value: value,
	}
}

// Generate returns a fixed string of the specified size.
func (fsg *FixedStringGenerator) Generate(size int) (string, error) {
	if size < 0 || size == 0 {
		return "", fmt.Errorf("invalid size %d: must be positive", size)
	}
	if size > len(fsg.value) {
		return fsg.value, nil
	}

	return fsg.value[:size], nil
}

// MustGenerate is a convenience method that panics
// if the generation fails.
func (fsg *FixedStringGenerator) MustGenerate(size int) string {
	str, err := fsg.Generate(size)
	if err != nil {
		panic(err)
	}
	return str
}

// randomString generates a cryptographically secure random string of the specified size.
func randomString(size int) (string, error) {
	if size < 0 || size == 0 {
		return "", fmt.Errorf("invalid size %d: must be positive", size)
	}
	const seed = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	letters, buffer := []rune(seed), make([]rune, size)

	for i := range size {
		// Generate a random index for the rune slice (letters)
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", fmt.Errorf("failed to generate random index: %w", err)
		}

		// Assign a random letter to the buffer at index i
		buffer[i] = letters[num.Int64()]
	}
	return string(buffer), nil
}
