package lightinggo

import (
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var fileName string = "env.yaml"

// EventOnChangeFunc the function that runs each time env change occurs.
type EventOnChangeFunc func(e fsnotify.Event)

type configFile struct {
	path    string
	watcher bool
	eventFn EventOnChangeFunc
}

// ConfigOption configures how we set up config file.
type ConfigOption func(e *configFile)

// WithConfigFile specifies the config file.
func WithConfigFile(filename string) ConfigOption {
	return func(e *configFile) {
		if v := strings.TrimSpace(filename); len(v) != 0 {
			e.path = filepath.Clean(v)
		}
	}
}

// WithConfigWatcher watching and re-reading config file.
func WithConfigWatcher(fn EventOnChangeFunc) ConfigOption {
	return func(e *configFile) {
		e.watcher = true
		e.eventFn = fn
	}
}

// LoadVariable will read your config file(s) and load them into ENV for this process.
// It will default to loading env.yaml in the current path if not specifies the filename.
func LoadVariable(options ...ConfigOption) error {
	env := &configFile{path: fileName}
	for _, f := range options {
		f(env)
	}

	filename, err := filepath.Abs(env.path)
	fileName = filename
	if err != nil {
		return err
	}

	statConfigFile(fileName)

	if env.watcher {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.Println("[lightinggo] env watcher panic", zap.Any("error", r), zap.String("config_file", filename), zap.ByteString("stack", debug.Stack()))
				}
			}()

			viper.OnConfigChange(func(e fsnotify.Event) {
				log.Println("Config file changed:", e.Name)
			})
			viper.WatchConfig()
		}()
	}

	return nil
}

func statConfigFile(fileName string) error {
	viper.AddConfigPath("./")
	viper.SetConfigFile(fileName)

	_, err := os.Stat(fileName)

	if err == nil {
		viper.ReadInConfig()
		return nil
	}

	viper.SetDefault("name", "dogger")
	viper.SetDefault("age", "18")
	viper.SetDefault("class", map[string]string{"class01": "01", "class02": "02"})
	viper.SetDefault("log", map[string]interface{}{"filename": "light.log", "stderr": true})

	return viper.WriteConfigAs(fileName)
}

func GetVariable(key string) interface{} {
	return viper.Get(key)
}

func SetVariable(key string, value interface{}) error {
	viper.Set(key, value)
	return viper.WriteConfigAs(fileName)
}
