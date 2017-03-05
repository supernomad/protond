// Copyright (c) 2017 Christian Saide <Supernomad>
// Licensed under the MPL-2.0, for details see https://github.com/Supernomad/protond/blob/master/LICENSE

package common

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net"
	"os"
	"path"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/Supernomad/protond/version"
	"gopkg.in/yaml.v2"
)

const (
	envPrefix = "PROTOND_"
)

/*
Config struct that handles marshalling in user supplied configuration data from cli arguments, environment variables, and configuration file entries.

The user supplied configuration is processed via a structured hierarchy:
	- Cli arguments override both environment variables and configuration file entries.
	- Environment variables will override file entries but can be overridden by cli arguments.
	- Configuration file entries will be overridden by both environment variables and cli arguments.
	- Defaults are used in the case that the use does not define a configuration argument.

The only exceptions to the above are the two special cli argments '-h'|'--help' or '-v'|'--version' which will output usage information or version information respectively and then exit the application.
*/
type Config struct {
	ConfFile        string            `skip:"false"  type:"string"    short:"c"    long:"conf-file"         default:""                              description:"The configuration file to use to configure protond."`
	Backlog         int               `skip:"false"  type:"int"       short:"b"    long:"backlog"           default:"1024"                          description:"The number of in flight events allowed per worker."`
	NumWorkers      int               `skip:"false"  type:"int"       short:"w"    long:"workers"           default:"0"                             description:"The number of protond workers to use, set to 0 for a worker per available cpu core."`
	FilterTimeout   time.Duration     `skip:"false"  type:"duration"  short:"t"    long:"filter-timeout"    default:"10s"                           description:"The maximum amount of time any filter can run before timing out and failing."`
	InputDirectory  string            `skip:"false"  type:"string"    short:"i"    long:"input-directory"   default:"/etc/protond/inputs.d"         description:"The directory containing arbitrary input filters for protond to use for ingesting events."`
	OutputDirectory string            `skip:"false"  type:"string"    short:"o"    long:"output-directory"  default:"/etc/protond/outputs.d"        description:"The directory containing arbitrary input filters for protond to use for ingesting events."`
	FilterDirectory string            `skip:"false"  type:"string"    short:"f"    long:"filter-directory"  default:"/etc/protond/filters.d"        description:"The directory containing arbitrary javascript filters for protond to use for event filtering."`
	DataDir         string            `skip:"false"  type:"string"    short:"d"    long:"data-dir"          default:"/var/lib/protond"              description:"The directory to store local protond state to."`
	PidFile         string            `skip:"false"  type:"string"    short:"p"    long:"pid-file"          default:"/var/run/protond/protond.pid"  description:"The pid file to use for tracking rolling restarts."`
	Log             *Logger           `skip:"true"` // The internal logger to use
	Inputs          []*PluginConfig   `skip:"true"` // The raw input configurations to use for event ingestion
	Outputs         []*PluginConfig   `skip:"true"` // The raw input configurations to use for event propagation
	Filters         []*FilterConfig   `skip:"true"` // The raw javascript filters to use during event filtering
	fileData        map[string]string `skip:"true"` // An internal map of data representing a passed in configuration file
}

func (config *Config) cliArg(short, long string, isFlag bool) (string, bool) {
	for i, arg := range os.Args {
		if arg == "-"+short ||
			arg == "--"+long {
			if !isFlag {
				return os.Args[i+1], true
			}
			return "true", true
		}
	}
	return "", false
}

func (config *Config) envArg(long string) (string, bool) {
	env := envPrefix + strings.ToUpper(strings.Replace(long, "-", "_", 10))
	output := os.Getenv(env)
	if output == "" {
		return output, false
	}
	return output, true
}

func (config *Config) fileArg(long string) (string, bool) {
	if config.fileData == nil {
		return "", false
	}
	value, ok := config.fileData[long]
	return value, ok
}

func (config *Config) usage(exit bool) {
	config.Log.Plain.Println("Usage of protond:")
	st := reflect.TypeOf(*config)

	numFields := st.NumField()
	for i := 0; i < numFields; i++ {
		field := st.Field(i)
		skip, fieldType, short, long, def, description := config.parseField(field.Tag)
		if skip == "true" {
			continue
		}

		config.Log.Plain.Printf("\t-%s|--%s  (%s)\n", short, long, fieldType)
		config.Log.Plain.Printf("\t\t%s (default: '%s')\n", description, def)
	}

	if exit {
		os.Exit(1)
	}
}

func (config *Config) version(exit bool) {
	config.Log.Plain.Printf("protond: v%s\n", version.VERSION)

	if exit {
		os.Exit(0)
	}
}

func (config *Config) parseFile() error {
	if config.ConfFile != "" {
		if !PathExists(config.ConfFile) {
			return errors.New("the supplied configuration file does not exist")
		}

		buf, err := ioutil.ReadFile(config.ConfFile)
		if err != nil {
			return err
		}

		data := make(map[string]string)
		ext := path.Ext(config.ConfFile)
		switch {
		case ".json" == ext:
			err = json.Unmarshal(buf, &data)
		case ".yaml" == ext || ".yml" == ext:
			err = yaml.Unmarshal(buf, &data)
		default:
			return errors.New("the supplied configuration file is not in a supported format, protond only supports 'json', or 'yaml' configuration files")
		}

		if err != nil {
			return err
		}

		config.fileData = data
	}
	return nil
}

