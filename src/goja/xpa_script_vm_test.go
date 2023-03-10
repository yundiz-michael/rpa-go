package goja

import (
	"fmt"
	"merkaba/common"
	"testing"
)

func TestJDScript(t *testing.T) {
	vm := New()
	common.InitEnviroment()
	printer := PrinterFunc(func(level string, s string) {
		fmt.Println("===>" + s)
	})
	InitConsole(vm, printer)
	_, err := vm.RunScript("3.js",
		`
		var client=new WebClient('jd.com');
		var url="https://passport.jd.com/new/login.aspx?ReturnUrl=https%3A%2F%2Fglobal.jd.com%2F";
		var page=client.load(url);
		page.click("div.login-tab-r a");
		page.setValue("#loginname", "xeach");
		page.setValue("#nloginpwd", "Startup.2013");
		page.click("div.login-btn a");
		/*处理登录图片*/
		var images = {}
		var qryBigImage = "div.JDJRV-bigimg img"
		var qrySmallImage = "div.JDJRV-smallimg img"
		var qryRefreshImage = "div.JDJRV-img-refresh"
		images[qryBigImage]=page.imageReady(qryBigImage)
		images[qrySmallImage]=page.imageReady(qrySmallImage)
		attrBig=page.readAttrs(qryBigImage, ["width","height"])
		var count = 3;
		while (count>=0) {
			var resp =client.decodeJDLogin(
				images[qrySmallImage].substring(22), 
				images[qryBigImage].substring(22))
			if (resp.isSuccess){
				var width =parseFloat(attrBig["width"]);
				var offsetX = (width / parseFloat(resp["w"])) * parseFloat(resp["x"]);
				page.mouseDrag("div.JDJRV-slide-btn", offsetX);
				resp= page.waitVisible("a.nickname")
			}    
			if (resp.isSuccess) break;
			page.click(qryRefreshImage)
			images[qryBigImage]=page.imageChanged(qryBigImage)
			images[qrySmallImage]=page.imageChanged(qrySmallImage)
			count--
		}
`)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	value := vm.Get("html")
	fmt.Println(value.String())
	text := vm.Get("text")
	fmt.Println(text.String())

}
