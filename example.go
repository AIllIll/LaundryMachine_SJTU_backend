package main

import (
	"fmt"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type ServerLog struct {
	time string
	test string
}

type Building struct {
	name string `json:"name"`
	pwd string `json:"age"`
}

type Machine struct{
	name string
	value float64
}

var session *mgo.Session
var database *mgo.Database
var collection_server_log *mgo.Collection
var collection_building *mgo.Collection
var collection_machine *mgo.Collection

func main() {
	// init mongo client
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(nil)
	}
	defer session.Clone()
	session.SetMode(mgo.Monotonic, true)
	database = session.DB("sjtu")
	collection_building := database.C("building")
	collection_machine := database.C("machine")
	collection_server_log := database.C("server_log")
	collection_statistic := database.C("statistic")
	collection_server_log.Insert(bson.M{"time": time.Now(), "test": '0'})
	fmt.Println(bson.M{"time": time.Now()})

	// init iris
	app := iris.New()
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*", "http://localhost:8080", "http://192.168.1.101:8080", "http://111.186.2.209:19034"},   //允许通过的主机名称
		AllowCredentials: true,
	})
	v1 := app.Party("/", crs).AllowMethods(iris.MethodOptions) // <- 对于预检很重要。
	{
		v1.Get("/ping", func(ctx iris.Context) {
			ctx.JSON(iris.Map{
				"serverTime": time.Now(),
			})
		})
		v1.Get("/check", func(ctx iris.Context) {
			id := ctx.URLParam("id")
			c := http.Client{Timeout: 5 * time.Second}
			res, _ := c.Get("https://www.weimaqi.net/Index_.aspx?device_name=" + id)
			defer res.Body.Close()
			body, _ := ioutil.ReadAll(res.Body)
			// fmt.Println(string(body))
			if strings.Contains(string(body), "离线") {
				ctx.JSON(iris.Map{"status": "offline"},)
			} else {
				ctx.JSON(iris.Map{"status": "online"},)
			}
		})
		v1.Post("/addBuilding", func(ctx iris.Context) {
			var formData iris.Map
			ctx.ReadJSON(&formData)
			name := formData["name"]
			pwd := formData["pwd"]
			fmt.Println(name, pwd)
			if pwd == "meiyoumima" {
				collection_building.Upsert(bson.M{"name": name}, bson.M{"name": name})
				ctx.JSON(iris.Map{"Msg": "success"})
			} else {
				ctx.JSON(iris.Map{"Msg": "invalid password"})
			}
		})
		v1.Get("/getBuildingList", func(ctx iris.Context) {
			query := collection_building.Find(nil)
			var result []iris.Map
			query.All(&result)
			ctx.JSON(result)
		})
		v1.Post("/addMachine", func(ctx iris.Context) {
			var formData iris.Map
			ctx.ReadJSON(&formData)
			building := formData["building"]
			name := formData["name"]
			code := formData["code"]
			pwd := formData["pwd"]
			fmt.Println(building, name, code, pwd)
			if pwd == "meiyoumima" {
				num, err := collection_machine.Upsert(bson.M{"code": code}, bson.M{"building": building, "name": name, "code": code})
				fmt.Println(num, err)
				ctx.JSON(iris.Map{"Msg": "success"})
			} else {
				ctx.JSON(iris.Map{"Msg": "invalid password"})
			}
		})
		v1.Get("/getMachineList", func(ctx iris.Context) {
			building := ctx.URLParam("building")
			if len(building)!= 0 {
				query := collection_machine.Find(bson.M{"building": building})
				var result []iris.Map
				query.All(&result)
				ctx.JSON(result)
			} else {
				ctx.JSON(iris.Map{"Msg": "invalid building"})
			}
		})
		v1.Get("test", func(ctx iris.Context) {
			fmt.Println(ctx.GetHeader)
		})
		v1.Get("getUserNumber", func(ctx iris.Context) {
			num, err := collection_statistic.Upsert(bson.M{"name": "accessNumber"}, bson.M{"$inc": bson.M{"accessNumber":1}})
			fmt.Println(num, err)
			query := collection_statistic.Find(bson.M{"name":"accessNumber"})
			var result iris.Map
			query.One(&result)
			ctx.JSON(result)
		})
	}
	app.Run(iris.Addr(":2334"))
}
