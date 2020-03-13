package router

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
)

type T struct {
}

func (t *T) Get(c *gin.Context) {
	c.Writer.WriteString("hello from T.Get")
}

func (t *T) Post(c *gin.Context) {
	c.Writer.WriteString("hello from T.Post")
}

func (t *T) Delete(c *gin.Context, id int) {
	fmt.Println("param id:", c.Param("id"))
	fmt.Println("arg id:", id)
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

	//r.Any("/*path", AutoRoute(&T{}))
	r.Any("/*path", REST(&T{}))
	r.Run(":9091")

}

func ExampleRouteAny() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	//r.Any("/article/*path", AutoRouteAny(&T{}))
	r.Any("/article/*path", RESTAny(&T{}))
	r.Run(":9092")
}

func TestFindHandlerFuncs(t *testing.T) {
	funcs := findHandlerFuncs(&T{}, true)
	if len(funcs) != 2 {
		t.Fatal("wrong func number")
	}
	exist := func(name string) bool {
		_, ok := funcs[name]
		return ok
	}
	if !exist("get") || !exist("post") {
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

type T2 struct{}

func (t *T2) F1(c *gin.Context) { c.Writer.WriteString("ok") }
func (t *T2) F2(c *gin.Context) { c.Writer.WriteString("ok") }
func (t *T2) F3(c *gin.Context) { c.Writer.WriteString("ok") }

func BenchmarkManualRoute(b *testing.B) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	t := &T2{}
	r.GET("/f1", t.F1)
	r.GET("/f2", t.F2)
	r.GET("/f3", t.F3)

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/f1", nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
	}

}

func BenchmarkAutoRoute(b *testing.B) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Any("/:path", AutoRoute(&T2{}))
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/f1", nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
	}
}
