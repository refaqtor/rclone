package cmd

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/artpar/rclone/fs"
	"github.com/artpar/rclone/fs/config/configflags"
	"github.com/artpar/rclone/fs/filter/filterflags"
	"github.com/artpar/rclone/fs/rc/rcflags"
	"github.com/artpar/rclone/lib/atexit"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Root is the main rclone command
var Root = &cobra.Command{
	Use:   "rclone",
	Short: "Show help for rclone commands, flags and backends.",
	Long: `
Rclone syncs files to and from cloud storage providers as well as
mounting them, listing them in lots of different ways.

See the home page (https://rclone.org/) for installation, usage,
documentation, changelog and configuration walkthroughs.

`,
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		fs.Debugf("rclone", "Version %q finishing with parameters %q", fs.Version, os.Args)
		atexit.Run()
	},
}

// root help command
var helpCommand = &cobra.Command{
	Use:   "help",
	Short: Root.Short,
	Long:  Root.Long,
	Run: func(command *cobra.Command, args []string) {
		Root.SetOutput(os.Stdout)
		_ = Root.Usage()
	},
}

// to filter the flags with
var flagsRe *regexp.Regexp

// Show the flags
var helpFlags = &cobra.Command{
	Use:   "flags [<regexp to match>]",
	Short: "Show the global flags for rclone",
	Run: func(command *cobra.Command, args []string) {
		if len(args) > 0 {
			re, err := regexp.Compile(args[0])
			if err != nil {
				log.Fatalf("Failed to compile flags regexp: %v", err)
			}
			flagsRe = re
		}
		Root.SetOutput(os.Stdout)
		_ = command.Usage()
	},
}

// Show the backends
var helpBackends = &cobra.Command{
	Use:   "backends",
	Short: "List the backends available",
	Run: func(command *cobra.Command, args []string) {
		showBackends()
	},
}

// Show a single backend
var helpBackend = &cobra.Command{
	Use:   "backend <name>",
	Short: "List full info about a backend",
	Run: func(command *cobra.Command, args []string) {
		if len(args) == 0 {
			Root.SetOutput(os.Stdout)
			_ = command.Usage()
			return
		}
		showBackend(args[0])
	},
}

// runRoot implements the main rclone command with no subcommands
func runRoot(cmd *cobra.Command, args []string) {
	if version {
		ShowVersion()
		resolveExitCode(nil)
	} else {
		_ = cmd.Usage()
		if len(args) > 0 {
			_, _ = fmt.Fprintf(os.Stderr, "Command not found.\n")
		}
		resolveExitCode(errorCommandNotFound)
	}
}

// setupRootCommand sets default usage, help, and error handling for
// the root command.
//
// Helpful example: http://rtfcode.com/xref/moby-17.03.2-ce/cli/cobra.go
func setupRootCommand(rootCmd *cobra.Command) {
	// Add global flags
	configflags.AddFlags(pflag.CommandLine)
	filterflags.AddFlags(pflag.CommandLine)
	rcflags.AddFlags(pflag.CommandLine)

	Root.Run = runRoot
	Root.Flags().BoolVarP(&version, "version", "V", false, "Print the version number")

	cobra.AddTemplateFunc("showGlobalFlags", func(cmd *cobra.Command) bool {
		return cmd.CalledAs() == "flags"
	})
	cobra.AddTemplateFunc("showCommands", func(cmd *cobra.Command) bool {
		return cmd.CalledAs() != "flags"
	})
	cobra.AddTemplateFunc("showLocalFlags", func(cmd *cobra.Command) bool {
		// Don't show local flags (which are the global ones on the root) on "rclone" and
		// "rclone help" (which shows the global help)
		return cmd.CalledAs() != "rclone" && cmd.CalledAs() != ""
	})
	cobra.AddTemplateFunc("backendFlags", func(cmd *cobra.Command, include bool) *pflag.FlagSet {
		backendFlagSet := pflag.NewFlagSet("Backend Flags", pflag.ExitOnError)
		cmd.InheritedFlags().VisitAll(func(flag *pflag.Flag) {
			matched := flagsRe == nil || flagsRe.MatchString(flag.Name)
			if _, ok := backendFlags[flag.Name]; matched && ok == include {
				backendFlagSet.AddFlag(flag)
			}
		})
		return backendFlagSet
	})
	rootCmd.SetUsageTemplate(usageTemplate)
	// rootCmd.SetHelpTemplate(helpTemplate)
	// rootCmd.SetFlagErrorFunc(FlagErrorFunc)
	rootCmd.SetHelpCommand(helpCommand)
	// rootCmd.PersistentFlags().BoolP("help", "h", false, "Print usage")
	// rootCmd.PersistentFlags().MarkShorthandDeprecated("help", "please use --help")

	rootCmd.AddCommand(helpCommand)
	helpCommand.AddCommand(helpFlags)
	helpCommand.AddCommand(helpBackends)
	helpCommand.AddCommand(helpBackend)

	cobra.OnInitialize(initConfig)

}

