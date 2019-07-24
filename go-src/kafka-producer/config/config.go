package config

import (
	"io/ioutil"

	"github.com/imdario/mergo"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var Default Config

type Config struct {
	Kafka Kafka `yaml:"kafka"`
}

type Kafka struct {
	Brokers []string `yaml:"brokers"`
	Topic   string   `yaml:"topic"`
}

func Init(defaultPath string, overriddenPaths ...string) error {

	//
	// read default config
	//

	log.Infof("read [%s]", defaultPath)

	defaultYaml, err := read(defaultPath)
	if err != nil {
		return err
	}

	for _, path := range overriddenPaths {

		log.Infof("read [%s]", path)

		overriddenYaml, err := read(path)
		if err != nil {
			return err
		}

		if err := mergo.Merge(&defaultYaml, overriddenYaml, mergo.WithOverride); err != nil {
			return err
		}

		log.Debugf("merged defaultYaml=[%#v]", defaultYaml)
	}

	mergedData, err := yaml.Marshal(defaultYaml)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(mergedData, &Default)
	if err != nil {
		return err
	}

	log.Debugln("===============")
	log.Debugf("%#v", Default)
	log.Debugln("===============")

	return nil
}

func read(yamlPath string) (map[string]interface{}, error) {

	log.Debugf("reading yaml [%s]", yamlPath)

	data, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		return nil, err
	}

	log.Debugf("read yaml with value [%s]", string(data))

	var yamlObject map[string]interface{}
	err = yaml.Unmarshal(data, &yamlObject)
	if err != nil {
		return nil, err
	}

	log.Debugf("read yaml object [%#v]", yamlObject)

	return yamlObject, nil
}
