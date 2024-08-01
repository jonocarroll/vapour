package main

import (
	"github.com/devOpifex/vapour/lexer"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
)

type vapour struct {
	name     string
	version  string
	handler  *protocol.Handler
	server   *server.Server
	root     *string
	files    lexer.Files
	combined []byte
}
