package goConfig

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type DataType int

const (
	String DataType = iota
	Int
	Date
	Bool
	Float
	Uint
)

func (d DataType) String() string {
	return [...]string{
		"String",
		"Int",
		"Date",
		"Bool",
		"Float",
		"Unsigned Int",
	}[d]
}

func init() {
	viper.SetDefault("config.filename", "config")
	viper.SetDefault("config.format", "yaml")
}

func Init(appName, version, envPrefix string) {
	return
}

func SetAppVersion(version string) {
	var appVersion *semver.Version
	var err error

	if appVersion, err = semver.NewVersion(version); err != nil {
		logrus.Fatal("Error setting App Version")
	}

	viper.SetDefault("app.version", appVersion)
}

func SetAppName(appName, appPrefix string) {
	viper.SetDefault("app.name", appName)
	viper.SetDefault("app.prefix", appPrefix)
}

func GetConfigs() {
	var helpRequested bool
	var versionRequested bool

	viper.SetConfigName(viper.GetString("config.filename"))
	viper.SetConfigType(viper.GetString("config.format"))

	viper.AddConfigPath(filepath.Join("/etc", viper.GetString("app.prefix")))
	viper.AddConfigPath(filepath.Join(os.Getenv("HOME"), "."+viper.GetString("app.prefix")))
	viper.AddConfigPath("./etc/")
	viper.WatchConfig()

	viper.SetEnvPrefix(viper.GetString("app.prefix"))
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	pflag.BoolVarP(&helpRequested, "help", "h", false, "This message")
	pflag.BoolVarP(&versionRequested, "version", "v", false, "Get Program Version")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		} else {
			logrus.Warn("Invalid Config File Format")
		}
	}

	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)

	if helpRequested {
		showHelp()
	}

	if versionRequested {
		showVersion()
	}
}

func showHelp() {
	fmt.Printf("Usage: %s:\n", os.Args[0])
	pflag.PrintDefaults()

	os.Exit(0)
}

func showVersion() {
	fmt.Printf("%s: %s\n", os.Args[0], viper.GetString("app.version"))

	os.Exit(0)
}
