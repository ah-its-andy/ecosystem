package config

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v3"
)

var _ ConfigReader = (*ymlConfigReader)(nil)

type ymlConfigReader struct {
	fileDiscoveryPathList []string
}

func NewYmlConfigReader(disconveryPathList []string) ConfigReader {
	return &ymlConfigReader{
		fileDiscoveryPathList: disconveryPathList,
	}
}

func (reader *ymlConfigReader) GetSection(uri string) ConfigSection {
	subSections := make([]ConfigSection, 0)
	for _, path := range reader.fileDiscoveryPathList {
		fileName := filepath.Join(path, uri)
		section, err := reader.ReadFromFile(fileName)
		if err == nil {
			subSections = append(subSections, section)
		}
	}
	return NewMultiConfigSection(subSections)
}

func (reader *ymlConfigReader) ReadFromFile(filePath string) (ConfigSection, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Unwrap(err)
	}
	defer f.Close()

	buffer, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, errors.Unwrap(err)
	}

	values := make(map[string]string)
	err = yaml.Unmarshal(buffer, &values)
	if err != nil {
		return nil, errors.Unwrap(err)
	}

	return NewYmlConfigSection(filePath, buffer, values), nil
}

type ymlConfigSection struct {
	FilePath string
	Values   map[string]string

	keys []string
	raw  []byte
}

func NewYmlConfigSection(filePath string, raw []byte, values map[string]string) ConfigSection {
	section := &ymlConfigSection{}
	section.FilePath = filePath
	section.Values = values
	section.raw = raw
	section.keys = make([]string, 0)
	for key := range section.Values {
		section.keys = append(section.keys, key)
	}
	return section
}

func (section *ymlConfigSection) GetString(key string) (string, bool) {
	ele, ok := section.Values[key]
	return ele, ok
}

func (section *ymlConfigSection) MustGetString(key string) string {
	if val, ok := section.GetString(key); ok {
		return val
	}
	return ""
}

func (section *ymlConfigSection) Keys() []string {
	keys := make([]string, len(section.keys))
	copy(keys, section.keys)
	return keys
}

func (section *ymlConfigSection) Raw() []byte {
	return section.raw
}

func (section *ymlConfigSection) HasChildren() bool {
	return false
}

func (section *ymlConfigSection) Children() []ConfigSection {
	results := make([]ConfigSection, 0)
	return results
}

func (section *ymlConfigSection) SectionType() string {
	return SECTION_TYPE_LOCAL_YML
}
