package goja

import (
	"go.uber.org/zap"
	"merkaba/chromedp"
	"merkaba/goja/unistring"
)

type webPageObject struct {
	baseObject
	m *chromedp.WebPage
}

func (mo *webPageObject) init(args []Value) {
	mo.baseObject.init()
	obj := args[1].ToObject(mo.val.runtime)
	client := obj.self.(*webClientObject).m
	mo.m = chromedp.NewPage(args[0].String(), client)
}

func (r *Runtime) CreateWebPageObject(e *chromedp.WebPage) *Object {
	o := &Object{runtime: r}
	wro := &webPageObject{m: e}
	wro.class = classWebPage
	wro.val = o
	wro.extensible = true
	o.self = wro
	wro.prototype = r.global.WebPagePrototype
	wro.baseObject.init()
	return o
}

func (r *Runtime) CreateErrorResponse(err string) *Object {
	o := r.CreateObject(r.global.ObjectPrototype)
	o.self.(*baseObject)._put("isSuccess", r.ToValue(false))
	o.self.(*baseObject)._put("error", r.ToValue(err))
	return o
}

func (r *Runtime) CreateSuccessResponse() *Object {
	o := r.CreateObject(r.global.ObjectPrototype)
	o.self.(*baseObject)._put("isSuccess", r.ToValue(true))
	return o
}

func (r *Runtime) checkResponse(err error) *Object {
	if err != nil {
		return r.CreateErrorResponse(err.Error())
	} else {
		return r.CreateSuccessResponse()
	}
}

func (r *Runtime) readWebPageObject(call FunctionCall) (wpo *webPageObject, ok bool) {
	thisObj := r.toObject(call.This)
	mo, ok := thisObj.self.(*webPageObject)
	if !ok {
		panic(r.NewTypeError("Method WebPage.prototype called on incompatible receiver %s", r.objectproto_toString(FunctionCall{This: thisObj})))
	}
	return mo, ok
}

func (r *Runtime) webPageProto_load(call FunctionCall) Value {
	mo, ok := r.readWebPageObject(call)
	if !ok {
		return valueNull{}
	}
	err := mo.m.Load(call.Argument(0).String())
	return r.checkResponse(err)
}

func (r *Runtime) webPageProto_url(call FunctionCall) Value {
	mo, ok := r.readWebPageObject(call)
	if !ok {
		return valueNull{}
	}
	return r.ToValue(mo.m.Url)
}

func (r *Runtime) webPageProto_html(call FunctionCall) Value {
	mo, ok := r.readWebPageObject(call)
	if !ok {
		return valueNull{}
	}
	return r.ToValue(mo.m.Html())
}

func (r *Runtime) webPageProto_waitVisible(call FunctionCall) Value {
	mo, ok := r.readWebPageObject(call)
	if !ok {
		return valueNull{}
	}
	second := int64(0)
	if len(call.Arguments) == 2 {
		second = call.Argument(1).ToInteger()
	}
	err := mo.m.WaitVisible(call.Argument(0).String(), second)
	return r.ToValue(err == nil)
}

func (r *Runtime) webPageProto_waitNotVisible(call FunctionCall) Value {
	mo, ok := r.readWebPageObject(call)
	if !ok {
		return valueNull{}
	}
	second := int64(0)
	if len(call.Arguments) == 2 {
		second = call.Argument(1).ToInteger()
	}
	err := mo.m.WaitNotVisible(call.Argument(0).String(), second)
	return r.ToValue(err == nil)
}

func (r *Runtime) webPageProto_waitNotPresent(call FunctionCall) Value {
	mo, ok := r.readWebPageObject(call)
	if !ok {
		return valueNull{}
	}
	second := int64(0)
	if len(call.Arguments) == 2 {
		second = call.Argument(1).ToInteger()
	}
	err := mo.m.WaitNotPresent(call.Argument(0).String(), second)
	return r.ToValue(err == nil)
}

