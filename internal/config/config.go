// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package config provides configuration for the application.
//
// Configuration is read from .env files.
package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Config is the configuration for the application.
type Config struct {
	// env and path can only be set by from the command line or environment variables
	env  Environment // the environment for the application
	path string      // path is the path to the configuration files.

	// Meta is the metadata for the configuration.
	Meta struct {
		ShowEnv      bool `json:"show-env,omitempty"`       // if set, dumps the environment to stdout
		ShowEnvFiles bool `json:"show-env-files,omitempty"` // if set, show the environment files that we are loading
		Verbose      bool `json:"verbose,omitempty"`        // verbose, if set, enables verbose logging.
	} `json:"meta,omitempty"`

	WorkingDir string `json:"working-directory,omitempty"`

	// Server configuration
	Server struct {
		Host string `json:"host,omitempty"`
		Port string `json:"port,omitempty"`
	} `json:"server,omitempty"`
}

// Environment is the environment the application is running in.
// It supports the following environment values:
//
// 1. development
// 2. test
// 3. production
//
// Each environment should have its own set of environment files.
type Environment int

const (
	Development Environment = iota
	Test
	Production
)

func (e Environment) String() string {
	switch e {
	case Development:
		return "development"
	case Test:
		return "test"
	case Production:
		return "production"
	default:
		panic(fmt.Errorf("%d: invalid environment", e))
	}
}

var (
	// environments must be sorted by ascending priority.
	environments = []Environment{Development, Test, Production}
)

// Default returns the environment the application is running in.
// It looks for the value on the command line first and then in the
// environment variables.
func Default(args []string) (*Config, error) {
	// initialize a configuration with default values.
	cfg := &Config{path: "."}
	if cwd, err := os.Getwd(); err != nil {
		return nil, err
	} else if cfg.WorkingDir, err = filepath.Abs(cwd); err != nil {
		return nil, err
	}
	cfg.Server.Host = "localhost"
	cfg.Server.Port = "8080"

	// check for values in the environment variables
	envSet := false
	if val, ok := os.LookupEnv("MOID_ENVIRONMENT"); ok {
		log.Printf("env: %-30s == %q\n", "MOID_ENVIRONMENT", val)
		switch strings.ToLower(val) {
		case "development":
			envSet, cfg.env = true, Development
		case "test":
			envSet, cfg.env = true, Test
		case "production":
			envSet, cfg.env = true, Production
		default:
			return nil, fmt.Errorf("%q: invalid environment", val)
		}
	}
	if val, ok := os.LookupEnv("MOID_CONFIG_PATH"); ok {
		log.Printf("env: %-30s == %q\n", "MOID_CONFIG_PATH", val)
		if val == "" {
			return nil, fmt.Errorf("%q: invalid path", val)
		}
		cfg.path = val
	}

	// parse the command line arguments
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--" {
			break
		} else if !strings.HasPrefix(arg, "--") {
			continue
		}
		opt, val, ok := strings.Cut(arg, "=")
		if opt == "--config-path" {
			log.Printf("env: %-30s == %q", "--config-path", val)
			if !ok || val == "" {
				return nil, fmt.Errorf("%q: invalid path", val)
			}
			cfg.path = val
		} else if opt == "--env" {
			log.Printf("env: %-30s == %q", "--env", val)
			switch strings.ToLower(val) {
			case "development":
				envSet, cfg.env = true, Development
			case "test":
				envSet, cfg.env = true, Test
			case "production":
				envSet, cfg.env = true, Production
			default:
				return nil, fmt.Errorf("%q: invalid environment", val)
			}
		} else if opt == "--help" {
			log.Fatalf("usage: moid [--env=development|test|production] [--config-path=...] [--help] [--show-env] [--show-env-files] [--verbose]\n")
		} else if opt == "--show-env" {
			cfg.Meta.ShowEnv = val == "" || val == "true" || val == "yes"
		} else if opt == "--show-env-files" {
			cfg.Meta.ShowEnvFiles = val == "" || val == "true" || val == "yes"
		} else if opt == "--verbose" {
			cfg.Meta.Verbose = val == "" || val == "true" || val == "yes"
		} else {
			log.Fatalf("unknown option: %q\n", arg)
		}
	}

	if !envSet { // we could not find the environment, so we return an error
		return nil, fmt.Errorf("missing environment")
	}

	// if the path is not absolute, then we need to convert it to an absolute path
	// that is relative to the working directory.
	if !filepath.IsAbs(cfg.path) {
		cfg.path = filepath.Join(cfg.WorkingDir, cfg.path)
		log.Printf("env: %-30s == %q\n", "MOID_CONFIG_PATH", cfg.path)
	}

	// path must be a valid directory that we can convert to an absolute path.
	if sb, err := os.Stat(cfg.path); os.IsNotExist(err) {
		return nil, fmt.Errorf("%q: does not exist", cfg.path)
	} else if !sb.IsDir() {
		return nil, fmt.Errorf("%q: is not a directory", cfg.path)
	}

	return cfg, nil
}

