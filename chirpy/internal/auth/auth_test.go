package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashPasswordAndCheck(t *testing.T) {
	const password = "correct horse battery staple"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}
	if hash == "" || hash == password {
		t.Fatalf("expected a hash distinct from the password, got %q", hash)
	}

	match, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash: %v", err)
	}
	if !match {
		t.Fatal("expected the password to match its hash")
	}

	match, err = CheckPasswordHash("wrong password", hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash wrong password: %v", err)
	}
	if match {
		t.Fatal("expected a different password to be rejected")
	}
}

func TestHashPasswordUsesAFreshSalt(t *testing.T) {
	const password = "same-password"

	first, err := HashPassword(password)
	if err != nil {
		t.Fatalf("first HashPassword: %v", err)
	}
	second, err := HashPassword(password)
	if err != nil {
		t.Fatalf("second HashPassword: %v", err)
	}
	if first == second {
		t.Fatal("expected two hashes of the same password to differ")
	}

	for _, hash := range []string{first, second} {
		match, err := CheckPasswordHash(password, hash)
		if err != nil {
			t.Fatalf("CheckPasswordHash: %v", err)
		}
		if !match {
			t.Fatalf("expected password to match hash %q", hash)
		}
	}
}

func TestCheckPasswordHashRejectsMalformedHash(t *testing.T) {
	_, err := CheckPasswordHash("password", "this-is-not-an-argon2-hash")
	if err == nil {
		t.Fatal("expected an error for a malformed hash")
	}
}

func TestMakeJWTAndValidate(t *testing.T) {
	userID := uuid.New()
	const secret = "token-secret"

	token, err := MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT: %v", err)
	}
	if token == "" {
		t.Fatal("expected a signed token")
	}

	got, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("ValidateJWT: %v", err)
	}
	if got != userID {
		t.Fatalf("user id: got %s, want %s", got, userID)
	}
}

func TestValidateJWTRejectsWrongSecret(t *testing.T) {
	token, err := MakeJWT(uuid.New(), "secret-a", time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT: %v", err)
	}

	if _, err := ValidateJWT(token, "secret-b"); err == nil {
		t.Fatal("expected validation to fail with the wrong secret")
	}
}

func TestValidateJWTRejectsExpiredToken(t *testing.T) {
	token, err := MakeJWT(uuid.New(), "secret", -time.Minute)
	if err != nil {
		t.Fatalf("MakeJWT: %v", err)
	}

	if _, err := ValidateJWT(token, "secret"); err == nil {
		t.Fatal("expected validation to fail for an expired token")
	}
}

func TestValidateJWTRejectsMalformedToken(t *testing.T) {
	if _, err := ValidateJWT("not-a-jwt", "secret"); err == nil {
		t.Fatal("expected validation to fail for a malformed token")
	}
}