func (r *Runtime) webPageProto_waitMoreThan(call FunctionCall) Value {
	mo, ok := r.readWebPageObject(call)
	if !ok {
		return valueNull{}
	}
	count := call.Argument(1).ToInteger()
	err := mo.m.WaitMoreThan(call.Argument(0).String(), count)
	return r.ToValue(err == nil)
}

func (r *Runtime) webPageProto_readValue(name string, call FunctionCall) Value {
	mo, ok := r.readWebPageObject(call)
	if !ok {
		return valueNull{}
	}
	result := ""
	var err error
	if name == "value" {
		err = mo.m.Value(call.Argument(0).String(), &result)
	} else if name == "text" {
		err = mo.m.Text(call.Argument(0).String(), &result)
	} else {
		return valueNull{}
	}
	if err != nil {
		return valueNull{}
	} else {
		return r.ToValue(result)
	}
}

func (r *Runtime) webPageProto_text(call FunctionCall) Value {
	return r.webPageProto_readValue("text", call)
}

func (r *Runtime) webPageProto_value(call FunctionCall) Value {
	return r.webPageProto_readValue("value", call)
}

func (r *Runtime) webPageProto_setValue(call FunctionCall) Value {
	return r.webPageProto_handleKeyValue("setValue", call)
}

func (r *Runtime) webPageProto_sendKeys(call FunctionCall) Value {
	return r.webPageProto_handleKeyValue("sendKeys", call)
}

func (r *Runtime) webPageProto_click(call FunctionCall) Value {
	return r.webPageProto_handleAction("click", call)
}

func (r *Runtime) webPageProto_mouseDrag(call FunctionCall) Value {
	return r.webPageProto_handleAction("mouseDrag", call)
}

func (r *Runtime) webPageProto_mouseOver(call FunctionCall) Value {
	return r.webPageProto_handleAction("mouseOver", call)
}

func (r *Runtime) webPageProto_scrollIntoView(call FunctionCall) Value {
	return r.webPageProto_handleAction("scrollIntoView", call)
}

func (r *Runtime) webPageProto_injectScript(call FunctionCall) Value {
	mo, ok := r.readWebPageObject(call)
	if !ok {
		return valueNull{}
	}
	mo.m.InjectScript(call.Argument(0).String())
	return valueNull{}
}

func (r *Runtime) webPageProto_readScriptVar(call FunctionCall) Value {
	mo, ok := r.readWebPageObject(call)
	if !ok {
		return valueNull{}
	}
	content, err := mo.m.ReadScriptVar(call.Argument(0).String())
	if err != nil {
		return valueNull{}
	} else {
		return r.ToValue(content)
	}
}

func (r *Runtime) webPageProto_handleAction(name string, call FunctionCall) Value {
	mo, ok := r.readWebPageObject(call)
	if !ok {
		return nil
	}
	var err error
	if name == "click" {
		millisecond := int64(100)
		if len(call.Arguments) == 2 {
			millisecond = call.Argument(1).ToInteger()
		}
		err = mo.m.Click(call.Argument(0).String(), millisecond)
	} else if name == "mouseDrag" {
		err = mo.m.MouseDrag(call.Argument(0).String(), call.Argument(1).ToFloat())
	} else if name == "mouseOver" {
		err = mo.m.MouseOver(call.Argument(0).String())
	} else if name == "scrollIntoView" {
		err = mo.m.ScrollIntoView(call.Argument(0).String())
	} else {
		err = nil
	}
	if err != nil {
		return nil
	} else {
		return nil
	}
}

func (r *Runtime) webPageProto_handleKeyValue(name string, call FunctionCall) Value {
	mo, ok := r.readWebPageObject(call)
	if !ok {
		return nil
	}
	var err error
	if name == "setValue" {
		err = mo.m.SetValue(call.Argument(0).String(), call.Argument(1).String())
	} else if name == "sendKeys" {
		err = mo.m.SendKeys(call.Argument(0).String(), call.Argument(1).String())
	}
	if err != nil {
		return nil
	} else {
		return nil
	}
}

