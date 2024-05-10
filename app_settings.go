package architect

import "embed"

// TODO: remove AppSettings and make all settings via Options or Flags.
type AppSettings struct {
	LogLevel    string
	Host        string
	PortHTTP    uint
	PortSwagger uint
	PortGRPC    uint
	SwaggerFS   embed.FS
}
