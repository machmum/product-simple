package cfg

import (
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

var Env = cfg{}

type cfg struct {
	Timeout           int    `mapstructure:"TIMEOUT"`
	Debug             bool   `mapstructure:"DEBUG"`
	ServerAddress     string `mapstructure:"SERVER_ADDRESS" default:"localhost"`
	BasicAuthUsername string `mapstructure:"BASIC_AUTH_USERNAME"`
	BasicAuthPassword string `mapstructure:"BASIC_AUTH_PASSWORD"`
	JwtPublicKey      string `mapstructure:"JWT_PUBLIC_KEY"`

	// database
	DBProvider        string `mapstructure:"DB_PROVIDER"`
	DBHost            string `mapstructure:"DB_HOST"`
	DBPort            int    `mapstructure:"DB_PORT"`
	DBUser            string `mapstructure:"DB_USER"`
	DBPass            string `mapstructure:"DB_PASS"`
	DBName            string `mapstructure:"DB_NAME"`
	DBConnMaxLifetime int    `mapstructure:"DB_MAX_LIFETIME" default:"60"`
	DBConnMaxOpenConn int    `mapstructure:"DB_MAX_OPEN_CONN" default:"5"`
	DBConnMaxIdleConn int    `mapstructure:"DB_MAX_IDLE_CONN" default:"2"`
}

// BindEnvs: bind envar to struct
// ref: https://github.com/spf13/viper/issues/188
func BindEnvs(field interface{}, parts ...string) {
	ifv := reflect.ValueOf(field)
	ift := reflect.TypeOf(field)
	for i := 0; i < ift.NumField(); i++ {
		v := ifv.Field(i)
		t := ift.Field(i)
		defTag, _ := t.Tag.Lookup("default")
		tv, ok := t.Tag.Lookup("mapstructure")
		if !ok {
			continue
		}
		switch v.Kind() {
		case reflect.Struct:
			BindEnvs(v.Interface(), append(parts, tv)...)
		default:
			_ = viper.BindEnv(strings.Join(append(parts, tv), "."))
			if defTag != "" {
				viper.SetDefault(tv, defTag)
			}
		}
	}
}

func LoadEnv() {
	_ = viper.ReadInConfig()
	BindEnvs(Env)
	_ = viper.Unmarshal(&Env)
}
