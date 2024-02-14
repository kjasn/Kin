package main

import (
	"Kjasn/Kin/kin"
	"fmt"
	"log"
)

type Engine struct {
}



func main() {
	router := kin.New()


	g1 := router.Group("/v1")

	g1.Use(func(ctx *kin.Context) {
		log.Printf("the path is %s\n", ctx.Path)
	})
	
	g1.Use(m1, m2, m3)


	{
		// parameters using ':'
		// g1.GET("/golang", func(ctx *kin.Context) {
		// 	ctx.String(http.StatusOK, "the route is fixed")
		// })




		// g1.GET("/golang/p", func(ctx *kin.Context) {
		// 	ctx.String(http.StatusOK, "hello world~")
		// })


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


func m1(ctx *kin.Context) {
	fmt.Println("start m1...")
	ctx.Next()
	fmt.Println("end m1----")
}


func m2(ctx *kin.Context) {
	fmt.Println("start m2...")
	ctx.Next()
	fmt.Println("end m2----")
}

func m3(ctx *kin.Context) {
	fmt.Println("start m3...")
	ctx.Next()
	fmt.Println("end m3----")
}