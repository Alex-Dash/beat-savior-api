module bsvapi/bsdweb

go 1.21.5

require (
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
)


require bsvapi/types v0.0.0

replace bsvapi/types => ../types

require bsvapi/bsddb v0.0.0

replace bsvapi/bsddb => ../bsddb