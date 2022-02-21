package data

import (
	"fmt"
	"os"
	"strings"

	"github.com/zurrty/cranky/data/fs"
	"github.com/zurrty/cranky/util"
)

type ConfigFile struct {
	Path string
	keys []string
	vals []string
}

type KeyInvalidError struct {
	key string
}

func (e KeyInvalidError) Error() string {
	return "Key \"" + e.key + "\" not found."
}

func OpenConfig(cfgPath string) ConfigFile {

	file, err := os.Open(cfgPath)
	if err != nil {
		file, err = os.Create(cfgPath)
		util.Assert(err)
	}
	file.Close() // haha fuck u

	keys := []string{}
	vals := []string{}
	for i, line := range fs.ReadLines(cfgPath) {
		kvp := strings.Split(line, "=")
		if len(kvp) != 2 {
			println(fmt.Sprintf("line %d isnt a key-value pair", i))
			continue
		}
		keys = append(keys, kvp[0])
		vals = append(vals, kvp[1])
	}

	cfg := ConfigFile{
		Path: cfgPath,
		keys: keys,
		vals: vals,
	}
	return cfg
}

func (cf *ConfigFile) indexOf(key string) (int, error) {
	for i, f := range cf.keys {
		if key == f {
			return i, nil
		}
	}
	return -1, KeyInvalidError{key}
}

func (cf *ConfigFile) GetString(key string) (string, error) {
	i, err := cf.indexOf(key)
	if err != nil {
		return "", err
	}
	return cf.vals[i], nil
}

func (cf *ConfigFile) SetString(key string, val string) error {
	i, err := cf.indexOf(key)
	if err != nil {
		return err
	}
	cf.vals[i] = val
	return nil
}
