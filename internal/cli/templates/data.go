package templates

type CommonData struct {
	ProjectName string
}

type ProtoServiceData struct {
	Module                             string
	ModuleForProto                     string
	ProjectNameSnakeCase               string
	ProjectNameCamelCaseWithFirstUpper string
	ProjectName                        string
}

type ServiceData struct {
	Module                             string
	ServiceName                        string
	ServiceNameCamelCaseWithFirstUpper string
}

type MainData struct {
	ProjectNameSnakeCase               string
	Module                             string
	ProjectNameCamelCaseWithFirstLower string
}

type EntityData struct {
	PkgName string
}

type TestingTestData struct {
	PkgName                      string
	FileDirPath                  string
	EntityTypeNameWithUpperFirst string
}
