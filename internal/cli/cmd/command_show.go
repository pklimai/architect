package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Base for other show sub commands",
	Long:  "Command is a root for various show sub commands.",
}

var architectureCmd = &cobra.Command{
	Use:   "architecture",
	Short: "Show architecture of the architect based application",
	Long:  "Show components of architect based application with short info.",
	Run: func(_ *cobra.Command, _ []string) {
		// nolint: forbidigo
		fmt.Print(architectureCmdOutput)
	},
}

const architectureCmdOutput = `Scheme of architect-application:
├── api	# proto contracts
│   ├── some_name_service # for one gRPC service
│   │   └── service.proto
│   └── types # tpyes that used in some gRPC services
│       └── types.proto
├── bin # for binaries (in .gitignore)
│   └── ...
├── cmd # run entry points
│   ├── helper  # entry for developer needs
│   │   └── ...
│   └── project # entry for app
│       └── main.go
├── config # env settings for app run
│   ├── config.go
│   ├── local_example.env
│   ├── prod.env
│   └── stg.env
├── internal # internal code of app
│   ├── adapter # data sources
│   │   ├── client # other microservices
│   │   │   ├──some_app_name
│   │   │   │   ├── client.go
│   │   │   │   └── some_method.go
│   │   │   └── provider.go
│   │   └── repository # db
│   │       └── some
│   │           ├── model.go
│   │           ├── repository.go
│   │           ├── some_method.go
│   │           └── sql.go
│   ├── api # handlers
│   │   ├── some_name_service_impl
│   │   │   ├── interface.go
│   │   │   ├── mapper.go
│   │   │   ├── service.go
│   │   │   └── some_handler.go
│   │   └── mapper.go
│   ├── business # logic
│   │   ├── manager # top logic element
│   │   │   └── some
│   │   │       ├── error.go
│   │   │       ├── interface.go
│   │   │       ├── manager.go
│   │   │       ├── some_method.go
│   │   │       └── testing_test.go
│   │   ├── sub_manager # lower logic element (code used in some managers)
│   │   │   └── some # same files as manager
│   │   │       └── ...
│   │   └── tool # different tools used in many places (future libs)
│   │       └── ...
│   ├── domain # entities in app 
│   │   └── ...
│   ├── generated # auto generated code (DO NOT EDIT)
│   │   ├── api # api of app
│   │   │   └── ...
│   │   ├── client # api of other apps
│   │   │   └── ...
│   │   └── swagger # dirty hack (will removed)
│   │       └── ...
│   └── init # init functions for main, (e.g. db)
│       └── ...
├── local # for local development
│   ├── docker # components for run
│   │   └── ...
│   └── example # commands, queries for developers
│       └── ...
├── migration # db migrations
│   └── ...
├── tests # integration tests (run with running app)
│   └── ...
├── script # developer scripts
│   └── ...
├── vendor.protogen # vendor for proto (in .gitignore)
│   └── ...
├── .gitattributes # customize action with git (e.g. diff)
├── .gitignore # list of ignore dirs & files
├── .gitlab-ci.yml # pipeline for gitlab
├── .golangci.yaml # linter configuration
├── Makefile # app make targets
├── Dockerfile # dockerfile to deploy app
├── go.mod # dependencies of app
├── go.sum # checksums for go dependencies
├── protodep.lock # checksums for proto vendoring (in .gitignore)
├── protodep.toml # config for proto vendoring
├── README.md # description of app
└── architect.mk # common make targets (DO NOT EDIT)
`
