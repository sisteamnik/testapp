package controllers

import (
	"database/sql"
	"github.com/coopernurse/gorp"
	_ "github.com/mattn/go-sqlite3"
	"github.com/revel/revel"
	"github.com/sisteamnik/guseful/geoobject"
	"github.com/sisteamnik/guseful/img"
	"github.com/sisteamnik/guseful/menu"
	"github.com/sisteamnik/guseful/pages"
	"log"
	"os"
	"strings"
)

func init() {
	revel.OnAppStart(Init)
	revel.TemplateFuncs["imgurl"] = func(a int64, b string) string {
		return strings.Replace(ImgUrl(a), ".jpg", b+".jpg", -1)
	}
	revel.TemplateFuncs["str"] = func(a []byte) string {
		return string(a)
	}
}

var Db *gorp.DbMap
var ImgApi *img.Api

func Init() {
	var dbpath = revel.BasePath + "/db"
	db, _ := sql.Open("sqlite3", dbpath)

	Db = &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	Db.TraceOn("[GORP]", log.New(os.Stdout, "myapp:", log.Lmicroseconds))
	t := Db.AddTable(geoobject.GeoObject{}).SetKeys(true, "Id")
	t.ColMap("SchemaData").SetMaxSize(3000)
	t.ColMap("Slug").SetUnique(true)
	Db.AddTable(geoobject.GeoPoint{}).SetKeys(true, "Id")
	t = Db.AddTable(pages.Page{}).SetKeys(true, "Id")
	t.ColMap("Slug").SetUnique(true)
	t = Db.AddTable(img.Img{}).SetKeys(true, "Id")
	t.ColMap("Name").SetUnique(true)

	Db.AddTable(menu.Menu{}).SetKeys(true, "Id")
	Db.AddTable(menu.MenuItem{}).SetKeys(true, "Id")

	Db.CreateTablesIfNotExists()
	ImgApi = img.NewApi(Db, revel.BasePath+"/public/img", 3, "Img",
		"/default.jpg", []img.Size{
			img.Size{
				100, 100, "thumb",
			},
			img.Size{
				900, 600, "fit",
			},
			img.Size{
				300, 300, "thumb",
			},
		})
	revel.InterceptMethod(Admin.CheckLogin, revel.BEFORE)
}
