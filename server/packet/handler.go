package packet

import "github.com/rauwekost/astrio/server"

type (
	Handler struct {
		protocol   int
		pressQ     bool
		pressW     bool
		pressSpace bool
		connection *server.Connection
	}
)

//NewHandler returns a new packet-handler
func NewHandler(c *server.Connection) *Handler {
	return &Handler{
		protocol:   0,
		pressQ:     false,
		pressW:     false,
		pressSpace: false,
		connection: c,
	}
}

func (h *Handler) OnMessage(m []byte) {
	//ignore empty messages
	if len(m) == 0 {
		return
	}
}
