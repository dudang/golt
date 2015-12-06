package runner

import (
	"github.com/dudang/golt/parser"
)

func generateGoltMap(goltTest parser.Golt) map[int][]parser.GoltJson {
	m := make(map[int][]parser.GoltJson)
	for _, element := range goltTest.Golt {
		array := m[element.Stage]
		if len(array) == 0 {
			m[element.Stage] = []parser.GoltJson{element}
		} else {
			m[element.Stage] = append(array, element)
		}
	}
	return m


}
