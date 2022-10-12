package alias

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func createExampleConfig() TopLevelConfig {
	testAlias := Alias{
		Name:        "test-command",
		Description: "test",
		Variables:   []string{"~/"},
	}
	testConfig := Config{
		Name:        "test-alias",
		Command:     "command ${} ${} ${} ${} ${}",
		Description: "Change directory",
		Alias:       []Alias{testAlias},
	}
	return TopLevelConfig{
		Config: []Config{testConfig},
	}
}

func TestExampleConfig(t *testing.T) {
	appFs = &afero.Afero{
		Fs: afero.NewMemMapFs(),
	}
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	AddExampleConfig()

	b, err := appFs.ReadFile(filepath.Join(home, ".galias"))
	if err != nil {
		panic(err)
	}
	assert.GreaterOrEqual(t, len(b), 1)
	var config TopLevelConfig
	err = json.Unmarshal(b, &config)
	assert.Nil(t, err)

	assert.Equal(t, config.Config[0].Name, "cd")
	assert.Equal(t, config.Config[0].Command, "cd ${}")
	assert.Equal(t, config.Config[0].Alias[0].Name, "home")

}
func TestProcessCorrectTemplate(t *testing.T) {
	template := "this is ${} ${}"
	args := []string{"a", "test"}

	processed, err := processTemplate(template, args, []string{})

	assert.Nil(t, err)
	assert.Equal(t, processed, fmt.Sprintf("this is %s %s", args[0], args[1]))

	extraVars := []string{"something", "extra"}
	processed, err = processTemplate(template, args, extraVars)

	assert.Nil(t, err)
	assert.Equal(t, processed, fmt.Sprintf("this is %s %s %s %s", args[0], args[1], extraVars[0], extraVars[1]))

}
func TestProcessTemplateIncorrectVariables(t *testing.T) {
	template := "this is ${}"

	args := []string{}
	empty, err := processTemplate(template, args, []string{})

	assert.Empty(t, empty)
	assert.NotNil(t, err)

	args = []string{"1", "2"}
	empty, err = processTemplate(template, args, []string{})

	assert.Empty(t, empty)
	assert.NotNil(t, err)
}
