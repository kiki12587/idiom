package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

var (
	wg   sync.WaitGroup
	lock sync.Mutex
)

//查询类型 1= 模糊查询 2=精确查询
func (c *MainController) Findcy(ty int, keyword string) {

	if ty == 1 { //模糊查询
		c.Vague02(keyword)
	}

	if ty == 2 { //精确查询
		c.Accurate(keyword)
	}
	return
}

//同步模糊查询
func (c *MainController) Vague01(keyword string) {
	var url = "https://route.showapi.com/1196-1?keyword=" + keyword + "&page=1&rows=10&showapi_appid=468266&showapi_sign=9ae2de37560b4eeeb3ef0a7a07d8dd7f"
	str, _ := c.GetJson(url)
	time.Sleep(time.Second)
	tempData := make(map[string]interface{})
	err := json.Unmarshal([]byte(str), &tempData)
	if err != nil {
		retMp = &RetMap{
			Code:    -1,
			Message: "查询失败",
		}
	}
	if tempData["showapi_res_code"].(float64) != 0 {
		retMp = &RetMap{
			Code:    -1,
			Message: tempData["showapi_res_body"].(map[string]interface{})["ret_message"].(string),
		}
	}

	dataSlice := tempData["showapi_res_body"].(map[string]interface{})["data"].([]interface{})
	for _, v := range dataSlice {
		title := v.(map[string]interface{})["title"].(string)
		accurate2, err := c.Accurate01(title)
		if err == nil {
			idiom := Idiom{
				Title:      title,
				Spell:      accurate2["spell"].(string),
				Content:    accurate2["content"].(string),
				Derivation: accurate2["derivation"].(string),
				Samples:    accurate2["samples"].(string),
			}
			idiomsMap[title] = idiom
			retMp = &RetMap{
				Code:    0,
				Message: "查询成功",
				Res:     idiomsMap,
			}

		} else {
			continue
		}
	}
	c.Data["json"] = retMp
}

//配合同步模糊查询
func (c *MainController) Accurate01(keyword string) (ret map[string]interface{}, err error) {
	//time.Sleep(time.Second)
	var url = "https://route.showapi.com/1196-2?keyword=" + keyword + "&page=1&rows=10&showapi_appid=468266&showapi_sign=9ae2de37560b4eeeb3ef0a7a07d8dd7f"
	str, _ := c.GetJson(url)
	tempData := make(map[string]interface{})
	err = json.Unmarshal([]byte(str), &tempData)
	ret = tempData["showapi_res_body"].(map[string]interface{})["data"].(map[string]interface{})

	if err != nil {
		fmt.Println("成语app报错信息", err)
		return nil, err
	} else {
		fmt.Println("ret", ret)
		return ret, nil
	}
}

//异步模糊查询
func (c *MainController) Vague02(keyword string) {
	var url = "https://route.showapi.com/1196-1?keyword=" + keyword + "&page=1&rows=10&showapi_appid=468266&showapi_sign=9ae2de37560b4eeeb3ef0a7a07d8dd7f"
	str, _ := c.GetJson(url)
	time.Sleep(time.Second)
	tempData := make(map[string]interface{})
	err := json.Unmarshal([]byte(str), &tempData)
	if err != nil {
		retMp = &RetMap{
			Code:    -1,
			Message: "查询失败",
		}
	}

	if tempData["showapi_res_code"].(float64) != 0 {
		retMp = &RetMap{
			Code:    -1,
			Message: tempData["showapi_res_body"].(map[string]interface{})["ret_message"].(string),
		}

	} else {
		dataSlice := tempData["showapi_res_body"].(map[string]interface{})["data"].([]interface{})
		for _, v := range dataSlice {
			title := v.(map[string]interface{})["title"].(string)
			wg.Add(1)
			go c.Accurate02(title)
		}

		wg.Wait()
		close(retCh)

		for i := 0; i < 10; i++ {
			v := <-retCh
			idiomsMap[v.Title] = v
		}

		retMp = &RetMap{
			Code:    0,
			Message: "查询成功",
			Res:     idiomsMap,
		}
	}

	c.Data["json"] = retMp
}

//配合异步模糊查询
func (c *MainController) Accurate02(keyword string) {
	lock.Lock()
	defer lock.Unlock()
	var url = "https://route.showapi.com/1196-2?keyword=" + keyword + "&page=1&rows=10&showapi_appid=468266&showapi_sign=9ae2de37560b4eeeb3ef0a7a07d8dd7f"
	str, _ := c.GetJson(url)
	//fmt.Println("查询返回数据str: ", str)
	tempData := make(map[string]interface{})
	json.Unmarshal([]byte(str), &tempData)
	if tempData["showapi_res_code"].(float64) != 0 {
		fmt.Println("调用失败原因: ", tempData)
	} else {
		ret := tempData["showapi_res_body"].(map[string]interface{})["data"].(map[string]interface{})
		//fmt.Println("查询返回数据: ", ret)
		idiom := Idiom{
			Title:      ret["title"].(string),
			Spell:      ret["spell"].(string),
			Content:    ret["content"].(string),
			Derivation: ret["derivation"].(string),
			Samples:    ret["samples"].(string),
		}
		retCh <- idiom
	}
	wg.Done()
	return

}

//精确查询
func (c *MainController) Accurate(keyword string) {
	var url = "https://route.showapi.com/1196-2?keyword=" + keyword + "&page=1&rows=10&showapi_appid=468266&showapi_sign=9ae2de37560b4eeeb3ef0a7a07d8dd7f"
	str, _ := c.GetJson(url)
	tempData := make(map[string]interface{})
	err := json.Unmarshal([]byte(str), &tempData)
	if err != nil {
		res = Ret{
			Code:    -1,
			Message: "查询失败",
		}
	}
	if tempData["showapi_res_body"].(map[string]interface{})["ret_code"].(float64) != 0 {
		res = Ret{
			Code:    -1,
			Message: tempData["showapi_res_body"].(map[string]interface{})["ret_message"].(string),
		}
	} else {
		ret := tempData["showapi_res_body"].(map[string]interface{})["data"].(map[string]interface{})
		fmt.Println("ret", ret)
		idiom := Idiom{
			Title:      ret["title"].(string),
			Spell:      ret["spell"].(string),
			Content:    ret["content"].(string),
			Derivation: ret["derivation"].(string),
			Samples:    ret["samples"].(string),
		}
		//idiomsMap[keyword] = idiom
		res = Ret{
			Code:    0,
			Message: "查询成功",
			Idiom:   idiom,
		}
	}

	c.Data["json"] = res
	return
}

//获取网络资源 转成json
func (c *MainController) GetJson(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("http 请求失败")
		return "", err
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("网球读取资源失败")
		return "", err
	}
	respStr := string(bytes)
	return respStr, nil
}
