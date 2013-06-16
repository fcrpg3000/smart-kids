package controllers

import (
	"github.com/robfig/revel"
	"smart-kids/ruler/app/models"
)

type Application struct {
	GorpController
}

func (c Application) Index() revel.Result {
	obj, err := c.Txn.Get(models.Admin{}, int64(1))
	if err != nil {
		panic(err)
	}
	admin := obj.(*models.Admin)
	return c.Render(admin)
}
