package main

import (
	"Kjasn/Kin/kin"
	"io"
	"net/http"
	"os"
)




func main() {
	router := kin.New()

	router.Use(kin.Logger())


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

	err := router.Run(":80")
	if err != nil {
		panic(err)
	}
}