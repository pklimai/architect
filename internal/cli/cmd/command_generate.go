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
	"gitlab.com/zigal0/architect/internal/cli/tool"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Base for other generate commands.",
	Long:  "Command is a root for various generate sub-commands.",
}

var servicesCmd = &cobra.Command{
	Use:   "services",
	Short: "Generate servises bases on given names.",
	Long: `Generate entitis for base application that responsible for connection between ptotoc generated code.
It generates code only for name that satisfies the pattern 'some-parts-name-service' with name SomePartsNameService.
`,
	Run: func(_ *cobra.Command, args []string) {
		logger.Info("Start generate services for ptotoc generated code.")

		moduleName, err := moduleFromGoMod()
		logger.FatalIfErr(err)

		wd, err := os.Getwd()
		logger.FatalIfErr(err)

		curProject := project.New(moduleName, wd)

		for _, rawServiceName := range args {
			if !serviceNameRegExp.MatchString(rawServiceName) {
				logger.Info(fmt.Sprintf("skip '%s' name", rawServiceName))

				continue
			}

			createService(curProject, strings.TrimSpace(rawServiceName))
		}

		logger.Info("Finish services generation.")
	},
}

func createService(curProject *project.Project, serviceName string) {
	filePath := filepath.Join(
		curProject.AbdPath(),
		"internal", "api", serviceName+"_impl", "service.go",
	)
	if checkFileExist(filePath) {
		return
	}

	data := templates.ServiceData{
		Module:                             curProject.Module(),
		ServiceName:                        serviceName,
		ServiceNameCamelCaseWithFirstUpper: tool.ToCamelCaseWithFirstUpper(serviceName),
	}

	content, err := createContentFromTemplate(templates.ServiceTemplate, data)
	logger.FatalIfErr(err)

	err = writeStringToFile(filePath, content)

	logger.FatalIfErr(err)
}
