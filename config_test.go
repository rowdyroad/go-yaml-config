package config

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/go-yaml/yaml"
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
	if configFile != "config.yaml" {
		t.Error()
	}
	if dst.Int != src.Int {
		t.Error()
	}
	if dst.String != src.String {
		t.Error()
	}
	if dst.Bool != src.Bool {
		t.Error()
	}

	if dst.Float != src.Float {
		t.Error()
	}
	if dst.Time != src.Time {
		t.Error()
	}
	if dst.Duration != src.Duration {
		t.Error()
	}
}

func TestDefault(t *testing.T) {
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
	if configFile != "" {
		t.Error()
	}
	if dst.Int != src.Int {
		t.Error()
	}
	if dst.String != src.String {
		t.Error()
	}
	if dst.Bool != src.Bool {
		t.Error()
	}

	if dst.Float != src.Float {
		t.Error()
	}
	if dst.Time != src.Time {
		t.Error()
	}
	if dst.Duration != src.Duration {
		t.Error()
	}
}

func TestNoDefaultError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("not panic")
		}
	}()
	var dst Config
	LoadConfigFromFile(&dst, "config.yaml", nil)
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
	if configFile != "config.custom.yaml" {
		t.Error()
	}

	if dst.Int != customSrc.Int {
		t.Error()
	}
	if dst.String != customSrc.String {
		t.Error()
	}
	if dst.Bool != src.Bool {
		t.Error()
	}

	if dst.Float != src.Float {
		t.Error()
	}
	if dst.Time != src.Time {
		t.Error()
	}
	if dst.Duration != src.Duration {
		t.Error()
	}
}

func TestBrokenFilePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error()
		}
	}()
	if err := ioutil.WriteFile("config.yaml", []byte("hello"), 0666); err != nil {
		panic(err)
	}
	defer os.Remove("config.yaml")
	var dst Config
	LoadConfigFromFile(&dst, "config.yaml", nil)
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
	LoadConfigFromFile(&dst, "config.yaml", src)
}
