package conf

import (
	"log"
	"path/filepath"
	"time"
	"errors"
	"strings"
	"github.com/BurntSushi/toml"
)

// DurationSecond represent a `time.Duration` with Seconds
type DurationSecond int

var errNotSupportedDuration = errors.New("go-isso DO NOT support duration unit y w d")

// UnmarshalText unmarshl the string to `time.Duration`
func (ds *DurationSecond) UnmarshalText(text []byte) error {
	durText := string(text)
	if strings.ContainsAny(durText, "ywd") {
		return errNotSupportedDuration
	}

	dur, err := time.ParseDuration(durText)
	if err != nil {
		return err
	}

	*ds = DurationSecond(dur.Seconds())
	return nil
}

// Generated with "https://toml-to-json.matiaskorhonen.fi/" & "https://app.quicktype.io/"

// Config save all config for this project
type Config struct {
	Name       string
	Hosts      []string
	Listen     string
	Notify     Notify
	Database   Database
	Admin      Admin
	Moderation Moderation
	Guard      Guard
	Markup     Markup
	Hash       Hash
}

// Admin control admin login stuff
type Admin struct {
	Enable   bool
	Password string
}

// Database contains all supported Database type
type Database struct {
	Dialect string
	Sqlite3  Sqlite3
}

// Sqlite3 contains config for sqlite3 database
type Sqlite3 struct {
	Path string
}

// Guard save the config for go-isso's guarder
type Guard struct {
	Enable        bool
	RateLimit     int64
	DirectReply   int64
	ReplyToSelf   bool
	RequireAuthor bool
	RequireEmail  bool
	EditMaxAge    DurationSecond
}

// Hash save the config for how to hash 
type Hash struct {
	Salt      string
	Algorithm string
}

type Markup struct {
	ExtensionsInt               int
	HTMLFlagsInt                int
	AdditionalAllowedElements   []string
	AdditionalAllowedAttributes []string
}

type Moderation struct {
	Enable     bool
	PurgeAfter DurationSecond
}

type Notify struct {
	Log      Log
	Email    Email
	Telegram Telegram
}

type Log struct {
	Enable   bool
	FilePath string
}

type Email struct {
	Enable                bool
	CanReplyNotifications bool
	To                    string
	From                  string
	SMTP                  SMTP
}

type SMTP struct {
	Username string
	Password string
	Host     string
	Port     string
	Security string
	Timeout  int64
}

type Telegram struct {
	Enable bool
	UserID int64
}

// Load config for isso
func Load(confPath string) (Config, error) {
	ConfigAbsPath, _ := filepath.Abs(confPath)
	log.Printf("[INFO] Load config from %v", ConfigAbsPath)

	var c Config

	if _, err := toml.DecodeFile(confPath, &c); err != nil {
		return Config{}, err
	}
	return c, nil
}
