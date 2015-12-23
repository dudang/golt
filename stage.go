package main

func generateGoltMap(goltTest Golts) map[int][]GoltThreadGroup {
	m := make(map[int][]GoltThreadGroup)
	for _, element := range goltTest.Golt {
		array := m[element.Stage]
		if len(array) == 0 {
			m[element.Stage] = []GoltThreadGroup{element}
		} else {
			m[element.Stage] = append(array, element)

		}
	}
	return m
}