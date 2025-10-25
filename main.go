package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
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

// ---> S 自定义中间件 <---

func Logger() gin.HandlerFunc {

	return func(c *gin.Context) {

		t := time.Now()

		// 设置 example 变量到 Context 的键中
		c.Set("example", "12345")

		// 请求前
		c.Next()

		// 请求后
		// 计算请求耗时
		latency := time.Since(t)
		log.Print("latency", latency)

		// 获取发送的status

		status := c.Writer.Status()
		log.Println("status", status)

	}
}

// ---> S 自定义验证器 <---

type Booking struct {
	CheckIn  time.Time `form:"check_in" binding:"required,bookabledate" time_format:"2006-01-02"`
	CheckOut time.Time `form:"check_out" binding:"required,gtfield=CheckIn,bookabledate" time_format:"2006-01-02"`
}

var bookableDate validator.Func = func(fl validator.FieldLevel) bool {
	date, ok := fl.Field().Interface().(time.Time)
	if ok {
		today := time.Now()
		if today.After(date) {
			return false
		}
	}
	return true
}

func getBookable(c *gin.Context) {
	var b Booking
	if err := c.ShouldBindWith(&b, binding.Query); err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Booking dates are valid!"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

// ---> E 自定义验证器 <---

// ---> S 错误处理中间件 <---
// ErrorHandler captures errors and returns a consistent JSON error response
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // Step1: Process the request first.

		// Step2: Check if any errors were added to the context
		if len(c.Errors) > 0 {
			// Step3: Use the last error
			err := c.Errors.Last().Err

			// Step4: Respond with a generic error message
			c.JSON(http.StatusInternalServerError, map[string]any{
				"success": false,
				"message": err.Error(),
			})
		}

		// Any other steps if no errors are found
	}
}

// ---> E 错误处理中间件 <---

// ---> S 在中间件中使用 Goroutine <---

// ---> E 在中间件中使用 Goroutine <---

func main() {

	r := gin.New()

	r.Use(Logger())

	// ---> S jsonp <---
	
	r.GET("/JSONP", func(c *gin.Context) {
    data := map[string]interface{}{
      "foo": "bar",
    }

    // /JSONP?callback=x
    // 将输出：x({\"foo\":\"bar\"})
    c.JSONP(http.StatusOK, data)
  })
	
	// ---> E jsonp <---

	// ---> S Html 渲染 <---

	// r.LoadHTMLGlob("templates/*")
	r.LoadHTMLGlob("templates/**/*")

	r.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main website",
		})

	})

	 r.GET("/posts/index", func(c *gin.Context) {
    c.HTML(http.StatusOK, "posts/index.tmpl", gin.H{
      "title": "Posts",
    })
  })

	r.GET("/users/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "users/index.tmpl", gin.H{
			"title": "Users",
		})
	})

	// ---> E Html 渲染 <---

	// ---> S 路由组 <---

	{
		v1 := r.Group("/v1")
		v1.GET("/login", func(ctx *gin.Context) {
			ctx.String(200, "v1 login endpoint")
		})
		v1.GET("/submit", func(ctx *gin.Context) {
			ctx.String(200, "v1 submit endpoint")
		})
		v1.GET("/read", func(ctx *gin.Context) {
			ctx.String(200, "v1 read endpoint")
		})
	}

	// ---> E 路由组 <---

	// ---> S 在中间件中使用 Goroutine <---
	r.GET("/long_async", func(c *gin.Context) {
		// 创建在 goroutine 中使用的副本
		cCp := c.Copy()
		go func() {
			// 用 time.Sleep() 模拟一个长任务。
			time.Sleep(5 * time.Second)

			// 请注意您使用的是复制的上下文 "cCp"，这一点很重要
			log.Println("Done! in path " + cCp.Request.URL.Path)
		}()
	})

	r.GET("/long_sync", func(c *gin.Context) {
		// 用 time.Sleep() 模拟一个长任务。
		time.Sleep(5 * time.Second)

		// 因为没有使用 goroutine，不需要拷贝上下文
		log.Println("Done! in path " + c.Request.URL.Path)
	})
	// ---> E 在中间件中使用 Goroutine <---

	// ---> S 错误处理中间件 <---
	// Attach the error-handling middleware
	r.Use(ErrorHandler())

	r.GET("/ok", func(c *gin.Context) {
		somethingWentWrong := false

		if somethingWentWrong {
			c.Error(errors.New("something went wrong"))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Everything is fine!",
		})
	})

	r.GET("/error", func(c *gin.Context) {
		somethingWentWrong := true

		if somethingWentWrong {
			c.Error(errors.New("something went wrong"))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Everything is fine!",
		})
	})
	// ---> E 错误处理中间件 <---

	// ---> S 自定义路由日志的格式 <---
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Printf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}

	r.POST("/foo", func(c *gin.Context) {
		c.JSON(http.StatusOK, "foo")
	})

	r.GET("/bar", func(c *gin.Context) {
		c.JSON(http.StatusOK, "bar")
	})

	r.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, "ok")
	})
	// ---> E 自定义路由日志的格式 <---

	// ---> S 自定义中间件 <---
	r.GET("/test", func(c *gin.Context) {

		example := c.MustGet("example").(string)

		log.Println(example)

	})
	// ---> E 自定义中间件 <---

	// ---> S 自定义验证器 <---

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("bookabledate", bookableDate)
	}

	r.GET("/bookable", getBookable)
	// ---> E 自定义验证器 <---

	r.Run(":8080")

}

