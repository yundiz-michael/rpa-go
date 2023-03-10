package goja

import (
	"merkaba/chromedp"
	"merkaba/common"
	"merkaba/goja/unistring"
)

type webClientObject struct {
	baseObject
	m *chromedp.WebClient
}

func (mo *webClientObject) init(r *Runtime, args []Value) {
	mo.baseObject.init()
	domain := args[0].String()
	webClientMap := WebClients[domain]
	if webClientMap == nil {
		webClientMap = make(map[string]*chromedp.WebClient)
		WebClients[domain] = webClientMap
	}
	webClient := webClientMap[r.Context.TaskName]
	if webClient != nil {
		mo.m = webClient
		r.WebClient = webClient
	} else {
		option := chromedp.WebClientOption{Domain: domain, Proxy: r.Context.Proxy, Headless: r.Context.Headless}
		webClient = chromedp.NewClient(option, r.Context)
		webClientMap[r.Context.TaskName] = webClient
		r.WebClient = webClient
		mo.m = webClient
	}
}

func (r *Runtime) CreateWebClientObject(e *chromedp.WebClient) *Object {
	o := &Object{runtime: r}
	wro := &webClientObject{m: e}
	wro.class = classWebClient
	wro.val = o
	wro.extensible = true
	o.self = wro
	wro.prototype = r.global.WebClientPrototype
	wro.baseObject.init()
	return o
}

func (r *Runtime) fromMap(resp map[string]any) *Object {
	o := r.CreateObject(r.global.ObjectPrototype)
	wro := o.self.(*baseObject)
	wro._put("isSuccess", r.ToValue(true))
	for k, v := range resp {
		wro._put(unistring.String(k), r.ToValue(v))
	}
	return o
}

func (r *Runtime) readWebClientObject(call FunctionCall) (wpo *webClientObject, ok bool) {
	thisObj := r.toObject(call.This)
	mo, ok := thisObj.self.(*webClientObject)
	if !ok {
		panic(r.NewTypeError("Method WebClient.prototype called on incompatible receiver %s", r.objectproto_toString(FunctionCall{This: thisObj})))
	}
	return mo, ok
}

func (r *Runtime) webClientProto_load(call FunctionCall) Value {
	mo, ok := r.readWebClientObject(call)
	if !ok {
		return valueNull{}
	}
	var page *chromedp.WebPage
	page = mo.m.Load(call.Argument(0).String(), call.Argument(1).String())
	return r.CreateWebPageObject(page)
}

func (r *Runtime) webClientProto_saveCookies(call FunctionCall) Value {
	mo, ok := r.readWebClientObject(call)
	if !ok {
		return nil
	}
	mo.m.SaveCookies()
	return nil
}

func (r *Runtime) webClientProto_pageBy(call FunctionCall) Value {
	mo, ok := r.readWebClientObject(call)
	if !ok {
		return valueNull{}
	}
	var page *chromedp.WebPage
	page = mo.m.PageBy(call.Argument(0).String())
	if page == nil {
		return nil
	} else {
		return r.CreateWebPageObject(page)
	}
}

func (r *Runtime) webClientProto_close(call FunctionCall) Value {
	mo, ok := r.readWebClientObject(call)
	if !ok {
		return valueNull{}
	}
	mo.m.Close()
	return nil
}

func (r *Runtime) webClientProto_decodeSlider1(call FunctionCall) Value {
	mo, ok := r.readWebClientObject(call)
	if !ok {
		return valueNull{}
	}
	background := call.Argument(0).String()
	templateWidth := call.Argument(1).ToInteger()
	resp, err := mo.m.DecodeSlider1(background, templateWidth)
	var o *Object
	if err != nil {
		return r.CreateErrorResponse(err.Error())
	} else {
		o = r.fromMap(resp)
	}
	return o
}

func (r *Runtime) webClientProto_decodeSlider2(call FunctionCall) Value {
	mo, ok := r.readWebClientObject(call)
	if !ok {
		return valueNull{}
	}
	background := call.Argument(0).String()
	template := call.Argument(1).String()
	resp, err := mo.m.DecodeSlider2(background, template)
	var o *Object
	if err != nil {
		return r.CreateErrorResponse(err.Error())
	} else {
		o = r.fromMap(resp)
	}
	return o
}

