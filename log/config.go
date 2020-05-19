package log

//Config log attributes in config file
type Config struct {
	Filename     string `mapstructure:"filename"`
	Level        int32  `mapstructure:"level"`
	MaxSize      int32  `mapstructure:"max_size"`
	MaxLines     int32  `mapstructure:"max_lines"`
	HourEnabled  bool   `mapstructure:"hour_enabled"`
	DailyEnabled bool   `mapstructure:"daily_enabled"`
	QueueSize    int    `mapstructure:"queue_size"`
}

//GetQueueSize get module queue size
func (c *Config) GetQueueSize() int {
	return c.QueueSize
}
