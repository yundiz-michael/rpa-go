package chromedp

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"image"
	"merkaba/chromedp/cdproto/browser"
	"merkaba/chromedp/cdproto/cdp"
	"merkaba/chromedp/cdproto/dom"
	"merkaba/chromedp/cdproto/network"
	"merkaba/chromedp/cdproto/page"
	"merkaba/chromedp/cdproto/runtime"
	"merkaba/chromedp/cdproto/target"
	"merkaba/common"
	"os"
	"strconv"
	"time"
)

type WebPage struct {
	runCtx     *common.RunContext
	ID         string
	Url        string
	Client     *WebClient
	Ctx        context.Context
	loggerStd  *zap.Logger
	loggerTask *zap.Logger
	cancel     func()
}

func NewPage(id string, client *WebClient) *WebPage {
	result := &WebPage{
		ID:     id,
		Client: client,
	}
	result.runCtx = client.Context
	result.Ctx, result.cancel = NewContext(client.AllocCtx)
	result.registerListener(result.Ctx)
	loggerCore := common.LoggerCoreBy(result.runCtx.TaskName)
	result.loggerStd = zap.New(common.LoggerStdCore, zap.AddCaller(), zap.AddCallerSkip(1), zap.Development())
	result.loggerTask = zap.New(loggerCore, zap.AddCaller(), zap.AddCallerSkip(1), zap.Development())

	client.AddPage(result)
	result.init()
	return result
}

func (p *WebPage) init() {
	ListenTarget(p.Ctx, func(ev interface{}) {
		switch ev := ev.(type) {
		case *page.EventFrameNavigated:
			p.Url = ev.Frame.URL
		}
	})
	ListenBrowser(p.Ctx, func(ev interface{}) {
		switch ev := ev.(type) {
		case *target.EventAttachedToTarget:
			p.Client.AddTabSession(ev.SessionID.String())
		case *target.EventDetachedFromTarget:
			p.Client.RemoveTabSession(ev.SessionID.String())
		}
	})
}

func (p *WebPage) CreateExecuteContext() (ctx context.Context, err error) {
	c, err := InitContextBrowser(p.Ctx)
	if err != nil {
		p.runCtx.Error("createExecuteContext", zap.NamedError("error", err))
		return nil, err
	}
	return cdp.WithExecutor(p.Ctx, c.Target), nil
}

func (p *WebPage) scriptHandler() common.ScriptHandler {
	return p.Client.ScriptHandler
}

func (p *WebPage) maxWaitTime() int64 {
	return p.runCtx.MaxWaitTime
}

func (p *WebPage) sendLog(level string, msg string, fields ...zap.Field) {
	p.Client.sendLog(level, msg, fields...)
}

func (p *WebPage) Logger() *zap.Logger {
	if p.runCtx.RunMode == 0 {
		return p.loggerStd
	} else {
		return p.loggerTask
	}
}

func (p *WebPage) info(msg string, fields ...zap.Field) {
	p.Logger().Info(msg, fields...)
}

func (p *WebPage) printError(msg string, err error) {
	p.Logger().Error(msg, zap.NamedError("error", err))
	p.sendLog("error", msg, zap.String("error", err.Error()))
}

func (p *WebPage) parseQueryOption(sel interface{}, all bool) QueryOption {
	return p.Client.parseQueryOption(sel, all)
}

func (p *WebPage) Load(url string) (err error) {
	p.Url = url
	p.registerListener(p.Ctx)
	err = Run(p.Ctx, p.loadCookies(), Navigate(url),
		WaitReady("body", ByQuery, SetRunParam(p.scriptHandler(), p, p.maxWaitTime())))
	if err != nil {
		p.printError("load:"+url, err)
	} else {
		p.info("load", zap.String("url", url))
	}
	p.RefreshVNC()
	return err
}

