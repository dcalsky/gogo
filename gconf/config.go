package gconf

import (
	"context"
	"errors"
	"fmt"
	"github.com/dcalsky/gogo/logs"
	"gopkg.in/yaml.v3"
	"os"
	"path"
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
func overrideConfigFromEnv(confPath string, target any) error {
	typ := reflect.TypeOf(target)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil
	}
	value := reflect.Indirect(reflect.ValueOf(target))
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		kind := f.Type.Kind()
		v := value.Field(i)
		if kind == reflect.Struct {
			err := overrideConfigFromEnv(confPath, v.Addr().Interface())
			if err != nil {
				return err
			}
			continue
		} else if kind == reflect.Ptr {
			if v.IsNil() {
				if f.Type.Elem().Kind() == reflect.Struct {
					continue
				}
				v.Set(reflect.New(f.Type.Elem()))
			}
			if f.Type.Elem().Kind() == reflect.Struct {
				err := overrideConfigFromEnv(confPath, v.Interface())
				if err != nil {
					return err
				}
				continue
			}
			v = v.Elem()
		}
		if relativeFilePath, ok := f.Tag.Lookup("file"); ok {
			filePath := path.Join(confPath, relativeFilePath)
			if !fileExist(filePath) {
				return fmt.Errorf("%s file not exist", filePath)
			}
			b, err := os.ReadFile(filePath)
			if err != nil {
				return err
			}
			if v.CanSet() {
				v.SetString(string(b))
			}
		}
		if envName, ok := f.Tag.Lookup("env"); ok {
			if envValue, exist := os.LookupEnv(envName); exist {
				if v.CanSet() {
					v.SetString(envValue)
				}
			}
		}
	}
	return nil
}

func UnmarshalConfFromDir(cluster string, env string, confDirPath string, target any) error {
	if confDirPath == "" {
		return errors.New("require a config dir")
	}
	paths := []string{
		filepath.Join(confDirPath, "base.yaml"),
	}
	if cluster != "" {
		paths = append(paths, filepath.Join(confDirPath, fmt.Sprintf("%s.yaml", cluster)))
		if env != "" {
			paths = append(paths, filepath.Join(confDirPath, fmt.Sprintf("%s.%s.yaml", cluster, env)))
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

	return overrideConfigFromEnv(confDirPath, target)
}
