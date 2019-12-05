package models

// API3.0 Response
type ResponseWrapper struct {
	Response ResponseCommon `json:"Response"`
}
type ResponseCommon struct {
	Error *ResponseError `json:"Error,omitempty"`
}
type ResponseError struct {
	Code string `json:"Code"`
}
