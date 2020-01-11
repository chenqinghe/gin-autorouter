package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"testing"
)

type T struct {
}

func (t *T) Get(c *gin.Context) {
	c.Writer.WriteString("hello from T.Get")
}

func (t *T) Post(c *gin.Context) {
	c.Writer.WriteString("hello from T.Post")
}

func (t *T) Delete(c *gin.Context) {
	c.Writer.WriteString("hello from T.Delete")
}

func (t *T) Hello(c *gin.Context, name string) {
	c.Writer.WriteString("hello " + name)
}

func (t *T) AddInt(c *gin.Context, a ...int) {
	var sum int
	for _, v := range a {
		sum += v
	}
	c.Writer.WriteString(strconv.Itoa(sum))
}

func ExampleAutoRouter() {

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Any("/*path", AutoRouter(&T{}))

	r.Run(":9091")

}

func ExampleRouterAny() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Any("/*path", RouterAny(&T{}))

	r.Run(":9092")
}

func TestFindHandlerFuncs(t *testing.T) {
	funcs := findHandlerFuncs(&T{}, true)
	if len(funcs) != 3 {
		t.Fatal("wrong func number")
	}
	exist := func(name string) bool {
		_, ok := funcs[name]
		return ok
	}
	if !exist("get") || !exist("post") || !exist("delete") {
		t.Fatal("wrong funcs")
	}

	funcs = findHandlerFuncs(&T{}, false)
	if len(funcs) != 5 {
		t.Fatal("wrong funcs number")
	}
	if !exist("addint") || !exist("hello") {
		t.Fatal("wrong funcs")
	}
}

func TestParseMethodAndArgs(t *testing.T) {
	var (
		method string
		args   []string
	)
	c := &gin.Context{
		Request: &http.Request{
			Method: http.MethodPost,
		},
		Params: gin.Params{
			{
				Key:   "path",
				Value: "",
			},
		},
	}

	setpath := func(path string) {
		c.Params[0].Value = path
	}

	method, args = parseMethodAndArgs(c)
	if method != http.MethodPost || len(args) > 0 {
		t.Fatal("method or args error")
	}

	setpath("/")
	method, args = parseMethodAndArgs(c)
	if method != http.MethodPost || len(args) > 0 {
		t.Fatal("method or args error")
	}

	setpath("/foo/bar/")
	method, args = parseMethodAndArgs(c)
	if method != "foo" || len(args) != 1 || args[0] != "bar" {
		t.Fatal("method or args error")
	}

}
