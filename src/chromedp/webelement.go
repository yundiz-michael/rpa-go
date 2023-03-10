package chromedp

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"merkaba/chromedp/cdproto/cdp"
	"merkaba/chromedp/cdproto/css"
	"merkaba/chromedp/cdproto/dom"
	"merkaba/chromedp/cdproto/page"
	"merkaba/common"
	"strconv"
	"time"
)

type WebElement struct {
	ctx        *common.RunContext
	Index      int
	Node       *cdp.Node
	Page       *WebPage
	loggerStd  *zap.Logger
	loggerTask *zap.Logger
}

func NewWebElement(index int, page *WebPage, node *cdp.Node) *WebElement {
	result := &WebElement{
		Index: index,
		Page:  page,
		Node:  node,
	}
	result.ctx = page.Client.Context
	loggerCore := common.LoggerCoreBy(result.ctx.TaskName)
	result.loggerStd = zap.New(common.LoggerStdCore, zap.AddCaller(), zap.AddCallerSkip(1), zap.Development())
	result.loggerTask = zap.New(loggerCore, zap.AddCaller(), zap.AddCallerSkip(1), zap.Development())
	return result
}

func (e *WebElement) Logger() *zap.Logger {
	if e.ctx.RunMode == 0 {
		return e.loggerStd
	} else {
		return e.loggerTask
	}
}

func (e *WebElement) sendLog(level string, msg string, fields ...zap.Field) {
	e.Page.Client.sendLog(level, msg, fields...)
}

func (e *WebElement) info(msg string, fields ...zap.Field) {
	e.Logger().Info(msg, fields...)
}

func (e *WebElement) printError(msg string, err error) {
	e.Logger().Error(msg, zap.NamedError("error", err))
	e.sendLog("error", msg, zap.String("error", err.Error()))
}

func (e *WebElement) parseQueryOption(sel interface{}, all bool) QueryOption {
	return e.Page.Client.parseQueryOption(sel, all)
}

func (e *WebElement) _createExecuteContext(page *WebPage) (ctx context.Context, err error) {
	return page.CreateExecuteContext()
}

func (e *WebElement) RunReadScript(name string) any {
	ctx, err := e._createExecuteContext(e.Page)
	if err != nil {
		return ""
	}
	var res any
	if name == "html" {
		err = CallFunctionOnNode(ctx,
			e.Node, attributeJS, &res, "outerHTML")
	} else if name == "text" {
		err = CallFunctionOnNode(ctx, e.Node, textJS, &res)
	} else if name == "value" {
		err = CallFunctionOnNode(ctx, e.Node, attributeJS, &res, "value")
	} else {
		err = CallFunctionOnNode(ctx, e.Node, attributeJS, &res, name)
	}
	if err != nil {
		msg := "element " + name
		e.printError(msg, err)
		return ""
	}
	return res

}

func (e *WebElement) Html() string {
	e.info(e.Node.NodeName + ".Html")
	return e.RunReadScript("html").(string)
}

func (e *WebElement) NodeId() string {
	if e.Node == nil {
		return ""
	} else {
		return strconv.FormatInt(int64(e.Node.NodeID), 10)
	}
}

func (e *WebElement) Text() string {
	e.info(e.Node.NodeName + ".Text")
	return e.RunReadScript("text").(string)
}

func (e *WebElement) Value() string {
	e.info(e.Node.NodeName + ".Value")
	return e.RunReadScript("value").(string)
}

func (e *WebElement) ScreenShot() string {
	e.info(e.Node.NodeName + ".ScreenShot")
	ctx, err := e._createExecuteContext(e.Page)
	if err != nil {
		return ""
	}
	dom.ScrollIntoViewIfNeeded().WithNodeID(e.Node.NodeID).Do(ctx)
	var bytes []byte
	ElementScreenshot(ctx, &bytes, e.Node)
	return e.Page.BASE64(bytes)
}

func (e *WebElement) Attrs() map[string]string {
	e.info(e.Node.NodeName + ".Attrs")
	e.Node.RLock()
	defer e.Node.RUnlock()
	result := make(map[string]string)
	attrs := e.Node.Attributes
	for i := 0; i < len(attrs); i += 2 {
		result[attrs[i]] = attrs[i+1]
	}
	return result
}

func (e *WebElement) ClientRect() (map[string]any, error) {
	ctx, _ := e._createExecuteContext(e.Page)
	result := make(map[string]any)
	var clip page.Viewport
	if err := CallFunctionOnNode(ctx, e.Node, getClientRectJS, &clip); err != nil {
		e.printError("ClientRect", err)
		return result, err
	}
	result["y"] = clip.Y
	result["x"] = clip.X
	result["width"] = clip.Width
	result["height"] = clip.Height
	return result, nil
}

