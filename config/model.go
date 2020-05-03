package config

// Config is the main config struct for go-isso
// All go-isso config related with it.
type Config struct {
	DBPath             string   `ini:"dbpath"`
	Name               string   `ini:"name"` // required to dispatch multiple websites, not used otherwise.
	Host               []string `ini:"host"`
	MaxAge             int
	Notify             []string `ini:"notify"`
	ReplyNotifications bool     `ini:"reply-notifications"`
	LogFilePath        string   `ini:"log-file"`
	Gravatar           bool     `ini:"gravatar"`
	GravatarURL        string   `ini:"gravatar-url"`
	Server             Server
	Admin              Admin
	Moderation         Moderation
	SMTP               SMTP
}

// Server store all HTTP server related config
type Server struct {
	Listen         string `ini:"listen"`
	PublicEndpoint string `ini:"public-endpoint"`
	Guard          Guard
}

// Admin interface config
type Admin struct {
	Password string `ini:"password"`
	Enable   bool   `ini:"enabled"`
}

// Moderation config
type Moderation struct {
	Enable              bool   `ini:"enabled"`
	PurgeAfter          string `ini:"purge-after"`
	ApproveAcquaintance bool   `ini:"approve-if-email-previously-approved"`
}

// Guard store basic spam protection config
type Guard struct {
	Enable        bool `ini:"enabled"`
	RateLimit     int  `ini:"ratelimit"`
	DirectReply   int  `ini:"direct-reply"`
	ReplyToSelf   bool `ini:"reply-to-self"`
	RequireAuthor bool `ini:"require-author"`
	RequireEmail  bool `ini:"require-email"`
	Markup        Markup
}

// Markup store basic Customize markup and sanitized HTML config
type Markup struct {
	AllowedElements   []string `ini:"allowed-elements"`
	AllowedAttributes []string `ini:"allowed-attributes"`
}

// SMTP save notify thought smtp config
type SMTP struct {
	Username string `ini:"username"`
	Password string `ini:"password"`
	Host     string `ini:"host"`
	Port     int    `ini:"port"`
	Security string `ini:"security"`
	To       string `ini:"to"`
	From     string `ini:"from"`
	Timeout  int    `ini:"timeout"`
}