var usageTemplate = `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if and (showCommands .) .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if and (showLocalFlags .) .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if and (showGlobalFlags .) .HasAvailableInheritedFlags}}

Global Flags:
{{(backendFlags . false).FlagUsages | trimTrailingWhitespaces}}

Backend Flags:
{{(backendFlags . true).FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}

Use "rclone [command] --help" for more information about a command.
Use "rclone help flags" for to see the global flags.
Use "rclone help backends" for a list of supported services.
`

// show all the backends
func showBackends() {
	fmt.Printf("All rclone backends:\n\n")
	for _, backend := range fs.Registry {
		fmt.Printf("  %-12s %s\n", backend.Prefix, backend.Description)
	}
	fmt.Printf("\nTo see more info about a particular backend use:\n")
	fmt.Printf("  rclone help backend <name>\n")
}

func quoteString(v interface{}) string {
	switch v.(type) {
	case string:
		return fmt.Sprintf("%q", v)
	}
	return fmt.Sprint(v)
}

// show a single backend
func showBackend(name string) {
	backend, err := fs.Find(name)
	if err != nil {
		log.Fatal(err)
	}
	var standardOptions, advancedOptions fs.Options
	done := map[string]struct{}{}
	for _, opt := range backend.Options {
		// Skip if done already (eg with Provider options)
		if _, doneAlready := done[opt.Name]; doneAlready {
			continue
		}
		if opt.Advanced {
			advancedOptions = append(advancedOptions, opt)
		} else {
			standardOptions = append(standardOptions, opt)
		}
	}
	optionsType := "standard"
	for _, opts := range []fs.Options{standardOptions, advancedOptions} {
		if len(opts) == 0 {
			continue
		}
		fmt.Printf("### %s Options\n\n", strings.Title(optionsType))
		fmt.Printf("Here are the %s options specific to %s (%s).\n\n", optionsType, backend.Name, backend.Description)
		optionsType = "advanced"
		for _, opt := range opts {
			done[opt.Name] = struct{}{}
			fmt.Printf("#### --%s\n\n", opt.FlagName(backend.Prefix))
			fmt.Printf("%s\n\n", opt.Help)
			fmt.Printf("- Config:      %s\n", opt.Name)
			fmt.Printf("- Env Var:     %s\n", opt.EnvVarName(backend.Prefix))
			fmt.Printf("- Type:        %s\n", opt.Type())
			fmt.Printf("- Default:     %s\n", quoteString(opt.GetValue()))
			if len(opt.Examples) > 0 {
				fmt.Printf("- Examples:\n")
				for _, ex := range opt.Examples {
					fmt.Printf("    - %s\n", quoteString(ex.Value))
					for _, line := range strings.Split(ex.Help, "\n") {
						fmt.Printf("        - %s\n", line)
					}
				}
			}
			fmt.Printf("\n")
		}
	}
}