func (p *WebPage) loadCookies() Action {
	return ActionFunc(func(ctx context.Context) error {
		sql := `
            select content from merkaba_cookie where ip=? and siteName=? and account=?
	   `
		var content []string
		common.Mysql.Select(&content, sql, common.LocalIP, p.runCtx.SiteName, p.runCtx.Account())
		if len(content) == 0 {
			return nil
		}
		cookieMap := make(map[string]any)
		json.Unmarshal([]byte(content[0]), &cookieMap)
		for name := range cookieMap {
			if name == "context" {
				continue
			}
			cm := cookieMap[name].(map[string]any)
			name := cm["name"].(string)
			value := cm["value"].(string)
			domain_ := cm["domain"].(string)
			expires := cm["expires"].(float64)
			var expr cdp.TimeSinceEpoch
			if expires < 0 {
				expr = cdp.TimeSinceEpoch(time.Now().Add(180 * 24 * time.Hour))
			} else {
				expr = cdp.TimeSinceEpoch(time.Unix(int64(expires), int64(1000)))
			}
			err := network.SetCookie(name, value).
				WithExpires(&expr).
				WithDomain(domain_).
				WithHTTPOnly(true).
				Do(ctx)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (p *WebPage) Close() {
	p.Cancel()
	for i := 0; i < len(p.Client.Pages); i++ {
		if p.Client.Pages[i] == p {
			p.Client.removePage(i)
			break
		}
	}
}

func (p *WebPage) Cancel() {
	p.cancel()
}

func (p *WebPage) BASE64(bytes []byte) string {
	return base64.StdEncoding.EncodeToString(bytes)
}

func (p *WebPage) Selects(sel interface{}, queryTimeout int64, mustVisible bool) (elements []*WebElement, err error) {
	p.info("select:" + sel.(string))
	var nodes []*cdp.Node
	options := make([]QueryOption, 0)

	options = append(options, p.Client.parseQueryOption(sel, true))
	if queryTimeout > 0 {
		options = append(options, SetScrollIntoView(), ActionTimeout(queryTimeout))
	}
	options = append(options, SetRunParam(p.scriptHandler(), p, p.maxWaitTime()))
	if mustVisible {
		err = Run(p.Ctx, VisibleNodes(sel, nil, &nodes, options...))
	} else {
		err = Run(p.Ctx, Nodes(sel, &nodes, options...))
	}
	if err == nil && nodes == nil {
		err = errors.New("not exist")
	}
	if err != nil {
		p.printError(sel.(string), err)
		return nil, err
	} else {
		p.info(sel.(string), zap.Int("count", len(nodes)))
	}
	elements = make([]*WebElement, len(nodes))
	for i, node := range nodes {
		elements[i] = NewWebElement(i, p, node)
	}
	return elements, err
}

func (p *WebPage) Select(sel interface{}, queryTimeout int64, mustVisible bool) (element *WebElement, err error) {
	elements, err := p.Selects(sel, queryTimeout, mustVisible)
	if err != nil || elements == nil {
		return nil, err
	}
	return elements[0], nil
}

func (p *WebPage) WaitVisible(sel interface{}, second int64) (err error) {
	if sel == nil {
		msg := "WaitVisible: select can't be null"
		p.printError(msg, err)
		return errors.New(msg)
	}
	p.info("waitVisible", zap.String("sel", sel.(string)), zap.Int64("second", second))
	err = Run(p.Ctx,
		WaitVisible(sel, p.Client.parseQueryOption(sel, false), SetScrollIntoView(), SetRunParam(p.scriptHandler(), p, p.maxWaitTime()), ActionTimeout(second)),
	)
	if err != nil {
		msg := "WaitVisible:" + sel.(string)
		p.printError(msg, err)
	} else {
		p.RefreshVNC()
	}
	return err
}

func (p *WebPage) WaitNotVisible(sel interface{}, second int64) (err error) {
	p.info("WaitNotVisible", zap.String("sel", sel.(string)), zap.Int64("second", second))
	err = Run(p.Ctx,
		WaitNotVisible(sel, p.Client.parseQueryOption(sel, false), SetScrollIntoView(), SetRunParam(p.scriptHandler(), p, p.maxWaitTime()), ActionTimeout(second)),
	)
	if err != nil {
		msg := "WaitNotVisible:" + sel.(string)
		p.printError(msg, err)
	} else {
		p.RefreshVNC()
	}
	return err
}

func (p *WebPage) WaitNotPresent(sel interface{}, second int64) (err error) {
	p.info("WaitNotPresent", zap.String("sel", sel.(string)), zap.Int64("second", second))
	err = Run(p.Ctx,
		WaitNotPresent(sel, p.Client.parseQueryOption(sel, false), SetScrollIntoView(), SetRunParam(p.scriptHandler(), p, p.maxWaitTime()), ActionTimeout(second)),
	)
	if err != nil {
		msg := "WaitNotPresent:" + sel.(string)
		p.printError(msg, err)
	} else {
		p.RefreshVNC()
	}
	return err
}

func (p *WebPage) WaitMoreThan(sel interface{}, count int64) (err error) {
	p.info("waitMoreThan", zap.String("sel", sel.(string)), zap.Int64("count", count))
	err = Run(p.Ctx,
		WaitMoreThan(sel, int(count), p.Client.parseQueryOption(sel, true), SetScrollIntoView(), SetRunParam(p.scriptHandler(), p, p.maxWaitTime())),
	)
	if err != nil {
		msg := "WaitMoreThan:" + sel.(string) + ">=" + strconv.Itoa(int(count))
		p.printError(msg, err)
	} else {
		p.RefreshVNC()
	}
	return err
}

func (p *WebPage) Text(sel interface{}, text *string) (err error) {
	p.info("readText", zap.String("sel", sel.(string)))
	err = Run(p.Ctx,
		Text(sel, text, p.Client.parseQueryOption(sel, false), SetRunParam(p.scriptHandler(), p, p.maxWaitTime())),
	)
	if err != nil {
		msg := "readText:" + sel.(string)
		p.printError(msg, err)
	} else {
		p.RefreshVNC()
	}
	return err
}

func (p *WebPage) Value(sel interface{}, value *string) (err error) {
	p.info("readText", zap.String("sel", sel.(string)))
	err = Run(p.Ctx,
		Value(sel, value, p.Client.parseQueryOption(sel, false), SetRunParam(p.scriptHandler(), p, p.maxWaitTime())),
	)
	if err != nil {
		msg := "readValue:" + sel.(string)
		p.printError(msg, err)
	} else {
		p.RefreshVNC()
	}
	return err
}

func (p *WebPage) SetValue(sel interface{}, value string) (err error) {
	p.info("setValue", zap.String("sel", sel.(string)), zap.String("value", value))
	err = Run(p.Ctx,
		SetValue(sel, value, p.Client.parseQueryOption(sel, false), SetRunParam(p.scriptHandler(), p, p.maxWaitTime())),
	)
	if err != nil {
		msg := "setValue:" + sel.(string)
		p.printError(msg, err)
	} else {
		p.RefreshVNC()
	}
	return err
}

func (p *WebPage) Styles(sel interface{}, attrMap map[string]string) (err error) {
	err = Run(p.Ctx,
		ReadAttributes(sel, attrMap, p.Client.parseQueryOption(sel, false), SetRunParam(p.scriptHandler(), p, p.maxWaitTime())),
	)
	if err != nil {
		msg := "Styles:" + sel.(string)
		p.printError(msg, err)
	}
	return err
}

func (p *WebPage) WaitContentChanged(sel interface{}, oldHtml string) (err error) {
	p.info("WaitChanged", zap.String("sel", sel.(string)))
	err = Run(p.Ctx,
		WaitContentChanged(sel, oldHtml, p.Client.parseQueryOption(sel, false), SetScrollIntoView(), SetRunParam(p.scriptHandler(), p, p.maxWaitTime())),
	)
	if err != nil {
		msg := "WaitChanged:" + sel.(string)
		p.printError(msg, err)
	} else {
		p.RefreshVNC()
	}
	return err
}

func (p *WebPage) ImageReady(sel interface{}, mapValue map[string]string) (err error) {
	p.info("ImageReady", zap.String("sel", sel.(string)))
	err = Run(p.Ctx,
		WaitImageSrcReady(sel, mapValue, p.Client.parseQueryOption(sel, false), SetScrollIntoView(), SetRunParam(p.scriptHandler(), p, p.maxWaitTime())),
	)
	if err != nil {
		msg := "imageReady:" + sel.(string)
		p.printError(msg, err)
	} else {
		p.RefreshVNC()
	}
	return err
}

func (p *WebPage) ImageChanged(sel interface{}, values map[string]string) (err error) {
	p.info("imageChange", zap.String("sel", sel.(string)))
	err = Run(p.Ctx,
		WaitImageSrcChanged(sel, values, p.Client.parseQueryOption(sel, false), SetScrollIntoView(), SetRunParam(p.scriptHandler(), p, p.maxWaitTime())),
	)
	if err != nil {
		msg := "ImageChanged:" + sel.(string)
		p.printError(msg, err)
	} else {
		p.RefreshVNC()
	}
	return err
}

func (p *WebPage) ScreenShot() (res string, err error) {
	bytes, e := p.WriteScreenShot()
	if e != nil {
		msg := "ScreenShot"
		p.printError(msg, e)
		return "", e
	} else {
		return p.BASE64(*bytes), nil
	}
}

func (p *WebPage) WriteScreenShot() (*[]byte, error) {
	var bytes []byte
	err := Run(p.Ctx,
		FullScreenshot(&bytes, 80))
	if err != nil {
		p.printError("WriteScreenShot", err)
	}
	return &bytes, err
}

func (p *WebPage) SendKeys(sel interface{}, value string) (err error) {
	p.info("sendKeys", zap.String("sel", sel.(string)), zap.String("value", value))
	err = Run(p.Ctx,
		SendKeys(sel, value, p.Client.parseQueryOption(sel, false), SetScrollIntoView(), SetRunParam(p.scriptHandler(), p, p.maxWaitTime())),
	)
	if err != nil {
		msg := "sendKeys:" + sel.(string)
		p.printError(msg, err)
	} else {
		p.RefreshVNC()
	}
	return err
}

func (p *WebPage) ClickDown(sel interface{}, maxSecond int64) (fileName string, err error) {
	p.info("Click " + sel.(string) + ",start download file")
	done := make(chan string, 1)
	ListenTarget(p.Ctx, func(v interface{}) {
		if ev, ok := v.(*browser.EventDownloadProgress); ok {
			if ev.State == browser.DownloadProgressStateCompleted {
				p.info("download finished")
				done <- ev.GUID
			} else if ev.State == browser.DownloadProgressStateInProgress {
				p.info("download file", zap.Float64("received", ev.ReceivedBytes), zap.Float64("total", ev.TotalBytes))
			}
		}
	})
	queryOption := p.Client.parseQueryOption(sel, false)
	e := Run(p.Ctx,
		browser.SetDownloadBehavior(browser.SetDownloadBehaviorBehaviorAllowAndName).
			WithDownloadPath(common.Env.Path.Temp).
			WithEventsEnabled(true),
		Click(sel, queryOption),
	)
	if e != nil {
		p.printError(sel.(string), err)
		return "", e
	}
	var guid string
	select {
	case <-time.After(time.Duration(maxSecond) * time.Second):
		e = errors.New("time out")
		p.printError(sel.(string), e)
		return "", e
	case guid = <-done:
		return common.Env.Path.Temp + guid, nil
	}
}

func (p *WebPage) Upload(sel interface{}, fileName string) (err error) {
	p.info("upload", zap.String("sel", sel.(string)), zap.String("fileName", fileName))
	err = Run(p.Ctx,
		SetUploadFiles(sel, []string{fileName}, p.Client.parseQueryOption(sel, false), SetScrollIntoView(), SetRunParam(p.scriptHandler(), p, p.maxWaitTime())),
	)
	if err != nil {
		msg := "upload:" + sel.(string)
		p.printError(msg, err)
	} else {
		p.RefreshVNC()
	}
	return err
}

func (p *WebPage) Wait(v int64) {
	p.info("upload", zap.Int64("time", v))
	time.Sleep(time.Duration(v) * time.Millisecond)
}

func (p *WebPage) Click(sel interface{}, millisecond int64) (err error) {
	p.info("Click", zap.String("sel", sel.(string)))
	queryOption := p.Client.parseQueryOption(sel, false)
	e := Run(p.Ctx,
		WaitVisible(sel, queryOption, SetScrollIntoView(), SetRunParam(p.scriptHandler(), p, p.maxWaitTime())),
		Click(sel, queryOption),
	)
	time.Sleep(time.Duration(millisecond) * time.Millisecond)
	if e != nil {
		msg := "click:" + sel.(string)
		p.printError(msg, err)
	} else {
		p.RefreshVNC()
	}
	return e
}

func (p *WebPage) ScrollIntoView(sel interface{}) (err error) {
	p.info("ScrollIntoView", zap.String("sel", sel.(string)))
	queryOption := p.Client.parseQueryOption(sel, false)
	err = Run(p.Ctx,
		ScrollIntoView(sel, queryOption, SetRunParam(p.scriptHandler(), p, p.maxWaitTime())),
	)
	if err != nil {
		msg := "ScrollIntoView:" + sel.(string)
		p.printError(msg, err)
	} else {
		p.RefreshVNC()
	}
	return err
}

func (p *WebPage) MouseDrag(sel interface{}, offsetX float64) (err error) {
	p.info("MouseDrag", zap.String("sel", sel.(string)), zap.Float64("offsetX", offsetX))
	err = Run(p.Ctx,
		MouseDrag(sel, offsetX, p.Client.parseQueryOption(sel, false), SetScrollIntoView(), SetRunParam(p.scriptHandler(), p, p.maxWaitTime())),
	)
	if err != nil {
		msg := "mouseDrag:" + sel.(string)
		p.printError(msg, err)
	} else {
		p.RefreshVNC()
	}
	return err
}

func (p *WebPage) MouseOver(sel interface{}) (err error) {
	p.info("mouseOver", zap.String("sel", sel.(string)))
	err = Run(p.Ctx,
		MouseOver(sel, p.Client.parseQueryOption(sel, false), SetScrollIntoView(), SetRunParam(p.scriptHandler(), p, p.maxWaitTime())),
	)
	if err != nil {
		msg := "mouseOver:" + sel.(string)
		p.printError(msg, err)
	} else {
		p.RefreshVNC()
	}
	return err
}

func (p *WebPage) _createMewPage(id string, newCtx context.Context, newCancel context.CancelFunc) (newPage *WebPage, err error) {
	newPage = &WebPage{
		Client: p.Client,
		Ctx:    newCtx,
		cancel: newCancel,
		ID:     id,
		runCtx: p.runCtx,
	}
	err = Run(newCtx,
		WaitReady("body", ByQuery, SetRunParam(p.scriptHandler(), newPage, p.maxWaitTime())),
	)
	if err != nil {
		return nil, err
	}
	p.Client.AddPage(newPage)
	return newPage, nil
}

func (p *WebPage) waitNewPage() (newPage *WebPage, err error) {
	var newUrl string
	ch := WaitNewTarget(p.Ctx, func(info *target.Info) bool {
		if info.URL != "" {
			newUrl = info.URL
			p.RefreshVNC()
			return true
		} else {
			return false
		}
	})
	newCtx, newCancel := NewContext(p.Ctx, WithTargetID(<-ch))
	p.registerListener(newCtx)
	_page, err := p._createMewPage(primitive.NewObjectID().Hex(), newCtx, newCancel)
	_page.Url = newUrl
	return _page, err
}

func (p *WebPage) InjectScript(script string) error {
	err := Run(p.Ctx,
		ActionFunc(func(ctx context.Context) error {
			_, exp, err := runtime.Evaluate(script).Do(ctx)
			if err != nil {
				p.runCtx.Error(script, zap.Error(err))
				return err
			}
			if exp != nil {
				p.runCtx.Error(script, zap.Error(exp))
				return exp
			}
			return nil
		}),
	)
	if err != nil {
		return err
	} else {
		p.RefreshVNC()
		return nil
	}
}

func (p *WebPage) ReadScriptVar(name string) (value string, err error) {
	ctx, err := p.CreateExecuteContext()
	if err != nil {
		return "", err
	}
	t := cdp.ExecutorFromContext(ctx).(*Target)
	_, _, execCtx, _ := t.ensureFrame()

	script := "function reaScriptVar() { return JSON.stringify(" + name + ",(key, value) =>{\n    if(value){\n       return value\n     } else {\n       return undefined\n   }\n});}"
	var res string
	err = CallFunctionOn(script, &res, func(p *runtime.CallFunctionOnParams) *runtime.CallFunctionOnParams {
		return p.WithExecutionContextID(execCtx)
	}).Do(ctx)
	if err != nil {
		msg := "ReadVariable " + name
		p.runCtx.Error(msg, zap.Error(err))
		return "", err
	}
	return res, nil
}

func (p *WebPage) registerListener(ctx context.Context) {
	ListenBrowser(ctx, func(ev interface{}) {
		switch ev := ev.(type) {
		case *runtime.EventExceptionThrown:
			s := ev.ExceptionDetails.Error()
			p.printError("Chrome Error", errors.New(s))
		case *page.EventJavascriptDialogOpening:
			Run(ctx, page.HandleJavaScriptDialog(true))
		}
	})
}

// Open  打开新页面
func (p *WebPage) Open(id string, url string) (newPage *WebPage, err error) {
	p.info("Open", zap.String("url", url))
	_page := p.Client.PageBy(id)
	if _page == nil {
		newCtx, newCancel := NewContext(p.Client.AllocCtx)
		defer newCancel()
		p.registerListener(newCtx)
		err = Run(newCtx, Navigate(url))
		if err != nil {
			return nil, err
		}
		_page, err = p._createMewPage(id, newCtx, newCancel)
		_page.Url = url
	} else {
		err = _page.Load(url)
	}
	if err != nil {
		msg := "Open:" + url
		p.printError(msg, err)
	} else {
		p.RefreshVNC()
	}
	return _page, err
}

func (p *WebPage) DownImage(url string) (content string, err error) {
	var newPage *WebPage
	newPage, err = p.Open("DownImage", url)
	if err != nil {
		return "", err
	}
	var bytes []byte
	err = Run(newPage.Ctx,
		FullScreenshot(&bytes, 80),
	)
	if err != nil {
		msg := "Download " + url
		p.printError(msg, err)
		newPage.Close()
		return "", err
	} else {
		newPage.RefreshVNC()
	}
	newPage.Close()
	content = p.BASE64(bytes)
	return content, nil
}

func (p *WebPage) RefreshVNC() {
	bytes_, err := p.WriteScreenShot()
	if err != nil {
		return
	}
	img, _, err := image.Decode(bytes.NewReader(*bytes_))
	if err != nil {
		return
	}
	common.RefreshVNC(p.Client.Context.TaskName, img)
}

func (p *WebPage) Html() string {
	var ids []cdp.NodeID
	var result string
	err := Run(p.Ctx,
		NodeIDs(`document`, &ids, ByJSPath),
		ActionFunc(func(ctx context.Context) error {
			var err error
			result, err = dom.GetOuterHTML().WithNodeID(ids[0]).Do(ctx)
			return err
		}),
	)
	if err != nil {
		return err.Error()
	} else {
		return result
	}
}

func (p *WebPage) SaveHtml(fileName string) {
	_fileName := common.RootPath + "tmp/" + fileName
	p.info("save html", zap.String("fileName", _fileName))
	content := p.Html()
	os.WriteFile(_fileName, []byte(content), 0644)
}
