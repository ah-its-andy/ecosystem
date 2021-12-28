package ecosystem

import (
	"errors"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v3"
)

var _ ConfigReader = (*YmlConfigReader)(nil)

type YmlConfigReader struct {
	fileDiscoveryPathList []string
}

func (reader *YmlConfigReader) GetSection(uri string) ConfigSection {
	results := &MultiConfigSection{
		sections: make([]ConfigSection, 0),
	}
	for _, path := range reader.fileDiscoveryPathList {
		fileName := filepath.Join(path, uri)
		section, err := reader.ReadFromFile(fileName)
		if err == nil {
			results.sections = append(results.sections, section)
		}
	}
	return results
}

func (reader *YmlConfigReader) ReadFromFile(filePath string) (ConfigSection, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Unwrap(err)
	}
	defer f.Close()
	decode := yaml.NewDecoder(f)
	section := &YmlConfigSection{}
	section.Values = make(map[string]string)
	section.FilePath = filePath
	err = decode.Decode(&section.Values)
	if err != nil {
		return nil, errors.Unwrap(err)
	}
	section.keys = make([]string, 0)
	for key := range section.Values {
		section.keys = append(section.keys, key)
	}
	return section, nil
}

type YmlConfigSection struct {
	FilePath string
	Values   map[string]string

	keys []string
}

func (section *YmlConfigSection) GetValue(key string) (string, bool) {
	ele, ok := section.Values[key]
	return ele, ok
}

func (section *YmlConfigSection) MustGetValue(key string) string {
	if val, ok := section.GetValue(key); ok {
		return val
	}
	return ""
}

func (section *YmlConfigSection) Keys() []string {
	keys := make([]string, len(section.keys))
	copy(keys, section.keys)
	return keys
}
