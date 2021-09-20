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
	Arr      []int
}
type CustomConfig struct {
	Int    int
	String string
	Arr    []int
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
		nil,
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
		nil,
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
		nil,
	}
	writeFile(src, "config.yaml")
	defer os.Remove("config.yaml")

	customSrc := CustomConfig{
		2,
		"Hello_Custom",
		nil,
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
		nil,
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

func TestCustom01(t *testing.T) {
	src := Config{
		1,
		"Hello",
		true,
		10.10,
		time.Unix(1000, 0),
		time.Second,
		nil,
	}
	writeFile(src, "config.yaml")
	defer os.Remove("config.yaml")

	customSrc := CustomConfig{
		2,
		"Hello_Custom",
		nil,
	}
	writeFile(customSrc, "config.custom.yaml")
	defer os.Remove("config.custom.yaml")

	custom01Src := CustomConfig{
		3,
		"Hello_Custom.01",
		nil,
	}
	writeFile(custom01Src, "config.custom.01.yaml")
	defer os.Remove("config.custom.01.yaml")

	var dst Config
	configFile := LoadConfigFromFile(&dst, "config.yaml", nil)
	assert.Equal(t, "config.custom.01.yaml", configFile)

	assert.Equal(t, custom01Src.Int, dst.Int)
	assert.Equal(t, custom01Src.String, dst.String)
	assert.Equal(t, src.Bool, dst.Bool)
	assert.Equal(t, src.Float, dst.Float)
	assert.Equal(t, src.Time, dst.Time)
	assert.Equal(t, src.Duration, dst.Duration)
}

func TestCustom02(t *testing.T) {
	src := Config{
		1,
		"Hello",
		true,
		10.10,
		time.Unix(1000, 0),
		time.Second,
		nil,
	}
	writeFile(src, "config.yaml")
	defer os.Remove("config.yaml")

	customSrc := CustomConfig{
		2,
		"Hello_Custom",
		nil,
	}
	writeFile(customSrc, "config.custom.yaml")
	defer os.Remove("config.custom.yaml")

	custom01Src := CustomConfig{
		3,
		"Hello_Custom.01",
		[]int{0, 1, 2},
	}
	writeFile(custom01Src, "config.custom.01.yaml")
	defer os.Remove("config.custom.01.yaml")

	custom02Src := CustomConfig{
		4,
		"Hello_Custom.02",
		[]int{3, 4},
	}
	writeFile(custom02Src, "config.custom.02.yaml")
	defer os.Remove("config.custom.02.yaml")

	var dst Config
	configFile := LoadConfigFromFile(&dst, "config.yaml", nil)
	assert.Equal(t, "config.custom.02.yaml", configFile)

	assert.Equal(t, custom02Src.Int, dst.Int)
	assert.Equal(t, custom02Src.String, dst.String)
	assert.Equal(t, custom02Src.Arr, dst.Arr)
	assert.Equal(t, src.Bool, dst.Bool)
	assert.Equal(t, src.Float, dst.Float)
	assert.Equal(t, src.Time, dst.Time)
	assert.Equal(t, src.Duration, dst.Duration)
}

func TestDumpYaml(t *testing.T) {
	c := struct {
		Z string
	}{Z: "ZZ"}
	r, e := DumpYaml(c)
	assert.NoError(t, e)
	assert.Equal(t, "z: ZZ\n", r)
}

func TestDumpJson(t *testing.T) {
	c := struct {
		Z string
	}{Z: "ZZ"}
	r, e := DumpJson(c)
	assert.NoError(t, e)
	assert.Equal(t, "{\"Z\":\"ZZ\"}", r)
}
