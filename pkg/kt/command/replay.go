package command

import (
	"github.com/gitlayzer/kt-connect/pkg/kt/command/general"
	opt "github.com/gitlayzer/kt-connect/pkg/kt/command/options"
	"github.com/gitlayzer/kt-connect/pkg/kt/command/replay"
	"github.com/spf13/cobra"
)

// NewReplayCommand return new replay command
func NewReplayCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "replay",
		Short: "Replay mirrored traffic logs to a target address",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return replay.Replay(opt.Get().Replay.LogPath, opt.Get().Replay.Target)
		},
		Example: "ktctl replay --log-path ./mirror-logs --target 127.0.0.1:8080",
	}

	cmd.SetUsageTemplate(general.UsageTemplate(true))
	opt.SetOptions(cmd, cmd.Flags(), opt.Get().Replay, opt.ReplayFlags())
	return cmd
}