// Load loads the configuration for the application.
//
// If path is empty, the default path is used.
// Environment File Loading Priority
// Load tries to emulate the priority list from the dotenv page at
// https://github.com/bkeepers/dotenv#what-other-env-files-can-i-use.
//
// This is essentially .env, .env.{environment}, .env.local, .env.local.{environment}.
// Local files take priority over environment files, which take priority over the
// global .env file.
//
// The environment files are loaded in the following order:
//
//	Filename________________  .gitignore?  safe to store secrets?
//	.env                      no           never use for secrets
//	.env.{environment}        no           never use for secrets
//	.env.local                yes          maybe
//	.env.local.{environment}  yes          maybe
//
// Take note of that `.gitignore` column.
// It's a reminder of which files are expected to contain
// If `.gitignore` is `no`, then the file should never contain sensitive
// information like credentials and tokens because it may be checked into Git.
// If the value is `yes`, then the file might contain that information, so it
// should never, ever be checked in to Git.
//
// It's important to note that `.env` is loaded in all environments.
// The `.env.local` file is loaded in all environments except for test.
func (cfg *Config) Load(args []string) error {
	log.Printf("env: %-30s == %q\n", "MOID_ENVIRONMENT", cfg.env.String())
	log.Printf("env: %-30s == %q\n", "MOID_CONFIG_PATH", cfg.path)

	// load the configuration from the environment files.
	// note that the order of the environment files is important because
	// each file will overwrite the values of the previous file.

	// .env is the lowest priority, so load it first.
	if err := cfg.load(filepath.Join(cfg.path, ".env.json")); err != nil {
		return err
	}

	// shared environment settings are next in priority.
	for _, shared := range environments {
		if shared == cfg.env {
			if err := cfg.load(filepath.Join(cfg.path, ".env."+shared.String()+".json")); err != nil {
				return err
			}
		}
	}

	// .env.local is loaded for all environments except test.
	if cfg.env != Test {
		if err := cfg.load(filepath.Join(cfg.path, ".env.local.json")); err != nil {
			return err
		}
	}

	// local environment files are the highest priority, so load them last.
	for _, local := range environments {
		if local == cfg.env {
			if err := cfg.load(filepath.Join(cfg.path, ".env.local."+local.String()+".json")); err != nil {
				return err
			}
		}
	}

	// finally, load the command line arguments.
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--" {
			break
		} else if !strings.HasPrefix(arg, "--") {
			log.Fatalf("unknown option: %q\n", arg)
		}
		opt, val, ok := strings.Cut(arg, "=")
		if opt == "--config-path" {
			// handled previously
		} else if opt == "--env" {
			// handled previously
		} else if opt == "--help" {
			// handled previously
		} else if opt == "--show-env" {
			// handled previously
		} else if opt == "--show-env-files" {
			// handled previously
		} else if opt == "--working-directory" && ok && val != "" {
			cfg.WorkingDir = val
		} else if opt == "--verbose" {
			// handled previously
		} else {
			log.Fatalf("unknown option: %q\n", arg)
		}
	}
	if cfg.Meta.ShowEnv {
		if data, err := json.MarshalIndent(cfg, "", "  "); err != nil {
			return err
		} else {
			log.Printf("env: configuration %s\n", string(data))
		}
	}

	// set def to the working directory now that we've loaded the configuration.
	if sb, err := os.Stat(cfg.WorkingDir); os.IsNotExist(err) {
		return fmt.Errorf("%q: does not exist", cfg.WorkingDir)
	} else if !sb.IsDir() {
		return fmt.Errorf("%q: is not a directory", cfg.WorkingDir)
	} else if err = os.Chdir(cfg.WorkingDir); err != nil {
		return err
	} else if wd, err := os.Getwd(); err != nil {
		return err
	} else {
		log.Printf("env: %-30s == %q\n", "working-directory", wd)
	}

	return nil
}

func (cfg *Config) load(path string) error {
	if cfg.Meta.ShowEnvFiles {
		log.Printf("env: trying:: %q\n", path)
	}
	if data, err := os.ReadFile(path); err == nil {
		if cfg.Meta.ShowEnvFiles {
			log.Printf("env: loading: %q\n", path)
		}
		if err = json.Unmarshal(data, &cfg); err != nil {
			return fmt.Errorf("%q: %v", path, err)
		}
		if cfg.Meta.ShowEnvFiles {
			log.Printf("env: loaded:: %q\n", path)
		}
	}
	return nil
}

type Option func(*Config) error

func ShowEnv() Option {
	return func(cfg *Config) error {
		cfg.Meta.ShowEnv = true
		return nil
	}
}

func Verbose() Option {
	return func(cfg *Config) error {
		cfg.Meta.Verbose = true
		return nil
	}
}

func WithWorkingDir(path string) Option {
	return func(cfg *Config) error {
		if err := os.Chdir(path); err != nil {
			return err
		} else if cwd, err := os.Getwd(); err != nil {
			return err
		} else if cfg.WorkingDir, err = filepath.Abs(cwd); err != nil {
			return err
		}
		log.Printf("env: %-30s == %q\n", "MOID_WORKING_DIR", cfg.WorkingDir)
		return nil
	}
}
