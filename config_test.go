package main

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

func Test_parseConfig(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    *config
		wantErr bool
	}{
		{
			name: "default config",
			args: args{r: bytes.NewBufferString("")},
			want: &config{
				Java:   configJava{Path: "/usr/bin/java", Memory: "1G"},
				Server: configServer{Dir: "/var/lib/minecraft", Path: "server.jar"},
			},
			wantErr: false,
		},
		{
			name: "set 2G memory",
			args: args{r: bytes.NewBufferString("java.memory = '2G'")},
			want: &config{
				Java:   configJava{Path: "/usr/bin/java", Memory: "2G"},
				Server: configServer{Dir: "/var/lib/minecraft", Path: "server.jar"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseConfig(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_config_serverCommand(t *testing.T) {
	tests := []struct {
		name string
		c    config
		want []string
	}{
		{
			name: "default config",
			c:    defaultConfig,
			want: []string{"/usr/bin/java", "-Xmx1G", "-Xms1G", "-jar", "server.jar", "--nogui"},
		},
		{
			name: "extra java args",
			c: config{
				Java:   configJava{Path: "/usr/bin/java", Memory: "1G", ExtraArgs: []string{"--javaArg1", "--javaArg2"}},
				Server: configServer{Dir: "/var/lib/minecraft", Path: "server.jar"},
			},
			want: []string{"/usr/bin/java", "-Xmx1G", "-Xms1G", "--javaArg1", "--javaArg2", "-jar", "server.jar", "--nogui"},
		},
		{
			name: "extra server args",
			c: config{
				Java:   configJava{Path: "/usr/bin/java", Memory: "1G"},
				Server: configServer{Dir: "/var/lib/minecraft", Path: "server.jar", ExtraArgs: []string{"--serverArg1", "--serverArg2"}},
			},
			want: []string{"/usr/bin/java", "-Xmx1G", "-Xms1G", "-jar", "server.jar", "--nogui", "--serverArg1", "--serverArg2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.serverCommand(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("config.serverCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}
