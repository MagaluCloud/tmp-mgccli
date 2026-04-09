package config

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"

	"github.com/magaluCloud/mgccli/cmd/common/structs"
	"github.com/magaluCloud/mgccli/cmd/common/validator"
	"github.com/magaluCloud/mgccli/cmd/common/workspace"
	"gopkg.in/yaml.v3"
)

type Config interface {
	Get(name string) (*ConfigItem, error)
	GetSchema(name string) (*ConfigSchema, error)
	Set(name string, value any) error
	Delete(name string) error
	List() (map[string]*ConfigItem, error)
	Value(name string) (Value, error)
	Write() error
}

type ConfigItem struct {
	Name        string  `json:"name"`
	Value       any     `json:"value"`
	Type        string  `json:"type"`
	Description string  `json:"description"`
	Validator   *string `json:"validator,omitempty"`
	Default     any     `json:"default"`
	Scope       string  `json:"scope"`
}

type ConfigSchema struct {
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Description string  `json:"description"`
	Validator   *string `json:"validator,omitempty"`
	Default     any     `json:"default"`
	Scope       string  `json:"scope"`
}

type CliConfig struct {
	Items map[string]*ConfigItem
}

type ConfigYaml struct {
	ChunkSize        int    `yaml:"chunk_size,omitempty"`
	Workers          int    `yaml:"workers,omitempty"`
	DefaultOutput    string `yaml:"default_output,omitempty"`
	Region           string `yaml:"region,omitempty"`
	Env              string `yaml:"env,omitempty"`
	Debug            bool   `yaml:"debug,omitempty"`
	NoConfirm        bool   `yaml:"no_confirm,omitempty"`
	RawOutput        bool   `yaml:"raw_output,omitempty"`
	Lang             string `yaml:"lang,omitempty"`
	ServerURL        string `yaml:"server_url,omitempty"`
	VersionLastCheck string `yaml:"version_last_check,omitempty"`
}

type config struct {
	cliConfig  CliConfig
	configYaml ConfigYaml
	workspace  workspace.Workspace
}

func StrToStrPtr(str string) *string {
	return &str
}

func NewConfig(workspace workspace.Workspace) Config {
	configFile := path.Join(workspace.Dir(), "cli.yaml")
	configYaml, err := structs.LoadFileToStruct[ConfigYaml](configFile)
	if err != nil {
		panic(err)
	}

	cliConfig := CliConfig{
		Items: make(map[string]*ConfigItem, 10),
	}

	cliConfig.Items[nameToKey("chunk_size")] = &ConfigItem{
		Name:        keyToName("chunk_size"),
		Value:       configYaml.ChunkSize,
		Type:        "int",
		Description: "Chunk size to consider when doing multipart requests. Specified in Mb",
		Validator:   StrToStrPtr("minimum=8,maximum=5120"),
		Default:     8,
		Scope:       "object-storage",
	}

	cliConfig.Items["workers"] = &ConfigItem{
		Name:        "workers",
		Value:       configYaml.Workers,
		Type:        "int",
		Description: "umber of routines that spawn to do parallel operations",
		Validator:   StrToStrPtr("minimum=1"),
		Default:     5,
		Scope:       "object-storage",
	}

	cliConfig.Items[nameToKey("default_output")] = &ConfigItem{
		Name:        keyToName("default_output"),
		Value:       configYaml.DefaultOutput,
		Type:        "string",
		Description: "Default output string to be used when no other is specified",
		Validator:   StrToStrPtr("oneof=json,table"),
		Default:     "json",
		Scope:       "global",
	}

	cliConfig.Items["region"] = &ConfigItem{
		Name:        "region",
		Value:       configYaml.Region,
		Type:        "string",
		Description: "Region to reach the service",
		Validator:   StrToStrPtr("oneof=br-se1,br-ne1,br-mgl1"),
		Default:     "br-se1",
		Scope:       "global",
	}

	cliConfig.Items["env"] = &ConfigItem{
		Name:        "env",
		Value:       configYaml.Env,
		Type:        "string",
		Description: "Environment",
		Validator:   StrToStrPtr("oneof=prod,pre-prod"),
		Default:     "prod",
		Scope:       "global",
	}
	cliConfig.Items["debug"] = &ConfigItem{
		Name:        "debug",
		Value:       configYaml.Debug,
		Type:        "bool",
		Description: "Debug",
		Default:     false,
		Scope:       "global",
	}

	cliConfig.Items[nameToKey("raw_output")] = &ConfigItem{
		Name:        keyToName("raw_output"),
		Value:       configYaml.RawOutput,
		Type:        "bool",
		Description: "Raw output",
		Default:     false,
		Scope:       "global",
	}
	cliConfig.Items["lang"] = &ConfigItem{
		Name:        "lang",
		Value:       configYaml.Lang,
		Type:        "string",
		Description: "Language",
		Validator:   StrToStrPtr("oneof=en-US,pt-BR"),
		Default:     "en-US",
		Scope:       "global",
	}
	cliConfig.Items[nameToKey("server_url")] = &ConfigItem{
		Name:        keyToName("server_url"),
		Value:       configYaml.ServerURL,
		Type:        "string",
		Description: "Server URL",
		Default:     "https://api.magalu.cloud",
		Scope:       "global",
	}

	return &config{workspace: workspace, cliConfig: cliConfig, configYaml: configYaml}
}