func (r *Runtime) webPageProto_selects(call FunctionCall) Value {
	mo, ok := r.readWebPageObject(call)
	if !ok {
		return valueNull{}
	}
	queryTimeout := int64(0)
	if len(call.Arguments) == 2 {
		queryTimeout = call.Argument(1).ToInteger()
	}
	mustVisible := false
	if len(call.Arguments) == 3 {
		mustVisible = call.Argument(2).ToBoolean()
	}
	elements, err := mo.m.Selects(call.Argument(0).String(), queryTimeout, mustVisible)
	if err != nil || elements == nil {
		return valueNull{}
	}
	values := make([]Value, len(elements))
	for i, e := range elements {
		values[i] = r.CreateWebElementObject(e)
	}
	return r.newArrayValues(values)
}

func (r *Runtime) webPageProto_select(call FunctionCall) Value {
	mo, ok := r.readWebPageObject(call)
	if !ok {
		return valueNull{}
	}
	queryTimeout := int64(0)
	if len(call.Arguments) == 2 {
		queryTimeout = call.Argument(1).ToInteger()
	}
	mustVisible := false
	if len(call.Arguments) == 3 {
		mustVisible = call.Argument(2).ToBoolean()
	}
	element, err := mo.m.Select(call.Argument(0).String(), queryTimeout, mustVisible)
	if err != nil || element == nil {
		return valueNull{}
	}
	o := r.CreateWebElementObject(element)
	return o
}

func (r *Runtime) webPageProto_styles(call FunctionCall) Value {
	mo, ok := r.readWebPageObject(call)
	if !ok {
		return valueNull{}
	}
	o := r.CreateObject(r.global.ObjectPrototype)
	wro := o.self.(*baseObject)
	var err error
	query := call.Argument(0).String()
	mapValue := make(map[string]string)
	names := call.Argument(1).ToObject(r).self.(*arrayObject)
	for _, name := range names.values {
		mapValue[name.String()] = ""
	}
	err = mo.m.Styles(query, mapValue)
	if err != nil {
		r.Context.Error(query, zap.Error(err))
		wro._put("isSuccess", r.ToValue(true))
		wro._put("error", r.ToValue(err))
	} else {
		wro._put("isSuccess", r.ToValue(true))
		for k, v := range mapValue {
			wro._put(unistring.String(k), r.ToValue(v))
		}
	}
	return o
}

func (r *Runtime) webPageProto_handleImage(name string, call FunctionCall) Value {
	mo, ok := r.readWebPageObject(call)
	if !ok {
		return valueNull{}
	}
	mapValue := make(map[string]string)
	var err error
	query := call.Argument(0).String()

	if name == "imageReady" {
		err = mo.m.ImageReady(query, mapValue)
	} else if name == "imageChanged" {
		mapValue[query] = call.Argument(1).String()
		err = mo.m.ImageChanged(query, mapValue)
	}
	if err != nil {
		r.Context.Error(query, zap.Error(err))
		return r.ToValue("")
	} else {
		return r.ToValue(mapValue[query])
	}
}

func (r *Runtime) webPageProto_imageReady(call FunctionCall) Value {
	return r.webPageProto_handleImage("imageReady", call)
}

func (r *Runtime) webPageProto_imageChanged(call FunctionCall) Value {
	return r.webPageProto_handleImage("imageChanged", call)
}

func (r *Runtime) builtin_newWebPage(args []Value, newTarget *Object) *Object {
	if newTarget == nil {
		panic(r.needNew("WebPage"))
	}
	proto := r.getPrototypeFromCtor(newTarget, r.global.WebPage, r.global.WebPagePrototype)
	o := &Object{runtime: r}
	mo := &webPageObject{}
	mo.class = classWebPage
	mo.val = o
	mo.extensible = true
	o.self = mo
	mo.prototype = proto
	mo.init(args)
	return o
}

