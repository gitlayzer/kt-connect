package options

func ReplayFlags() []OptionConfig {
	flags := []OptionConfig{
		{
			Target:       "LogPath",
			DefaultValue: "",
			Description:  "Path to mirror log file or directory",
		},
		{
			Target:       "Target",
			DefaultValue: "",
			Description:  "Target address to replay traffic to, e.g. 127.0.0.1:8080",
		},
	}
	return flags
}
