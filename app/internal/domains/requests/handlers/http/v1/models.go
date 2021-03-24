package v1

type Response struct {
	Success         bool          `json:"success"`
	Message         string        `json:"message"`
	UrlRespPairList []UrlRespPair `json:"list,omitempty"`
}

type UrlRespPair struct {
	Url  string `json:"url"`
	Resp string `json:"response"`
}

type Request struct {
	URLs []string `json:"urls"`
}
