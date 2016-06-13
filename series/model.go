package series

type Series struct {
	UUID          string `json:"uuid"`
	PrefLabel     string `json:"prefLabel"`
	TmeIdentifier string `json:"tmeIdentifier,omitempty"`
	Type          string `json:"type,omitempty"`
}
