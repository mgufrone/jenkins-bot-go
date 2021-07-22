package handlers

type Step struct {
	Class              string      `json:"_class"`
	Links              Links       `json:"_links"`
	Actions            []Actions   `json:"actions"`
	DisplayDescription interface{} `json:"displayDescription"`
	DisplayName        string      `json:"displayName"`
	DurationInMillis   int         `json:"durationInMillis"`
	ID                 string      `json:"id"`
	Input              *Input      `json:"input"`
	Result             string      `json:"result"`
	StartTime          string      `json:"startTime"`
	State              string      `json:"state"`
	Type               string      `json:"type"`
}
type Self struct {
	Class string `json:"_class"`
	Href  string `json:"href"`
}
type Actions struct {
	Class string `json:"_class"`
	Href  string `json:"href"`
}
type Links struct {
	Self    Self    `json:"self"`
	Actions Actions `json:"actions"`
}
type Input struct {
	Class      string        `json:"_class"`
	Links      Links         `json:"_links"`
	ID         string        `json:"id"`
	Message    string        `json:"message"`
	Ok         string        `json:"ok"`
	Parameters []interface{} `json:"parameters"`
	Submitter  interface{}   `json:"submitter"`
}
