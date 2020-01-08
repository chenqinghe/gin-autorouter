package router

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

func AutoRouter(handler interface{}) gin.HandlerFunc {
	funcs := findhandlerFuncs(handler, true)

	fmt.Println(funcs)

	return func(c *gin.Context) {
		queryPath := c.Param("path")
		fmt.Println("path:", queryPath)

		var funcName string
		if queryPath != "" {
			if strings.Contains(queryPath, "/") {
				c.Status(http.StatusNotFound)
				c.Abort()
				return
			}
			funcName = queryPath
		} else {
			funcName = c.Request.Method
		}

		fn, exist := funcs[funcName]
		if !exist {
			c.Status(http.StatusNotFound)
			c.Abort()
			return
		}
		fn.Call([]reflect.Value{reflect.ValueOf(c)})
	}
}

func RouterAny(handler interface{}) gin.HandlerFunc {
	funcs := findhandlerFuncs(handler, false)

	fmt.Println(funcs)

	return func(c *gin.Context) {
		queryPath := c.Param("path")
		fmt.Println("path:", queryPath)

		segments := strings.Split(queryPath, "/")

		funcName := segments[0]
		if funcName == "" {
			funcName = c.Request.Method
		}

		fn, exist := funcs[funcName]
		if !exist {
			c.Status(http.StatusNotFound)
			c.Abort()
			return
		}

		numIn := fn.Type().NumIn()
		if numIn > len(segments)+1 { // not enough arguments
			c.Status(http.StatusNotFound)
			c.Abort()
			return
		}

		// TODO: convert args type

		// TODO: handle variadic arguments
		//if fn.Type().IsVariadic() {
		//}

		arguments := []reflect.Value{reflect.ValueOf(c)}

		fn.Call(arguments)
	}
}

var ginContextType = reflect.TypeOf(&gin.Context{})

func findhandlerFuncs(handler interface{}, onlyOne bool) map[string]reflect.Value {
	funcs := make(map[string]reflect.Value)
	rv := reflect.ValueOf(handler)
	rt := reflect.TypeOf(handler)

	for i := 0; i < rt.NumMethod(); i++ {
		fn := rv.Method(i)
		ft := rt.Method(i)

		fmt.Println("method: ", ft.Name)
		fmt.Println("numin:", ft.Type.NumIn())

		if (!onlyOne || ft.Type.NumIn() == 2) &&
			ft.Type.In(1) == ginContextType {
			funcs[ft.Name] = fn
		}
	}
	return funcs
}
