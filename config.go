package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/BurntSushi/toml"
)

type config struct {
	Java   configJava   `toml:"java"`
	Server configServer `toml:"server"`
}

type configJava struct {
	// Path is the Path to the Java binary.
	Path string `toml:"path"`
	// Memory is the amount of Memory to launch the JVM with.
	Memory string `toml:"memory"`
	// ExtraArgs holds additional command line arguments for the JVM.
	ExtraArgs []string `toml:"extraArgs"`
}

type configServer struct {
	// Dir is the directory of the server containing the world data and config files.
	Dir string `toml:"dir"`
	// Path is the Path to the server JAR file.
	Path string `toml:"path"`
	// ExtraArgs holds additional command line arguments for the server.
	ExtraArgs []string `toml:"extraArgs"`
}

const configFile = "/etc/minecraftd.toml"

var defaultConfig = config{
	Java: configJava{
		Path:   "/usr/bin/java",
		Memory: "1G",
	},
	Server: configServer{
		Dir:  "/var/lib/minecraft",
		Path: "server.jar",
	},
}

func loadConfig() (*config, error) {
	f, err := os.Open(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			return &defaultConfig, nil
		}
		return nil, err
	}
	defer f.Close()
	return parseConfig(bufio.NewReader(f))
}

func parseConfig(r io.Reader) (*config, error) {
	cfg := defaultConfig
	md, err := toml.NewDecoder(r).Decode(&cfg)
	if err != nil {
		return nil, err
	}
	if len(md.Undecoded()) > 0 {
		return nil, errors.New("config files contains unknown keys")
	}
	return &cfg, nil
}

func (c config) serverCommand() []string {
	cmd := make([]string, 0)

	cmd = append(cmd, c.Java.Path)
	cmd = append(cmd,
		fmt.Sprintf("-Xmx%s", c.Java.Memory),
		fmt.Sprintf("-Xms%s", c.Java.Memory),
	)
	cmd = append(cmd, c.Java.ExtraArgs...)

	cmd = append(cmd, "-jar", c.Server.Path)
	cmd = append(cmd, "--nogui")
	cmd = append(cmd, c.Server.ExtraArgs...)

	return cmd
}