func (r *Runtime) webPageProto_open(call FunctionCall) Value {
	mo, ok := r.readWebPageObject(call)
	if !ok {
		return valueNull{}
	}
	newPage, err := mo.m.Open(call.Argument(0).String(), call.Argument(1).String())
	if err != nil {
		return valueNull{}
	}
	o := r.CreateWebPageObject(newPage)
	return o
}

func (r *Runtime) webPageProto_close(call FunctionCall) Value {
	mo, ok := r.readWebPageObject(call)
	if !ok {
		return valueNull{}
	}
	mo.m.Close()
	return nil
}
func (r *Runtime) webPageProto_upload(call FunctionCall) Value {
	mo, ok := r.readWebPageObject(call)
	if !ok {
		return nil
	}
	mo.m.Upload(call.Argument(0).String(), call.Argument(1).String())
	return valueNull{}
}

func (r *Runtime) webPageProto_saveHtml(call FunctionCall) Value {
	mo, ok := r.readWebPageObject(call)
	if !ok {
		return nil
	}
	mo.m.SaveHtml(call.Argument(0).String())
	return valueNull{}
}

func (r *Runtime) webPageProto_clickDown(call FunctionCall) Value {
	mo, ok := r.readWebPageObject(call)
	if !ok {
		return nil
	}
	second := int64(20)
	if len(call.Arguments) == 2 {
		second = call.Argument(1).ToInteger()
	}
	fileName, err := mo.m.ClickDown(call.Argument(0).String(), second)
	if err != nil {
		return valueNull{}
	} else {
		return r.ToValue(fileName)
	}
}

func (r *Runtime) webPageProto_downImage(call FunctionCall) Value {
	mo, ok := r.readWebPageObject(call)
	if !ok {
		return nil
	}
	content, err := mo.m.DownImage(call.Argument(0).String())
	if err != nil {
		return valueNull{}
	}
	return r.ToValue(content)
}

func (r *Runtime) webPageProto_client(call FunctionCall) Value {
	mo, ok := r.readWebPageObject(call)
	if !ok {
		return valueNull{}
	}
	return r.CreateWebClientObject(mo.m.Client)
}

func (r *Runtime) webPageProto_wait(call FunctionCall) Value {
	mo, ok := r.readWebPageObject(call)
	if !ok {
		return valueNull{}
	}
	mo.m.Wait(call.Argument(0).ToInteger())
	return valueNull{}
}

