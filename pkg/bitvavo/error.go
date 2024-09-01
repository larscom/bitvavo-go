package bitvavo

import (
	"fmt"

	"github.com/larscom/bitvavo-go/v2/internal/util"
)

// ApiError Complete list of errorCodes: https://docs.bitvavo.com/#tag/Error-messages
type ApiError = WebSocketError

// WebSocketError Complete list of errorCodes: https://docs.bitvavo.com/#tag/Error-messages
type WebSocketError struct {
	Code    int    `json:"errorCode"`
	Message string `json:"error"`
	Action  string `json:"action"`
}

func (b *WebSocketError) Error() string {
	msg := fmt.Sprintf("code %d: %s", b.Code, b.Message)
	return fmt.Sprint(util.IfOrElse(len(b.Action) > 0, func() string { return fmt.Sprintf("%s action: %s", msg, b.Action) }, msg))
}
