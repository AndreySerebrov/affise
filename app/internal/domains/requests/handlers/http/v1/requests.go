package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Handler struct {
	log *log.Logger
	ctx context.Context
	ex  Requester
}

func New(ctx context.Context, log *log.Logger, requester Requester) *Handler {
	return &Handler{
		log: log,
		ctx: ctx,
		ex:  requester,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		err := Response{Message: "only post requests are supported"}
		errorBytes, _ := json.Marshal(err)
		_, _ = w.Write(errorBytes)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		err := Response{Message: "can't read request body"}
		errorBytes, _ := json.Marshal(err)
		_, _ = w.Write(errorBytes)
		return
	}

	req := Request{}
	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Printf("Error parsing body: %s", string(body))
		err := Response{Message: "can't parse request body"}
		errorBytes, _ := json.Marshal(err)
		_, _ = w.Write(errorBytes)
		return
	}

	resp, err := h.ex.MakeRequest(r.Context(), req.URLs)
	if err != nil {
		log.Printf("Error during external request")
		err := Response{Message: fmt.Sprintf("Error during external request, %s", err.Error())}
		errorBytes, _ := json.Marshal(err)
		_, _ = w.Write(errorBytes)
		return
	}

	httpResp := Response{Success: true, Message: "OK", UrlRespPairList: make([]UrlRespPair, 0, len(resp))}

	for _, item := range resp {
		httpResp.UrlRespPairList = append(httpResp.UrlRespPairList, UrlRespPair{Url: item.Url, Resp: item.Response})
	}

	errorBytes, _ := json.Marshal(resp)
	_, _ = w.Write(errorBytes)
}
