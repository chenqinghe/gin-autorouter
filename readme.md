# gin-autorouter

gin-autorouter is a middleware which could automatic mapping request url to a handler method.

## Api
the project have four main functions:

* AutoRoute
* RouteAny
* REST
* RESTAny


## Basic Uage

```go
package main


type T struct{
	
}

func (t *T)Greet(c *gin.Context)  {
    c.Writer.WriteString("hello from *T.Greet")
}

func (t *T)Hello(c *gin.Context) {
    c.Writer.WriteString("hello from *T.Hello")
}


func main(){
	r:=gin.Default()
	r.Any("/*path",router.AutoRoute(&T{}))
	r.Run(":8080")	
}


```
you only need to register a router with pattern "/*path"

view http://localhost:8080/greet, you could see "hello from *T.Greet"

view http://localhost:8080/hello, you will see "hello from *T.Hello"

## RESTful Api
with AutoRouter, you can create restful api very easily.

```go
package main

func (h *Article) Get(c *gin.Context) {
	articleId := c.Param("id")
	// search artile stuff.....
	article := model.SearchArticle(articleId)
	c.JSONP(http.StatusOK,article)
}

func (h *Article)Delete(c *gin.Context) {
	articleId := c.Param("id")
	model.DeleteArticle(articleId)
	c.JSONP(http.StatusOK,"ok")
}

func main(){
	r := gin.Default()
	r.Any("/article/:id",router.REST(&Article{}))
}

//  * GET /article/123 => *Article.Get
//  * Delete /article/123 => *Article.Delete
```
also, you can use RESTAny, things will be extremely easy!!

```go
package main

func (h *Article)Get(c *gin.Context, id int) {
	fmt.Println("article:",id) // output: article: 123
	article := model.SearchArticle(id)
	c.JSONP(http.StatusOK,article)
}

func (h *Article)Delete(c *gin.Context, id int) {
	fmt.Println("article:",id) // output: article: 123
	model.DeleteArticle(id)
	c.JSONP(http.StatusOK,"ok")
}


func main(){
	r:= gin.Default()
	r.Any("/article/:path",router.RESTAny(&Article{}))
}

// GET /article/123 => *Article.Get(c, 123)
// DELETE /article/123 => *Article.Delete(c, 123)


```

## Mapping Rules 
the mapping is basic on *gin.Context.Param("path"). path will be exploded to several segments by '/'
* if path is empty, method is request http method
* the first segment is method
* others will be method arguments
* segments number MUST be equal or greater than method arguments number
* if method is variadic, the segments mapping to last argument could be zero
* otherwise, "404 not found" will be returned


some examples: 
<table>
    <tr>
        <td>path</td>
        <td>method</td>
        <td>arguments(exclude *gin.Context)</td>
    </tr>
    <tr>
        <td>/</td>
        <td>REQUEST METHOD</td>
        <td>nil</td>
    </tr> 
    <tr>
        <td>/foo</td>
        <td>foo</td>
        <td>nil</td>
    </tr>
    <tr>
        <td>/foo/bar</td>
        <td>foo</td>
        <td>[bar]</td>
    </tr>
    <tr>
        <td>/foo/bar/123</td>
        <td>foo</td>
        <td>[bar,123]</td>
    </tr>
</table>

## License

the project is under MIT license protected which you can find in [LICENSE](https://github.com/chenqinghe/gin-autorouter/blob/master/LICENSE) file.