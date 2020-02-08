package tester

import "io"

// TestRequest contains expected request from the API client being tested
type TestRequest struct {
	Method  string
	Path    string
	Headers map[string]string
	Body    io.ReadSeeker
}

// TestResponse contains mock response from the API
type TestResponse struct {
	Code    int
	Status  string
	Headers map[string]string
	Body    io.ReadSeeker
}

// TestScenario contains expected request and mock response
type TestScenario struct {
	Request  TestRequest
	Response TestResponse
}

// TestCategory contains the scenarios for particular category
type TestCategory struct {
	Scenarios map[string]TestScenario
}
