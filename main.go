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

	g1 := router.Group("/v1")

	g1.Use(func(ctx *kin.Context) {
		log.Printf("the path is %s\n", ctx.Path)
	})

	{
		// parameters using ':'
		// g1.GET("/golang", func(ctx *kin.Context) {
		// 	ctx.String(http.StatusOK, "the route is fixed")
		// })



		// g1.GET("/golang/p", func(ctx *kin.Context) {
		// 	ctx.String(http.StatusOK, "hello world~")
		// })

		g1.GET("/cpp", func(ctx *kin.Context) {
			ctx.String(http.StatusOK, "this is cpp url")
		})
		g1.GET("/:lang", func(ctx *kin.Context) {
			ctx.JSON(http.StatusOK, kin.H {
				"lang" : ctx.Param("lang"),
			})
			ctx.String(http.StatusOK, "this is a dynamic route")
		})


		// wildcard '*' match
		// g1.GET("/test/*filepath", func(ctx *kin.Context) {
		// 	ctx.JSON(http.StatusOK, kin.H {
		// 		"filepath" : ctx.Param("filepath"),
		// 	})
		// })

	}
	err := router.Run(":80")
	if err != nil {
		panic(err)
	}
}



func indexHandler(ctx *kin.Context) {
	ctx.HTML(http.StatusOK, "<h1> Welcome </h1>")
}