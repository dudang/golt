package main

// A Golts contains all the GoltThreadGroup generated from a configuration file.
type Golts struct {
	Golt []GoltThreadGroup
}

// A GoltThreadGroup contains the configuration of a single thread generated
// from a configuration file.
type GoltThreadGroup struct {
	Threads     int
	Timeout     int
	Repetitions int
	Stage       int
	Requests    []GoltRequest
}

// A GoltRequest contains the configuration of a single HTTP request.
type GoltRequest struct {
	URL     string
	Method  string
	Payload string
	Headers map[string]*string
	Assert  GoltAssert
	// TODO: Have the possibility to extract multiple values
	Extract GoltExtractor
}

// A GoltAssert contains the configuration of the assertions to be made for a
// GoltRequest.
type GoltAssert struct {
	Status int
	Type   string
}

// A GoltExtractor contains the configuration to extract information of the
// response of a GoltRequest.
type GoltExtractor struct {
	Var   string
	Field string
	Regex string
	// TODO: Have the possibility to extract the value of a JSON field from the headers/body
}
