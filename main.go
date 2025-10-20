package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type formA struct {
	Foo string `json:"foo"  xml:"foo" binding:"required" `
}

type formB struct {
	Bar string `json:"bar"  xml:"bar" binding:"required" `
}

func SomeHandler(c *gin.Context) {

	objA := formA{}

	objB := formB{}

	if err := c.ShouldBind(&objA); err == nil {

		c.String(http.StatusOK, `the body should be formA`)

	} else if err := c.ShouldBind(&objB); err == nil {

		c.String(http.StatusOK, `the body should be formB`)

	} else {
		c.String(http.StatusOK, `the body is valid`)
	}

}

// ---> S 绑定表单数据至自定义结构体 <---

type StructA struct {
	FiledA string `form:"field_a"`
}

type StructB struct {
	NestedStruct StructA
	FiledB       string `form:"field_b"`
	// 标签格式 └──┬──┘ └──┬────┘
	//         标签类型   参数名
}

type StructC struct {
	NestedStructPointer *StructA
	FieldC              string `form:"field_c"`
}

type StructD struct {
	NestedAnonyStruct struct {
		FieldX string `form:"field_x"`
	}
	FieldD string `form:"field_d"`
}

type Persion struct {
	Name string
	Age  int
}

func GetDataB(c *gin.Context) {

	// 声明接收请求数据的结构体实例
	var b StructB
	// 核心绑定操作：将请求数据自动填充到结构体
	c.Bind(&b) // 等效于 ShouldBindWith 的快捷方式
	// 打印 b 的内容以便调试
	fmt.Printf(" b >>> %+v\n", b)

	// gin.H 是 map[string]interface{} 的简写
	c.JSON(200, gin.H{
		"a": b.NestedStruct.FiledA,
		"b": b.FiledB,
	})

}

func GetDataC(c *gin.Context) {

	var b StructC

	c.Bind(&b)

	fmt.Printf(" b >>> %+v\n", b)

	c.JSON(200, gin.H{
		"a": b.NestedStructPointer,
		"c": b.FieldC,
	})

}

type myForm struct {
	Colors []string `form:"colors[]"`
}

func formHandler(c *gin.Context) {

	var fakeForm myForm
	c.Bind(&fakeForm)
	c.JSON(200, gin.H{
		"colors": fakeForm.Colors,
	})

}

func GetDataD(c *gin.Context) {
	var b StructD
	c.Bind(&b)
	c.JSON(200, gin.H{
		"x": b.NestedAnonyStruct,
		"d": b.FieldD,
	})
}

func GetDataE(c *gin.Context) {
	var p Persion
	c.Bind(&p)
	c.JSON(200, gin.H{
		"name": p.Name,
		"age":  p.Age,
	})
}

// -

// ---> E 绑定表单数据至自定义结构体 <---

// ---> S 绑定查询字符串或表单数据 <---

type Persion2 struct {
	Name     string    `from:"name"`
	Addres   string    `form:"address"`
	Birthday time.Time `form:"birthday" time_format:"2006-01-02"`
}

func startPage(c *gin.Context) {

	var person Persion2

	if c.ShouldBind(&person) == nil {

		fmt.Printf(" person >>> %+v\n", person)

		log.Println(person.Name)
		log.Println(person.Addres)
		log.Println(person.Birthday)

	}
	fmt.Printf(" person >>> %+v\n", person)

	c.String(200, "Success")

}

// ---> E 绑定查询字符串或表单数据 <---

// ---> S 绑定 Uri <---

type Person3 struct {
	ID   string `uri:"id" binding:"required,uuid"`
	Name string `uri:"name" binding:"required"`
}

// ---> E 绑定 Uri <---

func main() {
	// 创建带默认中间件（日志与恢复）的 Gin 路由器
	r := gin.Default()

	// ---> S 绑定 Uri <---
	r.GET("/:name/:id", func(c *gin.Context) {
		var person Person3
		if err := c.ShouldBindUri(&person); err != nil {
			c.JSON(400, gin.H{"msg": err.Error()})
			return
		}
		c.JSON(200, gin.H{"name": person.Name, "uuid": person.ID})
	})

	// ---> E 绑定 Uri <---

	// 定义简单的 GET 路由
	r.POST("/ping", func(c *gin.Context) {
		// 返回 JSON 响应
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})

	})
	r.POST("/test", func(c *gin.Context) {
		// 返回 JSON 响应
		SomeHandler(c)

	})

	// ---> S 绑定表单数据至自定义结构体 <---

	r.POST("/getb", GetDataB)
	r.POST("/getc", GetDataC)
	r.POST("/getd", GetDataD)
	r.POST("/gete", GetDataE)
	r.POST("/formHandler", formHandler)

	// ---> E 绑定表单数据至自定义结构体 <---

	// ---> S 静态资源 <---
	// 添加 GZIP 压缩

	// 配置静态资源
	r.Static("/assets", "./static/assets")
	r.StaticFile("/favicon.ico", "./static/favicon.ico")

	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})
	// 处理 SPA 路由
	// r.NoRoute(func(c *gin.Context) {
	// 	c.File("./static/index.html")
	// })

	// ---> E 静态资源 <---

	// 默认端口 8080 启动服务器
	// 监听 0.0.0.0:800（Windows 下为 localhost:8080）

	// ---> S 绑定查询字符串或表单数据 <---

	r.POST("/startPage", startPage)
	// ---> E 绑定查询字符串或表单数据 <---
	

   r.Run(":8088")
}
