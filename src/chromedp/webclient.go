package chromedp

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"merkaba/chromedp/cdproto/cdp"
	"merkaba/chromedp/cdproto/network"
	"merkaba/common"
	"os"
	"strconv"
	"time"
)

type WebClientOption struct { //record ,class
	Domain   string
	Proxy    bool
	Headless bool
}

type WebCookie struct {
	Path    string  `json:"path"`
	Name    string  `json:"name"`
	Domain  string  `json:"domain"`
	Value   string  `json:"value"`
	Expires float64 `json:"expires"`
}

type WebClient struct {
	Option           WebClientOption
	Pages            []*WebPage
	SessionIds       []string
	CurrentPageIndex int
	Context          *common.RunContext
	ScriptHandler    common.ScriptHandler
	AllocCtx         context.Context
	loggerStd        *zap.Logger
	loggerTask       *zap.Logger
	cancel           func()
}

func NewAttrMap(names ...string) map[string]string {
	result := make(map[string]string)
	for _, name := range names {
		result[name] = ""
	}
	return result
}

func ReadFloat(attrMap map[string]string, name string) float64 {
	value := attrMap[name]
	result, _ := strconv.ParseFloat(value[0:len(value)-2], 64)
	return result
}

func _initOptions(option WebClientOption) []ExecAllocatorOption {
	edpOptions := append(
		DefaultExecAllocatorOptions[3:],
		NoFirstRun,
		NoDefaultBrowserCheck,
	)
	if option.Proxy {
		edpOptions = append(edpOptions, ProxyServer("todo proxy pool"))
	}
	if option.Headless {
		edpOptions = append(edpOptions, Flag("headless", true))
	} else {
		/*2k分辨率*/
		edpOptions = append(edpOptions, WindowSize(2560, 1600))
	}
	return edpOptions
}

func NewClient(option WebClientOption, ctx *common.RunContext) *WebClient {
	result := &WebClient{
		Option:  option,
		Context: ctx,
	}
	loggerCore := common.LoggerCoreBy(ctx.TaskName)
	result.loggerStd = zap.New(common.LoggerStdCore, zap.AddCaller(), zap.AddCallerSkip(2), zap.Development())
	result.loggerTask = zap.New(loggerCore, zap.AddCaller(), zap.AddCallerSkip(2), zap.Development())
	result.AllocCtx, result.cancel = NewExecAllocator(context.Background(), _initOptions(option)...)
	return result
}

func (client *WebClient) Logger() *zap.Logger {
	if client.Context.RunMode == 0 {
		return client.loggerStd
	} else {
		return client.loggerTask
	}
}

func (client *WebClient) Load(pageId string, url string) *WebPage {
	page := client.PageBy(pageId)
	if page == nil {
		page = NewPage(pageId, client)
	}
	page.Load(url)
	return page
}

func (client *WebClient) sendLog(level string, msg string, fields ...zap.Field) {
	if client.Context.RunMode != common.RunModeBrowserRun {
		return
	}
	log := make(map[string]any)
	log["message"] = msg
	data := make(map[string]any)
	log["data"] = data
	for _, field := range fields {
		log[field.Key] = field.Interface
	}
	common.SendLog(level, client.Context, log)
}

func (client *WebClient) Info(msg string, fields ...zap.Field) {
	client.Logger().Info(msg, fields...)
	client.sendLog("info", msg, fields...)
}

func (client *WebClient) Error(msg string, fields ...zap.Field) {
	client.Logger().Error(msg, fields...)
	client.sendLog("error", msg, fields...)
}

func (client *WebClient) parseQueryOption(sel interface{}, all bool) QueryOption {
	queryMode := QueryMode(sel.(string))
	switch queryMode {
	case "xpath":
		return BySearch
	default:
		if all {
			return ByQueryAll
		} else {
			return ByQuery
		}
	}
}

func (client *WebClient) PageBy(pageId string) *WebPage {
	index := -1
	var page *WebPage
	var i int
	for i, page = range client.Pages {
		if page.ID == pageId {
			index = i
			break
		}
	}
	if index >= 0 {
		page = client.Pages[index]
	} else {
		page = nil
	}
	return page
}

