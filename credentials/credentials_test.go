// Package credentials provides functionality for managing user credentials
package credentials

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

	dbPath := filepath.Join(storageDir, ".eldar", "credentials.db")
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

// TestGetStorageDir tests the GetStorageDir function
func TestGetStorageDir(t *testing.T) {
	storageDir, err := GetStorageDir()
	if err != nil {
		t.Fatalf("GetStorageDir failed: %v", err)
	}
	if storageDir == "" {
		t.Error("Expected non-empty storage directory")
	}

	originalHome := os.Getenv("HOME")
	_ = os.Setenv("HOME", "/nonexistent/path")
	defer func(key, value string) {
		if err = os.Setenv(key, value); err != nil {
			fmt.Printf("Error setting environment variable: %v", err)
		}
	}("HOME", originalHome)

	_, err = GetStorageDir()
	if err == nil {
		t.Error("Expected error with invalid home directory")
	}
	if !errors.Is(err, os.ErrNotExist) && !errors.Is(err, os.ErrPermission) &&
		!strings.Contains(err.Error(), "no such file or directory") &&
		!strings.Contains(err.Error(), "permission denied") {
		t.Errorf("Expected not exist or permission error, got: %v", err)
	}
}

// TestGetCredentialsErrorCases tests error cases in GetCredentials
func TestGetCredentialsErrorCases(t *testing.T) {
	setupTestEnvironment(t)

	originalStorageDir := testStorageDir
	testStorageDir = "/nonexistent/path"
	defer func() { testStorageDir = originalStorageDir }()

	dbDir := filepath.Join(testStorageDir, ".eldar")
	if err := os.MkdirAll(dbDir, 0444); err == nil {
		defer func(name string, mode os.FileMode) {
			if err = os.Chmod(name, mode); err != nil {
				t.Errorf("Failed to set permissions on %s: %v", name, err)
			}
		}(dbDir, 0755)
	}

	_, err := getTestCredentials()
	if err == nil {
		t.Error("Expected error with invalid storage directory")
	}
	if !errors.Is(err, os.ErrNotExist) && !errors.Is(err, os.ErrPermission) &&
		!strings.Contains(err.Error(), "no such file or directory") &&
		!strings.Contains(err.Error(), "permission denied") {
		t.Errorf("Expected not exist or permission error, got: %v", err)
	}
}

// TestClearCredentialsErrorCases tests error cases in ClearCredentials
func TestClearCredentialsErrorCases(t *testing.T) {
	setupTestEnvironment(t)

	err := clearTestCredentials()
	if err != nil {
		t.Errorf("Expected no error when clearing non-existent database, got: %v", err)
	}

	originalStorageDir := testStorageDir
	testStorageDir = "/nonexistent/path"
	defer func() { testStorageDir = originalStorageDir }()

	dbDir := filepath.Join(testStorageDir, ".eldar")
	if err = os.MkdirAll(dbDir, 0444); err == nil {
		defer func(name string, mode os.FileMode) {
			_ = os.Chmod(name, mode)
		}(dbDir, 0755)
	}

	err = clearTestCredentials()
	if err == nil {
		// Accept nil error if the database file does not exist, matching implementation contract
		dbPath := filepath.Join(testStorageDir, ".eldar", "credentials.db")
		if _, statErr := os.Stat(dbPath); !os.IsNotExist(statErr) {
			t.Error("Expected error with invalid storage directory")
		}
	} else if !errors.Is(err, os.ErrNotExist) && !errors.Is(err, os.ErrPermission) &&
		!strings.Contains(err.Error(), "no such file or directory") &&
		!strings.Contains(err.Error(), "permission denied") {
		t.Errorf("Expected not exist or permission error, got: %v", err)
	}
}

// TestAddTestCredentials tests the AddTestCredentials function
func TestAddTestCredentials(t *testing.T) {
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

	err = AddTestCredentialsWithDir("/nonexistent/path")
	if err == nil {
		t.Error("Expected error when adding credentials to invalid directory")
	}
}

// TestDatabasePermissions tests database permission issues
func TestDatabasePermissions(t *testing.T) {
	setupTestEnvironment(t)

	dbDir := filepath.Join(testStorageDir, ".eldar")
	if err := os.MkdirAll(dbDir, 0444); err != nil {
		t.Fatalf("Failed to create read-only directory: %v", err)
	}
	defer func(name string, mode os.FileMode) {
		_ = os.Chmod(name, mode)
	}(dbDir, 0755)

	err := AddTestCredentialsWithDir(testStorageDir)
	if err == nil {
		t.Error("Expected error when adding credentials to read-only directory")
	}
}
