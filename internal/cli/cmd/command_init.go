package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gitlab.com/zigal0/architect/internal/cli/logger"
	"gitlab.com/zigal0/architect/internal/cli/project"
	"gitlab.com/zigal0/architect/internal/cli/templates"
	"golang.org/x/mod/module"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize architect application.",
	Long:  "Cteate new architect application with all necessary infrastructure and ready to start.",
	Run: func(_ *cobra.Command, args []string) {
		logger.Info("Start generating application with architect.")

		if len(args) == 0 {
			logger.Fatal("The 'module' is required argument.")
		}
		logger.FatalIfErr(validateModule(args[0]))

		wd, err := os.Getwd()
		logger.FatalIfErr(err)

		curProject := project.New(args[0], wd)
		logger.FatalIfErr(validateProjectName(curProject.Name()))

		if !checkFileExist(filepath.Join(curProject.AbdPath(), goModFileName)) {
			logger.FatalIfErr(execute("go", "mod", "init", curProject.Module()))
		}

		// create architect of application
		for _, info := range prepareProjectPartInfos(curProject) {
			createProjectPart(info)
		}

		logger.FatalIfErr(executeMake("generate", curProject.AbdPath()))

		logger.Info("Finish mudule initialization.")
	},
}

// --------------------//
// Creation
// --------------------//

func prepareProjectPartInfos(curProject *project.Project) []projectPartInfo {
	return []projectPartInfo{
		{
			curProject:     curProject,
			pathParts:      []string{".gitignore"},
			tmplt:          templates.GitIgnoreTemplate,
			tmpltData:      nil,
			needToOverride: false,
		},
		{
			curProject:     curProject,
			pathParts:      []string{".golangci.yaml"},
			tmplt:          templates.GolangCITemplate,
			tmpltData:      nil,
			needToOverride: false,
		},
		{
			curProject:     curProject,
			pathParts:      []string{".gitattributes"},
			tmplt:          templates.GitAttributesTemplate,
			tmpltData:      nil,
			needToOverride: false,
		},
		{
			curProject: curProject,
			pathParts:  []string{"Dockerfile"},
			tmplt:      templates.DockerfileTempalte,
			tmpltData: templates.CommonData{
				ProjectName: curProject.Name(),
			},
			needToOverride: true,
		},
		{
			curProject:     curProject,
			pathParts:      []string{".gitlab-ci.yml"},
			tmplt:          templates.GitlabCITemplate,
			tmpltData:      nil,
			needToOverride: true,
		},
		{
			curProject:     curProject,
			pathParts:      []string{"protodep.toml"},
			tmplt:          templates.ProtodepConfigTemplate,
			tmpltData:      nil,
			needToOverride: false,
		},
		{
			curProject: curProject,
			pathParts:  []string{"api", curProject.NameUnderscored() + "_service", "service.proto"},
			tmplt:      templates.ProtoServiceTemplate,
			tmpltData: templates.ProtoServiceData{
				Module:                             curProject.Module(),
				ModuleForProto:                     curProject.ModuleForProto(),
				ProjectNameUnderscored:             curProject.NameUnderscored(),
				ProjectNameCamelCaseWithFirstUpper: curProject.NameCamelCaseWithFirstUpper(),
				ProjectName:                        curProject.Name(),
			},
			needToOverride: false,
		},
		{
			curProject: curProject,
			pathParts:  []string{"architect.mk"},
			tmplt:      templates.ArchitectMKTemplate,
			tmpltData: templates.CommonData{
				ProjectName: curProject.Name(),
			},
			needToOverride: true,
		},
		{
			curProject:     curProject,
			pathParts:      []string{"Makefile"},
			tmplt:          templates.MakefileTemplate,
			tmpltData:      nil,
			needToOverride: false,
		},
		{
			curProject: curProject,
			pathParts:  []string{"cmd", curProject.Name(), "main.go"},
			tmplt:      templates.MainTemplate,
			tmpltData: templates.MainData{
				Module:                             curProject.Module(),
				ProjectNameSnakeCase:               curProject.NameSnakeCase(),
				ProjectNameCamelCaseWithFirstLower: curProject.NameCamelCaseWithFirstLower(),
			},
			needToOverride: false,
		},
		{
			curProject:     curProject,
			pathParts:      []string{"config", "config.go"},
			tmplt:          templates.ConfigTemplate,
			tmpltData:      nil,
			needToOverride: false,
		},
		{
			curProject:     curProject,
			pathParts:      []string{"config", "env_local_example.env"},
			tmplt:          templates.EnvLocalExampleTemplate,
			tmpltData:      nil,
			needToOverride: false,
		},
		// TODO: dirty hack for swagger, need to fix
		{
			curProject:     curProject,
			pathParts:      []string{"internal", "generated", "swagger", "embed.go"},
			tmplt:          templates.SwaggerHackTemplate,
			tmpltData:      nil,
			needToOverride: false,
		},
		{
			curProject:     curProject,
			pathParts:      []string{"script", "generate_swagger_ui.sh"},
			tmplt:          templates.GenerateSwaggerUITemplate,
			tmpltData:      nil,
			needToOverride: true,
		},
	}
}

// --------------------//
// Validation
// --------------------//

func validateModule(moduleName string) error {
	if err := module.CheckPath(moduleName); err != nil {
		return err
	}

	if !strings.HasPrefix(moduleName, "gitlab.com") {
		return fmt.Errorf("invalid module %s: must have prefix 'gitlab.com'", moduleName)
	}

	return nil
}

func validateProjectName(name string) error {
	if !projectNameRegExp.MatchString(name) {
		return fmt.Errorf("invalid application name '%s': starts with '-' symbol or contains forbidden symbols, only alphanumeric & '-' allowed", name) // nolint: lll
	}

	return nil
}
