package goja

import (
	"go.uber.org/zap"
	"merkaba/chromedp"
	"merkaba/goja/unistring"
)

type webElementObject struct {
	baseObject
	m *chromedp.WebElement
}

func (e *webElementObject) getOwnPropStr(name unistring.String) Value {
	if name == "text" {
		return newStringValue(e.m.Text())
	} else if name == "html" {
		return newStringValue(e.m.Html())
	} else {
		return valueNull{}
	}
}

func (r *Runtime) CreateWebElementObject(e *chromedp.WebElement) *Object {
	o := &Object{runtime: r}
	wro := &webElementObject{m: e}
	wro.class = classWebElement
	wro.val = o
	wro.extensible = true
	o.self = wro
	wro.prototype = r.global.WebElementPrototype
	wro.init()
	return o
}

func (r *Runtime) ReadWebElementObject(call FunctionCall) (wpo *webElementObject, ok bool) {
	thisObj := r.toObject(call.This)
	mo, ok := thisObj.self.(*webElementObject)
	if !ok {
		panic(r.NewTypeError("Method WebElement.prototype.text called on incompatible receiver %s", r.objectproto_toString(FunctionCall{This: thisObj})))
	}
	return mo, ok
}
func (r *Runtime) builtin_newWebElement(args []Value, newTarget *Object) *Object {
	if newTarget == nil {
		panic(r.needNew("WebElement"))
	}
	proto := r.getPrototypeFromCtor(newTarget, r.global.WebElement, r.global.WebElementPrototype)
	o := &Object{runtime: r}
	mo := &webElementObject{}
	mo.class = classWebElement
	mo.val = o
	mo.extensible = true
	o.self = mo
	mo.prototype = proto
	mo.init()
	return o
}

func (r *Runtime) webElementProto_text(call FunctionCall) Value {
	mo, ok := r.ReadWebElementObject(call)
	if !ok {
		return valueNull{}
	}
	v := mo.m.Text()
	return nilSafe(r.ToValue(v))
}

func (r *Runtime) webElementProto_html(call FunctionCall) Value {
	mo, ok := r.ReadWebElementObject(call)
	if !ok {
		return valueNull{}
	}
	v := mo.m.Html()
	return nilSafe(r.ToValue(v))
}

func (r *Runtime) webElementProto_nodeId(call FunctionCall) Value {
	mo, ok := r.ReadWebElementObject(call)
	if !ok {
		return valueNull{}
	}
	v := mo.m.NodeId()
	return nilSafe(r.ToValue(v))
}

func (r *Runtime) webElementProto_value(call FunctionCall) Value {
	mo, ok := r.ReadWebElementObject(call)
	if !ok {
		return valueNull{}
	}
	v := mo.m.Value()
	return nilSafe(r.ToValue(v))
}

func (r *Runtime) webElementProto_shot(call FunctionCall) Value {
	mo, ok := r.ReadWebElementObject(call)
	if !ok {
		return valueNull{}
	}
	v := mo.m.ScreenShot()
	return nilSafe(r.ToValue(v))
}

func (r *Runtime) webElementProto_attrs(call FunctionCall) Value {
	mo, ok := r.ReadWebElementObject(call)
	if !ok {
		return valueNull{}
	}
	mapValue := mo.m.Attrs()
	o := r.CreateObject(r.global.ObjectPrototype)
	for k, v := range mapValue {
		o.self.(*baseObject)._put(unistring.String(k), r.ToValue(v))
	}
	return o
}

func (r *Runtime) webElementProto_selects(call FunctionCall) Value {
	mo, ok := r.ReadWebElementObject(call)
	if !ok {
		return nil
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
		return nil
	}
	values := make([]Value, len(elements))
	for i, e := range elements {
		values[i] = r.CreateWebElementObject(e)
	}
	return r.newArrayValues(values)
}

