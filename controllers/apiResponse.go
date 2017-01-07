package controllers

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type EmptySuccessResponse struct {
	Status string `json:"status"`
}

type ApiController struct {
	BaseController
	Complete bool
}

const (
	StatusSuccess = "success"
	StatusFail    = "fail"
)

func (this *ApiController) SetError(err error) bool {
	if err == nil {
		return false
	}
	result := ErrorResponse{StatusFail, err.Error()}
	this.Data["json"] = &result
	this.ServeJSON()
	this.Complete = true
	return true
}

func (this *ApiController) SetSuccessResponse() {
	this.Data["json"] = &EmptySuccessResponse{StatusSuccess}
	this.ServeJSON()
}

func (this *ApiController) SetData(result interface{}) {
	this.Data["json"] = &result
	this.ServeJSON()
	this.Complete = true
}
