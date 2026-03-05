package config

type Config struct {
	Databases    []DatabaseConfig   `yaml:"databases"`
	Backup       BackupConfig       `yaml:"backup"`
	Storage      StorageConfig      `yaml:"storage"`
	Notification NotificationConfig `yaml:"notification"`
	Schedule     ScheduleConfig     `yaml:"schedule"`
}

type DatabaseConfig struct {
	Name     string `yaml:"name"`
	Type     string `yaml:"type"`
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Port     int    `yaml:"port"`
	AuthDB   string `yaml:"auth_db"` // MongoDB: authentication database (default "admin")
}

type BackupConfig struct {
	Type       string `yaml:"type"`       // "full", "incremental"
	Path       string `yaml:"path"`       // Backup storage path
	Retain     int    `yaml:"retain"`     // Number of backups to keep
	Compress   bool   `yaml:"compress"`   // Enable compression
	CompressionLevel int `yaml:"compression_level"` // Gzip compression level (1-9)
}

type StorageConfig struct {
	Type    string `yaml:"type"` // "local", "s3", "gcs"
	Path    string `yaml:"path"` // Only used for "local" storage
	Bucket  string `yaml:"bucket"`
	Region  string `yaml:"region"`
	Project string `yaml:"project"`
	Retain  int    `yaml:"retain"` // Number of backups to keep
}

// NotificationConfig holds all notification settings.
type NotificationConfig struct {
	Slack SlackConfig `yaml:"slack"`
}

// SlackConfig holds Slack Incoming Webhook settings.
type SlackConfig struct {
	Enabled    bool   `yaml:"enabled"`
	WebhookURL string `yaml:"webhook_url"`
	Channel    string `yaml:"channel"`    // optional display name
	OnSuccess  bool   `yaml:"on_success"`
	OnFailure  bool   `yaml:"on_failure"`
}

// ScheduleConfig holds cron scheduling settings.
type ScheduleConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Cron     string `yaml:"cron"`     // e.g. "0 2 * * *"
	TimeZone string `yaml:"timezone"` // e.g. "Asia/Tokyo"
}
