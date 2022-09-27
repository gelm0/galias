package alias

import (
	"testing"
	"os"
	"path/filepath"
	"encoding/json"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func createExampleConfig() TopLevelConfig {
	testAlias := Alias{
		Name: "test-command",
		Description: "test",
		Variables: []string{"~/"},
	}
	testConfig := Config{
		Name: "test-alias",
		Command: "command ${} ${} ${} ${} ${}",
		Description: "Change directory",
		Alias: []Alias{testAlias},

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

}
func TestProcessTemplateIncorrectVariables(t *testing.T) {

}
func TestProcessTemplateVariablesNoInterpolation(t *testing.T) {

}