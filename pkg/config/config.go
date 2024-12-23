package config

import (
	"bytes"
	"reflect"
	"runtime"
	"text/template"

	"github.com/go-viper/mapstructure/v2"
)

func TemplatingDecodeHook() mapstructure.DecodeHookFunc {
	return func(f reflect.Kind, t reflect.Kind, data interface{}) (interface{}, error) {
		if t != reflect.String {
			return data, nil
		}

		switch f {
		case reflect.String:
			return TemplateString(data.(string))
		}

		return data, nil
	}
}

func TemplateString(input string) (string, error) {
	templ, err := template.New("config").Funcs(FuncMap()).Parse(input)
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer

	if err := templ.Execute(&buffer, nil); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"arch": func() string { return runtime.GOARCH },
	}
}
