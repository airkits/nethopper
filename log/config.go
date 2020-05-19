package log

//Config log attributes in config file
type Config struct {
	Filename     string `mapstructure:"filename"`
	Level        int    `mapstructure:"level"`
	MaxSize      int    `mapstructure:"max_size"`
	MaxLines     int    `mapstructure:"max_lines"`
	HourEnabled  bool   `mapstructure:"hour_enabled"`
	DailyEnabled bool   `mapstructure:"daily_enabled"`
}
