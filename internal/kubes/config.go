package kubes

// Config is the Kubernetes config
type Config struct {
	APIVersion     string            `yaml:"apiVersion"`
	Clusters       []Cluster         `yaml:"clusters"`
	Contexts       []Context         `yaml:"contexts"`
	CurrentContext string            `yaml:"current-context"`
	Kind           string            `yaml:"kind"`
	Preferences    map[string]string `yaml:"preferences"`
	Users          []User            `yaml:"users"`
}

// Cluster is a kubernetes cluster config
type Cluster struct {
	Name    string `yaml:"name"`
	Cluster struct {
		Server string `yaml:"server"`
	} `yaml:"cluster"`
}

// Context is a kubernetes context config
type Context struct {
	Name    string `yaml:"name"`
	Context struct {
		Cluster   string `yaml:"cluster"`
		Namespace string `yaml:"namespace"`
		User      string `yaml:"user"`
	} `yaml:"context"`
}

// User is a kubernetes user config
type User struct {
	Name string `yaml:"name"`
	User struct {
		Token string `yaml:"token"`
	} `yaml:"user"`
}

// DefaultConfig is the default kubernetes configuration to be populated
var DefaultConfig = Config{
	APIVersion:     "v1",
	Clusters:       []Cluster{{Name: "cluster"}},
	Contexts:       []Context{{Name: "cluster"}},
	CurrentContext: "cluster",
	Kind:           "config",
	Preferences:    make(map[string]string),
	Users:          []User{{Name: "default"}},
}
