package postman

type postmanCollectionInfo struct {
	ID     string `json:"_postman_id"`
	Name   string `json:"name"`
	Schema string `json:"schema"`
}

type postmanHeader struct {
	Key   string `json:"key"`
	Name  string `json:"name"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

type postmanBody struct {
	Mode string `json:"mode"`
	Raw  string `json:"raw"`
}

type postmanURL struct {
	Raw  string   `json:"raw"`
	Host []string `json:"host"`
	Path []string `json:"path"`
}

type postmanRequest struct {
	Method  string          `json:"method"`
	Headers []postmanHeader `json:"header"`
	Body    postmanBody     `json:"body"`
	URL     postmanURL      `json:"url"`
}

type postmanResponse struct {
	Name            string          `json:"name"`
	OriginalRequest postmanRequest  `json:"originalRequest"`
	Headers         []postmanHeader `json:"header"`
	Status          string          `json:"status"`
	Code            int             `json:"code"`
	Body            string          `json:"body"`
}

type postmanItem struct {
	Name     string            `json:"name"`
	Request  postmanRequest    `json:"request"`
	Response []postmanResponse `json:"response"`
}

type postmanVariable struct {
	ID    string `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

type postmanCollection struct {
	Info      postmanCollectionInfo `json:"info"`
	Items     []postmanItem         `json:"item"`
	Variables []postmanVariable     `json:"variable"`
}
