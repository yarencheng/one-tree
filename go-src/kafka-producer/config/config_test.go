package config

import (
	"io/ioutil"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func TestRead(t *testing.T) {

	// arrange

	yamlFile, err := ioutil.TempFile("", "gotest")
	require.NoError(t, err)
	defer os.Remove(yamlFile.Name())

	err = ioutil.WriteFile(yamlFile.Name(), []byte(
		`
aaa: 123
`,
	), 0644)
	require.NoError(t, err)

	// action

	data, err := read(yamlFile.Name())

	// assert

	assert.NoError(t, err)
	assert.Equal(t, map[string]interface{}{
		"aaa": 123,
	}, data)
}

func TestInit(t *testing.T) {

	// arrange: default config

	defaultFile, err := ioutil.TempFile("", "gotest")
	require.NoError(t, err)
	defer os.Remove(defaultFile.Name())

	err = ioutil.WriteFile(defaultFile.Name(), []byte(
		`
kafka:
  brokers:
    - aa
    - bb
`,
	), 0644)
	require.NoError(t, err)

	// arrange: overide config

	overideFile, err := ioutil.TempFile("", "gotest")
	require.NoError(t, err)
	defer os.Remove(overideFile.Name())

	err = ioutil.WriteFile(overideFile.Name(), []byte(
		`
kafka:
  brokers:
    - cc
    - dd
`,
	), 0644)
	require.NoError(t, err)

	// action

	err = Init(defaultFile.Name(), overideFile.Name())

	// assert

	assert.NoError(t, err)
	assert.Equal(t, Config{
		Kafka: Kafka{
			Brokers: []string{"cc", "dd"},
		},
	}, Default)
}
