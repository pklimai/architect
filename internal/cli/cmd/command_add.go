package cmd

import (
	"path/filepath"

	"github.com/spf13/cobra"
	"gitlab.com/zigal0/architect/internal/cli/logger"
	"gitlab.com/zigal0/architect/internal/cli/project"
	"gitlab.com/zigal0/architect/internal/cli/templates"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Base for other add sub commands.",
	Long:  "Command is a root for various add sub commands.",
}

var managerCmd = &cobra.Command{
	Use:   "manager",
	Short: "Generate new manager, top logic entity, with given pkg name.",
	Long: `Create new manager, top logic entity in the specified path internal/business/manager/manager_name/manager.go.
Also adds file interfaces.go  with commands for minimock and testing_test.go in the same place if it do not exist.
Name should satisfy snake_case.
`,
	Run: func(_ *cobra.Command, args []string) {
		logger.Info("Start manager creation.")

		if len(args) == 0 {
			logger.Fatal("No manager pkg name was provided.")
		}
		logger.FatalIfErr(validateEntityPkgName(args[0]))

		moduleName, err := moduleFromGoMod()
		logger.FatalIfErr(err)

		curProject := project.New(moduleName)

		// create architect of application
		for _, info := range projectPartInfosForManagerAdd(curProject, args[0]) {
			createProjectPart(info)
		}

		logger.Info("Finish manager creation.")
	},
}

// --------------------//
// Creation
// --------------------//

func projectPartInfosForManagerAdd(curProject *project.Project, managerPkgName string) []projectPartInfo {
	baseParths := []string{"internal", "business", "manager", managerPkgName}

	return []projectPartInfo{
		{
			curProject: curProject,
			pathParts:  append(baseParths, "manager.go"),
			tmplt:      templates.TemplateManager,
			tmpltData: templates.EntityData{
				EntityPkgName: managerPkgName,
			},
			needToOverride: false,
		},
		{
			curProject: curProject,
			pathParts:  append(baseParths, "interface.go"),
			tmplt:      templates.TemplateInterface,
			tmpltData: templates.EntityData{
				EntityPkgName: managerPkgName,
			},
			needToOverride: false,
		},
		{
			curProject: curProject,
			pathParts:  append(baseParths, "testing_test.go"),
			tmplt:      templates.TemplateTestingTest,
			tmpltData: templates.TestingTestData{
				EntityPkgName: managerPkgName,
				FileDirPath: filepath.Join(
					append([]string{curProject.Module()}, baseParths...)...,
				),
				EntityTypeName:               "manager",
				EntityTypeNameWithUpperFirst: "Manager",
			},
			needToOverride: false,
		},
	}
}
