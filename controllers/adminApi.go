package controllers

import (
	ai "github.com/BabichMikhail/Hanabi/AI"
)

type ApiAdminController struct {
	ApiController
}

func (c *ApiAdminController) GetAINames() {
	c.SetData(ai.AINames)
}