func (r *Runtime) webElementProto_select(call FunctionCall) Value {
	mo, ok := r.ReadWebElementObject(call)
	if !ok {
		return nil
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

func (r *Runtime) webElementProto_clickPage(call FunctionCall) Value {
	mo, ok := r.ReadWebElementObject(call)
	if !ok {
		return valueNull{}
	}
	page, err := mo.m.ClickPage()
	if err != nil {
		return valueNull{}
	}
	o := r.CreateWebPageObject(page)
	return o
}

func (r *Runtime) webElementProto_click(call FunctionCall) Value {
	mo, ok := r.ReadWebElementObject(call)
	if !ok {
		return valueNull{}
	}
	millisecond := int64(100)
	if len(call.Arguments) == 2 {
		millisecond = call.Argument(1).ToInteger()
	}
	mo.m.Click(millisecond)
	return valueNull{}
}

func (r *Runtime) webElementProto_sendKeys(call FunctionCall) Value {
	mo, ok := r.ReadWebElementObject(call)
	if !ok {
		return valueNull{}
	}
	mo.m.SendKeys(call.Argument(0).String())
	return valueNull{}
}

func (r *Runtime) webElementProto_setValue(call FunctionCall) Value {
	mo, ok := r.ReadWebElementObject(call)
	if !ok {
		return valueNull{}
	}
	mo.m.SetValue(call.Argument(0).String())
	return valueNull{}
}

func (r *Runtime) webElementProto_setAttr(call FunctionCall) Value {
	mo, ok := r.ReadWebElementObject(call)
	if !ok {
		return valueNull{}
	}
	v := call.Argument(1)
	mo.m.SetAttr(call.Argument(0).String(), v.Export())
	return valueNull{}
}

func (r *Runtime) webElementProto_show(call FunctionCall) Value {
	mo, ok := r.ReadWebElementObject(call)
	if !ok {
		return valueNull{}
	}
	mo.m.Show()
	return valueNull{}
}

func (r *Runtime) webElementProto_attr(call FunctionCall) Value {
	mo, ok := r.ReadWebElementObject(call)
	if !ok {
		return valueNull{}
	}
	v := mo.m.Attr(call.Argument(0).String())
	return nilSafe(r.ToValue(v))
}

func (r *Runtime) webElementProto_handleImage(name string, call FunctionCall) Value {
	mo, ok := r.ReadWebElementObject(call)
	if !ok {
		return nil
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

func (r *Runtime) webElementProto_imageReady(call FunctionCall) Value {
	return r.webElementProto_handleImage("imageReady", call)
}

func (r *Runtime) webElementProto_imageChanged(call FunctionCall) Value {
	return r.webElementProto_handleImage("imageChanged", call)
}

func (r *Runtime) webElementProto_waitContentChanged(call FunctionCall) Value {
	mo, ok := r.ReadWebElementObject(call)
	if !ok {
		return nil
	}
	var err error
	query := call.Argument(0).String()
	oldHtml := call.Argument(1).String()
	timeout := int64(0)
	if len(call.Arguments) >= 3 {
		timeout = call.Argument(2).ToInteger()
	}
	err = mo.m.WaitChanged(query, oldHtml, timeout)
	if err != nil {
		r.Context.Error(query, zap.Error(err))
	}
	return valueNull{}
}

func (r *Runtime) webElementProto_styles(call FunctionCall) Value {
	mo, ok := r.ReadWebElementObject(call)
	if !ok {
		return nil
	}
	o := r.CreateObject(r.global.ObjectPrototype)
	wro := o.self.(*baseObject)
	var err error
	mapValue := make(map[string]string)
	err = mo.m.Styles(mapValue)
	if err != nil {
		wro._put("error", r.ToValue(err))
	} else {
		for k, v := range mapValue {
			wro._put(unistring.String(k), r.ToValue(v))
		}
	}
	return o
}

func (r *Runtime) webElementProto_clientRect(call FunctionCall) Value {
	mo, ok := r.ReadWebElementObject(call)
	if !ok {
		return nil
	}
	o := r.CreateObject(r.global.ObjectPrototype)
	wro := o.self.(*baseObject)
	mapValue, err := mo.m.ClientRect()
	if err != nil {
		wro._put("error", r.ToValue(err))
	} else {
		for k, v := range mapValue {
			wro._put(unistring.String(k), r.ToValue(v))
		}
	}
	return o
}

func (r *Runtime) webElementProto_mouseDrag(call FunctionCall) Value {
	mo, ok := r.ReadWebElementObject(call)
	if !ok {
		return valueNull{}
	}
	mo.m.MouseDrag(call.Argument(0).String(), call.Argument(1).ToFloat())
	return nil
}

func (r *Runtime) webElementProto_mouseOver(call FunctionCall) Value {
	mo, ok := r.ReadWebElementObject(call)
	if !ok {
		return valueNull{}
	}
	mo.m.MouseOver(call.Argument(0).String())
	return nil
}

func (r *Runtime) webElementProto_scrollIntoView(call FunctionCall) Value {
	mo, ok := r.ReadWebElementObject(call)
	if !ok {
		return valueNull{}
	}
	mo.m.ScrollIntoView()
	return nil
}

func (r *Runtime) webElementProto_waitVisible(call FunctionCall) Value {
	mo, ok := r.ReadWebElementObject(call)
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

func (r *Runtime) createWebElementProto(val *Object) objectImpl {
	o := newBaseObjectObj(val, r.global.WebElementPrototype, classWebElement)
	o._putProp("setAttr", r.newNativeFunc(r.webElementProto_setAttr, nil, "setAttr", nil, 2), true, false, true)
	o._putProp("attr", r.newNativeFunc(r.webElementProto_attr, nil, "attr", nil, 2), true, false, true)
	o._putProp("shot", r.newNativeFunc(r.webElementProto_shot, nil, "shot", nil, 0), true, false, true)
	o._putProp("waitVisible", r.newNativeFunc(r.webElementProto_waitVisible, nil, "waitVisible", nil, 2), false, true, true)
	o._putProp("waitChanged", r.newNativeFunc(r.webElementProto_waitContentChanged, nil, "waitChanged", nil, 3), false, true, true)

	o._putProp("selects", r.newNativeFunc(r.webElementProto_selects, nil, "selects", nil, 2), true, false, true)
	o._putProp("select", r.newNativeFunc(r.webElementProto_select, nil, "select", nil, 2), true, false, true)

	o._putProp("show", r.newNativeFunc(r.webElementProto_show, nil, "show", nil, 0), false, true, true)
	o._putProp("clickPage", r.newNativeFunc(r.webElementProto_clickPage, nil, "clickPage", nil, 1), true, false, true)
	o._putProp("click", r.newNativeFunc(r.webElementProto_click, nil, "click", nil, 1), true, false, true)
	o._putProp("sendKeys", r.newNativeFunc(r.webElementProto_sendKeys, nil, "sendKeys", nil, 1), true, false, true)
	o._putProp("setValue", r.newNativeFunc(r.webElementProto_setValue, nil, "setValue", nil, 1), true, false, true)
	o._putProp("imageReady", r.newNativeFunc(r.webElementProto_imageReady, nil, "imageReady", nil, 2), true, false, true)
	o._putProp("imageChanged", r.newNativeFunc(r.webElementProto_imageChanged, nil, "imageChanged", nil, 2), true, false, true)
	o._putProp("mouseDrag", r.newNativeFunc(r.webElementProto_mouseDrag, nil, "mouseDrag", nil, 3), true, false, true)
	o._putProp("mouseOver", r.newNativeFunc(r.webElementProto_mouseOver, nil, "mouseOver", nil, 3), true, false, true)
	o._putProp("scrollIntoView", r.newNativeFunc(r.webElementProto_scrollIntoView, nil, "mouseWheel", nil, 3), true, false, true)
	o.values["text"] = &valueProperty{getterFunc: r.newNativeFunc(r.webElementProto_text, nil, "get text", nil, 0), accessor: true, writable: true, configurable: true}
	o.values["html"] = &valueProperty{getterFunc: r.newNativeFunc(r.webElementProto_html, nil, "get html", nil, 0), accessor: true, writable: true, configurable: true}
	o.values["nodeId"] = &valueProperty{getterFunc: r.newNativeFunc(r.webElementProto_nodeId, nil, "get nodeId", nil, 0), accessor: true, writable: true, configurable: true}
	o.values["value"] = &valueProperty{getterFunc: r.newNativeFunc(r.webElementProto_value, nil, "get value", nil, 0), accessor: true, writable: true, configurable: true}
	o.values["attrs"] = &valueProperty{getterFunc: r.newNativeFunc(r.webElementProto_attrs, nil, "get attrs", nil, 0), accessor: true, writable: true, configurable: true}
	o.values["styles"] = &valueProperty{getterFunc: r.newNativeFunc(r.webElementProto_styles, nil, "get styles", nil, 0), accessor: true, writable: true, configurable: true}
	o.values["clientRect"] = &valueProperty{getterFunc: r.newNativeFunc(r.webElementProto_clientRect, nil, "get clientRect", nil, 0), accessor: true, writable: true, configurable: true}
	return o
}

func (r *Runtime) createWebElement(val *Object) objectImpl {
	o := r.newNativeConstructOnly(val, r.builtin_newWebElement, r.global.WebElementPrototype, "WebResponse", 0)
	r.putSpeciesReturnThis(o)
	return o
}

func (r *Runtime) initWebElement() {
	r.global.WebElementPrototype = r.newLazyObject(r.createWebElementProto)
	r.global.WebElement = r.newLazyObject(r.createWebElement)
	r.addToGlobal("WebElement", r.global.WebElement)
}
