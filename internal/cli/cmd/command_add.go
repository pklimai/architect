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

var (
	localPostgresFlag bool
)

func init() {
	addPotgresCmd.Flags().BoolVar(&localPostgresFlag, "local", false, "add postgres to docker-compose & make targets for local work") // nolint: lll
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Base for other add sub commands",
	Long:  "Command is a root for various add sub commands.",
}

var addManagerCmd = &cobra.Command{
	Use:   "manager",
	Short: "Add new manager, top logic entity, with given name",
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

var addSubManagerCmd = &cobra.Command{
	Use:   "sub-manager",
	Short: "Add new sub manager, lower logic entity, with given name",
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

var addRepositoryCmd = &cobra.Command{
	Use:   "repository",
	Short: "Add new repository with given name",
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
				dirNameInternal,
				dirNameAdapter,
				entityTypeNameRepository,
				args[0],
				entityTypeNameRepository + extensionGo,
			},
			tmplt: templates.TemplateRepository,
			tmpltData: templates.EntityData{
				PkgName: args[0] + "_" + entityTypeNameRepository,
			},
			needToOverride: false,
		})

		logger.FatalIfErr(executeGoModTidy())

		logger.Infof(formatLogFinishCreation, entityTypeNameRepository)
	},
}

var addProtoServiceCmd = &cobra.Command{
	Use:   "proto-service",
	Short: "Add proto contract for new service with given name",
	Long: `Create proto contract for new service in the specified path api/some_name_service/service.proto.
It contains example service that should be override. 
After changes you need to execute make target 'generate'.
`,
	Run: func(_ *cobra.Command, args []string) {
		logger.Infof(formatLogStartCreation, "proto-service")

		if len(args) == 0 {
			logger.Fatal(logNoEntityNameProvided)
		}

		serviceName := args[0]

		logger.FatalIfErr(validateServiceName(serviceName))

		moduleName, err := moduleFromGoMod()
		logger.FatalIfErr(err)

		curProject := project.New(moduleName)

		createProjectPart(projectPartInfo{
			absPath: curProject.AbsPath(),
			pathParts: []string{
				dirNameAPI,
				serviceName,
				entityTypeNameService + extensionProto,
			},
			tmplt: templates.TemplateProtoService,
			tmpltData: templates.ProtoServiceData{
				Module:                             curProject.Module(),
				ModuleForProto:                     curProject.ModuleForProto(),
				ServiceNameSnakeCase:               serviceName,
				ServiceNameCamelCaseWithFirstUpper: tool.ToCamelCaseWithFirstUpper(serviceName),
			},
			needToOverride: false,
		})

		logger.Infof(formatLogFinishCreation, "proto-service")
	},
}

var addClientCmd = &cobra.Command{
	Use:   "grpc-client",
	Short: "Add code to connect and interact with given client via gRPC",
	// nolint: lll
	Long: `Generate code for connect with given client via gRPC in the specified path internal/adapter/client/client_name/client.proto.
It adds contracnts to protodep.toml. Also it generates connection provider if it's not exist.
For execution requires 2 arguments: 
* client name in snake_case;
* path to proto contract with branch (e.g. gitalb.com/project/api/service/service.proto@main).
After changes you need to execute make target 'generate' and add code for gRPC connection in new client.go.
`,
	Run: func(_ *cobra.Command, args []string) {
		logger.Infof(formatLogStartCreation, entityTypeNameClient)

		if len(args) != 2 {
			logger.Fatal("Incorrect number of arguments, need client_name and path to proto contract.")
		}

		clientName := args[0]
		logger.FatalIfErr(validateEntityPkgName(clientName))

		fullPathToProto := args[1]
		logger.FatalIfErr(validateProtoContractsPath(fullPathToProto))

		moduleName, err := moduleFromGoMod()
		logger.FatalIfErr(err)

		curProject := project.New(moduleName)

		createProjectPart(projectPartInfo{
			absPath: curProject.AbsPath(),
			pathParts: []string{
				dirNameInternal,
				dirNameAdapter,
				entityTypeNameClient,
				fileNameProvider,
			},
			tmplt:          templates.TemplateProvider,
			tmpltData:      nil,
			needToOverride: false,
		})

		createProjectPart(projectPartInfo{
			absPath: curProject.AbsPath(),
			pathParts: []string{
				dirNameInternal,
				dirNameAdapter,
				entityTypeNameClient,
				clientName,
				entityTypeNameClient + extensionGo,
			},
			tmplt: templates.TemplateClient,
			tmpltData: templates.EntityData{
				PkgName: clientName + "_" + "client",
			},
			needToOverride: false,
		})

		appendToProjectPart(projectPartInfo{
			absPath:   curProject.AbsPath(),
			pathParts: []string{fileNameProtodep},
			tmplt:     templates.TemplateProtodepClient,
			tmpltData: constructProtodepClientData(fullPathToProto, clientName),
		})

		logger.Infof(formatLogFinishCreation, entityTypeNameClient)
	},
}

