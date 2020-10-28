package responses

// ListProject response
type ListProject struct {
	Projects []ListProjectEntry `json:"projects"`
}

type ListProjectEntry struct {
	Name string `json:"name"`
}

type Project struct {
	Name string `json:"name"`
}