func (config *Config) parseField(tag reflect.StructTag) (skip, fieldType, short, long, def, description string) {
	skip = tag.Get("skip")
	fieldType = tag.Get("type")
	short = tag.Get("short")
	long = tag.Get("long")
	def = tag.Get("default")
	description = tag.Get("description")
	return
}

func (config *Config) parseSpecial(args []string, exit bool) {
	for _, arg := range args {
		switch {
		case arg == "-h" || arg == "--help":
			config.usage(exit)
		case arg == "-v" || arg == "--version":
			config.version(exit)
		}
	}
}

func (config *Config) parseArgs() error {
	st := reflect.TypeOf(*config)
	sv := reflect.ValueOf(config).Elem()

	numFields := st.NumField()
	for i := 0; i < numFields; i++ {
		field := st.Field(i)
		fieldValue := sv.Field(i)
		skip, fieldType, short, long, def, _ := config.parseField(field.Tag)

		if skip == "true" || !fieldValue.CanSet() {
			continue
		}

		var raw string
		if value, ok := config.cliArg(short, long, fieldType == "bool"); ok {
			raw = value
		} else if value, ok := config.envArg(long); ok {
			raw = value
		} else if value, ok := config.fileArg(long); ok {
			raw = value
		} else {
			raw = def
		}

		switch fieldType {
		case "int":
			i, err := strconv.Atoi(raw)
			if err != nil {
				return errors.New("error parsing value for '" + long + "' got, '" + raw + "', expected an 'int'")
			}
			fieldValue.Set(reflect.ValueOf(i))
		case "duration":
			dur, err := time.ParseDuration(raw)
			if err != nil {
				return errors.New("error parsing value for '" + long + "' got, '" + raw + "', expected a 'duration' for example: '10s' or '2d'")
			}
			fieldValue.Set(reflect.ValueOf(dur))
		case "ip":
			ip := net.ParseIP(raw)
			if ip == nil && raw != "" {
				return errors.New("error parsing value for '" + long + "' got, '" + raw + "', expected an 'ip' for example: '10.0.0.1' or 'fd42:dead:beef::1'")
			}
			fieldValue.Set(reflect.ValueOf(ip))
		case "bool":
			b, err := strconv.ParseBool(raw)
			if err != nil {
				return errors.New("error parsing value for '" + long + "' got, '" + raw + "', expected a 'bool'")
			}
			fieldValue.Set(reflect.ValueOf(b))
		case "list":
			list := strings.Split(raw, ",")
			fieldValue.Set(reflect.ValueOf(list))
		case "string":
			fieldValue.Set(reflect.ValueOf(raw))
		default:
			return errors.New("build error unknown configuration type")
		}

		if field.Name == "ConfFile" {
			config.parseFile()
		}
	}

	return nil
}

func (config *Config) computeArgs() error {
	if numCPU := runtime.NumCPU(); config.NumWorkers == 0 || config.NumWorkers > numCPU {
		config.NumWorkers = numCPU
	}

	os.MkdirAll(config.DataDir, os.ModeDir)
	os.MkdirAll(path.Dir(config.PidFile), os.ModeDir)

	pid := os.Getpid()

	err := ioutil.WriteFile(config.PidFile, []byte(strconv.Itoa(pid)), os.ModePerm)
	if err != nil {
		return err
	}

	if PathExists(config.FilterDirectory) {
		filterFiles, err := ioutil.ReadDir(config.FilterDirectory)
		if err != nil {
			return err
		}

		config.Filters = make([]*FilterConfig, 0)
		for i := 0; i < len(filterFiles); i++ {
			name := filterFiles[i].Name()
			ext := path.Ext(name)
			switch ext {
			case ".js":
				fileData, err := ioutil.ReadFile(path.Join(config.FilterDirectory, name))
				if err != nil {
					return err
				}

				filterCfg := &FilterConfig{
					Type: ext[1:],
					Name: name,
					Code: string(fileData),
				}
				config.Filters = append(config.Filters, filterCfg)
			default:
				config.Log.Warn.Printf("Filter file '%s' is not one of the compatible filter types: 'js'.", name)
			}
		}
	} else {
		config.Log.Warn.Println("The specified FilterDirectory path does not exist, using Noop filter.")
	}

	inputConfigs, err := ParsePluginConfigs(config.InputDirectory, config.Log)
	if err != nil {
		return err
	}
	config.Inputs = inputConfigs

	outputConfigs, err := ParsePluginConfigs(config.OutputDirectory, config.Log)
	if err != nil {
		return err
	}
	config.Outputs = outputConfigs

	return nil
}

// NewConfig creates a new Config struct based on user supplied input.
func NewConfig(log *Logger) (*Config, error) {
	config := &Config{
		Log: log,
	}

	// Handle the help and version commands if the exist
	config.parseSpecial(os.Args, true)

	// Handle parsing user supplied configuration data
	if err := config.parseArgs(); err != nil {
		return nil, err
	}

	// Compute internal configuration based on the user supplied configuration data
	if err := config.computeArgs(); err != nil {
		return nil, err
	}

	return config, nil
}
