package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/gophercloud/gophercloud/internal/cli/util"
	"gopkg.in/ini.v1"
	"gopkg.in/urfave/cli.v1"
)

func configure(c *cli.Context) error {
	fmt.Println("\nThis interactive session will walk you through creating\n" +
		"a profile in your configuration file. You may fill in all or none of the\n" +
		"values.\n")
	reader := bufio.NewReader(os.Stdin)
	m := map[string]string{
		"username": "",
		"password": "",
		"region":   "",
	}

	fmt.Print("Username: ")
	username, _ := reader.ReadString('\n')
	m["username"] = strings.TrimSpace(username)

	fmt.Print("Password: ")
	pwd, err := terminal.ReadPassword(syscall.Stdin)
	if err != nil {
		return err
	}
	m["password"] = string(pwd)

	fmt.Print("Region: ")
	region, _ := reader.ReadString('\n')
	m["region"] = strings.ToUpper(strings.TrimSpace(region))

	fmt.Print("Profile Name (leave blank to create a default profile): ")
	profile, _ := reader.ReadString('\n')
	profile = strings.TrimSpace(profile)

	configFile, err := util.ConfigFile()
	if err != nil {
		return err
	}

	var cfg *ini.File
	cfg, err = ini.Load(configFile)
	if err != nil {
		// fmt.Printf("Error loading config file: %s\n", err)
		cfg = ini.Empty()
	}

	if strings.TrimSpace(profile) == "" || strings.ToLower(profile) == "default" {
		profile = "DEFAULT"
	}

	var section *ini.Section
	for {
		if section, err = cfg.GetSection(profile); err == nil && len(section.Keys()) != 0 {
			fmt.Printf("\nA profile named %s already exists. Overwrite? (y/n): ", profile)
			choice, _ := reader.ReadString('\n')
			choice = strings.TrimSpace(choice)
			switch strings.ToLower(choice) {
			case "y", "yes":
				break
			case "n", "no":
				fmt.Print("Profile Name: ")
				profile, _ = reader.ReadString('\n')
				profile = strings.TrimSpace(profile)
				continue
			default:
				continue
			}
			break
		}
		break
	}

	section, err = cfg.NewSection(profile)
	if err != nil {
		//fmt.Printf("Error creating new section [%s] in config file: %s\n", profile, err)
		return err
	}

	for key, val := range m {
		section.NewKey(key, val)
	}

	err = cfg.SaveTo(configFile)
	if err != nil {
		//fmt.Printf("Error saving config file: %s\n", err)
		return err
	}

	if profile == "DEFAULT" {
		fmt.Printf("\nCreated new default profile for username %s", username)
	} else {
		fmt.Printf("\nCreated profile %s with username %s", profile, username)
	}

	return nil

}
