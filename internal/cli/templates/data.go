package templates

type CommonData struct {
	ProjectName string
	Module      string
}

type ProtoAppServiceData struct {
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

type LogicEntityData struct {
	PkgName                               string
	EntityTypeNameCamelCaseWithFirstUpper string
}

type TestingTestData struct {
	PkgName                               string
	FileDirPath                           string
	EntityTypeNameCamelCaseWithFirstUpper string
}

type ProtoServiceData struct {
	Module                             string
	ModuleForProto                     string
	ServiceNameSnakeCase               string
	ServiceNameCamelCaseWithFirstUpper string
}

type ProtodepClientData struct {
	ModuleWithPathToProtoDir string
	Branch                   string
	ClientNameSnakeCase      string
	PtotoFileName            string
}
