package model

type Role struct {
	//id决定了Role的优先级，从0到3，优先级增高
	//当角色拥有多个优先级，采用最高优先级
	Id          int `yaml:"id"`
	*Resource   `yaml:"resource"`
	*Permission `yaml:"permission"`
}

type Resource struct {
	//资源的host
	Host string `json:"host" yaml:"host"`
	//资源的path
	Path string `json:"path" yaml:"path"`
	//资源的method
	Method string `json:"method" yaml:"method"`
}

type Permission struct {
	//允许访问资源的角色
	AuthorizedRoles []string `json:"authorized_roles" yaml:"authorized_roles"`
	//不允许访问资源的角色
	ForbiddenRoles []string `json:"forbidden_roles" yaml:"forbidden_roles"`
	//是否允许所有人访问，在Permission中优先级最高
	AllowAnyone bool `json:"allow_anyone" yaml:"allow_anyone"`
}
