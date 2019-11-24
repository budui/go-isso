package config

// Config is the main config struct for go-isso
// All go-isso config related with it.
type Config struct {
	Debug bool `ini:"debug"`
	DBPath string `ini:"dbpath"`
	Name string `ini:"name"` // required to dispatch multiple websites, not used otherwise.
	Host []string `ini:"host"`
	//MaxAge string
	Notify string `ini:"notify"`
	ReplyNotifications bool `ini:"reply-notifications"`
	LogFilePath string `ini:"log-file"`
	// Gravatar string
	// GravatarURL string
	Server Server
	Admin Admin
	Moderation Moderation
}

// Server store all HTTP server related config
type Server struct {
	Listen string `ini:"listen"`
	PublicEndpoint string `ini:"public-endpoint"`
}

// Admin interface config
type Admin struct {
	Password string `ini:"password"`
	Enable bool `ini:"enabled"`
}

// Moderation config
type Moderation struct {
	Enable bool `ini:"enabled"`
	PurgeAfter string `ini:"purge-after"`
}