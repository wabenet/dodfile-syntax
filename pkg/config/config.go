package config

import (
	"bytes"
	"errors"
	"reflect"
	"regexp"
	"runtime"
	"text/template"

	"github.com/go-viper/mapstructure/v2"
)

var (
	ErrUnsupportedArch = errors.New("unsupported GOARCH value")
	ErrUnsupportedOS   = errors.New("unsupported GOOS value")
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
		"arch":     func() string { return runtime.GOARCH },
		"archLike": ArchLike,
		"os":       func() string { return runtime.GOOS },
		"osLike":   OsLike,
	}
}

func ArchLike(archs ...string) (string, error) {
	re, ok := MagicArchMap()[runtime.GOARCH]
	if !ok {
		return "", ErrUnsupportedArch
	}

	for _, arch := range archs {
		if re.MatchString(arch) {
			return arch, nil
		}
	}

	return "", ErrUnsupportedArch
}

func MagicArchMap() map[string]*regexp.Regexp {
	return map[string]*regexp.Regexp{
		"amd64":   regexp.MustCompile(`(?i)(x64|amd64|x86(-|_)?64)`),
		"386":     regexp.MustCompile(`(?i)(x32|amd32|x86(-|_)?32|i?386)`),
		"arm":     regexp.MustCompile(`(?i)(arm32|armv6|arm\b)`),
		"arm64":   regexp.MustCompile(`(?i)(arm64|armv8|aarch64|arm)`),
		"riscv64": regexp.MustCompile(`(?i)(riscv64)`),
	}
}

func OsLike(oss ...string) (string, error) {
	re, ok := MagicOSMap()[runtime.GOOS]
	if !ok {
		return "", ErrUnsupportedOS
	}

	for _, os := range oss {
		if re.MatchString(os) {
			return os, nil
		}
	}

	return "", ErrUnsupportedOS
}

func MagicOSMap() map[string]*regexp.Regexp {
	return map[string]*regexp.Regexp{
		"darwin":  regexp.MustCompile(`(?i)(darwin|mac.?(os)?|osx)`),
		"windows": regexp.MustCompile(`(?i)([^r]win|windows)`),
		"linux":   regexp.MustCompile(`(?i)(linux|ubuntu)`),
	}
}
