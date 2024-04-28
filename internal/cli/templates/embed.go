package templates

import _ "embed"

//go:embed template_gitignore
var TemplateGitIgnore string

//go:embed template_gitattributes
var TemplateGitAttributes string

//go:embed template_architect.mk
var TemplateArchitectMK string

//go:embed template_makefile
var TemplateMakefile string

//go:embed template_golangci.yaml
var TemplateGolangCI string

//go:embed template_dockerfile
var TemplateDockerfile string

//go:embed template_gitlab-ci.yml
var TemplateGitlabCI string

//go:embed template_protodep.toml
var TemplateProtodepConfig string

//go:embed template_proto_app_service.proto
var TemplateProtoAppService string

//go:embed template_service.txt
var TemplateService string

//go:embed template_swagger_hack.txt
var TemplateSwaggerHack string

//go:embed template_main.txt
var TemplateMain string

//go:embed template_config.txt
var TemplateConfig string

//go:embed template_env_local_example.env
var TemplateEnvLocalExample string

//go:embed template_generate_swagger_ui.sh
var TemplateGenerateSwaggerUI string

//go:embed template_logic_entity.txt
var TemplateLogicEntity string

//go:embed template_interface.txt
var TemplateInterface string

//go:embed template_testing_test.txt
var TemplateTestingTest string

//go:embed template_repository.txt
var TemplateRepository string

//go:embed template_proto_service.proto
var TemplateProtoService string
