package log

//Config log attributes in config file
type Config struct {
	Filename     string `yarm:"filename"`
	Level        int    `yarm:"level"`
	MaxSize      uint32 `yarm:"max_size"`
	MaxLines     uint32 `yarm:"max_lines"`
	HourEnabled  bool   `yarm:"hour_enabled"`
	DailyEnabled bool   `yarm:"daily_enabled"`
}