func (r *Runtime) createWebPageProto(val *Object) objectImpl {
	o := newBaseObjectObj(val, r.global.WebPagePrototype, classWebPage)
	o._putProp("constructor", r.global.WebPage, true, false, true)
	o._putProp("styles", r.newNativeFunc(r.webPageProto_styles, nil, "styles", nil, 3), true, false, true)
	o._putProp("load", r.newNativeFunc(r.webPageProto_load, nil, "load", nil, 1), true, false, true)
	o._putProp("open", r.newNativeFunc(r.webPageProto_open, nil, "load", nil, 1), true, false, true)
	o._putProp("close", r.newNativeFunc(r.webPageProto_close, nil, "load", nil, 0), true, false, true)
	o._putProp("waitVisible", r.newNativeFunc(r.webPageProto_waitVisible, nil, "waitVisible", nil, 2), false, true, true)
	o._putProp("waitNotVisible", r.newNativeFunc(r.webPageProto_waitNotVisible, nil, "waitVisible", nil, 2), false, true, true)
	o._putProp("waitNotPresent", r.newNativeFunc(r.webPageProto_waitNotPresent, nil, "waitNotPresent", nil, 2), false, true, true)
	o._putProp("waitMoreThan", r.newNativeFunc(r.webPageProto_waitMoreThan, nil, "waitMoreThan", nil, 2), false, true, true)
	o._putProp("text", r.newNativeFunc(r.webPageProto_text, nil, "text", nil, 1), true, false, true)
	o._putProp("value", r.newNativeFunc(r.webPageProto_value, nil, "value", nil, 1), true, false, true)
	o._putProp("setValue", r.newNativeFunc(r.webPageProto_setValue, nil, "setValue", nil, 3), true, false, true)
	o._putProp("sendKeys", r.newNativeFunc(r.webPageProto_sendKeys, nil, "sendKeys", nil, 3), true, false, true)
	o._putProp("imageReady", r.newNativeFunc(r.webPageProto_imageReady, nil, "imageReady", nil, 3), true, false, true)
	o._putProp("imageChanged", r.newNativeFunc(r.webPageProto_imageChanged, nil, "imageChanged", nil, 3), true, false, true)
	o._putProp("click", r.newNativeFunc(r.webPageProto_click, nil, "click", nil, 3), true, false, true)
	o._putProp("mouseDrag", r.newNativeFunc(r.webPageProto_mouseDrag, nil, "mouseDrag", nil, 3), true, false, true)
	o._putProp("mouseOver", r.newNativeFunc(r.webPageProto_mouseOver, nil, "mouseOver", nil, 3), true, false, true)
	o._putProp("scrollIntoView", r.newNativeFunc(r.webPageProto_scrollIntoView, nil, "scrollIntoView", nil, 1), true, false, true)
	o._putProp("injectScript", r.newNativeFunc(r.webPageProto_injectScript, nil, "injectScript", nil, 1), true, false, true)
	o._putProp("readScriptVar", r.newNativeFunc(r.webPageProto_readScriptVar, nil, "readScriptVar", nil, 1), true, false, true)

	o._putProp("selects", r.newNativeFunc(r.webPageProto_selects, nil, "selects", nil, 2), true, false, true)
	o._putProp("select", r.newNativeFunc(r.webPageProto_select, nil, "select", nil, 2), true, false, true)

	o._putProp("upload", r.newNativeFunc(r.webPageProto_upload, nil, "upload", nil, 2), true, false, true)
	o._putProp("downImage", r.newNativeFunc(r.webPageProto_downImage, nil, "downImage", nil, 2), true, false, true)
	o._putProp("clickDown", r.newNativeFunc(r.webPageProto_clickDown, nil, "clickDown", nil, 2), true, false, true)
	o._putProp("saveHtml", r.newNativeFunc(r.webPageProto_saveHtml, nil, "saveHtml", nil, 2), true, false, true)
	o._putProp("wait", r.newNativeFunc(r.webPageProto_wait, nil, "wait", nil, 1), true, false, true)
	o.values["client"] = &valueProperty{getterFunc: r.newNativeFunc(r.webPageProto_client, nil, "get client", nil, 0), accessor: true, writable: true, configurable: true}
	o.values["url"] = &valueProperty{getterFunc: r.newNativeFunc(r.webPageProto_url, nil, "get url", nil, 0), accessor: true, writable: true, configurable: true}
	o.values["html"] = &valueProperty{getterFunc: r.newNativeFunc(r.webPageProto_html, nil, "get html", nil, 0), accessor: true, writable: true, configurable: true}
	return o
}

func (r *Runtime) createWebPage(val *Object) objectImpl {
	o := r.newNativeConstructOnly(val, r.builtin_newWebPage, r.global.WebPagePrototype, "WebPage", 0)
	r.putSpeciesReturnThis(o)
	return o
}

func (r *Runtime) initWebPage() {
	r.global.WebPagePrototype = r.newLazyObject(r.createWebPageProto)
	r.global.WebPage = r.newLazyObject(r.createWebPage)
	r.addToGlobal("WebPage", r.global.WebPage)
}
