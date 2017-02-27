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
	Filters         []*FilterConfig   `skip:"true"` // The raw javascript filters to use during event filtering
	fileData        map[string]string `skip:"true"` // An internal map of data representing a passed in configuration file
}

func (cfg *Config) cliArg(short, long string, isFlag bool) (string, bool) {
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

func (cfg *Config) envArg(long string) (string, bool) {
	env := envPrefix + strings.ToUpper(strings.Replace(long, "-", "_", 10))
	output := os.Getenv(env)
	if output == "" {
		return output, false
	}
	return output, true
}

func (cfg *Config) fileArg(long string) (string, bool) {
	if cfg.fileData == nil {
		return "", false
	}
	value, ok := cfg.fileData[long]
	return value, ok
}

func (cfg *Config) usage(exit bool) {
	cfg.Log.Plain.Println("Usage of protond:")
	st := reflect.TypeOf(*cfg)

	numFields := st.NumField()
	for i := 0; i < numFields; i++ {
		field := st.Field(i)
		skip, fieldType, short, long, def, description := cfg.parseField(field.Tag)
		if skip == "true" {
			continue
		}

		cfg.Log.Plain.Printf("\t-%s|--%s  (%s)\n", short, long, fieldType)
		cfg.Log.Plain.Printf("\t\t%s (default: '%s')\n", description, def)
	}

	if exit {
		os.Exit(1)
	}
}

func (cfg *Config) version(exit bool) {
	cfg.Log.Plain.Printf("protond: v%s\n", version.VERSION)

	if exit {
		os.Exit(0)
	}
}

func (cfg *Config) parseFile() error {
	if cfg.ConfFile != "" {
		if !PathExists(cfg.ConfFile) {
			return errors.New("the supplied configuration file does not exist")
		}

		buf, err := ioutil.ReadFile(cfg.ConfFile)
		if err != nil {
			return err
		}

		data := make(map[string]string)
		ext := path.Ext(cfg.ConfFile)
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

		cfg.fileData = data
	}
	return nil
}

func (cfg *Config) parseField(tag reflect.StructTag) (skip, fieldType, short, long, def, description string) {
	skip = tag.Get("skip")
	fieldType = tag.Get("type")
	short = tag.Get("short")
	long = tag.Get("long")
	def = tag.Get("default")
	description = tag.Get("description")
	return
}

func (cfg *Config) parseSpecial(args []string, exit bool) {
	for _, arg := range args {
		switch {
		case arg == "-h" || arg == "--help":
			cfg.usage(exit)
		case arg == "-v" || arg == "--version":
			cfg.version(exit)
		}
	}
}

func (cfg *Config) parseArgs() error {
	st := reflect.TypeOf(*cfg)
	sv := reflect.ValueOf(cfg).Elem()

	numFields := st.NumField()
	for i := 0; i < numFields; i++ {
		field := st.Field(i)
		fieldValue := sv.Field(i)
		skip, fieldType, short, long, def, _ := cfg.parseField(field.Tag)

		if skip == "true" || !fieldValue.CanSet() {
			continue
		}

		var raw string
		if value, ok := cfg.cliArg(short, long, fieldType == "bool"); ok {
			raw = value
		} else if value, ok := cfg.envArg(long); ok {
			raw = value
		} else if value, ok := cfg.fileArg(long); ok {
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
			cfg.parseFile()
		}
	}

	return nil
}

func (cfg *Config) computeArgs() error {
	if numCPU := runtime.NumCPU(); cfg.NumWorkers == 0 || cfg.NumWorkers > numCPU {
		cfg.NumWorkers = numCPU
	}

	os.MkdirAll(cfg.DataDir, os.ModeDir)
	os.MkdirAll(path.Dir(cfg.PidFile), os.ModeDir)

	pid := os.Getpid()

	err := ioutil.WriteFile(cfg.PidFile, []byte(strconv.Itoa(pid)), os.ModePerm)
	if err != nil {
		return err
	}

	if !PathExists(cfg.FilterDirectory) {
		cfg.Log.Warn.Println("The specified FilterDirectory path does not exist, using Noop filter.")
		return nil
	}

	files, err := ioutil.ReadDir(cfg.FilterDirectory)
	if err != nil {
		return err
	}

	cfg.Filters = make([]*FilterConfig, 0)
	for i := 0; i < len(files); i++ {
		name := files[i].Name()
		ext := path.Ext(name)
		switch ext {
		case ".js":
			fileData, err := ioutil.ReadFile(path.Join(cfg.FilterDirectory, name))
			if err != nil {
				return err
			}

			filterCfg := &FilterConfig{
				Type: ext[1:],
				Name: name,
				Code: string(fileData),
			}
			cfg.Filters = append(cfg.Filters, filterCfg)
		default:
			cfg.Log.Warn.Printf("Filter file '%s' is not one of the compatible filter types: 'js'.", name)
		}
	}

	return nil
}

// NewConfig creates a new Config struct based on user supplied input.
func NewConfig(log *Logger) (*Config, error) {
	cfg := &Config{
		Log: log,
	}

	// Handle the help and version commands if the exist
	cfg.parseSpecial(os.Args, true)

	// Handle parsing user supplied configuration data
	if err := cfg.parseArgs(); err != nil {
		return nil, err
	}

	// Compute internal configuration based on the user supplied configuration data
	if err := cfg.computeArgs(); err != nil {
		return nil, err
	}

	return cfg, nil
}
