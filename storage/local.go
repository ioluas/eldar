package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"fyne.io/fyne/v2"
	"go.etcd.io/bbolt"
)

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
			return currentApp.Storage().RootURI().Path(), nil
		}
		// Fallback to UserConfigDir which is more reliable on Android than UserHomeDir
		storageDir, err = os.UserConfigDir()
		if err != nil {
			return "", fmt.Errorf("failed to get config directory: %w", err)
		}
		return storageDir, nil
	}
	// Use home directory for desktop platforms
	storageDir, err = os.UserConfigDir()

	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	if _, err = os.Stat(storageDir); err != nil {
		return "", fmt.Errorf("home directory is not accessible: %w", err)
	}

	return storageDir, nil
}

func InitDatabase() (*Database, error) {
	var err error
	storageDir, err := GetStorageDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get storage directory: %w", err)
	}
	// Create the directory for our database if it doesn't exist
	dbDir := filepath.Join(storageDir, "eldar")
	if err = os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	db, err := bbolt.Open(filepath.Join(dbDir, "eldar.db"), 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	database := &Database{
		Db:          db,
		Credentials: &Credentials{},
		Config:      &Config{},
	}
	ensureBucketExists := func(name string) error {
		return db.Update(func(tx *bbolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte(name))
			return err
		})
	}

	if err = ensureBucketExists("config"); err != nil {
		return nil, err
	}
	if err = ensureBucketExists("credentials"); err != nil {
		return nil, err
	}
	return database, nil
}

func (db *Database) LoadConfig() error {
	if err := db.Db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("config"))
		if b == nil {
			return fmt.Errorf("config bucket does not exist")
		}
		endpoint := b.Get([]byte("endpoint"))
		if endpoint == nil {
			db.Config.Endpoint = ""
		} else {
			db.Config.Endpoint = string(endpoint)
		}
		anonKey := b.Get([]byte("anonKey"))
		if anonKey == nil {
			db.Config.AnonKey = ""
		} else {
			db.Config.AnonKey = string(anonKey)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	return nil
}

func (db *Database) SaveConfig(endpoint, anonKey string) error {
	if err := db.Db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("config"))
		if b == nil {
			return fmt.Errorf("config bucket does not exist")
		}
		if err := b.Put([]byte("endpoint"), []byte(endpoint)); err != nil {
			return fmt.Errorf("failed to set endpoint: %w", err)
		}
		if err := b.Put([]byte("anonKey"), []byte(anonKey)); err != nil {
			return fmt.Errorf("failed to set Anonmous Key: %w", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}
	db.Config.Endpoint = endpoint
	db.Config.AnonKey = anonKey
	return nil
}

func (db *Database) LoadCredentials() error {
	if err := db.Db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("credentials"))
		if b == nil {
			return fmt.Errorf("credentials bucket does not exist")
		}
		username := b.Get([]byte("username"))
		if username == nil {
			db.Credentials.Username = ""
		} else {
			db.Credentials.Username = string(username)
		}
		accessToken := b.Get([]byte("access_token"))
		if accessToken == nil {
			db.Credentials.AccessToken = ""
		} else {
			db.Credentials.AccessToken = string(accessToken)
		}
		refreshToken := b.Get([]byte("refresh_token"))
		if refreshToken == nil {
			db.Credentials.RefreshToken = ""
		} else {
			db.Credentials.RefreshToken = string(refreshToken)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to load credentials: %w", err)
	}
	return nil
}

func (db *Database) SaveCredentials(username, accessToken, refreshToken string) error {
	if err := db.Db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("credentials"))
		if b == nil {
			return fmt.Errorf("credentials bucket does not exist")
		}

		if err := b.Put([]byte("username"), []byte(username)); err != nil {
			return fmt.Errorf("failed to set username: %w", err)
		}
		if err := b.Put([]byte("access_token"), []byte(accessToken)); err != nil {
			return fmt.Errorf("failed to set access token: %w", err)
		}
		if err := b.Put([]byte("refresh_token"), []byte(refreshToken)); err != nil {
			return fmt.Errorf("failed to set refresh token: %w", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to save credentials: %w", err)
	}
	db.Credentials.Username = username
	db.Credentials.AccessToken = accessToken
	db.Credentials.RefreshToken = refreshToken
	return nil
}
