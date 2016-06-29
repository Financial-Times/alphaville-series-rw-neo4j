package alphavilleseries

// AlphavilleSeries json data
type AlphavilleSeries struct {
	UUID          string `json:"uuid"`
	Description   string `json:"description,omitempty"`
	PrefLabel     string `json:"prefLabel"`
	TmeIdentifier string `json:"tmeIdentifier,omitempty"`
	Type          string `json:"type,omitempty"`
}
