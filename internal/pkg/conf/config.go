package conf

import (
	"github.com/BurntSushi/toml"
	"log"
	"path/filepath"
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
	Options                     []string
	AdditionalAllowedElements   []string
	AdditionalAllowedAttributes []string
}

type Moderation struct {
	Enable     bool
	PurgeAfter string
}

type Notify struct {
	Type     string
	Log      Log
	Email    Email
	Telegram Telegram
}

type Log struct {
	Path string
}

type Email struct {
	CanReplyNotifications bool
	SMTP                  SMTP
}

type SMTP struct {
	Username string
	Password string
	Host     string
	Port     string
	Security string
	To       string
	From     string
	Timeout  int64
}

type Telegram struct {
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
