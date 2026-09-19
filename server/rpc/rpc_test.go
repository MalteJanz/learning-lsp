package rpc

import "testing"

type MyMsg struct {
	Hello bool
}

func TestEncode(t *testing.T) {
	expected := "Content-Length: 14\r\n\r\n{\"Hello\":true}"
	actual := EncodeMsg(MyMsg{Hello: true})

	if expected != actual {
		t.Fatalf("Expected: %s, Actual: %s", expected, actual)
	}
}

func TestDecode(t *testing.T) {
	incomingMsg := "Content-Length: 15\r\n\r\n{\"method\":\"hi\"}"
	method, content, err := DecodeMsg([]byte(incomingMsg))
	if err != nil {
		t.Fatal(err)
	}

	if len(content) != 15 {
		t.Fatalf("Expected: %d, Got: %d", 15, len(content))
	}

	if method != "hi" {
		t.Fatalf("Exptected: '%s', Got: '%s', with raw content %s", "hi", method, string(content))
	}
}
