package main

import (
	"Kjasn/Kin/kin"
	"net/http"
)

type Engine struct {
}



func main() {
	router := kin.New()

	router.GET("/", indexHandler)
	router.GET("/ping", func(ctx *kin.Context) {
		ctx.JSON(http.StatusOK, kin.H{
			"name": "kjasn", 
			"opt": "test",
		})
	})
	router.GET("/hello", func(ctx *kin.Context) {
		ctx.String(http.StatusOK, "hello %s, the path is %s\n", 
		ctx.Query("name"), ctx.Path)
	})

	// parameters using ':'
	router.GET("/hello/:lang", func(ctx *kin.Context) {
		ctx.JSON(http.StatusOK, kin.H {"lang" : ctx.Param("lang")})
	})

	// wildcard '*' match
	router.GET("/test/*filepath", func(ctx *kin.Context) {
		ctx.JSON(http.StatusOK, kin.H {
			"filepath" : ctx.Param("filepath"),
		})
	})

	err := router.Run(":80")
	if err != nil {
		panic(err)
	}
}



func indexHandler(ctx *kin.Context) {
	ctx.HTML(http.StatusOK, "<h1> Welcome </h1>")
}