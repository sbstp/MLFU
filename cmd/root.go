package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/sbstp/MLFU/drivers"
	"github.com/spf13/cobra"
)

type Args struct {
	ConfigPath string
	LogPath    string
}

var (
	args = &Args{}

	rootCmd = &cobra.Command{
		Use:   "mlfu",
		Short: "Magnet Link Forwarding Utility",
	}

	setupCmd = &cobra.Command{
		Use:   "setup",
		Short: "register MLFU as the magnet link handler",
		RunE: func(cmd *cobra.Command, args []string) error {
			return setup()
		},
	}

	openCmd = &cobra.Command{
		Use:   "open",
		Short: "open magnet link and send to configured client",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return open(args[0])
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&args.ConfigPath, "config", "/etc/mlfu/config.json", "path to config file")
	rootCmd.PersistentFlags().StringVar(&args.LogPath, "log", "/var/log/mlfu.log", "path to log file")

	rootCmd.AddCommand(setupCmd)
	rootCmd.AddCommand(openCmd)
}

func Execute() {
	rootCmd.Execute()
}

func setup() error {
	execPath, err := os.Executable()
	if err != nil {
		return err
	}

	absConfigPath, _ := filepath.Abs(args.ConfigPath)

	desktopFilePath := fmt.Sprintf("%s/.local/share/applications/mlfu.desktop", os.Getenv("HOME"))
	desktopFileData := fmt.Sprintf(`[Desktop Entry]
Type=Application
Name=Magnet Link Relay Utility
Exec=%s --config %s open %%u
StartupNotify=false
MimeType=x-scheme-handler/magnet;
`, execPath, absConfigPath)

	err = os.WriteFile(desktopFilePath, []byte(desktopFileData), 0644)
	if err != nil {
		return err
	}
	log.Printf("wrote Desktop Entry to %s", desktopFilePath)

	cmd := exec.Command("xdg-mime", "default", "mlfu.desktop", "x-scheme-handler/magnet")
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func open(magnet string) error {
	config, err := drivers.LoadConfig(args.ConfigPath)
	if err != nil {
		return err
	}

	driver := drivers.GetDriver(config.Driver)
	if driver == nil {
		return fmt.Errorf("driver %q not found", config.Driver)
	}

	return driver.AddMagnetURL(config, magnet)
}
