package controllers

import (
	"github.com/astaxie/beego"
	"strconv"
	"strings"
)

type MainController struct {
	beego.Controller
}

type Ret struct {
	Code    int
	Message string
	Idiom
}

type RetMap struct {
	Code    int
	Message string
	Res     map[string]Idiom
}

type Idiom struct {
	Title      string
	Spell      string
	Content    string
	Derivation string
	Samples    string
}

var (
	idiomsMap map[string]Idiom
	res       Ret     //精确查询 数据返回
	retMp     *RetMap //模糊查询数据返回
	retCh     chan Idiom
)

func (c *MainController) Get() {
	c.TplName = "index.html"
}

func (c *MainController) Post() {
	idiomsMap = make(map[string]Idiom)
	res = Ret{}
	retMp = &RetMap{}
	retCh = make(chan Idiom, 10)
	tp, err := strconv.Atoi(c.GetString("type"))
	if err != nil {
		c.Ctx.WriteString("类型选择错误")
		return
	}
	keyword := c.GetString("keyword")
	keyword = strings.Replace(keyword, " ", "", -1)
	if keyword == "" {
		c.Ctx.WriteString("请输入成语关键字")
		return
	}
	c.Findcy(tp, keyword)
	c.ServeJSON()
	return
}
