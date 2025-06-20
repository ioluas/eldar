package storage

import "go.etcd.io/bbolt"

// Credentials structure to hold the user credentials
type Credentials struct {
	Username     string
	AccessToken  string
	RefreshToken string
}

type Config struct {
	Endpoint string
	AnonKey  string
}

type Database struct {
	Db          *bbolt.DB
	Config      *Config
	Credentials *Credentials
}
