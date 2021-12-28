package ecosystem

import "strings"

var _ ConfigSection = (*MultiConfigSection)(nil)

type MultiConfigSection struct {
	sections []ConfigSection
	keys     []string
}

func NewMultiConfigSection(s []ConfigSection) {
	r := &MultiConfigSection{
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
}

func indexOf(a []string, e string) int {
	for i, v := range a {
		if strings.Compare(v, e) == 0 {
			return i
		}
	}
	return -1
}

func (section *MultiConfigSection) Keys() []string {
	keys := make([]string, len(section.keys))
	copy(keys, section.keys)
	return keys
}

func (section *MultiConfigSection) GetValue(key string) (string, bool) {
	if len(section.sections) == 0 {
		return "", false
	}
	retVal := ""
	retOk := false
	for _, sub := range section.sections {
		if val, ok := sub.GetValue(key); ok {
			retVal = val
			retOk = true
			break
		}
	}
	return retVal, retOk
}

func (section *MultiConfigSection) MustGetValue(key string) string {
	if val, ok := section.GetValue(key); ok {
		return val
	}
	return ""
}

var EmptyConfigSection ConfigSection = &emptyConfigSection{}

type emptyConfigSection struct {
}

func (section *emptyConfigSection) GetValue(key string) (string, bool) {
	return "", false
}
func (section *emptyConfigSection) MustGetValue(key string) string {
	return ""
}

var emptyConfigSectionKeys = make([]string, 0)

func (section *emptyConfigSection) Keys() []string {
	return emptyConfigSectionKeys
}
