package cmd

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/mod/module"
)

// RegExp
var (
	serviceNameRegExp = regexp.MustCompile(`^([a-zA-Z]+_)+service$`)
	kebabCaseRegExp   = regexp.MustCompile(`^[a-zA-Z]+(-[a-zA-Z]+)*$`)
	snakeCaseRegExp   = regexp.MustCompile(`^[a-zA-Z]+(_[a-zA-Z]+)*$`)
)

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
	if !kebabCaseRegExp.MatchString(name) {
		return fmt.Errorf("invalid application name '%s': not kebab-case or contains forbidden symbols, only alphanumeric & '-' allowed", name) // nolint: lll
	}

	return nil
}

func validateEntityPkgName(name string) error {
	if !snakeCaseRegExp.MatchString(name) {
		return fmt.Errorf("invalid entity pkg name '%s': not snake-case or contains forbidden symbols, only alphanumeric & '-' allowed", name) // nolint: lll
	}

	return nil
}

func validateServiceName(name string) error {
	if !serviceNameRegExp.MatchString(name) {
		return fmt.Errorf("invalid service name '%s': not snake-case, no postfix _service or contains forbidden symbols, only alphanumeric & '_' allowed", name) // nolint: lll
	}

	return nil
}
