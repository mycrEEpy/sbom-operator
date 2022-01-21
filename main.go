package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"

	"github.com/ckotzbauer/sbom-operator/internal"
	"github.com/ckotzbauer/sbom-operator/internal/daemon"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Version sets the current Operator version
	Version = "0.0.1"
	Commit  = "main"
	Date    = ""
	BuiltBy = ""

	verbosity  string
	daemonCron string

	rootCmd = &cobra.Command{
		Use:              "sbom-operator",
		Short:            "An operator for cataloguing all k8s-cluster-images to git.",
		PersistentPreRun: internal.BindFlags,
		Run: func(cmd *cobra.Command, args []string) {
			internal.SetUpLogs(os.Stdout, verbosity)
			printVersion()

			daemon.Start(viper.GetString(internal.ConfigKeyCron))

			logrus.Info("Webserver is running at port 8080")
			http.HandleFunc("/health", health)
			logrus.WithError(http.ListenAndServe(":8080", nil)).Fatal("Starting webserver failed!")
		},
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&verbosity, internal.ConfigKeyVerbosity, "v", logrus.InfoLevel.String(), "Log-level (debug, info, warn, error, fatal, panic)")
	rootCmd.PersistentFlags().StringVarP(&daemonCron, internal.ConfigKeyCron, "c", "@hourly", "Backround-Service interval (CRON)")
	rootCmd.PersistentFlags().String(internal.ConfigKeyFormat, "json", "SBOM-Format.")
	rootCmd.PersistentFlags().String(internal.ConfigKeyGitWorkingTree, "/work", "Directory to place the git-repo.")
	rootCmd.PersistentFlags().String(internal.ConfigKeyGitRepository, "", "Git-Repository-URL (HTTPS).")
	rootCmd.PersistentFlags().String(internal.ConfigKeyGitBranch, "main", "Git-Branch to checkout.")
	rootCmd.PersistentFlags().String(internal.ConfigKeyGitPath, "", "Folder-Path inside the Git-Repository.")
	rootCmd.PersistentFlags().String(internal.ConfigKeyGitAccessToken, "", "Git-Access-Token.")
	rootCmd.PersistentFlags().String(internal.ConfigKeyGitAuthorName, "", "Author name to use for Git-Commits.")
	rootCmd.PersistentFlags().String(internal.ConfigKeyGitAuthorEmail, "", "Author email to use for Git-Commits.")
	rootCmd.PersistentFlags().String(internal.ConfigKeyPodLabelSelector, "", "Kubernetes Label-Selector for pods.")
	rootCmd.PersistentFlags().String(internal.ConfigKeyNamespaceLabelSelector, "", "Kubernetes Label-Selector for namespaces.")
}

func initConfig() {
	viper.SetEnvPrefix("SGO")
	viper.AutomaticEnv()
}

func printVersion() {
	logrus.Info(fmt.Sprintf("Version: %s", Version))
	logrus.Info(fmt.Sprintf("Commit: %s", Commit))
	logrus.Info(fmt.Sprintf("Buit at: %s", Date))
	logrus.Info(fmt.Sprintf("Buit by: %s", BuiltBy))
	logrus.Info(fmt.Sprintf("Go Version: %s", runtime.Version()))
}

func health(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
	fmt.Fprint(w, "Running!")
}

func main() {
	rootCmd.Execute()
}
