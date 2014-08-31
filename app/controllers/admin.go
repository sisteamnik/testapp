package controllers

import (
	"github.com/revel/revel"
	"github.com/sisteamnik/guseful/chpu"
	"github.com/sisteamnik/guseful/geoobject"
	"github.com/sisteamnik/guseful/menu"
	"github.com/sisteamnik/guseful/pages"
)

type (
	Admin struct {
		*revel.Controller
	}
)

func (c Admin) Index() revel.Result {
	return c.Render()
}

func (c Admin) CheckLogin() revel.Result {
	if c.Session["u"] == "admin" {
		return nil
	}
	c.Flash.Error("Please log in first")
	return c.Redirect(App.Login)
}

func (c Admin) Pages() revel.Result {
	p, err := pages.GetPages(Db, 0, 10)
	if err != nil {
		panic(err)
	}
	c.RenderArgs["pages"] = p
	return c.Render()
}

func (c Admin) CreatePage(page pages.Page) revel.Result {
	_, err := pages.CreatePage(Db, &page)
	if err != nil {
		panic(err)
	}
	return c.Redirect(Admin.Pages)
}

func (c Admin) UpdatePage(page pages.Page) revel.Result {
	err := pages.UpdatePage(Db, &page)
	if err != nil {
		panic(err)
	}
	return c.Redirect(Admin.Pages)
}

func (c Admin) EditPage(slug string) revel.Result {
	p := &pages.Page{}
	c.RenderArgs["new_page"] = true
	var err error
	if slug != "" {
		p, err = pages.GetPage(Db, slug)
		if err != nil {
			panic(err)
		}
		c.RenderArgs["new_page"] = false
	}
	c.RenderArgs["page"] = p
	return c.Render()
}

func (c Admin) Menu() revel.Result {
	menus, err := menu.GetAllMenu(Db)
	if err != nil {
		panic(err)
	}
	return c.Render(menus)
}

func (c Admin) CreateMenu(title string) revel.Result {
	_, err := menu.CreateMenu(Db, title)
	if err != nil {
		panic(err)
	}
	return c.RenderJson("ok")
}

func (c Admin) Villages() revel.Result {
	villages, err := geoobject.GetGeoObjects(Db, "village", 0, 10)
	if err != nil {
		panic(err)
	}
	return c.Render(villages)
}

func (c Admin) CreateVillage(village *geoobject.GeoObject) revel.Result {
	err := geoobject.CreateGeoObject(Db, village)
	if err != nil {
		panic(err)
	}
	return c.Redirect(Admin.Villages)
}

func (c Admin) UpdateVillage(village *geoobject.GeoObject, longtext string) revel.Result {
	village.LongText = []byte(longtext)
	err := geoobject.UpdateGeoObject(Db, village)
	if err != nil {
		panic(err)
	}
	return c.Redirect(Admin.Villages)
}

func (c Admin) EditVillage(slug string) revel.Result {
	c.RenderArgs["new_village"] = true
	if slug != "" {
		o, err := geoobject.GetGeoObject(Db, slug)
		if err != nil {
			panic(err)
		}
		c.RenderArgs["village"] = o
		c.RenderArgs["new_village"] = false
		in := ImgName(o.SchemaId)
		if len(in) != 0 {
			c.RenderArgs["villageSchemaUrl"] = "/public/img/" + in[0:1] + "/" + in[1:2] + "/" +
				in[2:3] + "/" + in + "_100x100.jpg"
		}
		in = ImgName(o.PhotoId)
		if len(in) != 0 {
			c.RenderArgs["villagePhotoUrl"] = "/public/img/" + in[0:1] + "/" + in[1:2] + "/" +
				in[2:3] + "/" + in + "_100x100.jpg"
		}

	}
	return c.Render()
}

func (c Admin) Upload(imgdata []byte, name, descr string) revel.Result {
	name = chpu.Chpu(name)
	im, err := ImgApi.Create(imgdata, name, descr)
	if err != nil {
		panic(err)
	}
	return c.RenderJson(im)
}

func ImgName(id int64) string {
	name, err := Db.SelectStr("select Name from Img where Id = ?", id)
	if err != nil {
		panic(err)
	}
	return name
}

func ImgUrl(id int64) string {
	in := ImgName(id)
	if len(in) > 2 {
		return "/public/img/" + in[0:1] + "/" + in[1:2] + "/" +
			in[2:3] + "/" + in + ".jpg"
	}
	return ""
}
