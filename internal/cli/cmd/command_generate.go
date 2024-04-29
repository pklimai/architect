package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.com/zigal0/architect/internal/cli/logger"
	"gitlab.com/zigal0/architect/internal/cli/project"
	"gitlab.com/zigal0/architect/internal/cli/templates"
	"gitlab.com/zigal0/architect/internal/cli/tool"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Base for other generate sub commands.",
	Long:  "Command is a root for various generate sub commands.",
}

var servicesCmd = &cobra.Command{
	Use:   "services",
	Short: "Generate servises bases on given names.",
	Long: `Generate entitis for base application that responsible for connection between ptotoc generated code.
It generates code only for name that satisfies snake_case_name_service with name SnakeCaseNameService.
`,
	Run: func(_ *cobra.Command, args []string) {
		const servicesString = "services"
		logger.Infof(formatLogStartCreation, servicesString)

		if len(args) == 0 {
			logger.Info("No services names were provided.")
		}

		moduleName, err := moduleFromGoMod()
		logger.FatalIfErr(err)

		curProject := project.New(moduleName)

		for _, serviceName := range args {
			if !serviceNameRegExp.MatchString(serviceName) {
				logger.Info(fmt.Sprintf("skip '%s' name", serviceName))

				continue
			}

			createProjectPart(projectPartInfo{
				absPath: curProject.AbsPath(),
				pathParts: []string{
					layerNameInternal,
					layerNameAPI,
					serviceName + "_impl",
					"service.go",
				},
				tmplt: templates.TemplateService,
				tmpltData: templates.ServiceData{
					Module:                             curProject.Module(),
					ServiceName:                        serviceName,
					ServiceNameCamelCaseWithFirstUpper: tool.ToCamelCaseWithFirstUpper(serviceName),
				},
				needToOverride: false,
			})
		}

		logger.Infof(formatLogFinishCreation, servicesString)
	},
}
