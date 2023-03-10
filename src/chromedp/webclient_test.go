package chromedp

import (
	"fmt"
	"merkaba/common"
	"testing"
	"time"
)

func TestWebClientLoadHomePage(t *testing.T) {
	common.InitEnviroment()
	client := NewClient(WebClientOption{Domain: "jd.com"}, &common.RunContext{})
	url := "https://www.jd.com"
	page := client.Load("home", url)
	element, _ := page.Select("#areamini", 0, false)
	fmt.Println(element.Html())
	fmt.Println(element.Text())
	values := element.Attrs()
	fmt.Println(values["href"])
}

func TestWebClientSearch(t *testing.T) {
	common.InitEnviroment()
	client := NewClient(WebClientOption{Domain: "jd.com"}, &common.RunContext{})
	page := client.Load("search", "https://search.jd.com/Search?keyword=iphone&enc=utf-8&wq=iphone&pvid=010a8e5fd9cb4d5eb5458cd6d1d5dd9e")
	elements, _ := page.Selects("#J_goodsList li.gl-item", 0, false)
	for index, e := range elements {
		img, _ := e.Select("div.p-img img", 0, false)
		fmt.Printf("%d ====> %s\n\n", index, e.Text())
		newPage, _ := img.ClickPage()
		ee, _ := newPage.Select("div.sku-name", 0, false)
		fmt.Println(ee.Text())
		break
	}
	time.Sleep(3 * time.Minute)
}

func TestWebClientLogin(t *testing.T) {
	common.InitEnviroment()
	client := NewClient(WebClientOption{Domain: "jd.com", Proxy: true, Headless: true}, &common.RunContext{})
	url := "https://passport.jd.com/new/login.aspx?ReturnUrl=https%3A%2F%2Fglobal.jd.com%2F"
	page := client.Load("login", url)
	page.Click("div.login-tab-r a", 0)
	page.SetValue("#loginname", "xeach")
	page.SetValue("#nloginpwd", "Startup.2013")
	page.Click("div.login-btn a", 0)
	imageMap := make(map[string]string)
	qryBigImage := "div.JDJRV-bigimg img"
	qrySmallImage := "div.JDJRV-smallimg img"
	page.ImageReady(qryBigImage, imageMap)
	page.ImageReady(qrySmallImage, imageMap)

	attrBig := NewAttrMap("width", "height")

	qryRefreshImage := "div.JDJRV-img-refresh"
	var resp map[string]any
	var err error
	count := 3
	for {
		if count <= 0 {
			break
		}
		resp, err = client.DecodeSlider2(imageMap[qrySmallImage][22:], imageMap[qryBigImage][22:])
		if err == nil && resp["isSuccess"].(bool) {
			width := ReadFloat(attrBig, "width")
			offsetX := (width / resp["w"].(float64)) * resp["x"].(float64)
			page.MouseDrag("div.JDJRV-slide-btn", offsetX)
			err = page.WaitVisible("a.nickname", 11)
		}
		if err == nil {
			break
		}
		page.Click(qryRefreshImage, 0)
		page.ImageChanged(qryBigImage, imageMap)
		page.ImageChanged(qrySmallImage, imageMap)
		count--
	}
}
