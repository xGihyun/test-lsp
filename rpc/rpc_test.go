package rpc_test

import (
	"fmt"
	"test-lsp/rpc"
	"testing"
)

type EncodingTest struct {
	Method string `json:"method"`
}

func TestEncodeMessage(t *testing.T) {
	testContent := "{\"method\":\"hello\"}"
	expected := fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(testContent), testContent)
	actual := rpc.EncodeMessage(EncodingTest{Method: "hello"})

	if expected != actual {
		t.Fatalf("Expected: %s, Got: %s", expected, actual)
	}
}

func TestDecodeMessage(t *testing.T) {
	testContent := "{\"method\":\"hello\"}"
	message := fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(testContent), testContent)
	method, content, err := rpc.DecodeMessage([]byte(message))
  contentLength := len(content)

	if err != nil {
		t.Fatal(err)
	}

	if contentLength != len(content) {
		t.Fatalf("Expected: %d, Got: %d", len(content), contentLength)
	}

	if method != "hello" {
		t.Fatalf("Expected: %s, Got: %s", "hello", method)
	}
}
