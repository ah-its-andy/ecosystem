package config

import (
	"reflect"
	"strings"
)

const SECTION_TYPE_LOCAL_YML = "local_yml"
const SECTION_TYPE_MULTI = "multi"
const SECTION_TYPE_EMPTY = "empty"
const SECTION_TYPE_STRONGTYPE = "strong_type"

type ConfigSection interface {
	GetString(key string) (string, bool)
	MustGetString(key string) string
	Keys() []string
	Raw() []byte
	HasChildren() bool
	Children() []ConfigSection
	SectionType() string
}

var _ ConfigSection = (*multiConfigSection)(nil)

type multiConfigSection struct {
	sections []ConfigSection
	keys     []string
}

func NewMultiConfigSection(s []ConfigSection) ConfigSection {
	r := &multiConfigSection{
		sections: make([]ConfigSection, len(s)),
		keys:     make([]string, 0),
	}
	for i, sub := range s {
		r.sections[i] = sub
		if i == 0 {
			r.keys = sub.Keys()
		} else {
			for _, key := range sub.Keys() {
				if indexOf(r.keys, key) == -1 {
					r.keys = append(r.keys, key)
				}
			}
		}
	}
	return r
}

func indexOf(a []string, e string) int {
	for i, v := range a {
		if strings.Compare(v, e) == 0 {
			return i
		}
	}
	return -1
}

func (section *multiConfigSection) Keys() []string {
	keys := make([]string, len(section.keys))
	copy(keys, section.keys)
	return keys
}

func (section *multiConfigSection) GetString(key string) (string, bool) {
	if len(section.sections) == 0 {
		return "", false
	}
	retVal := ""
	retOk := false
	for _, sub := range section.sections {
		if val, ok := sub.GetString(key); ok {
			retVal = val
			retOk = true
			break
		}
	}
	return retVal, retOk
}

func (section *multiConfigSection) MustGetString(key string) string {
	if val, ok := section.GetString(key); ok {
		return val
	}
	return ""
}

func (section *multiConfigSection) Raw() []byte {
	return make([]byte, 0)
}

func (section *multiConfigSection) HasChildren() bool {
	return len(section.sections) > 0
}

func (section *multiConfigSection) Children() []ConfigSection {
	results := make([]ConfigSection, len(section.sections))
	copy(results, section.sections)
	return results
}

func (section *multiConfigSection) SectionType() string {
	return SECTION_TYPE_MULTI
}

var EmptyConfigSection ConfigSection = &emptyConfigSection{}

type emptyConfigSection struct {
}

func (section *emptyConfigSection) GetString(key string) (string, bool) {
	return "", false
}
func (section *emptyConfigSection) MustGetString(key string) string {
	return ""
}

func (section *emptyConfigSection) Keys() []string {
	return make([]string, 0)
}

func (section *emptyConfigSection) Raw() []byte {
	return make([]byte, 0)
}

func (section *emptyConfigSection) HasChildren() bool {
	return false
}

func (section *emptyConfigSection) Children() []ConfigSection {
	results := make([]ConfigSection, 0)
	return results
}
func (section *emptyConfigSection) SectionType() string {
	return SECTION_TYPE_EMPTY
}

type StrongTypeConfigSection struct {
	reflect.Type
	instance interface{}
	raw      []byte
}

func (section *StrongTypeConfigSection) Raw() []byte {
	return section.raw
}

func (section *StrongTypeConfigSection) SectionType() string {
	return SECTION_TYPE_STRONGTYPE
}

func (section *StrongTypeConfigSection) GetInstance() interface{} {
	return section.instance
}
