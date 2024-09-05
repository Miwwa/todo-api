package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/crypto/argon2"
	"strings"
)

type argonOptions struct {
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
	saltLen uint32
}

func argonDefault() argonOptions {
	return argonOptions{
		time:    1,
		memory:  60 * 1024,
		threads: 4,
		keyLen:  32,
		saltLen: 16,
	}
}

// HashPassword returns the hashed version of the provided plaintext password and an error, if any.
// It uses the argon2 crypto algorithm
func HashPassword(plaintextPassword string) (string, error) {
	a := argonDefault()

	salt, err := getRandomBytes(a.saltLen)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(plaintextPassword), salt, a.time, a.memory, a.threads, a.keyLen)

	b64salt := base64.RawStdEncoding.EncodeToString(salt)
	b64hash := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf("$argon2id$v=%d,m=%d,t=%d,p=%d$%s$%s", argon2.Version, a.memory, a.time, a.threads, b64salt, b64hash)

	return encoded, err
}

// ComparePassword compares a plaintext password with a hashed password and returns true if they match.
func ComparePassword(plaintextPassword string, hashedPassword string) (bool, error) {
	a := argonDefault()

	parts := strings.Split(hashedPassword, "$")
	if len(parts) != 5 {
		return false, errors.New("hashed password invalid format")
	}

	var version int
	_, err := fmt.Sscanf(parts[2], "v=%d,m=%d,t=%d,p=%d", &version, &a.memory, &a.time, &a.threads)
	if err != nil {
		return false, err
	}
	if version != argon2.Version {
		return false, errors.New("argon2 versions mismatch")
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[3])
	if err != nil {
		return false, err
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}
	a.keyLen = uint32(len(decodedHash))

	computedHash := argon2.IDKey([]byte(plaintextPassword), salt, a.time, a.memory, a.threads, a.keyLen)

	return subtle.ConstantTimeCompare(decodedHash, computedHash) == 1, nil
}

func getRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
