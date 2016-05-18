package controllers

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/deepglint/glog"
	"github.com/deepglint/muses/eventserver/models"
)

type ApiController struct {
	name string
	Apis map[string]models.Api
}

func NewApiController(name string) *ApiController {
	output := new(ApiController)
	output.name = name

	output.Apis = make(map[string]models.Api)
	output.SetDocuments()
	return output
}

func (this *ApiController) AppendDocument(api *models.Api) {
	this.Apis[api.Name] = *api
}

func (this *ApiController) GetDocument(name string) models.Api {
	if api, ok := this.Apis[name]; ok {
		return api
	}
	return models.Api{}
}

func (this *ApiController) GetDocumentList() []string {
	var names []string
	for name, _ := range this.Apis {
		names = append(names, name)
	}
	return names
}

func (this *ApiController) GetAllDocument() []models.Api {
	var rets []models.Api
	for _, api := range this.Apis {
		rets = append(rets, api)
	}
	return rets
}

func (this *ApiController) Document(rw http.ResponseWriter, req *http.Request) {
	parametermap, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	args, ok := parametermap["name"]
	var ret interface{}
	// name, ok := params["name"]
	if !ok {
		// return document list
		ret = this.GetDocumentList()
	} else {
		name := args[0]
		// return document for name
		glog.Infoln(name)
		if name == "all" {
			ret = this.GetAllDocument()
		} else {
			ret = this.GetDocument(name)
		}
	}
	body, err := json.Marshal(ret)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.Write(body)
	rw.WriteHeader(http.StatusOK)
}
