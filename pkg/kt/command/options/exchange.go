package options

import "github.com/gitlayzer/kt-connect/pkg/kt/util"

func ExchangeFlags() []OptionConfig {
	flags := []OptionConfig{
		{
			Target:       "Expose",
			DefaultValue: "",
			Description:  "Ports to expose, use ',' separated, in [port] or [local:remote] format, e.g. 7001,8080:80",
			Required:     true,
		},
		{
			Target:       "Mode",
			DefaultValue: util.ExchangeModeSelector,
			Description:  "Exchange method 'selector', 'scale' or 'ephemeral'(experimental)",
		},
		{
			Target:       "SkipPortChecking",
			DefaultValue: false,
			Description:  "Do not check whether specified local ports are listened",
		},
		{
			Target:       "RecoverWaitTime",
			DefaultValue: 120,
			Description:  "(scale method only) Seconds to wait for original deployment recover before turn off the shadow pod",
		},
		{
			Target:       "MirrorTarget",
			DefaultValue: "",
			Description:  "Mirror traffic to the specified address, e.g. 127.0.0.1:18080",
		},
		{
			Target:       "MirrorSampleRate",
			DefaultValue: 100,
			Description:  "Mirror sample rate in percentage (0-100)",
		},
		{
			Target:       "MirrorRedactRules",
			DefaultValue: "",
			Description:  "Mirror redact rules in 'pattern=replacement' format, separated by ';'",
		},
		{
			Target:       "MirrorLogPath",
			DefaultValue: "",
			Description:  "Directory to write mirror request logs",
		},
	}
	return flags
}
