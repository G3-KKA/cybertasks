package config

import (
	"log"
	"time"
)

// Environment variables validated automatically.
var (
	environment = [...]string{
		// Every path in service works around WORKSPACE,
		// Removing this will break every env-based path in service.
		"WORKSPACE",
		"CONFIG_FILE",
		"PG_PASSWORD",
	}

	// Other viper options.
	elses = [...]elseSetter{}

	// Uses	viper.Set.
	override = [...]overrideContainer{}
)

// Configuration Constraints
//
// # Enivronment variables
//   - Must be defined, otherwise application shouldn't start
//   - Constant, shouldnt be overridden in runtime
//   - Should not have a default value
//
// # Configuration File
//   - Must exist, have same structure as config.Config, otherwise application shouldn't start
//   - May be overridden in runtime or exist in multiple variants across application parts
//   - Should not have a default value

// Signalizes that config field may contain env signature,
// and it must be replaced with value of the env.
type EnvString string

// Example:
// WORKSPACE = '~/user/goapp'
//
// ${WORKSPACE}/file/path   -->    ~/user/goapp/file/path.

// Represents expected contents of configuration file.
type Config struct {
	L  Logger      `mapstructure:"Logger"`
	C  Controller  `mapstructure:"Controller"`
	DB Database    `mapstructure:"Database"`
	TM TaskMetrics `mapstructure:"TaskMetrics"`
}
type (
	Logger struct {
		Writers []LogWriter `mapstructure:"Writers"`
	}
	LogWriter struct {
		Dst        EnvString `mapstructure:"Dst"`
		Type       string    `mapstructure:"Type"`
		Level      int8      `mapstructure:"Level"` // might be negative.
		MustCreate bool      `mapstructure:"MustCreate"`
	}
	Controller struct {
		//GRPCServer `mapstructure:"GRPCServer"`
		HTTPServer `mapstructure:"HTTPServer"`
	}
	/* 	GRPCServer struct {
		Address string `mapstructure:"Address"`
	} */
	HTTPServer struct {
		Address string `mapstructure:"Address"`
	}
	Database struct {
		User    string `mapstructure:"User"`
		Address string `mapstructure:"Address"`
		Port    string `mapstructure:"Port"`
		Dbname  string `mapstructure:"Dbname"`
	}
	TaskMetrics struct {
		ConnectionReties    uint          `mapstructure:"ConnectionReties"`
		RetryAfter          time.Duration `mapstructure:"RetryAfter"`
		PerServiceTableSize uint          `mapstructure:"PerServiceTableSize"`
		ExtServices         []ExtService  `mapstructure:"ExtServices"`
	}
	ExtService struct {
		Name       string        `mapstructure:"Name"`
		Autoupdate time.Duration `mapstructure:"Autoupdate"`
		Address    string        `mapstructure:"Address"`
	}
)

// Get reads from CONFIG_FILE.
// Return config or zero value config and error.
func Get() (Config, error) {
	err := initConfig()
	if err != nil {
		return Config{}, err
	}

	return c, nil
}

// MustGet reads from CONFIG_FILE.
// Return config or panics, if any error happened.
func MustGet() Config {
	cfg, err := Get()
	if err != nil {
		log.Fatal(err)
	}

	return cfg
}
