package serve

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/elazarl/goproxy"
	"github.com/elazarl/goproxy/ext/auth"
	"go.zoe.im/proxman/cmd"
	"go.zoe.im/x"
	"go.zoe.im/x/cli"
)

type options struct {
	Addr     string `opts:"help=Address for listen"`
	Config   string `opts:"short=c,help=Configuration file"`
	Verbose  bool   `opts:"short=v,help=Verbose outout more information"`
	Username string `opts:"short=u,help=Username for proxy authorization"`
	Password string `opts:"short=p,help=Password for proxy authorization"`
}

type server struct {
	opts  *options
	proxy *goproxy.ProxyHttpServer
}

func (s server) Run() error {
	s.proxy.Verbose = s.opts.Verbose

	fmt.Println("注册八方城会员数据修改工具")
	s.proxy.OnResponse(
		goproxy.ReqConditionFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
			if req != nil && req.URL != nil && req.URL.Path == "/web/customer/queryMemberData" {
				return true
			}
			return false
		}),
	).DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		b, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			// 替换手机号
			b = bytes.Replace(b, x.Str2Bytes(`"mobile":""`), x.Str2Bytes(`"mobile":"13812911234"`), -1)
			b = bytes.Replace(b, x.Str2Bytes(`"mobile":null`), x.Str2Bytes(`"mobile":"13812911234"`), -1)
			// create a new reader from b
			resp.Body = ioutil.NopCloser(bytes.NewReader(b))
		} else {
			fmt.Println("解析错误:", err)
		}
		return resp
	})

	auth.ProxyBasic(s.proxy, "", func(user, passwd string) bool {
		if user == s.opts.Username && passwd == s.opts.Password {
			return true
		}
		return false
	})

	fmt.Println("代理监听:", s.opts.Addr)
	return http.ListenAndServe(s.opts.Addr, s.proxy)
}

func newServer() *server {
	return &server{
		opts: &options{
			Addr:     ":8080",
			Username: "liya-zoe",
			Password: "zoe-liya",
		},
		proxy: goproxy.NewProxyHttpServer(),
	}
}

func init() {

	svr := newServer()

	cmd.Register(
		cli.New(
			cli.Name("serve"),
			cli.Short("serve the http serve"),
			cli.Config(svr.opts),
			cli.Run(func(c *cli.Command, args ...string) {
				err := svr.Run()
				if err != nil {
					log.Fatal(err)
				}
			}),
		),
	)
}

/**
modify:
- match:
	url: ^http(s)://www.baidu.com
  request:
  response:
	body: {{repace response.body}}
*/

// type match struct {
// 	URL     string   `json:"url"` // regex
// 	Methods []string `json:"methods"`
// 	Custom  string   `json:"custom"` // rule
// }

// type modify struct {
// 	URL     string `json:"url"`
// 	Request `json:request`
// }

// type config struct {
// 	modifies []modify
// }
