/*
Package Name: common
File Name: hasher.go
Abstract: Hasher provides helper functions for encoding/decoding
strings with the argon2 algorithm.

Author: Alejandro Modroño <alex@sureservice.es>
Created: 07/12/2023
Last Updated: 07/24/2023

# MIT License

# Copyright 2023 Alejandro Modroño Vara

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package common

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// ======== TYPES ========

// Extracted from https://www.alexedwards.net/blog/how-to-hash-and-verify-passwords-with-argon2-in-go
//
// The Argon2 algorithm accepts a number of configurable parameters:
//   - Memory: The amount of memory used by the algorithm (in kibibytes).
//   - Iterations: The number of iterations (or passes) over the memory.
//   - Parallelism: The number of threads (or lanes) used by the algorithm.
//   - Salt length: Length of the random salt. 16 bytes is recommended for
//     password hashing.
//   - Key length: Length of the generated key (or password hash).
//     16 bytes or more is recommended.
//
// The memory and iterations parameters control the computational cost of
// hashing the password. The higher these figures are, the greater the cost
// of generating the hash. It also follows that the greater the cost will be
// for any attacker trying to guess the password.
//
// But there's a balance that you need to strike. As you increase the cost,
// the time taken to generate the hash also increases. If you're generating
// the hash in response to a user action (like signing up or logging in to
// a website) then you probably want to keep the runtime to less than 500ms
// to avoid a negative user experience.
//
// If the Argon2 algorithm is running on a machine with multiple cores, then
// one way to decrease the runtime without reducing the cost is to increase
// the parallelism parameter. This controls the number of threads that the
// work is spread across. There's an important thing to note here though:
// changing the value of the parallelism parameter changes the output of the
// algorithm. So — for example — running Argon2 with a parallelism parameter
// of 2 will result in a different password hash to running it with a parallelism
// parameter of 4.
type Parameters struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

// ======== NAMESPACES ========

// hasherT is used for creating a namespace
type hasherT struct{}

// the Hasher namespace
var Hasher hasherT

// ======== ERRORS ========
var (
	InvalidHashException         = errors.New("The encoded hash is not in the correct format.")
	IncompatibleVersionException = errors.New("Incompatible version of argon2.")
)

// ======== PUBLIC METHODS ========

// Hasher.hash returns a hash from a string.
func (hasherT) Hash(from_string string) (encodedHash string, err error) {
	// Set the parameters to be used by the argon2 algorithm.
	params := &Parameters{
		memory:      64 * 1024,
		iterations:  3,
		parallelism: 2,
		saltLength:  16,
		keyLength:   32,
	}

	// Generate a random salt to be appended to the hash.
	salt, err := generateRandomBytes(params.saltLength)
	if err != nil {
		return "", err
	}

	// Generate the hash
	hash := argon2.IDKey(
		[]byte(from_string),
		salt,
		params.iterations,
		params.memory,
		params.parallelism,
		params.keyLength,
	)

	// Base64 encode the salt and hashed password.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Return a string using the standard encoded hash representation.
	encodedHash = fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		params.memory,
		params.iterations,
		params.parallelism,
		b64Salt,
		b64Hash,
	)

	return encodedHash, nil
}

// Hasher.compare compares a plaintext string with a hash and returns whether
// they match or not.
func (hasherT) Compare(plaintext, encodedHash string) (matches bool, err error) {
	// Decodes the salt and hash from base64 and extracts the parameters,
	// salt and derived key from the encoded password hash.
	params, salt, hash, err := decode(encodedHash)
	if err != nil {
		return false, err
	}

	// Derive the key from the other password using the same parameters.
	otherHash := argon2.IDKey(
		[]byte(plaintext),
		salt,
		params.iterations,
		params.memory,
		params.parallelism,
		params.keyLength,
	)

	// Check that the contents of the hashed passwords are identical. Note
	// that we are using the subtle.ConstantTimeCompare() function for this
	// to help prevent timing attacks.
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}

// ======== PRIVATE METHODS ========

// GenerateRandomBytes generates a random salt that will be appended
// to the hash.
func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}

	return b, nil
}

// Hasher.decodeHash decodes the salt and hash and extracts the parameters
// from an argon2 hash.
func decode(encodedHash string) (params *Parameters, salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, InvalidHashException
	}

	var version int
	if _, err = fmt.Sscanf(vals[2], "v=%d", &version); err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, IncompatibleVersionException
	}

	params = &Parameters{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &params.memory, &params.iterations, &params.parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	params.saltLength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	params.keyLength = uint32(len(hash))

	return params, salt, hash, nil
}
