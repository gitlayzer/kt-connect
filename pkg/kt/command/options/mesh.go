package options

import (
	"fmt"
	"github.com/gitlayzer/kt-connect/pkg/kt/util"
)

func MeshFlags() []OptionConfig {
	flags := []OptionConfig{
		{
			Target:       "Expose",
			DefaultValue: "",
			Description:  "Ports to expose, use ',' separated, in [port] or [local:remote] format, e.g. 7001,8080:80",
			Required:     true,
		},
		{
			Target:       "Mode",
			DefaultValue: util.MeshModeAuto,
			Description:  "Mesh method 'auto' or 'manual'",
		},
		{
			Target:       "VersionMark",
			DefaultValue: "",
			Description:  "Specify the version of mesh service, e.g. '0.0.1' or 'mark:local'",
		},
		{
			Target:       "SkipPortChecking",
			DefaultValue: false,
			Description:  "Do not check whether specified local ports are listened",
		},
		{
			Target:       "RouterImage",
			DefaultValue: fmt.Sprintf("%s:v%s", util.ImageKtRouter, Store.Version),
			Description:  "(auto method only) Customize router image",
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
