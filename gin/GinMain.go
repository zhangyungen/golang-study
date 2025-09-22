package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
	t "zyj.com/golang-study/gin/test"
)

func main() {
	router := gin.Default()
	gin.ForceConsoleColor()
	//gin中间件拦截器日志格式
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// 你的自定义格式
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	//路由日志格式内容
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Printf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}
	//重定向
	router.GET("/test", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "http://www.google.com/")
	})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	err := router.SetTrustedProxies([]string{"192.168.1.2"})
	if err != nil {
		log.Println(err)
		return
	}
	// 使用预定义的 gin.PlatformXXX 头
	// Google App Engine
	router.TrustedPlatform = gin.PlatformGoogleAppEngine
	// Cloudflare
	router.TrustedPlatform = gin.PlatformCloudflare
	// Fly.io
	router.TrustedPlatform = gin.PlatformFlyIO
	// 或者，你可以设置自己的可信请求头。但要确保你的 CDN
	// 能防止用户传递此请求头！例如，如果你的 CDN 将客户端
	// IP 放在 X-CDN-Client-IP 中：
	router.TrustedPlatform = "X-CDN-Client-IP"
	router.GET("/", func(c *gin.Context) {
		// 如果设置了 TrustedPlatform，ClientIP() 将解析
		// 对应的请求头并直接返回 IP
		fmt.Printf("ClientIP: %s\n", c.ClientIP())
	})

	router.GET("/getb", t.GetDataB)

	router.GET("/getb2", t.GetDataB)

	router.GET("/getc", t.GetDataC)

	router.GET("/getd", t.GetDataD)

	router.GET("/testing", t.StartPage)

	router.GET("/testing2", t.StartPage)

	router.GET("/testing3", t.StartPage)

	router.GET("/testing4", t.StartPage)

	router.GET("/testing5", t.StartPage)

	router.GET("/testing6", t.StartPage)

	router.StaticFS("/public", gin.Dir("/static/xxxapp", false))

	{
		v1 := router.Group("/v1/user")
		v1.POST("/login", t.StartPage)
		v1.POST("/submit", t.StartPage)
		v1.POST("/read", t.StartPage)
	}

	// 绑定 JSON ({"user": "manu", "password": "123"})
	router.POST("/loginJSON", func(c *gin.Context) {
		var json t.Login
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if json.User != "manu" || json.Password != "123" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
	})

	router.GET("/cookie", func(c *gin.Context) {

		cookie, err := c.Cookie("gin_cookie")

		if err != nil {
			cookie = "NotSet"
			c.SetCookie("gin_cookie", "test", 3600, "/", "localhost", false, true)
		}

		fmt.Printf("Cookie value: %s \n", cookie)
	})

	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.POST("/upload-batch", func(c *gin.Context) {
		// Multipart form
		form, _ := c.MultipartForm()
		files := form.File["upload[]"]

		for _, file := range files {
			log.Println(file.Filename)
			// 上传文件至指定目录
			err := c.SaveUploadedFile(file, "./files/"+file.Filename)
			if err != nil {
				log.Println("文件保存失败 ", err)
				c.String(http.StatusBadRequest, fmt.Sprintf("%d files error upload!", len(files)))
				return
			}
		}
		c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
	})
	router.POST("/upload", func(c *gin.Context) {
		// 单文件
		file, _ := c.FormFile("file")
		log.Println(file.Filename)

		dst := "./" + file.Filename
		// 上传文件至指定的完整文件路径
		err := c.SaveUploadedFile(file, dst+file.Filename)
		if err != nil {
			log.Println("文件保存失败 ", err)
			c.String(http.StatusBadRequest, "files error upload!")
			return
		}
		c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	})

	s := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	err = s.ListenAndServe()
	if err != nil {
		log.Println("启动失败 ", err)
		return
	}

	//router.Run(":8080") // 默认监听 0.0.0.0:8080
}
