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

func (c *ApiController) SetError(err error) bool {
	if err == nil {
		return false
	}
	result := ErrorResponse{StatusFail, err.Error()}
	c.Data["json"] = &result
	c.ServeJSON()
	c.Complete = true
	return true
}

func (c *ApiController) SetSuccessResponse() {
	c.Data["json"] = &EmptySuccessResponse{StatusSuccess}
	c.ServeJSON()
}

func (c *ApiController) SetData(result interface{}) {
	c.Data["json"] = &result
	c.ServeJSON()
	c.Complete = true
}
