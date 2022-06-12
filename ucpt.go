package uscp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/wexder/uscp/encoding"
	"github.com/wexder/uscp/merge"
	"gopkg.in/yaml.v3"
)

type config struct {
	content  []byte
	fileType string
}

type Uscp struct {
	FileName     string
	ConfigPaths  []string
	FileTypes    []string
	FileContents []config
	Decoders     map[string]encoding.Decoder
	AutoParseEnv bool
}

func New() Uscp {
	return Uscp{
		AutoParseEnv: true,
		Decoders: map[string]encoding.Decoder{
			"yaml": encoding.YamlDecoder{},
			"json": encoding.JsonDecoder{},
		},
	}
}

func (u *Uscp) SetConfigName(name string) {
	u.FileName = name
}

func (u *Uscp) SetConfigPaths(paths []string) {
	u.ConfigPaths = paths
}

func (u *Uscp) SetFileTypes(fileTypes []string) {
	u.FileTypes = fileTypes
}

func (u *Uscp) AddConfigPath(path string) {
	u.ConfigPaths = append(u.ConfigPaths, path)
}

func (u *Uscp) AddFileType(fileType string) {
	u.FileTypes = append(u.FileTypes, fileType)
}

func (u *Uscp) ReadInConfiguration() error {
	for _, path := range u.ConfigPaths {
		files, err := os.ReadDir(path)
		if err != nil {
			return err
		}
		for _, file := range files {
			if file.IsDir() {
				continue
			}

			name, ext := getNameAndExtenstion(file.Name())
			if !sliceContains(u.FileTypes, ext) {
				continue
			}
			if name != u.FileName {
				continue
			}

			fileContent, err := os.ReadFile(fmt.Sprintf("%s/%s", path, file.Name()))
			if err != nil {
				return err
			}
			u.FileContents = append(u.FileContents, config{
				content:  fileContent,
				fileType: ext,
			})
		}
	}

	return nil
}

func (u *Uscp) Unmarshal(out interface{}) error {
	configs := []map[string]any{}
	for _, file := range u.FileContents {
		rawConfig := map[string]any{}
		err := u.Decoders[file.fileType].Decode(bytes.NewBuffer(file.content), &rawConfig)
		if err != nil {
			return err
		}

		configs = append(configs, rawConfig)
	}

	result := map[string]any{}
	if len(configs) > 0 {
		result = configs[0]
		for _, conf := range configs {
			merge.Merge(result, conf)
		}
	}

	jsonBytes, err := yaml.Marshal(result)
	if err != nil {
		return err
	}

	fmt.Println(string(jsonBytes))
	err = yaml.Unmarshal(jsonBytes, out)
	if err != nil {
		return err
	}

	return bindEnvs(out, u.FileName)
}

func bindEnvs(i interface{}, path string) error {
	t := reflect.ValueOf(i)
	for t.Kind() == reflect.Ptr {
		t = reflect.Indirect(t)
	}

	return iterateStruct(t, []string{path})
}

var (
	unmarshalerType = reflect.TypeOf((*json.Unmarshaler)(nil)).Elem()
	durationType    = reflect.TypeOf((*time.Duration)(nil)).Elem()
)

func iterateStruct(t reflect.Value, path []string) error {
	for i := 0; i < t.NumField(); i++ {
		f := t.Type().Field(i)
		n := f.Name
		value, present := os.LookupEnv(strings.Join(append(path, n), "_"))
		switch f.Type.Kind() {
		case reflect.Struct:
			// We want to check if we can unmarshal the value directly with unmarshaller
			impl := reflect.PointerTo(f.Type).Implements(unmarshalerType)
			if impl {
				if present {
					unmarsh, ok := t.Field(i).Addr().Interface().(json.Unmarshaler)
					if ok {
						err := unmarsh.UnmarshalJSON([]byte(`"` + value + `"`))
						if err != nil {
							return err
						}
					}
				}
			} else {
				iterateStruct(t.Field(i), append(path, n))
			}
			break
		default:
			if present {
				if f.Type == durationType {
					dur, err := time.ParseDuration(value)
					if err != nil {
						return err
					}
					t.Field(i).Set(reflect.ValueOf(dur))
				} else {
					t.Field(i).Set(reflect.ValueOf(value))
				}
			}
			break
		}

		// check if the field has tag for binding envs
		envsTags := strings.Split(f.Tag.Get("uscp_env"), ",")
		if len(envsTags) > 0 {
			for _, env := range envsTags {
				value, present := os.LookupEnv(env)
				if present {
					// TODO refactor
					if f.Type == durationType {
						dur, err := time.ParseDuration(value)
						if err != nil {
							return err
						}
						t.Field(i).Set(reflect.ValueOf(dur))
					} else {
						t.Field(i).Set(reflect.ValueOf(value))
					}
				}
			}
		}

		// check if the field has tag and if it's non zero if required
		tags := strings.Split(f.Tag.Get("uscp"), ",")
		if sliceContains(tags, "required") {
			isZero := t.Field(i).IsZero()
			if isZero {
				return fmt.Errorf("Required value %s is nil or zero value", strings.Join(append(path, n), "."))
			}
		}
	}

	return nil
}

func sliceContains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func getNameAndExtenstion(filename string) (string, string) {
	split := strings.Split(filename, ".")
	return split[0], split[1]
}
