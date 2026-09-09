package auth_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/ijon6k/song-moodboard/internal/auth"
)

func TestAuthManager_TokenLifecycle(t *testing.T) {
	secret := "test-secret-twilight-key"
	authMgr := auth.NewAuthManager(secret, nil)

	userID := uuid.New()
	username := "aimer_listener"
	email := "test@song-moodboard.io"

	token, err := authMgr.GenerateToken(userID, username, email)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	claims, err := authMgr.ValidateToken(token)
	if err != nil {
		t.Fatalf("failed to validate token: %v", err)
	}

	if claims.UserID != userID || claims.Username != username || claims.Email != email {
		t.Errorf("token claims mismatch: got %+v", claims)
	}
}

func TestAuthManager_PasswordHashing(t *testing.T) {
	authMgr := auth.NewAuthManager("test-secret", nil)
	pass := "myQuietHaze2026!"

	hash, err := authMgr.HashPassword(pass)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	if !authMgr.CheckPassword(hash, pass) {
		t.Errorf("password check failed for valid password")
	}

	if authMgr.CheckPassword(hash, "wrongPass") {
		t.Errorf("password check succeeded for invalid password")
	}
}
