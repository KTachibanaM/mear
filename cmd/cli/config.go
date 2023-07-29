package main

import (
	"fmt"
	"os"
	"os/user"
	"path"

	"github.com/pelletier/go-toml/v2"
)

func ask_for_config_value(prompt string) (string, error) {
	var value string
	for {
		fmt.Printf("%s: ", prompt)
		_, err := fmt.Scanln(&value)
		if err != nil {
			return "", err
		}
		if value != "" {
			return value, nil
		}
	}
}

func get_config(section, key, prompt string) string {
	// Get home dir
	current_user, err := user.Current()
	if err != nil {
		fmt.Printf("failed to get current user: %v", err)
		return ""
	}
	home_dir := current_user.HomeDir

	// Create .config dir
	config_dir := path.Join(home_dir, ".config")
	err = os.MkdirAll(config_dir, os.ModePerm)
	if err != nil {
		fmt.Printf("failed to create config dir: %v", err)
		return ""
	}

	config_path := path.Join(config_dir, "mear.toml")
	if _, err := os.Stat(config_path); os.IsNotExist(err) {
		// Create empty config file if it does not exist
		empty_config_bytes, err := toml.Marshal(map[string]interface{}{})
		if err != nil {
			fmt.Printf("failed to marshal empty config: %v", err)
			return ""
		}
		err = os.WriteFile(config_path, empty_config_bytes, os.ModePerm)
		if err != nil {
			fmt.Printf("failed to write empty config file: %v", err)
			return ""
		}
	}

	// Read config file
	config_bytes, err := os.ReadFile(config_path)
	if err != nil {
		fmt.Printf("failed to read config file: %v", err)
		return ""
	}
	var config map[string]interface{}
	err = toml.Unmarshal(config_bytes, &config)
	if err != nil {
		fmt.Printf("failed to unmarshal config file: %v", err)
		return ""
	}

	// Read config value or ask for it
	if config == nil {
		config = map[string]interface{}{}
	}
	var value string
	if _, ok := config[section]; ok {
		if _, ok := config[section].(map[string]interface{})[key]; ok {
			return config[section].(map[string]interface{})[key].(string)
		} else {
			value, err = ask_for_config_value(prompt)
			if err != nil {
				fmt.Printf("failed to ask for config value: %v", err)
				return ""
			}
		}
	} else {
		config[section] = map[string]interface{}{}
		value, err = ask_for_config_value(prompt)
		if err != nil {
			fmt.Printf("failed to ask for config value: %v", err)
			return ""
		}
	}

	// Write back config file
	config[section].(map[string]interface{})[key] = value
	config_bytes, err = toml.Marshal(config)
	if err != nil {
		fmt.Printf("failed to marshal config file: %v", err)
		return ""
	}
	err = os.WriteFile(config_path, config_bytes, os.ModePerm)
	if err != nil {
		fmt.Printf("failed to write back config file: %v", err)
		return ""
	}

	return value
}
