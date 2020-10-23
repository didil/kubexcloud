package responses

// ListApp response
type ListApp struct {
	Apps []ListAppEntry `json:"apps"`
}

type ListAppEntry struct {
	Name                string `json:"name"`
	AvailableReplicas   int32  `json:"availableReplicas"`
	UnavailableReplicas int32  `json:"unavailableReplicas"`
}
