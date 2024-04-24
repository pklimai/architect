package architect

import "embed"

type AppSettings struct {
	LogLevel    string
	Host        string
	PortHTTP    uint
	PortSwagger uint
	PortGRPC    uint
	SwaggerFS   embed.FS
}
