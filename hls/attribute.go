package hls

import (
	"fmt"
	"net/url"
)

type Attribute struct {
	Name  string
	Value string
}

func (a Attribute) Format() string {
	if _, err := url.ParseRequestURI(a.Value); err == nil {
		return fmt.Sprintf("%s=\"%s\"", a.Name, a.Value)
	}

	return fmt.Sprintf("%s=%s", a.Name, a.Value)
}

func BuildAttributeTagLine(tag string, attributes []Attribute) string {
	if len(attributes) == 0 {
		return ""
	}

	attributeList := attributes[0].Format()
	for _, attribute := range attributes[1:] {
		attributeList = fmt.Sprintf("%s,%s", attributeList, attribute.Format())
	}

	return fmt.Sprintf("#%s:%s", tag, attributeList)
}
