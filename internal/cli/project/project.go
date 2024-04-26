package project

import (
	"os"
	"strings"

	"gitlab.com/zigal0/architect/internal/cli/logger"
	"gitlab.com/zigal0/architect/internal/cli/tool"
)

type Project struct {
	absPath string
	module  string
	name    string
}

func New(module string) *Project {
	parts := strings.Split(module, "/")
	name := parts[len(parts)-1]

	wd, err := os.Getwd()
	logger.FatalIfErr(err)

	return &Project{
		module:  module,
		absPath: wd,
		name:    name,
	}
}

func (p *Project) AbdPath() string {
	return p.absPath
}

func (p *Project) Name() string {
	return p.name
}

func (p *Project) NameCamelCaseWithFirstUpper() string {
	return tool.ToCamelCaseWithFirstUpper(p.name)
}

func (p *Project) NameCamelCaseWithFirstLower() string {
	return tool.ToCamelCaseWithFirstLower(p.name)
}

func (p *Project) NameSnakeCase() string {
	return strings.ReplaceAll(p.name, "-", "_")
}

func (p *Project) Module() string {
	return p.module
}

func (p *Project) ModuleForProto() string {
	moduleParts := strings.Split(p.module, "/")

	modulePartsWithoutName := moduleParts[:len(moduleParts)-1]

	moduleWithDots := strings.Join(modulePartsWithoutName, ".")

	return moduleWithDots + "." + p.NameSnakeCase()
}
