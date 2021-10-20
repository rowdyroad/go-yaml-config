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

	configFile := LoadConfigFromFile(&dst, "config.yaml", nil)
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
	configFile := LoadConfigFromFile(&dst, "config.yaml", src)
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
	assert.Panics(t, func() {
		LoadConfigFromFile(&dst, "config.yaml", nil)
	})

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
	configFile := LoadConfigFromFile(&dst, "config.yaml", nil)
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
	assert.Panics(t, func() { LoadConfigFromFile(&dst, "config.yaml", nil) })
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
	assert.NotPanics(t, func() {
		LoadConfigFromFile(&dst, "config.yaml", src)
	})
}

func TestEnv(t *testing.T) {

	err := ioutil.WriteFile("config.yaml", []byte(
		`string: ${STR_ENV}
int: ${INT_ENV}
bool: ${BOOL_ENV}
float: ${FLOAT_ENV}
time: ${TIME_ENV}
duration: ${DURATION_ENV}
`), 0777)
	assert.NoError(t, err)
	defer os.Remove("config.yaml")
	tt := time.Now()

	os.Setenv("STR_ENV", "str")
	os.Setenv("INT_ENV", "1")
	os.Setenv("BOOL_ENV", "true")
	os.Setenv("FLOAT_ENV", "10.10")
	os.Setenv("TIME_ENV", tt.Format(time.RFC3339Nano))
	os.Setenv("DURATION_ENV", time.Second.String())

	var cfg Config
	assert.NotPanics(t, func() { LoadConfig(&cfg, "config.yaml", nil) })

	assert.Equal(t, 1, cfg.Int)
	assert.Equal(t, "str", cfg.String)
	assert.Equal(t, true, cfg.Bool)
	assert.Equal(t, 10.10, cfg.Float)
	assert.True(t, tt.Equal(cfg.Time))
	assert.Equal(t, time.Second, cfg.Duration)

}

func TestGetEnvWithDefault(t *testing.T) {
	os.Remove("config.yaml")

	type configCustom struct {
		String1 string
		String2 string
		String3 string
		String4 string
		String5 string
	}

	err := ioutil.WriteFile("config.yaml", []byte(
		`string1: ${TEST_ENV1}
string2: ${TEST_ENV2=hello}
string3: ${TEST_ENV3=}
string4: ${TEST_ENV4}
string5: ${TEST_ENV5=hello}
`), 0777)
	assert.NoError(t, err)
	defer os.Remove("config.yaml")

	os.Setenv("TEST_ENV1", "test_val")
	defer os.Unsetenv("TEST_ENV1")

	os.Setenv("TEST_ENV5", "test_val")
	defer os.Unsetenv("TEST_ENV5")

	var cfg configCustom
	assert.NotPanics(t, func() { LoadConfigFromFile(&cfg, "config.yaml", nil) })

	assert.Equal(t, "test_val", cfg.String1)
	assert.Equal(t, "hello", cfg.String2)
	assert.Equal(t, "", cfg.String3)
	assert.Equal(t, "", cfg.String4)
	assert.Equal(t, "test_val", cfg.String5)
}
