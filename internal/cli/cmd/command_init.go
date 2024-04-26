package cmd

import (
	"path/filepath"

	"github.com/spf13/cobra"
	"gitlab.com/zigal0/architect/internal/cli/logger"
	"gitlab.com/zigal0/architect/internal/cli/project"
	"gitlab.com/zigal0/architect/internal/cli/templates"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize architect application.",
	Long: `Cteate new architect application with all necessary infrastructure and ready to start.
Name of application (the last part of mudule) shuould be in kebab-case.
`,
	Run: func(_ *cobra.Command, args []string) {
		logger.Info("Start module initialization.")

		if len(args) == 0 {
			logger.Fatal("The 'module' is required argument.")
		}
		logger.FatalIfErr(validateModule(args[0]))

		curProject := project.New(args[0])
		logger.FatalIfErr(validateProjectName(curProject.Name()))

		if !checkFileExist(filepath.Join(curProject.AbdPath(), goModFileName)) {
			logger.FatalIfErr(execute("go", "mod", "init", curProject.Module()))
		}

		// create architect of application
		for _, info := range projectPartInfosForInit(curProject) {
			createProjectPart(info)
		}

		logger.FatalIfErr(executeMake("generate", curProject.AbdPath()))

		logger.Info("Finish mudule initialization.")
	},
}

// --------------------//
// Creation
// --------------------//

func projectPartInfosForInit(curProject *project.Project) []projectPartInfo {
	return []projectPartInfo{
		{
			curProject:     curProject,
			pathParts:      []string{".gitignore"},
			tmplt:          templates.TemplateGitIgnore,
			tmpltData:      nil,
			needToOverride: false,
		},
		{
			curProject:     curProject,
			pathParts:      []string{".gitattributes"},
			tmplt:          templates.TemplateGitAttributes,
			tmpltData:      nil,
			needToOverride: false,
		},
		{
			curProject:     curProject,
			pathParts:      []string{".gitlab-ci.yml"},
			tmplt:          templates.TemplateGitlabCI,
			tmpltData:      nil,
			needToOverride: true,
		},
		{
			curProject: curProject,
			pathParts:  []string{"Dockerfile"},
			tmplt:      templates.TemplateDockerfile,
			tmpltData: templates.CommonData{
				ProjectName: curProject.Name(),
			},
			needToOverride: true,
		},
		{
			curProject:     curProject,
			pathParts:      []string{".golangci.yaml"},
			tmplt:          templates.TemplateGolangCI,
			tmpltData:      nil,
			needToOverride: false,
		},

		{
			curProject:     curProject,
			pathParts:      []string{"protodep.toml"},
			tmplt:          templates.TemplateProtodepConfig,
			tmpltData:      nil,
			needToOverride: false,
		},
		{
			curProject: curProject,
			pathParts:  []string{"architect.mk"},
			tmplt:      templates.TemplateArchitectMK,
			tmpltData: templates.CommonData{
				ProjectName: curProject.Name(),
			},
			needToOverride: true,
		},
		{
			curProject:     curProject,
			pathParts:      []string{"config", "config.go"},
			tmplt:          templates.TemplateConfig,
			tmpltData:      nil,
			needToOverride: false,
		},
		{
			curProject:     curProject,
			pathParts:      []string{"config", "env_local_example.env"},
			tmplt:          templates.TemplateEnvLocalExample,
			tmpltData:      nil,
			needToOverride: false,
		},
		{
			curProject:     curProject,
			pathParts:      []string{"Makefile"},
			tmplt:          templates.TemplateMakefile,
			tmpltData:      nil,
			needToOverride: false,
		},
		{
			curProject: curProject,
			pathParts:  []string{"api", curProject.NameSnakeCase() + "_service", "service.proto"},
			tmplt:      templates.TemplateProtoService,
			tmpltData: templates.ProtoServiceData{
				Module:                             curProject.Module(),
				ModuleForProto:                     curProject.ModuleForProto(),
				ProjectNameSnakeCase:               curProject.NameSnakeCase(),
				ProjectNameCamelCaseWithFirstUpper: curProject.NameCamelCaseWithFirstUpper(),
				ProjectName:                        curProject.Name(),
			},
			needToOverride: false,
		},
		{
			curProject: curProject,
			pathParts:  []string{"cmd", curProject.Name(), "main.go"},
			tmplt:      templates.TemplateMain,
			tmpltData: templates.MainData{
				Module:                             curProject.Module(),
				ProjectNameSnakeCase:               curProject.NameSnakeCase(),
				ProjectNameCamelCaseWithFirstLower: curProject.NameCamelCaseWithFirstLower(),
			},
			needToOverride: false,
		},
		// TODO: dirty hack for swagger, need to fix
		{
			curProject:     curProject,
			pathParts:      []string{"internal", "generated", "swagger", "embed.go"},
			tmplt:          templates.TemplateSwaggerHack,
			tmpltData:      nil,
			needToOverride: false,
		},
		{
			curProject:     curProject,
			pathParts:      []string{"script", "generate_swagger_ui.sh"},
			tmplt:          templates.TemplateGenerateSwaggerUI,
			tmpltData:      nil,
			needToOverride: true,
		},
	}
}