func (c *config) Value(name string) (Value, error) {
	item, err := c.Get(name)
	if err != nil {
		return nil, err
	}
	return NewValue(item.Value), nil
}

func (c *config) Get(name string) (*ConfigItem, error) {
	value, ok := c.cliConfig.Items[nameToKey(name)]
	if !ok {
		return nil, fmt.Errorf("config %s not found", name)
	}

	if value.Value == nil {
		value.Value = value.Default
		return value, nil
	}
	if reflect.ValueOf(value.Value).IsZero() {
		value.Value = value.Default
		return value, nil
	}
	return value, nil
}

func (c *config) Set(name string, value any) error {
	item, err := c.Get(name)
	if err != nil {
		return err
	}

	if value != nil {
		updatedValue, err := formattedValue(name, item.Type, value)
		if err != nil {
			return err
		}

		item.Value = updatedValue
	}

	if item.Validator != nil {
		err = validator.NewValidator(item.Value, *item.Validator).Validate()
		if err != nil {
			return err
		}
	}

	err = structs.Set(&c.configYaml, keyToName(name), value)
	if err != nil {
		return err
	}
	err = c.Write()
	if err != nil {
		return err
	}
	return nil
}

func (c *config) Delete(name string) error {
	return c.Set(name, nil)
}

func (c *config) Write() error {
	data, err := yaml.Marshal(c.configYaml)
	if err != nil {
		return err
	}
	err = os.WriteFile(path.Join(c.workspace.Dir(), "cli.yaml"), data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (c *config) List() (map[string]*ConfigItem, error) {
	return c.cliConfig.Items, nil
}

func (c *config) GetSchema(name string) (*ConfigSchema, error) {
	config := c.cliConfig.Items[name]
	if config == nil {
		return nil, fmt.Errorf(`config "%s" not found`, name)
	}

	return &ConfigSchema{
		Name:        config.Name,
		Type:        config.Type,
		Description: config.Description,
		Validator:   config.Validator,
		Default:     config.Default,
		Scope:       config.Scope,
	}, nil
}

// name_to_key -> name-to-key
func nameToKey(name string) string {
	return strings.ToLower(strings.ReplaceAll(name, "_", "-"))
}

// key-to-name -> key_to_name
func keyToName(key string) string {
	return strings.ToLower(strings.ReplaceAll(key, "-", "_"))
}

func formattedValue(configName string, configType string, value any) (any, error) {
	switch configType {
	case "string":
		return anyToString(value), nil
	case "int":
		val, err := strconv.Atoi(anyToString(value))
		if err != nil {
			return nil, err
		}
		return val, nil
	case "bool":
		val, err := strconv.ParseBool(anyToString(value))
		if err != nil {
			return nil, err
		}
		return val, nil
	default:
		return nil, fmt.Errorf("unsupported type for config %s", configName)
	}
}
