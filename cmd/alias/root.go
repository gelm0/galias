package alias

import (
	"fmt"
	"os"
	"path/filepath"

	alias "github.com/gelm0/go-alias/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var reservedCommands = map[string]int{"add-example-config": 0}

var (
	config *alias.TopLevelConfig
)

var rootCmd = &cobra.Command{
	Use:   "galias",
	Short: "galias - an interpolating alias command line interpreter",
}

func initConfig() *alias.TopLevelConfig {
	home, err := os.UserHomeDir()
	alias.ExitIfErr(err)

	viper.AddConfigPath(filepath.Join(home))
	viper.SetConfigType("json")
	viper.SetConfigName(".galias")
	err = viper.ReadInConfig()
	alias.ExitIfErr(err)

	conf := &alias.TopLevelConfig{}
	err = viper.Unmarshal(conf)
	alias.ExitIfErr(err)

	return conf
}

func addReservedCommands() {
	addCommand := &cobra.Command{
		Use:   "add-example-config",
		Short: "generate a example config",
		Run: func(cmd *cobra.Command, args []string) {
			alias.AddExampleConfig()
		},
	}
	rootCmd.AddCommand(addCommand)
}

func isReserved(c alias.Config) bool {
	_, ok := reservedCommands[c.Name]
	return ok
}

func addConfigCommands() {
	for _, c := range config.Config {
		if isReserved(c) {
			continue;
		}
		var description string
		genDescription := "user-generated"
		if c.Description == "" {
			description = genDescription
		} else {
			description = c.Description
		}
		command := &cobra.Command{
			Use:   c.Name,
			Short: description,
		}
		for _, a := range c.Alias {
			alias := &cobra.Command{
				Use:   a.Name,
				Short: a.Description,
				Run: func(cmd *cobra.Command, args []string) {
					fmt.Fprintf(os.Stdout, "%s", args)
					fmt.Println()
					alias.RunCommand(c.Command, a.Variables, args)
				},
			}
			command.AddCommand(alias)
		}
		rootCmd.AddCommand(command)
	}
}

func init() {
	config = initConfig()
	addReservedCommands()
	addConfigCommands()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "'%s'", err)
		os.Exit(1)
	}
}
