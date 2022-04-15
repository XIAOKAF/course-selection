package model

// Rule 用于解析yaml配置文件
type Rule struct {
	Id              string `yaml:"id"`
	Host            string `yaml:"host"`
	Path            string `yaml:"path"`
	Method          string `yaml:"method"`
	AuthorizedRoles string `yaml:"authorized_roles"`
	ForbiddenRoles  string `yaml:"forbidden_roles"`
	AllowAnyone     bool   `yaml:"allow_anyone"`
}
