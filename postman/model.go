package postman

type PostmanCollectionInfo struct {
	ID     string `json:"_postman_id"`
	Name   string `json:"name"`
	Schema string `json:"schema"`
}

type PostmanHeader struct {
	Key   string `json:"key"`
	Name  string `json:"name"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

type PostmanBody struct {
	Mode string `json:"mode"`
	Raw  string `json:"raw"`
}

type PostmanURL struct {
	Raw  string   `json:"raw"`
	Host []string `json:"host"`
	Path []string `json:"path"`
}

type PostmanRequest struct {
	Method  string          `json:"method"`
	Headers []PostmanHeader `json:"header"`
	Body    PostmanBody     `json:"body"`
	URL     PostmanURL      `json:"url"`
}

type PostmanResponse struct {
	Name            string          `json:"name"`
	OriginalRequest PostmanRequest  `json:"originalRequest"`
	Headers         []PostmanHeader `json:"header"`
	Status          string          `json:"status"`
	Code            int             `json:"code"`
	Body            string          `json:"body"`
}

type PostmanItem struct {
	Name     string            `json:"name"`
	Request  PostmanRequest    `json:"request"`
	Response []PostmanResponse `json:"response"`
}

type PostmanVariable struct {
	Id    string `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

type PostmanCollection struct {
	Info      PostmanCollectionInfo `json:"info"`
	Items     []PostmanItem         `json:"item"`
	Variables []PostmanVariable     `json:"variable"`
}
