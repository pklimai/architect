package templates

import _ "embed"

//go:embed template_gitignore.txt
var TemplateGitIgnore string

//go:embed template_gitattributes.txt
var TemplateGitAttributes string

//go:embed template_architect.txt
var TemplateArchitectMK string

//go:embed template_makefile.txt
var TemplateMakefile string

//go:embed template_golangci.txt
var TemplateGolangCI string

//go:embed template_dockerfile.txt
var TemplateDockerfile string

//go:embed template_gitlab-ci.txt
var TemplateGitlabCI string

//go:embed template_protodep.txt
var TemplateProtodepConfig string

//go:embed template_proto_app_service.txt
var TemplateProtoAppService string

//go:embed template_service.txt
var TemplateService string

//go:embed template_swagger_hack.txt
var TemplateSwaggerHack string

//go:embed template_main.txt
var TemplateMain string

//go:embed template_config.txt
var TemplateConfig string

//go:embed template_env_local_example.txt
var TemplateEnvLocalExample string

//go:embed template_generate_swagger_ui.txt
var TemplateGenerateSwaggerUI string

//go:embed template_logic_entity.txt
var TemplateLogicEntity string

//go:embed template_interface.txt
var TemplateInterface string

//go:embed template_testing_test.txt
var TemplateTestingTest string

//go:embed template_repository.txt
var TemplateRepository string

//go:embed template_proto_service.txt
var TemplateProtoService string

//go:embed template_provider.txt
var TemplateProvider string

//go:embed template_client.txt
var TemplateClient string

//go:embed template_protodep_client.txt
var TemplateProtodepClient string

//go:embed template_postgres_connect.txt
var TemplatePostgresConnect string

//go:embed template_postgres_docker.txt
var TemplatePostgresDocker string

//go:embed template_postgres_env.txt
var TemplatePostgresEnv string

//go:embed template_postgres_makefile.txt
var TemplatePostgresMakefile string
