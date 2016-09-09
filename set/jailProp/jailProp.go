package jailProp

import (
	"fmt"
	"os"
	"os/exec"
)

// jailProp accepts a path where the config file lives and the property.
// A default boolean governs if it's a default property or a jail property.
func Args(path string, prop, value string, def bool) error {
	var config string

	if def {
		config = path + "/.default"
	} else {
		config = path + "/config"
	}

	_, err := exec.Command("/usr/local/bin/uclcmd", "set", "-u", "-f",
		config, "-t", "string", "--", prop, value).Output() // TODO: Use C Lib, support different mountpoints for /ioc

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return nil
}