func (r *Runtime) webClientProto_ocr(call FunctionCall) Value {
	mo, ok := r.readWebClientObject(call)
	if !ok {
		return valueNull{}
	}
	image := call.Argument(0).String()
	resp, err := mo.m.Ocr(image)
	var o *Object
	if err != nil {
		return r.CreateErrorResponse(err.Error())
	} else {
		o = r.fromMap(resp)
	}
	return o
}

func (r *Runtime) convertValue(v Value) interface{} {
	if v.ExportType() == nil {
		r.Context.Error("Data is automatically converted to empty str,because data is undefined")
		return ""
	}
	typeName := v.ExportType().Kind().String()
	switch typeName {
	case "int64":
		return v.ToInteger()
	case "bool":
		return v.ToBoolean()
	case "slice":
		object := v.ToObject(r)
		values := object.self.(*arrayObject).values
		results := make([]interface{}, 0)
		for _, v1 := range values {
			if v1.ExportType() == nil {
				continue
			}
			results = append(results, r.convertValue(v1))
		}
		return results
	case "map":
		object := v.ToObject(r)
		result := make(map[string]any)
		for k, v2 := range object.self.(*baseObject).values {
			if v2.ExportType() == nil {
				continue
			}
			result[k.String()] = r.convertValue(v2)
		}
		return result
	default:
		return v.String()
	}
}

func (r *Runtime) convertToMap(obj *baseObject, result map[string]any) {
	for k, v := range obj.values {
		key := k.String()
		result[key] = r.convertValue(v)
	}
}

func (r *Runtime) webClientProto_httpPost(call FunctionCall) Value {
	uri := call.Argument(0).String()
	obj := call.Argument(1).ToObject(r)

	mapValue := make(map[string]any)
	if obj != _undefined {
		iObj := obj.self.(*baseObject)
		r.convertToMap(iObj, mapValue)
	}
	resp, err := common.HttpPost(uri, mapValue)
	var o *Object
	if err != nil {
		return r.CreateErrorResponse(err.Error())
	} else {
		o = r.fromMap(resp)
	}
	return o
}

func (r *Runtime) builtin_newWebClient(args []Value, newTarget *Object) *Object {
	if newTarget == nil {
		panic(r.needNew("WebClient"))
	}
	proto := r.getPrototypeFromCtor(newTarget, r.global.WebClient, r.global.WebClientPrototype)
	o := &Object{runtime: r}

	mo := &webClientObject{}
	mo.class = classWebClient
	mo.val = o
	mo.extensible = true
	o.self = mo
	mo.prototype = proto
	mo.init(r, args)
	mo.m.ScriptHandler = r.ScriptHandler
	return o
}

func (r *Runtime) createWebClientProto(val *Object) objectImpl {
	o := newBaseObjectObj(val, r.global.WebClientPrototype, classWebClient)
	o._putProp("constructor", r.global.WebClient, true, false, true)
	o._putProp("load", r.newNativeFunc(r.webClientProto_load, nil, "load", nil, 2), true, false, true)
	o._putProp("saveCookies", r.newNativeFunc(r.webClientProto_saveCookies, nil, "saveCookies", nil, 0), true, false, true)
	o._putProp("pageBy", r.newNativeFunc(r.webClientProto_pageBy, nil, "pageBy", nil, 2), true, false, true)
	o._putProp("close", r.newNativeFunc(r.webClientProto_close, nil, "close", nil, 1), true, false, true)
	o._putProp("decodeSlider1", r.newNativeFunc(r.webClientProto_decodeSlider1, nil, "decodeSlider1", nil, 2), true, false, true)
	o._putProp("decodeSlider2", r.newNativeFunc(r.webClientProto_decodeSlider2, nil, "decodeSlider2", nil, 2), true, false, true)
	o._putProp("ocr", r.newNativeFunc(r.webClientProto_ocr, nil, "ocr", nil, 2), true, false, true)
	o._putProp("httpPost", r.newNativeFunc(r.webClientProto_httpPost, nil, "httpPost", nil, 2), true, false, true)
	return o
}

func (r *Runtime) createWebClient(val *Object) objectImpl {
	o := r.newNativeConstructOnly(val, r.builtin_newWebClient, r.global.WebClientPrototype, "WebClient", 0)
	r.putSpeciesReturnThis(o)
	return o
}

func (r *Runtime) initWebClient() {
	r.global.WebClientPrototype = r.newLazyObject(r.createWebClientProto)
	r.global.WebClient = r.newLazyObject(r.createWebClient)
	r.addToGlobal("WebClient", r.global.WebClient)
}
