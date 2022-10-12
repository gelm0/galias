package alias

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/riywo/loginshell"
	"github.com/spf13/afero"
	"github.com/valyala/fasttemplate"
)

var appFs = &afero.Afero{
	Fs: afero.NewOsFs(),
}

type TopLevelConfig struct {
	Config []Config `json:"config"`
}

type Config struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Command     string  `json:"command"`
	Alias       []Alias `json:"alias"`
}

type Alias struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Variables   []string `json:"variables"`
}

func ExitIfErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func AddExampleConfig() {
	homeAlias := Alias{
		Name:        "home",
		Description: "home directory",
		Variables:   []string{"~/"},
	}
	firstConfig := Config{
		Name:        "cd",
		Command:     "cd ${}",
		Description: "Change directory",
		Alias:       []Alias{homeAlias},
	}
	overlay := TopLevelConfig{
		Config: []Config{firstConfig},
	}
	file, err := json.MarshalIndent(overlay, "", " ")
	if err != nil {
		fmt.Println(err)
	}
	home, err := os.UserHomeDir()
	ExitIfErr(err)
	err = appFs.WriteFile(filepath.Join(home, ".galias"), file, 0644)
	ExitIfErr(err)
}

func processTemplate(command string, vars, args []string) (exec string, err error) {
	template := fasttemplate.New(command, "${", "}")
	ii := 0
	maxLen := len(vars)
	exec, err = template.ExecuteFuncStringWithErr(func(w io.Writer, tag string) (int, error) {
		if ii >= maxLen {
			return 0, fmt.Errorf("cannot match template to supplied variables\n, ")
		}
		n, err := w.Write([]byte(vars[ii]))
		ii += 1
		return n, err
	})
	if err != nil {
		return "", err
	}
	// Check that template matches arguments
	if ii != maxLen {
		return "", fmt.Errorf("supplied %d arguments, but want to interpolate %d argument(s)", maxLen, ii)
	}
	if len(args) > 0 {
		for _, arg := range args {
			exec += " " + arg
		}
	}
	return exec, nil
}

func execShell(command string) error {
	var shellCommandString string
	switch runtime.GOOS {
	case "linux":
	case "darwin":
	case "freebsd":
	case "openbsd":
		/*
			This will most likely will not work for all shells,
			 but supported for bash,zsh,ksh,fish which should cover most cases
			 if anyone wishes to extend this or have a link to some extensive documentation
			 around this issue it would be much appreciated
		*/
		shellCommandString = "-c"
	default:
		return fmt.Errorf("environment not supported %s", runtime.GOOS)

	}
	shell, err := loginshell.Shell()
	ExitIfErr(err)
	cmd := exec.Command(shell, shellCommandString, command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
	return nil

}

func RunCommand(command string, vars []string, args []string) {
	exec, err := processTemplate(command, vars, args)
	ExitIfErr(err)
	err = execShell(exec)
	ExitIfErr(err)
}
