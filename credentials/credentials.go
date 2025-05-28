package credentials

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"fyne.io/fyne/v2"
	"go.etcd.io/bbolt"
)

// Credentials structure to hold the user credentials
type Credentials struct {
	Username     string
	AccessToken  string
	RefreshToken string
}

// GetStorageDir returns the appropriate storage directory based on the platform
// For Android, it uses the app's storage directory if available, or falls back to UserConfigDir
// For other platforms, it uses the user's home directory
func GetStorageDir() (string, error) {
	var storageDir string
	var err error

	if runtime.GOOS == "android" || runtime.GOOS == "ios" {
		// Try to use Fyne's storage API for mobile platforms
		currentApp := fyne.CurrentApp()
		if currentApp != nil {
			storageDir = currentApp.Storage().RootURI().Path()
		} else {
			// Fallback to UserConfigDir which is more reliable on Android than UserHomeDir
			if storageDir, err = os.UserConfigDir(); err != nil {
				return "", fmt.Errorf("failed to get config directory: %w", err)
			}
		}
	} else {
		// Use home directory for desktop platforms
		if storageDir, err = os.UserHomeDir(); err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}

		// Verify the directory exists and is accessible
		if _, err = os.Stat(storageDir); err != nil {
			return "", fmt.Errorf("home directory is not accessible: %w", err)
		}
	}

	return storageDir, nil
}

// AddTestCredentials adds test credentials to the database for testing purposes
// This function is used by the main application to add test credentials
func AddTestCredentials() error {
	// Get the appropriate storage directory for the platform
	storageDir, err := GetStorageDir()
	if err != nil {
		return fmt.Errorf("failed to get storage directory: %w", err)
	}

	// Call the implementation with the storage directory
	return AddTestCredentialsWithDir(storageDir)
}

// AddTestCredentialsWithDir adds test credentials to the database for testing purposes
// with the specified storage directory
func AddTestCredentialsWithDir(storageDir string) error {
	// Check if the storage directory exists
	if _, err := os.Stat(storageDir); err != nil {
		return fmt.Errorf("storage directory is not accessible: %w", err)
	}

	dbDir := filepath.Join(storageDir, ".eldar")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return fmt.Errorf("failed to create database directory: %w", err)
	}

	// Check if we have write permissions to the directory
	testFile := filepath.Join(dbDir, ".test")
	f, err := os.Create(testFile)
	if err != nil {
		return fmt.Errorf("directory is not writable: %w", err)
	}
	_ = f.Close()
	_ = os.Remove(testFile)
	dbPath := filepath.Join(dbDir, "credentials.db")
	db, err := bbolt.Open(dbPath, 0600, nil)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer func(db *bbolt.DB) {
		_ = db.Close()
	}(db)

	err = db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("credentials"))
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}

		if err := b.Put([]byte("username"), []byte("testuser")); err != nil {
			return fmt.Errorf("failed to set username: %w", err)
		}

		if err := b.Put([]byte("access_token"), []byte("test-access-token-123")); err != nil {
			return fmt.Errorf("failed to set access token: %w", err)
		}

		if err := b.Put([]byte("refresh_token"), []byte("test-refresh-token-456")); err != nil {
			return fmt.Errorf("failed to set refresh token: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to add test credentials: %w", err)
	}

	return nil
}

// GetCredentials initializes and checks the bbolt database for credentials
func GetCredentials() (*Credentials, error) {
	// Get the appropriate storage directory for the platform
	storageDir, err := GetStorageDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get storage directory: %w", err)
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

// ClearCredentials removes all stored credentials from the database
func ClearCredentials() error {
	// Get the appropriate storage directory for the platform
	storageDir, err := GetStorageDir()
	if err != nil {
		return fmt.Errorf("failed to get storage directory: %w", err)
	}

	// Check if the database file exists
	dbPath := filepath.Join(storageDir, ".eldar", "credentials.db")
	if _, err = os.Stat(dbPath); os.IsNotExist(err) {
		fmt.Println("Database file does not exist. Nothing to clear.")
		return nil
	}

	// Open the bbolt database
	db, err := bbolt.Open(dbPath, 0600, nil)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer func(db *bbolt.DB) {
		if err := db.Close(); err != nil {
			fmt.Printf("Error closing database: %v", err)
		}
	}(db)

	// Clear the credentials
	if err = db.Update(func(tx *bbolt.Tx) error {
		// Check if the bucket exists
		b := tx.Bucket([]byte("credentials"))
		if b == nil {
			return fmt.Errorf("credentials bucket does not exist")
		}

		// Delete the credentials
		if err := b.Delete([]byte("username")); err != nil {
			return fmt.Errorf("failed to delete username: %w", err)
		}
		if err := b.Delete([]byte("access_token")); err != nil {
			return fmt.Errorf("failed to delete access token: %w", err)
		}
		if err := b.Delete([]byte("refresh_token")); err != nil {
			return fmt.Errorf("failed to delete refresh token: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("failed to clear credentials: %w", err)
	}

	fmt.Println("Credentials cleared successfully!")
	return nil
}
