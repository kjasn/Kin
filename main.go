package main

import (
	"Kjasn/Kin/kin"
	"log"
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

	g1 := router.Group("/v1")

	g1.Use(func(ctx *kin.Context) {
		log.Printf("the path is %s\n", ctx.Path)
	})

	{
		g1.GET("/hello", func(ctx *kin.Context) {
			ctx.String(http.StatusOK, "hello %s, the path is %s\n", 
			ctx.Query("name"), ctx.Path)
		})

		// parameters using ':'
		g1.GET("/hello/:lang", func(ctx *kin.Context) {
			ctx.JSON(http.StatusOK, kin.H {"lang" : ctx.Param("lang")})
		})

		// wildcard '*' match
		g1.GET("/test/*filepath", func(ctx *kin.Context) {
			ctx.JSON(http.StatusOK, kin.H {
				"filepath" : ctx.Param("filepath"),
			})
		})

	}
	err := router.Run(":80")
	if err != nil {
		panic(err)
	}
}



func indexHandler(ctx *kin.Context) {
	ctx.HTML(http.StatusOK, "<h1> Welcome </h1>")
}