// ---> E 自定义中间件 <---

// func main() {
// 	强制日志颜色化
// 	gin.ForceConsoleColor()
// 	创建带默认中间件（日志与恢复）的 Gin 路由器
// 	r := gin.Default()

// 	---> S 设置和获取 Cookie <---
// 	r.GET("/cookie", func(ctx *gin.Context) {

// 		cookie, err := ctx.Cookie("gin_cookie")
// 		if err != nil {

// 			cookie = "NotSet"
// 			设置一个新的 Cookie
// 			ctx.SetCookie("gin_cookie", "test_cookie_value", 3600, "/", "localhost", false, true)

// 			fmt.Printf("Cookie value: %s \n", cookie)
// 		}

// 	})
// 	---> E 设置和获取 Cookie <---

// 	---> S 绑定 Uri <---
// 	r.GET("/:name/:id", func(c *gin.Context) {
// 		var person Person3
// 		if err := c.ShouldBindUri(&person); err != nil {
// 			c.JSON(400, gin.H{"msg": err.Error()})
// 			return
// 		}
// 		c.JSON(200, gin.H{"name": person.Name, "uuid": person.ID})
// 	})

// 	---> E 绑定 Uri <---

// 	定义简单的 GET 路由
// 	r.POST("/ping", func(c *gin.Context) {
// 		返回 JSON 响应
// 		c.JSON(http.StatusOK, gin.H{
// 			"message": "pong",
// 		})

// 	})
// 	r.POST("/test", func(c *gin.Context) {
// 		返回 JSON 响应
// 		SomeHandler(c)

// 	})

// 	---> S 绑定表单数据至自定义结构体 <---

// 	r.POST("/getb", GetDataB)
// 	r.POST("/getc", GetDataC)
// 	r.POST("/getd", GetDataD)
// 	r.POST("/gete", GetDataE)
// 	r.POST("/formHandler", formHandler)

// 	---> E 绑定表单数据至自定义结构体 <---

// 	---> S 静态资源 <---
// 	添加 GZIP 压缩

// 	配置静态资源
// 	r.Static("/assets", "./static/assets")
// 	r.StaticFile("/favicon.ico", "./static/favicon.ico")

// 	r.GET("/", func(c *gin.Context) {
// 		c.File("./static/index.html")
// 	})
// 	处理 SPA 路由
// 	r.NoRoute(func(c *gin.Context) {
// 		c.File("./static/index.html")
// 	})

// 	---> E 静态资源 <---

// 	默认端口 8080 启动服务器
// 	监听 0.0.0.0:800（Windows 下为 localhost:8080）

// 	---> S 绑定查询字符串或表单数据 <---

// 	r.POST("/startPage", startPage)
// 	---> E 绑定查询字符串或表单数据 <---

// 	---> S z自定义日志文件 <---

// 	---> E z自定义日志文件 <---
// 	LoggerWithFormatter 中间件会写入日志到 gin.DefaultWriter
// 	默认 gin.DefaultWriter = os.Stdout
// 	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
// 		你的自定义格式
// 		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
// 			param.ClientIP,
// 			param.TimeStamp.Format(time.RFC1123),
// 			param.Method,
// 			param.Path,
// 			param.Request.Proto,
// 			param.StatusCode,
// 			param.Latency,
// 			param.Request.UserAgent(),
// 			param.ErrorMessage,
// 		)
// 	}))
// 	r.Use(gin.Recovery())
// 	r.GET("/ping2", func(c *gin.Context) {
// 		c.String(200, "pong")
// 	})

// 	---> S 自定义 HTTP 配置 <---
// 	http.ListenAndServe(":8088", r)

// 	s := &http.Server{
// 		Addr:           ":8088",
// 		Handler:        r,
// 		ReadTimeout:    10 * time.Second,
// 		WriteTimeout:   10 * time.Second,
// 		MaxHeaderBytes: 1 << 20,
// 	}
// 	s.ListenAndServe()
// 	r.Run(":8088")

// 	---> E 自定义 HTTP 配置 <---

// }
