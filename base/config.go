package base

import (
	"context"
	"errors"
	"fmt"
	"github.com/dcalsky/gogo/common/logs"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"reflect"
)

func fileExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// overrideConfigFromEnv parses each fields type of `target`, if the field has `env` tag and is one of [string, int, float64, bool]:
// replace the value with the value from env
func overrideConfigFromEnv(target any) {
	typ := reflect.TypeOf(target)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return
	}
	value := reflect.Indirect(reflect.ValueOf(target))
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		kind := f.Type.Kind()
		v := value.Field(i)
		if kind == reflect.Struct {
			overrideConfigFromEnv(v.Addr().Interface())
			continue
		} else if kind == reflect.Ptr {
			if v.IsNil() {
				if f.Type.Elem().Kind() == reflect.Struct {
					continue
				}
				v.Set(reflect.New(f.Type.Elem()))
			}
			if f.Type.Elem().Kind() == reflect.Struct {
				overrideConfigFromEnv(v.Interface())
				continue
			}
			v = v.Elem()

		}
		if envName, ok := f.Tag.Lookup("env"); ok {
			if envValue, exist := os.LookupEnv(envName); exist {
				if v.CanSet() {
					v.SetString(envValue)
				}
			}
		}
	}
}

func UnmarshalConfFromDir(confDirPath string, target any) error {
	if confDirPath == "" {
		return errors.New("require a config dir")
	}
	paths := []string{
		filepath.Join(confDirPath, "base.yaml"),
	}
	cluster := GetCluster()
	envName := GetEnv()
	if cluster != "" {
		paths = append(paths, filepath.Join(confDirPath, fmt.Sprintf("%s.yaml", cluster)))
		if envName != "" {
			paths = append(paths, filepath.Join(confDirPath, fmt.Sprintf("%s.%s.yaml", cluster, envName)))
		}
	}
	ctx := context.Background()
	for _, path := range paths {
		if !fileExist(path) {
			logs.Infof(ctx, "[Conf] %s doesn't exist", path)
			continue
		}
		logs.Infof(ctx, "[Conf] use %s", path)
		fileBody, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if err = yaml.Unmarshal(fileBody, target); err != nil {
			return err
		}
	}

	overrideConfigFromEnv(target)
	return nil
}
