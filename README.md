# 简介
ginfra是基于gin搭建的一个服务框架，也可以认为是一个gin服务的模板。
ginfra提供了以下几个常用功能的封装：
- requestid
- prometheus metric
- structured logging and rotating, based on zap
- timeout
- gorm with context timeout
- config, based on viper

# 目的
* 提供服务开发模板，提高开发效率，统一代码结构；
* 提供多个常用中间件：生成requestid，封装日志组件，支持metric统计，支持超时处理等；
* 提供GORM支持context timeout的方案；
* 提供单元测试demo，以进一步促进单元测试的覆盖面，提升代码质量；

# 中间件
## requestid
requestid主要用于日志染色标记，便于日志检索。

生成requestid的逻辑:
* 首先，检查请求的header是否存在“X-Request-Id”，存在则复用该requestid；
* header中无“X-Request-Id”，则本地利用uuid.NewV4()生成一个新的requestid。

## metric
支持Prometheus指标。

## 日志
日记组件基于zap进行封装，支持日志切割，插入自定义关键字段，如requestid，client_ip，userid，productid等。

默认插入了requestid、client_ip字段，每次请求可以根据业务逻辑，插入其他关键字段，如：
```go
func Ping(c *gin.Context) {
    // set custom fileds into logger
    mw.SetLoggerField(c, map[string]string{
        mw.CtxProductID: "cbd271dec6133d7065bb5391a105f6ea",
        mw.CtxUserID:    "0qkkoqm22idmnmsno203u4nljdsf9",
    })  
 
    // mw.Log(c) returns the zap logger with custom fields
    mw.Log(c).Info("ping...pong")
    
    ...
}
```
日志内容如下：
```shell
{"level":"info","time":"2019-10-15T11:17:14.007+0800","caller":"handler/ping.go:17","msg":"ping...pong","client_ip":"127.0.0.1","request_id":"4ce2ee1d-5534-480c-a5b9-adc66af6b3fb","X-User-ID":"0qkkoqm22idmnmsno203u4nljdsf9","X-Product-ID":"cbd271dec6133d7065bb5391a105f6ea"}
{"level":"info","time":"2019-10-15T11:17:14.007+0800","caller":"middleware/metrics.go:65","msg":"/ping","client_ip":"127.0.0.1","request_id":"4ce2ee1d-5534-480c-a5b9-adc66af6b3fb","X-User-ID":"0qkkoqm22idmnmsno203u4nljdsf9","X-Product-ID":"cbd271dec6133d7065bb5391a105f6ea","status":200,"method":"GET","path":"/ping","query":"","ip":"127.0.0.1","user-agent":"curl/7.29.0","etime":"2019-10-15T11:17:14+08:00","latency":0.000080627}
```

## 超时处理
请求超时处理使用的是context.WithTimeout机制，在超时情况下，快速释放相关goroutine资源。

handler以及业务创建的goroutine，需要监听context的事件（不监听则超时控制不起效）：
```go
func TimedHandler(c *gin.Context) {

        ctx := c.Request.Context()

        // create the response data type to use as a channel type
        type responseData struct {
                status int
                body   map[string]interface{}
        }

        // create a done channel to tell the request it's done
        doneChan := make(chan responseData)

        go func(c *gin.Context) {
                select {

                // if the context is done it timed out or was cancelled
                // so don't return anything
                case <-ctx.Done():
                        c.AbortWithStatus(http.StatusGatewayTimeout)
                        mw.Log(c).Info("timeout, terminate sub goroutine...")
                        return

                // use timer to simulate I/O or logical opertion
                case <-time.After(time.Second * 30):
                        doneChan <- responseData{
                                status: 200,
                                body:   gin.H{"hello": "world"},
                        }
                }
        }(c)

        // non-blocking select on two channels see if the request
        // times out or finishes
        select {

        // if the context is done it timed out or was cancelled
        // so don't return anything
        case <-ctx.Done():
                c.AbortWithStatus(http.StatusGatewayTimeout)
                mw.Log(c).Info("timeout")
                return

        // if the request finished then finish the request by
        // writing the response
        case res := <-doneChan:
                c.JSON(res.status, res.body)
        }
}
```
日志如下：
```shell
{"level":"info","time":"2019-10-15T11:57:31.777+0800","caller":"handler/timeout.go:55","msg":"timeout","client_ip":"127.0.0.1","request_id":"d6e4ee5a-b5a9-4389-a6b3-48c679a6b7a4"}
{"level":"info","time":"2019-10-15T11:57:31.777+0800","caller":"handler/timeout.go:35","msg":"timeout, terminate sub goroutine...","client_ip":"127.0.0.1","request_id":"d6e4ee5a-b5a9-4389-a6b3-48c679a6b7a4"}
{"level":"info","time":"2019-10-15T11:57:31.777+0800","caller":"middleware/metrics.go:65","msg":"/timeout","client_ip":"127.0.0.1","request_id":"d6e4ee5a-b5a9-4389-a6b3-48c679a6b7a4","status":504,"method":"GET","path":"/timeout","query":"","ip":"127.0.0.1","user-agent":"curl/7.29.0","etime":"2019-10-15T11:57:31+08:00","latency":2.000209287}
```
# GORM
通过将sql db与context绑定，传入gorm DB中，即可实现gorm调用sql引擎层的时候，将相应api替换为withContext的api，从而实现context timeout能力：
```shell
    db, err := datasource.GormWithContext(ctx)
    if err != nil {
        return nil, err 
    } 
```

# 单元测试
## hander测试
利用httptest和httpexpect测试handler。

例：
```go
func Test_TagGet(t *testing.T) {

        // create a server for testing
        server := httptest.NewServer(g)
        defer server.Close()

        // create a test engine from server
        e := httpexpect.New(t, server.URL)

        // expect get / status is 200
        e.GET("/TagGet").
                Expect().
                Status(http.StatusOK).
                JSON().Object().ContainsKey("Response")
}
```

利用gomonkey进行mock：
```go
// Http Client
func Test_UseHttpClient(t *testing.T) {
        patches := ApplyFunc(req.Get, func(_ string, _ ...interface{}) (*req.Resp, error) {
                return nil, errors.New("failed")
        })
        defer patches.Reset()

        // create a server for testing
        server := httptest.NewServer(g)
        defer server.Close()

        // create a test engine from server
        e := httpexpect.New(t, server.URL)

        // expect get / status is 200
        e.GET("/UseHttpClient").
                Expect().
                Status(http.StatusInternalServerError)
}
```

## db相关单元测试
利用sqlmock测试db相关代码。

例：
```go
func Test_GetPostById(t *testing.T) {
        var (
                id    = 1
                title = "post title"
                body  = "blabla..."
                view  = 10
        )

        mock.ExpectQuery(regexp.QuoteMeta(
                "SELECT * FROM `posts` WHERE `posts`.`deleted_at` IS NULL AND ((id = ?)) ORDER BY `posts`.`id` ASC LIMIT 1")).
                WithArgs(id).
                WillReturnRows(sqlmock.NewRows([]string{"title", "body", "view"}).
                        AddRow(title, body, view))

        res, err := GetPostById(strconv.Itoa(id))
        convey.Convey("models.GetPostById", t, func() {
                convey.So(err, convey.ShouldEqual, nil)
        })
        convey.Convey("models.GetPostById", t, func() {
                convey.So(res, convey.ShouldResemble, &Post{Title: title, Body: body, View: view})
        })
}
```

## 执行单元测试
```shell
go test -cover -v ./...
```
