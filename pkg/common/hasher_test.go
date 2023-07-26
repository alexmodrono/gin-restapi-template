/*
Package Name: common
File Name: hasher_test.go
Abstract: Tests for the hasher functions
Author: Alejandro Modroño <alex@sureservice.es>
Created: 07/26/2023
Last Updated: 07/26/2023

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
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHasher_HashAndCompare(t *testing.T) {
	// Test case 1: Valid hash and comparison
	password := "mySecretPassword"

	// Generate the hash for the password
	hashedPassword, err := Hasher.Hash(password)
	require.NoError(t, err)

	// Test that the password matches the hash
	matches, err := Hasher.Compare(password, hashedPassword)
	assert.True(t, matches)
	assert.NoError(t, err)

	// Test case 2: Invalid comparison
	wrongPassword := "wrongPassword"

	// Test that the wrong password does not match the hash
	matches, err = Hasher.Compare(wrongPassword, hashedPassword)
	assert.False(t, matches)
	assert.NoError(t, err)
}

func TestHasher_Decode(t *testing.T) {
	// Test case 1: Valid encoded hash
	encodedHash := "$argon2id$v=19$m=65536,t=3,p=2$Zm9v$MTIzNDU2"

	params, salt, hash, err := decode(encodedHash)
	require.NoError(t, err)
	assert.NotNil(t, params)
	assert.Equal(t, uint32(65536), params.memory)
	assert.Equal(t, uint32(3), params.iterations)
	assert.Equal(t, uint8(2), params.parallelism)
	assert.Equal(t, uint32(3), params.saltLength)
	assert.Equal(t, uint32(6), params.keyLength)

	expectedSalt, _ := base64.RawStdEncoding.Strict().DecodeString("Zm9v")
	expectedHash, _ := base64.RawStdEncoding.Strict().DecodeString("MTIzNDU2")

	assert.Equal(t, expectedSalt, salt)
	assert.Equal(t, expectedHash, hash)

	// Test case 2: Invalid encoded hash
	invalidHash := "invalidHash"

	_, _, _, err = decode(invalidHash)
	assert.Error(t, err)
}

func TestHasher_Decode_InvalidHash(t *testing.T) {
	// Test case 3: Invalid encoded hash with less than 6 parts
	encodedHash := "$argon2id$v=19$m=65536,t=3,p=2$Zm9v"

	_, _, _, err := decode(encodedHash)
	assert.EqualError(t, err, InvalidHashException.Error())

	// Test case 4: Invalid encoded hash with incompatible version
	encodedHash = "$argon2id$v=18$m=65536,t=3,p=2$Zm9v$MTIzNDU2"

	_, _, _, err = decode(encodedHash)
	assert.EqualError(t, err, IncompatibleVersionException.Error())
}
