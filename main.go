package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/Kjasn/Kin/kin"
)
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

	router.GET("/index/:lang/doc", func(ctx *kin.Context) {
		lang, ok := ctx.Param("lang")
		if !ok {
			log.Panicln("field not exists")
		}
		log.Println("matches " + ctx.Path, ",lang is ", lang)
	})

	router.GET("/index/go/doc", func(ctx *kin.Context) {
		log.Println("matches " + ctx.Path)
	})

	err := router.Run(":80")
	if err != nil {
		panic(err)
	}
}