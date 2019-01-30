package conf

import (
	"errors"
	"github.com/BurntSushi/toml"
	"log"
	"path/filepath"
	"time"
)

// Generated with "https://toml-to-json.matiaskorhonen.fi/" & "https://app.quicktype.io/"

// Configure save all config for this project
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

type Admin struct {
	Enable   bool
	Password string
}

type Database struct {
	Sqlite Sqlite
}

type Sqlite struct {
	Path string
}

type Guard struct {
	Enable        bool
	RateLimit     int64
	DirectReply   int64
	ReplyToSelf   bool
	RequireAuthor bool
	RequireEmail  bool
	EditMaxAge    string
}

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
	PurgeAfter string
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

	if _, err := time.ParseDuration(c.Guard.EditMaxAge); err != nil {
		return Config{}, errors.New("Guard.EditMaxAge can't be parsed as Duration")
	}
	if _, err := time.ParseDuration(c.Moderation.PurgeAfter); err != nil {
		return Config{}, errors.New("Moderation.PurgeAfter can't be parsed as Duration")
	}
	return c, nil
}

func DurationSeconds(duration string) int {
	Duration, _ := time.ParseDuration(duration)
	return int(Duration.Seconds())
}