func (e *WebElement) Styles(attrMap map[string]string) error {
	ctx, err := e._createExecuteContext(e.Page)
	if err != nil {
		return err
	}
	computed, err := css.GetComputedStyleForNode(e.Node.NodeID).Do(ctx)
	if err != nil {
		e.printError("Styles", err)
		return err
	}
	for _, prop := range computed {
		name := prop.Name
		attrMap[name] = prop.Value
	}
	return nil
}

func (e *WebElement) Selects(sel interface{}, queryTimeout int64, mustVisible bool) (elements []*WebElement, err error) {
	e.info(e.Node.NodeName, zap.String("select", sel.(string)))
	var nodes []*cdp.Node
	options := make([]QueryOption, 0)
	options = append(options, e.parseQueryOption(sel, true))
	options = append(options, FromNode(e.Node))
	if queryTimeout > 0 {
		options = append(options, ActionTimeout(queryTimeout))
	}
	options = append(options, SetRunParam(e.Page.Client.ScriptHandler, e.Page, e.Page.maxWaitTime()))
	if mustVisible {
		err = Run(e.Page.Ctx, VisibleNodes(sel, e.Node, &nodes, options...))
	} else {
		err = Run(e.Page.Ctx, Nodes(sel, &nodes, options...))
	}
	if err == nil && nodes == nil {
		err = errors.New("not exist")
	}
	if err != nil {
		e.printError(sel.(string), err)
		return nil, err
	} else {
		e.info(sel.(string), zap.Int("count", len(nodes)))
	}
	elements = make([]*WebElement, len(nodes))
	for i, node := range nodes {
		elements[i] = NewWebElement(i, e.Page, node)
	}
	return elements, err
}

func (e *WebElement) Select(sel interface{}, queryTimeout int64, mustVisible bool) (element *WebElement, err error) {
	elements, err := e.Selects(sel, queryTimeout, mustVisible)
	if err != nil || elements == nil {
		return nil, err
	}
	return elements[0], nil
}

func (e *WebElement) Click(millisecond int64) error {
	e.info(e.Node.NodeName + " click")
	err := Run(e.Page.Ctx,
		MouseClickNode(e.Node),
	)
	time.Sleep(time.Duration(millisecond) * time.Millisecond)
	if err != nil {
		msg := e.Node.NodeName + " click"
		e.printError(msg, err)
	} else {
		e.Page.RefreshVNC()
	}
	return err
}

func (e *WebElement) ClickPage() (newPage *WebPage, err error) {
	msg := e.Node.NodeName + ".clickPage"
	e.info(msg)
	err = Run(e.Page.Ctx,
		MouseClickNode(e.Node),
	)
	newPage, err = e.Page.waitNewPage()
	if err != nil {
		e.printError(msg, err)
	} else {
		e.info("clickPage", zap.String("url", newPage.Url))
		e.Page.RefreshVNC()
	}
	return newPage, err
}

func (e *WebElement) SendKeys(value string) (err error) {
	e.info("sendKeys", zap.String("value", value))
	err = Run(e.Page.Ctx,
		SendElementKeys(e.Node, value),
	)
	if err != nil {
		msg := "sendKeys:" + value
		e.printError(msg, err)
	} else {
		e.Page.RefreshVNC()
	}
	return err
}

func (e *WebElement) SetValue(value string) (err error) {
	e.info("setValue", zap.String("value", value))
	err = Run(e.Page.Ctx,
		SetElementJavascriptAttribute(e.Node, "value", value),
	)
	if err != nil {
		msg := "SetValue:" + value
		e.printError(msg, err)
	} else {
		e.Page.RefreshVNC()
	}
	return err
}

func (e *WebElement) SetHtml(value string) (err error) {
	e.info("SetHtml", zap.String("html", value))
	_err := Run(e.Page.Ctx,
		SetElementJavascriptAttribute(e.Node, "innerHTML", value),
	)
	if _err != nil {
		msg := "SetHtml:" + value
		e.printError(msg, err)
	} else {
		e.Page.RefreshVNC()
	}
	return _err
}

func (e *WebElement) Attr(name string) any {
	e.info(e.Node.NodeName + "." + name)
	return e.RunReadScript(name)
}

func (e *WebElement) SetAttr(key string, value any) (err error) {
	e.info("attr", zap.String("key", key), zap.String("value", fmt.Sprintf("%v", value)))
	err = Run(e.Page.Ctx,
		SetElementJavascriptAttribute(e.Node, key, value),
	)
	if err != nil {
		msg := "Set:" + key + ":" + fmt.Sprintf("%v", value)
		e.printError(msg, err)
	} else {
		e.Page.RefreshVNC()
	}
	return err
}

