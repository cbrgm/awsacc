package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type Config struct {
	Accounts []Account `json:"Accounts"`
}

type Account struct {
	ID   string `json:"Id"`
	Name string `json:"Name"`
}

const (
	globalConfigDir        = ".aws/"
	globalConfigWindowsDir = "aws\\"
	globalConfigFile       = "accounts.json"
)

var (
	// ErrUnableLoadConfig is returned when the config cannot be loaded from config dir
	ErrUnableLoadConfig = errors.New("failed to load config from config directory")
	// ErrUnableLoadFile is returned when the input file cannot be loaded
	ErrUnableLoadFile = errors.New("failed to load input file")
	// ErrUnableReadFile is returned when the input file cannot be read
	ErrUnableReadFile = errors.New("unable to read input file")
	// ErrUnableLocateHomeDir is returned when the home dir cannot be located
	ErrUnableLocateHomeDir = errors.New("unable to locate home directory")
	// ErrInvalidArgument is returned when invalid arguments have been passed to function
	ErrInvalidArgument = errors.New("argument is invalid or was nil")
)

// customConfigDir contains the whole path to config dir. Only access via get/set functions.
var customConfigDir string

// SetConfigDir sets a custom config folder.
func SetConfigDir(configDir string) {
	customConfigDir = configDir
}

// GetConfigDir constructs config folder.
func GetConfigDir() (string, error) {
	if customConfigDir != "" {
		return customConfigDir, nil
	}
	homeDir, e := homedir.Dir()
	if e != nil {
		return "", ErrUnableLocateHomeDir
	}
	var configDir string
	// For windows the path is slightly different
	if runtime.GOOS == "windows" {
		configDir = filepath.Join(homeDir, globalConfigWindowsDir)
	} else {
		configDir = filepath.Join(homeDir, globalConfigDir)
	}
	return configDir, nil
}

// GetConfigPath constructs the configuration path.
func GetConfigPath() (string, error) {
	if customConfigDir != "" {
		return filepath.Join(customConfigDir, globalConfigFile), nil
	}
	dir, err := GetConfigDir()
	if err != nil {
		return "", ErrUnableLocateHomeDir
	}
	return filepath.Join(dir, globalConfigFile), nil
}

// LoadConfig loads the config file
func LoadConfig() (*Config, error) {
	file, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadFile(file)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("%s, %s", ErrUnableLoadConfig, err)
	}

	var list = Config{
		Accounts: []Account{},
	}
	err = json.Unmarshal(b, &list)
	if err != nil {
		return nil, fmt.Errorf("%s, %s", ErrUnableLoadConfig, err)
	}

	return &list, nil
}

func MapIDsToName(c *Config) map[string]string {
	result := make(map[string]string)
	for _, item := range c.Accounts {
		result[item.ID] = strings.ToLower(item.Name)
	}
	return result
}

func MapNameToIDs(c *Config) map[string]string {
	result := make(map[string]string)
	for _, item := range c.Accounts {
		result[strings.ToLower(item.Name)] = item.ID
	}
	return result
}
