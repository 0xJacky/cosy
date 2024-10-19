package settings

import (
	"gopkg.in/ini.v1"
	"log"
	"reflect"
)

var (
	Conf     *ini.File
	ConfPath string
)

type section struct {
	Name string
	Ptr  any
}

var sections = []section{
	{
		Name: "app",
		Ptr:  AppSettings,
	},
	{
		Name: "server",
		Ptr:  ServerSettings,
	},
	{
		Name: "database",
		Ptr:  DataBaseSettings,
	},
	{
		Name: "redis",
		Ptr:  RedisSettings,
	},
	{
		Name: "sonyflake",
		Ptr:  SonyflakeSettings,
	},
}

// Register the setting, this should be called before Init
func Register(name string, ptr any) {
	sections = append(sections, section{name, ptr})
}

// Init the settings
func Init(confPath string) {
	ConfPath = confPath
	setup()
}

// Load the settings
func load() {
	var err error
	Conf, err = ini.LoadSources(ini.LoadOptions{
		Loose:        true,
		AllowShadows: true,
	}, ConfPath)

	if err != nil {
		log.Fatalf("setting.init, fail to parse 'app.ini': %v", err)
	}
}

// Reload the settings
func Reload() {
	load()
}

// Set up the settings
func setup() {
	load()

	for _, s := range sections {
		mapTo(s.Name, s.Ptr)
	}
}

// MapTo the settings
func mapTo(section string, v any) {
	err := Conf.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("setting.mapTo %s err: %v", section, err)
	}
}

// ReflectFrom the settings
func ReflectFrom(section string, v any) {
	err := Conf.Section(section).ReflectFrom(v)
	if err != nil {
		log.Fatalf("Cfg.ReflectFrom %s err: %v", section, err)
	}
}

// ProtectedFill fill the target settings with new settings
func ProtectedFill(targetSettings interface{}, newSettings interface{}) {
	s := reflect.TypeOf(targetSettings).Elem()
	vt := reflect.ValueOf(targetSettings).Elem()
	vn := reflect.ValueOf(newSettings).Elem()

	// copy the values from new to target settings if it is not protected
	for i := 0; i < s.NumField(); i++ {
		if s.Field(i).Tag.Get("protected") != "true" {
			vt.Field(i).Set(vn.Field(i))
		}
	}
}

// Save the settings
func Save() (err error) {
	for _, s := range sections {
		ReflectFrom(s.Name, s.Ptr)
	}
	err = Conf.SaveTo(ConfPath)
	if err != nil {
		return
	}
	setup()
	return
}