func (client *WebClient) SetCookies(cookies ...string) error {
	expr := cdp.TimeSinceEpoch(time.Now().Add(180 * 24 * time.Hour))
	for i := 0; i < len(cookies); i += 2 {
		for _, page := range client.Pages {
			err := network.SetCookie(cookies[i], cookies[i+1]).
				WithExpires(&expr).
				WithDomain(client.Option.Domain).
				WithHTTPOnly(true).
				Do(page.Ctx)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (client *WebClient) Close() {
	if len(client.Pages) == 0 {
		/*如果页面为空，标示已经关闭了*/
		return
	}
	for _, page := range client.Pages {
		page.Cancel()
	}
	client.Pages = make([]*WebPage, 0)
	client.cancel()
	client.Context.OnClientClose(client)
}

func (client *WebClient) AddPage(page *WebPage) {
	client.Pages = append(client.Pages, page)
	client.CurrentPageIndex = len(client.Pages) - 1
}

func (client *WebClient) removePage(index int) {
	client.Pages = append(client.Pages[:index], client.Pages[index+1:]...)
	client.CurrentPageIndex = len(client.Pages) - 1
}

func (client *WebClient) AddTabSession(id string) {
	client.SessionIds = append(client.SessionIds, id)
}

func (client *WebClient) RemoveTabSession(id string) {
	for i := 0; i < len(client.SessionIds); i++ {
		if client.SessionIds[i] == id {
			client.SessionIds = append(client.SessionIds[:i], client.SessionIds[i+1:]...)
			i--
		}
	}
	if len(client.SessionIds) == 0 && client.Context.OnClientClose != nil {
		client.Context.OnClientClose(client)
	}
}

func (client *WebClient) DecodeSlider2(background string, template string) (resp map[string]any, err error) {
	param := make(map[string]any)
	param["background"] = background
	param["template"] = template
	param["debug"] = true
	resp, err = common.DecodeSlider(2, param)
	return resp, err
}

func (client *WebClient) DecodeSlider1(background string, templateWidth int64) (resp map[string]any, err error) {
	param := make(map[string]any)
	param["background"] = background
	param["templateWidth"] = templateWidth
	param["debug"] = true
	resp, err = common.DecodeSlider(1, param)
	return resp, err
}

func (client *WebClient) Ocr(image string) (resp map[string]any, err error) {
	param := make(map[string]any)
	param["image"] = image
	param["debug"] = false
	resp, err = common.HttpOCRImage(param)
	return resp, err
}

func (client *WebClient) RefreshVNC() {
	if client.CurrentPageIndex >= len(client.Pages) {
		return
	}
	page := client.Pages[client.CurrentPageIndex]
	page.RefreshVNC()
}

func (client *WebClient) Snapshot(fileName string) {
	if client.CurrentPageIndex >= len(client.Pages) {
		return
	}
	page := client.Pages[client.CurrentPageIndex]
	bytes, err := page.WriteScreenShot()
	if err != nil {
		return
	}
	err = os.WriteFile(fileName, *bytes, 0644)
	if err == nil {
		client.Info("write snapshot ", zap.String("fileName", fileName))
	}
}

func (client *WebClient) SaveCookies() {
	if len(client.Pages) == 0 {
		return
	}
	cookieMap := make(map[string]any)
	task := Tasks{
		network.Enable(),
		ActionFunc(func(ctx context.Context) error {
			cookies, err := network.GetAllCookies().Do(ctx)
			if err != nil {
				client.Error("SaveCookies ", zap.Error(err))
				return nil
			}
			for _, cookie := range cookies {
				cookieMap[cookie.Name] = map[string]any{
					"path":    cookie.Path,
					"name":    cookie.Name,
					"domain":  cookie.Domain,
					"value":   cookie.Value,
					"expires": cookie.Expires,
				}
			}
			return err
		}),
	}
	err := Run(client.Pages[0].Ctx, task)
	if err == nil {
		data, _ := json.Marshal(cookieMap)
		sql := `
            insert into merkaba_cookie(ip,siteName,account,content,lastModified) values(?,?,?,?,?)
            on duplicate key update 
        	content = values(content),lastModified = values(lastModified)                         
	   `
		var _, err = common.Mysql.Update(sql, common.LocalIP, client.Context.SiteName, client.Context.Account(),
			string(data), time.Now().UnixMilli())
		if err != nil {
			common.LoggerStd.Error(sql, zap.Error(err))
			return
		}
	}
}
