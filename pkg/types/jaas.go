package types

type JaaSServerAuth struct {
	Address  string `yaml:"address,omitempty"`
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
}

type JaaSConfig struct {
	Auths []JaaSServerAuth `yaml:"auths"`
}

type TaskCreateStatus struct {
	ID string
}

type Task struct {
	Name     string
	Replicas uint64
	Status   string
}
