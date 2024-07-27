package message

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRawFromBytes(t *testing.T) {
	payload := `
{"payload":"eyJuYW1lIjoiZ2VuZXJhbCIsIm93bmVyIjoic3lzdGVtIiwiY3JlYXRlZEF0IjoiMjAyNC0wNy0yN1QxNzo1NzoyMi42NDY0MDUrMDI6MDAiLCJjdXJyZW50VXNlcnMiOjUsInRvdGFsTWVzc2FnZXMiOjQsIndlbGNvbWVNZXNzYWdlIjoiV2VsY29tZSB0byB0aGUgSnVuZ2xlIiwiUGFzc3dvcmQiOm51bGx9","user_id":"","type":"channel_join_response"}`

	for i := 0; i < 1000000; i++ {
		resp, err := RawFromBytes([]byte(payload))
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	}
}
