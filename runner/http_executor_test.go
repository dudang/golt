package runner

import (
	"testing"
	"fmt"
)

const testKey1, testValue1 = "TEST_VALUE_1", "test-string-1"
const testKey2, testValue2 = "TEST_VALUE_2", "test-string-2"

func TestReplaceRegex(t *testing.T) {
	testingMap := make(map[string]string)
	testingMap[testKey1] = testValue1
	testingMap[testKey2] = testValue2

	replaceString1 := fmt.Sprintf("This is a test to replace the regex $(%s)", testKey1)
	replaceRegex(r, &replaceString1, testingMap)
	expectedString1 := fmt.Sprintf("This is a test to replace the regex %s", testValue1)
	if replaceString1 != expectedString1 {
		t.Fail()
	}

	replaceString2 := fmt.Sprintf("This is a test $(%s) to replace the regex $(%s)", testKey2, testKey2)
	replaceRegex(r, &replaceString2, testingMap)
	expectedString2 := fmt.Sprintf("This is a test %s to replace the regex %s", testValue2, testValue2)
	if replaceString2 != expectedString2 {
		t.Fail()
	}
}
