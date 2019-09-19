package main

import (
	"context"
	"fmt"
	"github.com/kataras/iris"
	"go.mongodb.org/mongo-driver/bson"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"
	"github.com/iris-contrib/middleware/cors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Machine struct{
	name string
	value float64
}


func main() {
	db_ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(db_ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	fmt.Println(err)
	collection := client.Database("laundryMachineSJTU").Collection("machines")
	app := iris.New()
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*", "http://localhost:8080", "http://192.168.1.101:8080", "http://111.186.2.209:19034"},   //允许通过的主机名称
		AllowCredentials: true,
	})
	v1 := app.Party("/", crs).AllowMethods(iris.MethodOptions) // <- 对于预检很重要。
	{
		v1.Get("/ping", func(ctx iris.Context) {
			db_ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
			res, err := collection.InsertOne(db_ctx, bson.M{"name": "pi", "value": 3.14159})
			fmt.Println(err)
			id := res.InsertedID
			ctx.JSON(iris.Map{
				"message": id,
			})
		})
		v1.Get("/test", func(ctx iris.Context) {
			id := ctx.URLParam("id")
			fmt.Println(reflect.TypeOf(id))
			c := http.Client{Timeout: 5 * time.Second}
			res, _ := c.Get("https://www.weimaqi.net/Index_.aspx?device_name=" + id)
			defer res.Body.Close()
			body, _ := ioutil.ReadAll(res.Body)
			fmt.Println(string(body))
			if strings.Contains(string(body), "离线") {
				ctx.HTML("<h1> 离线 </h1>")
			} else {
				ctx.HTML("<h1> 可以用 </h1>")
			}
		})
		v1.Get("/check", func(ctx iris.Context) {
			id := ctx.URLParam("id")
			fmt.Println(reflect.TypeOf(id))
			c := http.Client{Timeout: 5 * time.Second}
			res, _ := c.Get("https://www.weimaqi.net/Index_.aspx?device_name=" + id)
			defer res.Body.Close()
			body, _ := ioutil.ReadAll(res.Body)
			fmt.Println(string(body))
			if strings.Contains(string(body), "离线") {
				ctx.JSON(iris.Map{"status": "offline"},)
			} else {
				ctx.JSON(iris.Map{"status": "online"},)
			}
		})
		v1.Get("/getList", func(ctx iris.Context) {
			building := ctx.URLParam("building") // the name of the building
			fmt.Println(reflect.TypeOf(building))

			db_ctx, _ = context.WithTimeout(context.Background(), 30*time.Second)
			cur, err := collection.Find(db_ctx, bson.D{})
			if err != nil { log.Fatal(err) }
			defer cur.Close(db_ctx)
			for cur.Next(db_ctx) {
				var result bson.M
				err := cur.Decode(&result)
				if err != nil { log.Fatal(err) }
				fmt.Println(result["name"])
			}
			if err := cur.Err(); err != nil {
				log.Fatal(err)
			}
		})
	}
	app.Run(iris.Addr(":8082"))
}
