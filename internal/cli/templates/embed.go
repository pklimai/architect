package templates

import _ "embed"

//go:embed .gitignore
var GitIgnoreTemplate string

//go:embed .gitattributes
var GitAttributesTemplate string

//go:embed architect.mk
var ArchitectMKTemplate string

//go:embed Makefile
var MakefileTemplate string

//go:embed .golangci.yaml
var GolangCITemplate string

//go:embed Dockerfile
var DockerfileTempalte string

//go:embed .gitlab-ci.yml
var GitlabCITemplate string

//go:embed protodep.toml
var ProtodepConfigTemplate string

//go:embed service.proto
var ProtoServiceTemplate string

//go:embed generated_service.txt
var ServiceTemplate string

//go:embed swagger_hack.txt
var SwaggerHackTemplate string

//go:embed main.txt
var MainTemplate string

//go:embed config.txt
var ConfigTemplate string

//go:embed env_local_example.env
var EnvLocalExampleTemplate string

//go:embed generate_swagger_ui.sh
var GenerateSwaggerUITemplate string
