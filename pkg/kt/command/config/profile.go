package config

import (
	"fmt"
	"github.com/gitlayzer/kt-connect/pkg/kt/command/general"
	opt "github.com/gitlayzer/kt-connect/pkg/kt/command/options"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
	"sort"
)

type profileAddOptions struct {
	KubeContext              string
	Namespace                string
	ProxyMode                string
	ExchangeMode             string
	ExchangeExpose           string
	ExchangeSkipPortChecking bool
	MeshMode                 string
	MeshExpose               string
	MeshSkipPortChecking     bool
	PreviewExpose            string
	PreviewExternal          bool
	PreviewSkipPortChecking  bool
	Overwrite                bool
}

var profileAddFlag profileAddOptions

func NewProfileCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profile",
		Short: "Manage ktctl configuration profiles",
		RunE: func(cmd *cobra.Command, args []string) error {
			opt.HideGlobalFlags(cmd)
			return cmd.Help()
		},
	}

	cmd.AddCommand(newProfileAddCommand())
	cmd.AddCommand(general.SimpleSubCommand("list", "List saved profiles", ProfileList, nil))
	cmd.AddCommand(general.SimpleSubCommand("use", "Set current profile", ProfileUse, ProfileUseHandle))
	cmd.AddCommand(general.SimpleSubCommand("delete", "Delete a profile", ProfileDelete, ProfileDeleteHandle))
	cmd.SetUsageTemplate(general.UsageTemplate(false))
	return cmd
}

func newProfileAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			opt.HideGlobalFlags(cmd)
			return ProfileAdd(cmd, args)
		},
	}
	cmd.Flags().StringVar(&profileAddFlag.KubeContext, "kube-context", "", "Kubernetes context name")
	cmd.Flags().StringVar(&profileAddFlag.Namespace, "namespace", "", "Kubernetes namespace")
	cmd.Flags().StringVar(&profileAddFlag.ProxyMode, "proxy-mode", "", "Connect proxy mode, e.g. tun2socks or sshuttle")
	cmd.Flags().StringVar(&profileAddFlag.ExchangeMode, "exchange-mode", "", "Default exchange mode")
	cmd.Flags().StringVar(&profileAddFlag.ExchangeExpose, "exchange-expose", "", "Default exchange expose ports")
	cmd.Flags().BoolVar(&profileAddFlag.ExchangeSkipPortChecking, "exchange-skip-port-checking", false, "Skip port checking for exchange")
	cmd.Flags().StringVar(&profileAddFlag.MeshMode, "mesh-mode", "", "Default mesh mode")
	cmd.Flags().StringVar(&profileAddFlag.MeshExpose, "mesh-expose", "", "Default mesh expose ports")
	cmd.Flags().BoolVar(&profileAddFlag.MeshSkipPortChecking, "mesh-skip-port-checking", false, "Skip port checking for mesh")
	cmd.Flags().StringVar(&profileAddFlag.PreviewExpose, "preview-expose", "", "Default preview expose ports")
	cmd.Flags().BoolVar(&profileAddFlag.PreviewExternal, "preview-external", false, "Create external preview service")
	cmd.Flags().BoolVar(&profileAddFlag.PreviewSkipPortChecking, "preview-skip-port-checking", false, "Skip port checking for preview")
	cmd.Flags().BoolVar(&profileAddFlag.Overwrite, "overwrite", false, "Overwrite if profile already exists")
	cmd.SetUsageTemplate(general.UsageTemplate(false))
	return cmd
}

func ProfileAdd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("must specify a profile name")
	}
	name := args[0]
	if !profileNamePattern.MatchString(name) {
		return fmt.Errorf("invalid profile name, must only contains letter, number, underline, hyphen or dot")
	}
	if exist, err := opt.ProfileExists(name); err != nil {
		return err
	} else if exist && !profileAddFlag.Overwrite {
		return fmt.Errorf("profile '%s' already exists, use --overwrite to replace it", name)
	}

	profile := opt.Profile{Name: name}
	if cmd.Flags().Changed("kube-context") {
		profile.KubeContext = profileAddFlag.KubeContext
	}
	if cmd.Flags().Changed("namespace") {
		profile.Namespace = profileAddFlag.Namespace
	}
	if cmd.Flags().Changed("proxy-mode") {
		profile.ProxyMode = profileAddFlag.ProxyMode
	}
	if cmd.Flags().Changed("exchange-mode") {
		profile.ExchangeMode = profileAddFlag.ExchangeMode
	}
	if cmd.Flags().Changed("exchange-expose") {
		profile.ExchangeExpose = profileAddFlag.ExchangeExpose
	}
	if cmd.Flags().Changed("exchange-skip-port-checking") {
		profile.ExchangeSkipPortChecking = &profileAddFlag.ExchangeSkipPortChecking
	}
	if cmd.Flags().Changed("mesh-mode") {
		profile.MeshMode = profileAddFlag.MeshMode
	}
	if cmd.Flags().Changed("mesh-expose") {
		profile.MeshExpose = profileAddFlag.MeshExpose
	}
	if cmd.Flags().Changed("mesh-skip-port-checking") {
		profile.MeshSkipPortChecking = &profileAddFlag.MeshSkipPortChecking
	}
	if cmd.Flags().Changed("preview-expose") {
		profile.PreviewExpose = profileAddFlag.PreviewExpose
	}
	if cmd.Flags().Changed("preview-external") {
		profile.PreviewExternal = &profileAddFlag.PreviewExternal
	}
	if cmd.Flags().Changed("preview-skip-port-checking") {
		profile.PreviewSkipPortChecking = &profileAddFlag.PreviewSkipPortChecking
	}

	if err := opt.SaveProfile(profile); err != nil {
		return err
	}
	log.Info().Msgf("Profile '%s' saved", name)
	return nil
}

func ProfileList(args []string) error {
	names, err := opt.ListProfileNames()
	if err != nil {
		return err
	}
	current, _, err := opt.CurrentProfileName()
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	sort.Strings(names)
	for _, name := range names {
		if name == current {
			fmt.Printf("* %s\n", name)
		} else {
			fmt.Printf("  %s\n", name)
		}
	}
	return nil
}

func ProfileUse(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("must specify a profile name")
	}
	name := args[0]
	exist, err := opt.ProfileExists(name)
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("profile '%s' not exists", name)
	}
	if err := opt.SetCurrentProfileName(name); err != nil {
		return err
	}
	log.Info().Msgf("Profile '%s' activated", name)
	return nil
}

func ProfileUseHandle(cmd *cobra.Command) {
	cmd.ValidArgsFunction = profileNameValidator
}

func ProfileDelete(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("must specify a profile name")
	}
	name := args[0]
	exist, err := opt.ProfileExists(name)
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("profile '%s' not exists", name)
	}
	if err := opt.DeleteProfile(name); err != nil && !os.IsNotExist(err) {
		return err
	}
	current, ok, err := opt.CurrentProfileName()
	if err != nil {
		return err
	}
	if ok && current == name {
		if err := opt.ClearCurrentProfile(); err != nil {
			return err
		}
	}
	log.Info().Msgf("Profile '%s' removed", name)
	return nil
}

func ProfileDeleteHandle(cmd *cobra.Command) {
	cmd.ValidArgsFunction = profileNameValidator
}

func profileNameValidator(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	names, err := opt.ListProfileNames()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}
