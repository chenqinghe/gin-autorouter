package router

import (
	"errors"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// parseMethodAndArgs parse method and arguments from request context param "path"
// it will trim both left and right '/' first.
// if path empty, we will treat request http method  as method name and arguments empty.
func parseMethodAndArgs(c *gin.Context) (method string, args []string) {
	queryPath := c.Param("path")
	queryPath = strings.Trim(queryPath, "/")
	if queryPath == "" {
		return c.Request.Method, nil
	}

	segments := strings.Split(queryPath, "/")
	method = segments[0]
	segments = segments[1:]

	if len(segments) == 0 {
		return
	}

	return method, segments
}

// parseRESTMethodAndArgs is like parseMethodAndArgs, but all segments are parsed to arguments.
func parseRESTMethodAndArgs(c *gin.Context) (method string, args []string) {
	method = c.Request.Method
	path := c.Param("path")
	path = strings.Trim(path, "/")
	return method, strings.Split(path, "/")
}

func AutoRoute(handler interface{}) gin.HandlerFunc {
	return findAndCall(handler, true, parseMethodAndArgs)
}

func AutoRouteAny(handler interface{}) gin.HandlerFunc {
	return findAndCall(handler, false, parseMethodAndArgs)
}

func convertType(val string, inType reflect.Kind) (reflect.Value, error) {
	switch inType {
	case reflect.Int:
		tmp, err := strconv.Atoi(val)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(tmp), nil
	case reflect.Int32:
		tmp, err := strconv.Atoi(val)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(int32(tmp)), nil
	case reflect.Int64:
		tmp, err := strconv.Atoi(val)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(int64(tmp)), nil
	case reflect.Uint:
		i, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(uint(i)), nil
	case reflect.Uint32:
		i, err := strconv.ParseUint(val, 10, 32)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(uint32(i)), nil
	case reflect.Uint64:
		i, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(uint64(i)), nil
	case reflect.Float32:
		f, err := strconv.ParseFloat(val, 32)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(float32(f)), nil
	case reflect.Float64:
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(f), nil
	case reflect.String:
		return reflect.ValueOf(val), nil
	}

	return reflect.Value{}, errors.New("unsupport args type")
}

var ginContextType = reflect.TypeOf(&gin.Context{})

// findHandlerFuncs find valid methods of handler.
// if onlyCtx is true, valid means the method only have one parameter(except receiver) and its type is *gin.Context
// otherwise, valid means the method could have more than one parameters(except receiver too), but the first argument type
// must be *gin.Context, and others arguments type must be basic type, such as any int, any float, or string.
func findHandlerFuncs(handler interface{}, onlyCtx bool) map[string]reflect.Value {
	funcs := make(map[string]reflect.Value)
	rv := reflect.ValueOf(handler)
	rt := reflect.TypeOf(handler)

	for i := 0; i < rt.NumMethod(); i++ {
		fn := rv.Method(i)
		ft := rt.Method(i)

		if (!onlyCtx || ft.Type.NumIn() == 2) &&
			ft.Type.In(1) == ginContextType {
			funcs[strings.ToLower(ft.Name)] = fn
		}
	}
	return funcs
}

func REST(handler interface{}) gin.HandlerFunc {
	return findAndCall(handler, true, parseRESTMethodAndArgs)
}

func RESTAny(handler interface{}) gin.HandlerFunc {
	return findAndCall(handler, false, parseRESTMethodAndArgs)
}

func findAndCall(handler interface{}, onlyCtx bool, finder func(c *gin.Context) (string, []string)) gin.HandlerFunc {
	funcs := findHandlerFuncs(handler, onlyCtx)
	return func(c *gin.Context) {

		method, args := finder(c)
		fn, exist := funcs[strings.ToLower(method)]
		if !exist {
			http.NotFound(c.Writer, c.Request)
			c.Abort()
			return
		}

		numIn := fn.Type().NumIn() // include method receiver
		if (fn.Type().IsVariadic() && numIn-2 > len(args)) ||
			numIn-1 > len(args) { // not enough arguments
			http.NotFound(c.Writer, c.Request)
			c.Abort()
			return
		}

		arguments := make([]reflect.Value, 1, numIn)
		arguments[0] = reflect.ValueOf(c) // *gin.Context

		if !onlyCtx {
			t := numIn - 1 // non-variadic arguments number
			isVariadic := fn.Type().IsVariadic()
			if isVariadic {
				t--
			}

			popArg := func() (string, bool) {
				if len(args) > 0 {
					t := args[0]
					args = args[1:]
					return t, true
				}
				return "", false
			}
			for i := 0; i < t; i++ {
				argStr, ok := popArg()
				if !ok {
					break
				}
				arg, err := convertType(argStr, fn.Type().In(i+1).Kind())
				if err != nil {
					http.NotFound(c.Writer, c.Request)
					c.Abort()
					return
				}
				arguments = append(arguments, arg)
			}

			if isVariadic {
				if len(args) > 0 {
					variadicKind := fn.Type().In(numIn - 1).Elem().Kind()
					for {
						argStr, ok := popArg()
						if !ok {
							break
						}
						arg, err := convertType(argStr, variadicKind)
						if err != nil {
							http.NotFound(c.Writer, c.Request)
							c.Abort()
							return
						}
						arguments = append(arguments, arg)
					}
				}
			}
		}

		fn.Call(arguments)

	}

}