var addPotgresCmd = &cobra.Command{
	Use:   "postgres",
	Short: "Add code to connect to postgres and stuff for local work by flag",
	// nolint: lll
	Long: "Add code to connect to postgres. With local flag add tools for local work (docker, env, make-targets, migration mechanism)",
	Run: func(_ *cobra.Command, _ []string) {
		logger.Infof(formatLogStartCreation, "postgres")

		moduleName, err := moduleFromGoMod()
		logger.FatalIfErr(err)

		curProject := project.New(moduleName)

		createProjectPart(projectPartInfo{
			absPath: curProject.AbsPath(),
			pathParts: []string{
				"database",
				"postgres.go",
			},
			tmplt: templates.TemplatePostgresConnect,
			tmpltData: templates.CommonData{
				Module: curProject.Module(),
			},
		})

		if localPostgresFlag {
			// TODO: need to append, not create.
			createProjectPart(projectPartInfo{
				absPath: curProject.AbsPath(),
				pathParts: []string{
					"local",
					"docker",
					"docker-compose.yaml",
				},
				tmplt: templates.TemplatePostgresDocker,
				tmpltData: templates.CommonData{
					ProjectName: curProject.Name(),
				},
			})

			appendToProjectPart(projectPartInfo{
				absPath: curProject.AbsPath(),
				pathParts: []string{
					dirNameConfig,
					"env_local_example.env",
				},
				tmplt: templates.TemplatePostgresEnv,
				tmpltData: templates.CommonData{
					ProjectName: curProject.Name(),
				},
			})

			// TODO: Need to change place for append.
			appendToProjectPart(projectPartInfo{
				absPath: curProject.AbsPath(),
				pathParts: []string{
					"Makefile",
				},
				tmplt: templates.TemplatePostgresMakefile,
				tmpltData: templates.CommonData{
					ProjectName: curProject.Name(),
				},
			})
		}

		logger.Infof(formatLogFinishCreation, "postgres")
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
	baseParths := []string{dirNameInternal, dirNameBusiness, entityTypeName, entityName}

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

func constructProtodepClientData(
	fullPathToProto, clientName string,
) templates.ProtodepClientData {
	pathToProtoAndBranch := strings.Split(fullPathToProto, "@")
	if len(pathToProtoAndBranch) != 2 {
		logger.Fatalf("no branch name was porvided in '%s'", fullPathToProto)
	}

	pathToProto := pathToProtoAndBranch[0]

	branch := pathToProtoAndBranch[1]

	separatorIndex := strings.LastIndex(pathToProto, "/")

	if separatorIndex == -1 {
		logger.Fatalf("incorrect path to proto contracts in '%s'", fullPathToProto)
	}

	moduleWithPathToProtoDir := pathToProto[0:separatorIndex]
	protoFileName := pathToProto[separatorIndex+1:]

	return templates.ProtodepClientData{
		ModuleWithPathToProtoDir: moduleWithPathToProtoDir,
		ClientNameSnakeCase:      clientName,
		Branch:                   branch,
		PtotoFileName:            protoFileName,
	}
}
