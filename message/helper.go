package message

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
)

func unmarshal(contentType string, rawMsg []byte, msg *Message) error {
	if contentType == "application/json" || contentType == "text/json" {
		return json.Unmarshal(rawMsg, msg)
	} else if contentType == "application/xml" || contentType == "text/xml" {
		return xml.Unmarshal(rawMsg, msg)
	}

	return fmt.Errorf("unsupport content type: %s", contentType)
}
