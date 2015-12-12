package runner

import (
	"github.com/dudang/golt/parser"
)

func generateGoltMap(goltTest parser.Golts) map[int][]parser.GoltThreadGroup {
	m := make(map[int][]parser.GoltThreadGroup)
	for _, element := range goltTest.Golt {
		array := m[element.Stage]
		if len(array) == 0 {
			m[element.Stage] = []parser.GoltThreadGroup{element}
		} else {
			m[element.Stage] = append(array, element)
		}
	}
	return m
}