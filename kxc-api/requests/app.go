package requests

// CreateApp request
type CreateApp struct {
	Name     string `json:"name"`
	Replicas int32  `json:"replicas"`

	Containers []Container `json:"containers"`
}

// Container object
type Container struct {
	Image   string   `json:"image"`
	Name    string   `json:"name"`
	Command []string `json:"command,omitempty"`
	Ports   []Port   `json:"ports,omitempty"`
}

// Port Object
type Port struct {
	Number   int32  `json:"number"`
	Protocol string `json:"protocol"`
}
