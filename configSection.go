package ecosystem

var _ ConfigSection = (*MultiConfigSection)(nil)

type MultiConfigSection struct {
	sections []ConfigSection
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
