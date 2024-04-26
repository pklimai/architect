package cmd

import (
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gitlab.com/zigal0/architect/internal/cli/logger"
	"gitlab.com/zigal0/architect/internal/cli/project"
	"gitlab.com/zigal0/architect/internal/cli/templates"
	"gitlab.com/zigal0/architect/internal/cli/tool"
)

const (
	formatLogStartCreation  = "Start %s creation."
	formatLogFinishCreation = "Finish %s creation."

	logNoEntityNameProvided = "No entity name was provided."
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Base for other add sub commands.",
	Long:  "Command is a root for various add sub commands.",
}

var managerCmd = &cobra.Command{
	Use:   entityTypeNameManager,
	Short: "Generate new manager, top logic entity, with given name.",
	// nolint: lll
	Long: `Create new manager, top logic entity in the specified path internal/business/manager/manager_pkg_name/manager.go.
Also adds file interfaces.go  with commands for minimock and testing_test.go in the same place if it do not exist.
Name should satisfy snake_case.
`,
	Run: func(_ *cobra.Command, args []string) {
		logger.Infof(formatLogStartCreation, entityTypeNameManager)

		if len(args) == 0 {
			logger.Fatal(logNoEntityNameProvided)
		}
		logger.FatalIfErr(validateEntityPkgName(args[0]))

		moduleName, err := moduleFromGoMod()
		logger.FatalIfErr(err)

		curProject := project.New(moduleName)

		for _, info := range projectPartInfosToAddLogicEntity(
			curProject,
			entityTypeNameManager,
			args[0],
		) {
			createProjectPart(info)
		}

		logger.FatalIfErr(executeGoModTidy())

		logger.Infof(formatLogFinishCreation, entityTypeNameManager)
	},
}

var subManagerCmd = &cobra.Command{
	Use:   entityTypeNameSubManager,
	Short: "Generate new sub manager, bottom logic entity, with given name.",
	// nolint: lll
	Long: `Create new sub manager, bottom logic entity in the specified path internal/business/sub_manager/sub_manager_pkg_name/manager.go.
Also adds file interfaces.go  with commands for minimock and testing_test.go in the same place if it do not exist.
Name should satisfy snake_case.
`,
	Run: func(_ *cobra.Command, args []string) {
		logger.Infof(formatLogStartCreation, entityTypeNameSubManager)

		if len(args) == 0 {
			logger.Fatal(logNoEntityNameProvided)
		}
		logger.FatalIfErr(validateEntityPkgName(args[0]))

		moduleName, err := moduleFromGoMod()
		logger.FatalIfErr(err)

		curProject := project.New(moduleName)

		for _, info := range projectPartInfosToAddLogicEntity(
			curProject,
			entityTypeNameSubManager,
			args[0],
		) {
			createProjectPart(info)
		}

		logger.FatalIfErr(executeGoModTidy())

		logger.Infof(formatLogFinishCreation, entityTypeNameSubManager)
	},
}

var repositoryCmd = &cobra.Command{
	Use:   entityTypeNameRepository,
	Short: "Generate new repository with given name.",
	// nolint: lll
	Long: `Create new repositoy based on sqlx in the specified path internal/adapter/repository/repository_pkg_name/repository.go.
Also adds sql.go for quieries & model.go for data.
Name should satisfy snake_case.
`,
	Run: func(_ *cobra.Command, args []string) {
		logger.Infof(formatLogStartCreation, entityTypeNameRepository)

		if len(args) == 0 {
			logger.Fatal(logNoEntityNameProvided)
		}
		logger.FatalIfErr(validateEntityPkgName(args[0]))

		moduleName, err := moduleFromGoMod()
		logger.FatalIfErr(err)

		curProject := project.New(moduleName)

		createProjectPart(projectPartInfo{
			absPath: curProject.AbsPath(),
			pathParts: []string{
				layerNameiInternal,
				layerNameAdapter,
				entityTypeNameRepository,
				args[0],
				entityTypeNameRepository + extensionGo,
			},
			tmplt: templates.TemplateRepository,
			tmpltData: templates.EntityData{
				PkgName: args[0],
			},
			needToOverride: false,
		})

		logger.FatalIfErr(executeGoModTidy())

		logger.Infof(formatLogFinishCreation, entityTypeNameRepository)
	},
}

// --------------------//
// Creation
// --------------------//

func projectPartInfosToAddLogicEntity(
	curProject *project.Project,
	entityTypeName string,
	entityName string,
) []projectPartInfo {
	baseParths := []string{layerNameiInternal, layerNameBusiness, entityTypeName, entityName}

	pkgName := strings.Join([]string{entityName, entityTypeName}, "_")

	entityTypeNameCamelCase := tool.ToCamelCaseWithFirstUpper(entityTypeName)

	return []projectPartInfo{
		{
			absPath:   curProject.AbsPath(),
			pathParts: append(baseParths, entityTypeName+extensionGo),
			tmplt:     templates.TemplateLogicEntity,
			tmpltData: templates.LogicEntityData{
				PkgName:                               pkgName,
				EntityTypeNameCamelCaseWithFirstUpper: entityTypeNameCamelCase,
			},
			needToOverride: false,
		},
		{
			absPath:   curProject.AbsPath(),
			pathParts: append(baseParths, fileNameInterface),
			tmplt:     templates.TemplateInterface,
			tmpltData: templates.EntityData{
				PkgName: pkgName,
			},
			needToOverride: false,
		},
		{
			absPath:   curProject.AbsPath(),
			pathParts: append(baseParths, fileNameTestingTest),
			tmplt:     templates.TemplateTestingTest,
			tmpltData: templates.TestingTestData{
				PkgName: pkgName,
				FileDirPath: filepath.Join(
					append([]string{curProject.Module()}, baseParths...)...,
				),
				EntityTypeNameCamelCaseWithFirstUpper: entityTypeNameCamelCase,
			},
			needToOverride: false,
		},
	}
}
