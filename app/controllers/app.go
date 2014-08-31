package controllers

import (
	"github.com/sisteamnik/guseful/geoobject"

	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) Login(password string) revel.Result {
	if val, found := c.Session["u"]; found && val == "admin" {
		return c.Redirect(Admin.Index)
	}
	if c.Request.Method == "POST" {
		if password == "zmsk" {
			c.Session["u"] = "admin"
			return c.Redirect(Admin.Index)
		}
	}
	return c.Render()
}

func (c App) Village(slug string) revel.Result {
	o, err := geoobject.GetGeoObject(Db, slug)
	if err != nil {
		panic(err)
	}
	c.RenderArgs["village"] = o
	return c.Render()
}

func (c App) Communication(slug string) revel.Result {
	o, err := geoobject.GetGeoObject(Db, slug)
	if err != nil {
		panic(err)
	}
	c.RenderArgs["village"] = o
	return c.Render()
}

func (c App) Villages() revel.Result {
	villages, err := geoobject.GetGeoObjects(Db, "village", 0, 20)
	if err != nil {
		panic(err)
	}
	return c.Render(villages)
}

func (c App) Kontakti() revel.Result {
	return c.Render()
}

func metaVillages() map[string]string {
	type res struct {
		Slug, Name string
	}
	var vs []res
	var result = map[string]string{}
	_, err := Db.Select(&vs, "select Slug,Name from GeoObject")
	if err != nil {
		return result
	}
	for _, v := range vs {
		result[v.Name] = v.Slug
	}
	return result
}

func (c App) Stat() revel.Result {
	return c.RenderJson("ok")
}

func TemplateAddInfoFilter(c *revel.Controller, fc []revel.Filter) {
	c.RenderArgs["villages_meta"] = metaVillages()
	if c.Request.URL.Path == "" {
		c.RenderArgs["path"] = "/"
	} else {
		c.RenderArgs["path"] = c.Request.URL.Path
	}
	fc[0](c, fc[1:])
}
