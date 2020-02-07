package common

import "io"

type TestRequest struct {
	Method  string
	Path    string
	Headers map[string]string
	Body    io.ReadSeeker
}

type TestResponse struct {
	Code    int
	Status  string
	Headers map[string]string
	Body    io.ReadSeeker
}

type TestScenario struct {
	Request  TestRequest
	Response TestResponse
}

type TestCategory struct {
	Scenarios map[string]TestScenario
}
