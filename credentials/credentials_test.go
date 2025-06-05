// Package credentials provides functionality for managing user credentials
package credentials

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"go.etcd.io/bbolt"
)

// testStorageDir is used to store the path to the temporary directory for tests
var testStorageDir string

// getTestStorageDir returns the test storage directory for tests
func getTestStorageDir() (string, error) {
	return testStorageDir, nil
}

// setupTestEnvironment creates a temporary directory for testing
func setupTestEnvironment(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "eldar-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	testStorageDir = tempDir
	t.Cleanup(func() {
		_ = os.RemoveAll(tempDir)
	})
}

// getTestCredentials is a test version of GetCredentials that uses the test storage directory
func getTestCredentials() (*Credentials, error) {
	// Get the test storage directory
	storageDir, err := getTestStorageDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get test storage directory: %w", err)
	}

	// Check if the storage directory exists
	if _, err = os.Stat(storageDir); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("storage directory does not exist: %w", err)
		}
		return nil, fmt.Errorf("storage directory is not accessible: %w", err)
	}

	// Create the directory for our database if it doesn't exist
	dbDir := filepath.Join(storageDir, ".eldar")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// Open the bbolt database
	dbPath := filepath.Join(dbDir, "credentials.db")
	db, err := bbolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	defer func(db *bbolt.DB) {
		_ = db.Close()
	}(db)

	// Create the bucket if it doesn't exist
	if err := db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("credentials"))
		return err
	}); err != nil {
		return nil, fmt.Errorf("failed to create bucket: %w", err)
	}

	// Check for the three values
	var creds Credentials
	err = db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("credentials"))
		if b == nil {
			return fmt.Errorf("credentials bucket does not exist")
		}

		username := b.Get([]byte("username"))
		if username != nil {
			creds.Username = string(username)
		}

		accessToken := b.Get([]byte("access_token"))
		if accessToken != nil {
			creds.AccessToken = string(accessToken)
		}

		refreshToken := b.Get([]byte("refresh_token"))
		if refreshToken != nil {
			creds.RefreshToken = string(refreshToken)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to read from database: %w", err)
	}

	return &creds, nil
}

// clearTestCredentials is a test version of ClearCredentials that uses the test storage directory
func clearTestCredentials() error {
	storageDir, err := getTestStorageDir()
	if err != nil {
		return fmt.Errorf("failed to get test storage directory: %w", err)
	}

	// Check if the storage directory exists
	if _, err = os.Stat(storageDir); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("storage directory does not exist: %w", err)
		}
		return fmt.Errorf("storage directory is not accessible: %w", err)
	}

	// If the database directory doesn't exist, there's nothing to clear
	dbDir := filepath.Join(storageDir, ".eldar")
	if _, err = os.Stat(dbDir); os.IsNotExist(err) {
		fmt.Println("Database directory does not exist. Nothing to clear.")
		return nil
	} else if err != nil {
		return fmt.Errorf("database directory is not accessible: %w", err)
	}

	dbPath := filepath.Join(dbDir, "credentials.db")
	if _, err = os.Stat(dbPath); os.IsNotExist(err) {
		fmt.Println("Database file does not exist. Nothing to clear.")
		return nil
	}

	db, err := bbolt.Open(dbPath, 0600, nil)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer func(db *bbolt.DB) {
		if err = db.Close(); err != nil {
			fmt.Printf("Error closing database: %v", err)
		}
	}(db)

	if err = db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("credentials"))
		if b == nil {
			return fmt.Errorf("credentials bucket does not exist")
		}

		if err = b.Delete([]byte("username")); err != nil {
			return fmt.Errorf("failed to delete username: %w", err)
		}
		if err = b.Delete([]byte("access_token")); err != nil {
			return fmt.Errorf("failed to delete access token: %w", err)
		}
		if err = b.Delete([]byte("refresh_token")); err != nil {
			return fmt.Errorf("failed to delete refresh token: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("failed to clear credentials: %w", err)
	}

	return nil
}

// TestGetCredentials tests the GetCredentials function
func TestGetCredentials(t *testing.T) {
	setupTestEnvironment(t)

	err := AddTestCredentialsWithDir(testStorageDir)
	if err != nil {
		t.Fatalf("Failed to add test credentials: %v", err)
	}

	creds, err := getTestCredentials()
	if err != nil {
		t.Fatalf("getTestCredentials failed: %v", err)
	}

	if creds.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", creds.Username)
	}
	if creds.AccessToken != "test-access-token-123" {
		t.Errorf("Expected access token 'test-access-token-123', got '%s'", creds.AccessToken)
	}
	if creds.RefreshToken != "test-refresh-token-456" {
		t.Errorf("Expected refresh token 'test-refresh-token-456', got '%s'", creds.RefreshToken)
	}
}

// TestClearCredentials tests the ClearCredentials function
func TestClearCredentials(t *testing.T) {
	setupTestEnvironment(t)

	err := AddTestCredentialsWithDir(testStorageDir)
	if err != nil {
		t.Fatalf("Failed to add test credentials: %v", err)
	}

	creds, err := getTestCredentials()
	if err != nil {
		t.Fatalf("getTestCredentials failed: %v", err)
	}
	if creds.Username == "" || creds.AccessToken == "" || creds.RefreshToken == "" {
		t.Fatalf("Test credentials were not properly added")
	}

	err = clearTestCredentials()
	if err != nil {
		t.Fatalf("clearTestCredentials failed: %v", err)
	}

	creds, err = getTestCredentials()
	if err != nil {
		t.Fatalf("getTestCredentials failed after clearing: %v", err)
	}
	if creds.Username != "" {
		t.Errorf("Username should be empty after clearing, got '%s'", creds.Username)
	}
	if creds.AccessToken != "" {
		t.Errorf("AccessToken should be empty after clearing, got '%s'", creds.AccessToken)
	}
	if creds.RefreshToken != "" {
		t.Errorf("RefreshToken should be empty after clearing, got '%s'", creds.RefreshToken)
	}
}
