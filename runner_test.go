package main
import (
	"testing"
	"time"
)

var testPlan = Golts{
	Golt: []GoltThreadGroup{
		GoltThreadGroup{Stage: 1, Timeout: 100},
		GoltThreadGroup{Stage: 3, Timeout: 300},
		GoltThreadGroup{Stage: 2, Timeout: 400},
		GoltThreadGroup{Stage: 1, Timeout: 200},
	},
}

func TestGenerateStageMap(t *testing.T) {
	m := generateStageMap(testPlan)
	if len(m[1]) != 2 || len(m[2]) != 1 || len(m[3]) != 1 {
		t.Error("The stage map was not generated properly")
	}
}

func TestGenerateHttpClient(t *testing.T) {
	for _, entry := range testPlan.Golt {
		client := generateHttpClient(entry)
		if client.Timeout != time.Duration(time.Millisecond * time.Duration(entry.Timeout)) {
			t.Error("The http client was not generated properly")
		}
	}
}
