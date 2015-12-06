package runner

import (
	"github.com/dudang/golt/parser"
)

func generateGoltMap(goltTest parser.Golts) map[int][]parser.GoltItem {
	m := make(map[int][]parser.GoltItem)
	for _, element := range goltTest.Golt {
		array := m[element.Stage]
		if len(array) == 0 {
			m[element.Stage] = []parser.GoltItem{element}
		} else {
			m[element.Stage] = append(array, element)
		}
	}
	return m
}
