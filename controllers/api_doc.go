package controllers

import (
	"github.com/deepglint/muses/eventserver/models"
)

func (this *ApiController) SetDocuments() {
	//
	a1 := models.NewApi("getdocument", "doc", this.name, "Document", models.REQUEST_GET, false)
	a1.SetDescription(
		"get the API document by a parameter \"name\" or ignore the parameter to return a list of all api",
		"取得以\"name\"命名的API文档，如不指定参数，将获得API列表",
	)
	a1.AddParameter(
		"name", models.TYPE_STRING, "getdocument",
		"the name of the API",
		"API名称",
	)
	a1.SetRequest(models.BODY_TYPE_NONE)
	a1.SetResponse(models.BODY_TYPE_JSON)
	a1.AddResponseBodyParameter("Name", models.TYPE_STRING, "getdocument",
		"the name of the API",
		"API名称",
	)
	a1.AddResponseBodyParameter("Path", models.TYPE_STRING, "doc",
		"the url path of the API",
		"API的URL路径",
	)
	a1.AddResponseBodyParameter("RequestType", models.TYPE_STRING, "get",
		"the http request type of the API, get or post",
		"API的http请求类型，get或post",
	)
	a1.AddResponseBodyParameter("Parameters", models.TYPE_STRING, models.Parameter{},
		"",
		"参数列表，包括类型、描述和示范值",
	)
	a1.AddResponseBodyParameter("RequestBodyType", models.TYPE_STRING, models.BODY_TYPE_NONE,
		"",
		"请求Body的类型",
	)
	a1.AddResponseBodyParameter("ResponseBodyType", models.TYPE_STRING, models.BODY_TYPE_JSON,
		"",
		"响应Body的类型",
	)
	a1.AddResponseBodyParameter("RequestBodyParameters", models.TYPE_STRING, models.Parameter{},
		"",
		"当请求Body的类型为Json时，Body中的参数列表",
	)
	a1.AddResponseBodyParameter("ResponseBodyParameters", models.TYPE_STRING, models.Parameter{},
		"",
		"当响应Body的类型为Json时，Body中的参数列表",
	)
	a1.AddResponseBodyParameter("RequestBodyDemo", models.TYPE_STRING, "",
		"",
		"请求Body的示范",
	)
	a1.AddResponseBodyParameter("ResponseBodyDemo", models.TYPE_STRING, "",
		"",
		"响应Body的示范",
	)
	this.AppendDocument(a1)

	// getDateApi := models.NewApi("getdate", "/api/get/date", "", "GetDate", models.REQUEST_GET, false)

	// this.AppendDocument(getDateApi)
}
