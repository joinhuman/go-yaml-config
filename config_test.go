package config

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/go-yaml/yaml"
	"github.com/stretchr/testify/assert"
)

type Config struct {
	Int      int
	String   string
	Bool     bool
	Float    float64
	Time     time.Time
	Duration time.Duration
}
type CustomConfig struct {
	Int    int
	String string
}

func writeFile(config interface{}, filename string) {
	b, _ := yaml.Marshal(config)
	if err := ioutil.WriteFile(filename, b, 0666); err != nil {
		panic(err)
	}
}
func TestMain(t *testing.T) {
	src := Config{
		1,
		"Hello",
		true,
		10.10,
		time.Unix(1000, 0),
		time.Second,
	}
	writeFile(src, "config.yaml")
	defer os.Remove("config.yaml")

	var dst Config

	configFile, err := LoadConfigFromFile(&dst, "config.yaml", nil)
	assert.NoError(t, err)
	assert.Equal(t, "config.yaml", configFile)
	assert.Equal(t, src.Int, dst.Int)
	assert.Equal(t, src.String, dst.String)
	assert.Equal(t, src.Bool, dst.Bool)
	assert.Equal(t, src.Float, dst.Float)
	assert.Equal(t, src.Time, dst.Time)
	assert.Equal(t, src.Duration, dst.Duration)
}

func TestDefault(t *testing.T) {
	os.Remove("config.yaml")
	src := Config{
		1,
		"Hello",
		true,
		10.10,
		time.Unix(1000, 0),
		time.Second,
	}
	var dst Config
	configFile, err := LoadConfigFromFile(&dst, "config.yaml", src)
	assert.NoError(t, err)
	assert.Empty(t, configFile)
	assert.Equal(t, src.Int, dst.Int)
	assert.Equal(t, src.String, dst.String)
	assert.Equal(t, src.Bool, dst.Bool)
	assert.Equal(t, src.Float, dst.Float)
	assert.Equal(t, src.Time, dst.Time)
	assert.Equal(t, src.Duration, dst.Duration)
}

func TestNoDefaultError(t *testing.T) {
	var dst Config
	_, err := LoadConfigFromFile(&dst, "config.yaml", nil)
	assert.Error(t, err)

}
func TestCustom(t *testing.T) {
	src := Config{
		1,
		"Hello",
		true,
		10.10,
		time.Unix(1000, 0),
		time.Second,
	}
	writeFile(src, "config.yaml")
	defer os.Remove("config.yaml")

	customSrc := CustomConfig{
		2,
		"Hello_Custom",
	}
	writeFile(customSrc, "config.custom.yaml")
	defer os.Remove("config.custom.yaml")

	var dst Config
	configFile, err := LoadConfigFromFile(&dst, "config.yaml", nil)
	assert.NoError(t, err)
	assert.Equal(t, "config.custom.yaml", configFile)

	assert.Equal(t, customSrc.Int, dst.Int)
	assert.Equal(t, customSrc.String, dst.String)
	assert.Equal(t, src.Bool, dst.Bool)
	assert.Equal(t, src.Float, dst.Float)
	assert.Equal(t, src.Time, dst.Time)
	assert.Equal(t, src.Duration, dst.Duration)
}

func TestBrokenFilePanic(t *testing.T) {
	assert.NoError(t, ioutil.WriteFile("config.yaml", []byte("hello"), 0666))
	defer os.Remove("config.yaml")
	var dst Config
	_, err := LoadConfigFromFile(&dst, "config.yaml", nil)
	assert.Error(t, err)
}

func TestBrokenFileDefault(t *testing.T) {
	if err := ioutil.WriteFile("config.yaml", []byte("hello"), 0666); err != nil {
		panic(err)
	}
	defer os.Remove("config.yaml")
	src := Config{
		1,
		"Hello",
		true,
		10.10,
		time.Unix(1000, 0),
		time.Second,
	}
	var dst Config

	_, err := LoadConfigFromFile(&dst, "config.yaml", src)
	assert.NoError(t, err)
}
