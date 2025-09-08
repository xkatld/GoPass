package crypto

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "testpassword123"
	
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}
	
	if hash == "" {
		t.Error("Hash should not be empty")
	}
	
	if hash == password {
		t.Error("Hash should not equal original password")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "testpassword123"
	wrongPassword := "wrongpassword"
	
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}
	
	// Test correct password
	if !CheckPasswordHash(password, hash) {
		t.Error("Password verification should succeed with correct password")
	}
	
	// Test wrong password
	if CheckPasswordHash(wrongPassword, hash) {
		t.Error("Password verification should fail with wrong password")
	}
}

func TestEncryptDecrypt(t *testing.T) {
	plaintext := "This is a secret message"
	key := GenerateKey("test-password")
	
	// Test encryption
	ciphertext, err := Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("Failed to encrypt: %v", err)
	}
	
	if ciphertext == "" {
		t.Error("Ciphertext should not be empty")
	}
	
	if ciphertext == plaintext {
		t.Error("Ciphertext should not equal plaintext")
	}
	
	// Test decryption
	decrypted, err := Decrypt(ciphertext, key)
	if err != nil {
		t.Fatalf("Failed to decrypt: %v", err)
	}
	
	if decrypted != plaintext {
		t.Errorf("Decrypted text should equal original plaintext. Got: %s, Expected: %s", decrypted, plaintext)
	}
}

func TestEncryptDecryptWithWrongKey(t *testing.T) {
	plaintext := "This is a secret message"
	key1 := GenerateKey("password1")
	key2 := GenerateKey("password2")
	
	// Encrypt with key1
	ciphertext, err := Encrypt(plaintext, key1)
	if err != nil {
		t.Fatalf("Failed to encrypt: %v", err)
	}
	
	// Try to decrypt with key2 (should fail)
	_, err = Decrypt(ciphertext, key2)
	if err == nil {
		t.Error("Decryption should fail with wrong key")
	}
}

func TestGenerateKey(t *testing.T) {
	password1 := "password1"
	password2 := "password2"
	
	key1 := GenerateKey(password1)
	key2 := GenerateKey(password2)
	key1Again := GenerateKey(password1)
	
	// Keys should be 32 bytes (256 bits)
	if len(key1) != 32 {
		t.Errorf("Key length should be 32 bytes, got %d", len(key1))
	}
	
	// Same password should generate same key
	if string(key1) != string(key1Again) {
		t.Error("Same password should generate same key")
	}
	
	// Different passwords should generate different keys
	if string(key1) == string(key2) {
		t.Error("Different passwords should generate different keys")
	}
}

func TestDecryptInvalidData(t *testing.T) {
	key := GenerateKey("test-password")
	
	// Test with invalid base64
	_, err := Decrypt("invalid-base64!", key)
	if err == nil {
		t.Error("Should fail with invalid base64")
	}
	
	// Test with too short data
	_, err = Decrypt("dGVzdA==", key) // "test" in base64, too short for GCM
	if err == nil {
		t.Error("Should fail with too short data")
	}
}