func (e *WebElement) Show() (err error) {
	err = Run(e.Page.Ctx,
		ShowNode(e.Node),
	)
	if err != nil {
		msg := "Show"
		e.printError(msg, err)
	} else {
		e.Page.RefreshVNC()
	}
	return err
}

func (e *WebElement) WaitVisible(sel interface{}, second int64) (err error) {
	e.info("waitVisible", zap.String("sel", sel.(string)), zap.Int64("second", second))
	p := e.Page
	err = Run(p.Ctx,
		WaitVisible(sel, p.Client.parseQueryOption(sel, false),
			FromNode(e.Node),
			SetRunParam(p.scriptHandler(), p, p.maxWaitTime()),
			SetScrollIntoView(),
			ActionTimeout(second)),
	)
	if err != nil {
		msg := "WaitVisible:" + sel.(string)
		e.printError(msg, err)
	} else {
		e.Page.RefreshVNC()
	}
	return err
}

func (e *WebElement) WaitChanged(sel interface{}, oldHtml string, timeout int64) (err error) {
	e.info("WaitChanged", zap.String("sel", sel.(string)), zap.Int64("timeout", timeout))
	if timeout == 0 {
		timeout = e.Page.maxWaitTime()
	}
	err = Run(e.Page.Ctx,
		WaitContentChanged(sel, oldHtml, FromNode(e.Node),
			e.Page.Client.parseQueryOption(sel, false),
			SetScrollIntoView(),
			ActionTimeout(timeout),
			SetRunParam(e.Page.scriptHandler(), e.Page, e.Page.maxWaitTime())),
	)
	if err != nil {
		msg := "WaitChanged:" + sel.(string)
		e.printError(msg, err)
	} else {
		e.Page.RefreshVNC()
	}
	return err
}

func (e *WebElement) ImageReady(sel interface{}, mapValue map[string]string) (err error) {
	e.info("imageReady", zap.String("sel", sel.(string)))
	err = Run(e.Page.Ctx,
		WaitImageSrcReady(sel, mapValue, FromNode(e.Node), ByQuery, SetRunParam(e.Page.Client.ScriptHandler, e.Page, e.Page.maxWaitTime())),
	)
	if err != nil {
		msg := "imageReady:" + sel.(string)
		e.printError(msg, err)
	} else {
		e.Page.RefreshVNC()
	}
	return err
}

func (e *WebElement) ImageChanged(sel interface{}, oldValues map[string]string) (err error) {
	e.info("imageChange", zap.String("sel", sel.(string)))
	err = Run(e.Page.Ctx,
		WaitImageSrcChanged(sel, oldValues, FromNode(e.Node), ByQuery, SetRunParam(e.Page.Client.ScriptHandler, e.Page, e.Page.maxWaitTime())),
	)
	if err != nil {
		msg := "imageChange:" + sel.(string)
		e.printError(msg, err)
	} else {
		e.Page.RefreshVNC()
	}
	return err
}

func (e *WebElement) MouseDrag(sel interface{}, offsetX float64) (err error) {
	e.info("mouseDrag", zap.String("sel", sel.(string)), zap.Float64("offsetX", offsetX))
	err = Run(e.Page.Ctx,
		MouseDrag(sel, offsetX, FromNode(e.Node), ByQuery, SetRunParam(e.Page.Client.ScriptHandler, e.Page, e.Page.maxWaitTime())),
	)
	if err != nil {
		msg := "mouseDrag:" + sel.(string)
		e.printError(msg, err)
	} else {
		e.Page.RefreshVNC()
	}
	return err
}

func (e *WebElement) MouseOver(sel interface{}) (err error) {
	e.info("mouseOver", zap.String("sel", sel.(string)))
	err = Run(e.Page.Ctx,
		MouseOver(sel, FromNode(e.Node), ByQuery, SetRunParam(e.Page.Client.ScriptHandler, e.Page, e.Page.maxWaitTime())),
	)
	if err != nil {
		msg := "mouseOver:" + sel.(string)
		e.printError(msg, err)
	} else {
		e.Page.RefreshVNC()
	}
	return err
}

func (e *WebElement) ScrollIntoView() (err error) {
	e.info("ScrollIntoView")
	ctx, _ := e._createExecuteContext(e.Page)
	if err := dom.ScrollIntoViewIfNeeded().WithNodeID(e.Node.NodeID).Do(ctx); err != nil {
		e.printError("ScrollIntoView", err)
		return err
	} else {
		e.Page.RefreshVNC()
	}
	return nil
}
