package elasticclient

type ElasticResponse[T any] struct {
	Hits struct {
		Total struct {
			Value int64 `json:"value"`
		} `json:"total"`
		Hits []struct {
			Source T `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}
type ElasticResultResponse struct {
	Deleted int64 `json:"deleted"`
	Updated int64 `json:"updated"`
}
