package main

import (
	"io"
	"net/http"
	"os"

	"github.com/Kjasn/Kin/kin"
)


type demo1 struct {
	a int8
	b int16
	c int32
}

type demo2 struct {
	a int8
	c int32
	b int16
}


func main() {
	router := kin.Default()

	router.LoadHTMLGlob("./static/*")
	router.Static("/assets", "./static")

	router.GET("/hello", func(ctx *kin.Context) {
		ctx.HTML(http.StatusOK, "template.html", nil)
	})

	router.GET("/holo", func(ctx *kin.Context) {
		pic, err := os.Open("./static/file1.jpg")
		if err != nil {
			ctx.Fail(500, "file not exist")
		}
		defer pic.Close()

		data, err := io.ReadAll(pic)
		if err != nil {
			ctx.Fail(500, "read failed")
		}
		ctx.Data(http.StatusOK, data)
	})

	router.GET("/panic", func(ctx *kin.Context) {
		ctx.String(http.StatusOK, "something occurred error~\n")
		names := []string{"hello everyone"}
		ctx.String(http.StatusOK, names[100])
	})

	err := router.Run(":80")
	if err != nil {
		panic(err)
	}
}