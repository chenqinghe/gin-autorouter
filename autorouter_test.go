package router

import (
	"github.com/gin-gonic/gin"
	"testing"
)

func TestAutoRouter(t *testing.T) {

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Any("/", AutoRouter(&T{}))

	r.Run(":9091")

}

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
