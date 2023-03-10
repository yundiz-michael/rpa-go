package goja

import (
	"encoding/hex"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"math"
	"merkaba/chromedp"
	"merkaba/common"
	"merkaba/goja/unistring"
	"merkaba/rpc"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

const hexUpper = "0123456789ABCDEF"

var (
	WebClients = make(map[string]map[string]*chromedp.WebClient)
)

func RemoveWebClient(domain string, taskName string) {
	webClientMap := WebClients[domain]
	if webClientMap == nil {
		return
	}
	delete(webClientMap, taskName)
}

var (
	parseFloatRegexp = regexp.MustCompile(`^([+-]?(?:Infinity|[0-9]*\.?[0-9]*(?:[eE][+-]?[0-9]+)?))`)
)

func (r *Runtime) builtin_isNaN(call FunctionCall) Value {
	if math.IsNaN(call.Argument(0).ToFloat()) {
		return valueTrue
	} else {
		return valueFalse
	}
}

func (r *Runtime) builtin_parseInt(call FunctionCall) Value {
	str := call.Argument(0).toString().toTrimmedUTF8()
	radix := int(toInt32(call.Argument(1)))
	v, _ := parseInt(str, radix)
	return v
}

func (r *Runtime) builtin_parseFloat(call FunctionCall) Value {
	m := parseFloatRegexp.FindStringSubmatch(call.Argument(0).toString().toTrimmedUTF8())
	if len(m) == 2 {
		if s := m[1]; s != "" && s != "+" && s != "-" {
			switch s {
			case "+", "-":
			case "Infinity", "+Infinity":
				return _positiveInf
			case "-Infinity":
				return _negativeInf
			default:
				f, err := strconv.ParseFloat(s, 64)
				if err == nil || isRangeErr(err) {
					return floatToValue(f)
				}
			}
		}
	}
	return _NaN
}

func (r *Runtime) builtin_isFinite(call FunctionCall) Value {
	f := call.Argument(0).ToFloat()
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return valueFalse
	}
	return valueTrue
}

func (r *Runtime) _encode(uriString valueString, unescaped *[256]bool) valueString {
	reader := uriString.reader()
	utf8Buf := make([]byte, utf8.UTFMax)
	needed := false
	l := 0
	for {
		rn, _, err := reader.ReadRune()
		if err != nil {
			if err != io.EOF {
				panic(r.newError(r.global.URIError, "Malformed URI"))
			}
			break
		}

		if rn >= utf8.RuneSelf {
			needed = true
			l += utf8.EncodeRune(utf8Buf, rn) * 3
		} else if !unescaped[rn] {
			needed = true
			l += 3
		} else {
			l++
		}
	}

	if !needed {
		return uriString
	}

	buf := make([]byte, l)
	i := 0
	reader = uriString.reader()
	for {
		rn, _, err := reader.ReadRune()
		if err == io.EOF {
			break
		}

		if rn >= utf8.RuneSelf {
			n := utf8.EncodeRune(utf8Buf, rn)
			for _, b := range utf8Buf[:n] {
				buf[i] = '%'
				buf[i+1] = hexUpper[b>>4]
				buf[i+2] = hexUpper[b&15]
				i += 3
			}
		} else if !unescaped[rn] {
			buf[i] = '%'
			buf[i+1] = hexUpper[rn>>4]
			buf[i+2] = hexUpper[rn&15]
			i += 3
		} else {
			buf[i] = byte(rn)
			i++
		}
	}
	return asciiString(buf)
}

func (r *Runtime) _decode(sv valueString, reservedSet *[256]bool) valueString {
	s := sv.String()
	hexCount := 0
	for i := 0; i < len(s); {
		switch s[i] {
		case '%':
			if i+2 >= len(s) || !ishex(s[i+1]) || !ishex(s[i+2]) {
				panic(r.newError(r.global.URIError, "Malformed URI"))
			}
			c := unhex(s[i+1])<<4 | unhex(s[i+2])
			if !reservedSet[c] {
				hexCount++
			}
			i += 3
		default:
			i++
		}
	}

	if hexCount == 0 {
		return sv
	}

	t := make([]byte, len(s)-hexCount*2)
	j := 0
	isUnicode := false
	for i := 0; i < len(s); {
		ch := s[i]
		switch ch {
		case '%':
			c := unhex(s[i+1])<<4 | unhex(s[i+2])
			if reservedSet[c] {
				t[j] = s[i]
				t[j+1] = s[i+1]
				t[j+2] = s[i+2]
				j += 3
			} else {
				t[j] = c
				if c >= utf8.RuneSelf {
					isUnicode = true
				}
				j++
			}
			i += 3
		default:
			if ch >= utf8.RuneSelf {
				isUnicode = true
			}
			t[j] = ch
			j++
			i++
		}
	}

	if !isUnicode {
		return asciiString(t)
	}

	us := make([]rune, 0, len(s))
	for len(t) > 0 {
		rn, size := utf8.DecodeRune(t)
		if rn == utf8.RuneError {
			if size != 3 || t[0] != 0xef || t[1] != 0xbf || t[2] != 0xbd {
				panic(r.newError(r.global.URIError, "Malformed URI"))
			}
		}
		us = append(us, rn)
		t = t[size:]
	}
	return unicodeStringFromRunes(us)
}

func ishex(c byte) bool {
	switch {
	case '0' <= c && c <= '9':
		return true
	case 'a' <= c && c <= 'f':
		return true
	case 'A' <= c && c <= 'F':
		return true
	}
	return false
}

func unhex(c byte) byte {
	switch {
	case '0' <= c && c <= '9':
		return c - '0'
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10
	}
	return 0
}

func (r *Runtime) builtin_decodeURI(call FunctionCall) Value {
	uriString := call.Argument(0).toString()
	return r._decode(uriString, &uriReservedHash)
}

func (r *Runtime) builtin_decodeURIComponent(call FunctionCall) Value {
	uriString := call.Argument(0).toString()
	return r._decode(uriString, &emptyEscapeSet)
}

func (r *Runtime) builtin_encodeURI(call FunctionCall) Value {
	uriString := call.Argument(0).toString()
	return r._encode(uriString, &uriReservedUnescapedHash)
}

func (r *Runtime) builtin_encodeURIComponent(call FunctionCall) Value {
	uriString := call.Argument(0).toString()
	return r._encode(uriString, &uriUnescaped)
}

func (r *Runtime) builtin_escape(call FunctionCall) Value {
	s := call.Argument(0).toString()
	var sb strings.Builder
	l := s.length()
	for i := 0; i < l; i++ {
		r := uint16(s.charAt(i))
		if r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z' || r >= '0' && r <= '9' ||
			r == '@' || r == '*' || r == '_' || r == '+' || r == '-' || r == '.' || r == '/' {
			sb.WriteByte(byte(r))
		} else if r <= 0xff {
			sb.WriteByte('%')
			sb.WriteByte(hexUpper[r>>4])
			sb.WriteByte(hexUpper[r&0xf])
		} else {
			sb.WriteString("%u")
			sb.WriteByte(hexUpper[r>>12])
			sb.WriteByte(hexUpper[(r>>8)&0xf])
			sb.WriteByte(hexUpper[(r>>4)&0xf])
			sb.WriteByte(hexUpper[r&0xf])
		}
	}
	return asciiString(sb.String())
}

func (r *Runtime) builtin_unescape(call FunctionCall) Value {
	s := call.Argument(0).toString()
	l := s.length()
	var asciiBuf []byte
	var unicodeBuf []uint16
	_, u := devirtualizeString(s)
	unicode := u != nil
	if unicode {
		unicodeBuf = make([]uint16, 1, l+1)
		unicodeBuf[0] = unistring.BOM
	} else {
		asciiBuf = make([]byte, 0, l)
	}
	for i := 0; i < l; {
		r := s.charAt(i)
		if r == '%' {
			if i <= l-6 && s.charAt(i+1) == 'u' {
				c0 := s.charAt(i + 2)
				c1 := s.charAt(i + 3)
				c2 := s.charAt(i + 4)
				c3 := s.charAt(i + 5)
				if c0 <= 0xff && ishex(byte(c0)) &&
					c1 <= 0xff && ishex(byte(c1)) &&
					c2 <= 0xff && ishex(byte(c2)) &&
					c3 <= 0xff && ishex(byte(c3)) {
					r = rune(unhex(byte(c0)))<<12 |
						rune(unhex(byte(c1)))<<8 |
						rune(unhex(byte(c2)))<<4 |
						rune(unhex(byte(c3)))
					i += 5
					goto out
				}
			}
			if i <= l-3 {
				c0 := s.charAt(i + 1)
				c1 := s.charAt(i + 2)
				if c0 <= 0xff && ishex(byte(c0)) &&
					c1 <= 0xff && ishex(byte(c1)) {
					r = rune(unhex(byte(c0))<<4 | unhex(byte(c1)))
					i += 2
				}
			}
		}
	out:
		if r >= utf8.RuneSelf && !unicode {
			unicodeBuf = make([]uint16, 1, l+1)
			unicodeBuf[0] = unistring.BOM
			for _, b := range asciiBuf {
				unicodeBuf = append(unicodeBuf, uint16(b))
			}
			asciiBuf = nil
			unicode = true
		}
		if unicode {
			unicodeBuf = append(unicodeBuf, uint16(r))
		} else {
			asciiBuf = append(asciiBuf, byte(r))
		}
		i++
	}
	if unicode {
		return unicodeString(unicodeBuf)
	}

	return asciiString(asciiBuf)
}

func (r *Runtime) fillMapByContext(target map[string]any) {
	target["cookieId"] = r.Context.CookieId
	target["taskName"] = r.Context.TaskName
	target["appServerPort"] = r.Context.AppServerPort
	target["appServerIP"] = r.Context.AppServerIP
	target["scriptUri"] = r.Context.ScriptUri
	target["runMode"] = r.Context.RunMode
	target["request"] = r.Context.Parameters
}

// builtin_success  记录任务完成的信息
func (r *Runtime) builtin_success(call FunctionCall) Value {
	dataObject := call.Argument(0).ToObject(r).self.(*baseObject)
	data := make(map[string]any)
	r.convertToMap(dataObject, data)
	r.Context.Info("🚥🚥🚥progress", zap.Any("content", data))
	msg := make(map[string]any)
	r.fillMapByContext(msg)
	msg["data"] = data
	common.Notify("success", msg)
	return valueTrue
}

// builtin_sendData  发送数据到pulsar队列
func (r *Runtime) builtin_sendData(call FunctionCall) Value {
	format := call.Argument(0).String()
	dataObject := call.Argument(1).ToObject(r).self.(*baseObject)
	msg := make(map[string]any)
	r.convertToMap(dataObject, msg)
	r.fillMapByContext(msg)
	err := common.SendData("data", format, r.Context, msg)
	if err != nil {
		r.Context.Error("sendData", zap.Error(err))
		return valueFalse
	} else {
		return valueTrue
	}
}

// builtin_paramBy  读取参数
func (r *Runtime) builtin_paramBy(call FunctionCall) Value {
	name := call.Argument(0).String()
	if r.Context.Parameters == nil {
		return nil
	}
	value := r.Context.Parameters[name]
	if value != nil {
		return r.ToValue(value)
	} else if len(call.Arguments) == 2 {
		return r.ToValue(call.Argument(1))
	} else {
		return nil
	}
}

func (r *Runtime) builtin_remoteCall(call FunctionCall) Value {
	funcNames := strings.Split(call.Argument(0).String(), ".")
	if len(funcNames) != 2 {
		err := errors.New("first parameter must be 'serviceName.methodName' format")
		r.Context.Error("invoke", zap.Error(err))
		common.SendMessage("error", r.Context, err.Error())
		return nil
	}
	param := make(map[string]any)
	if len(call.Arguments) == 2 {
		dataObject := call.Argument(1).ToObject(r).self.(*baseObject)
		r.convertToMap(dataObject, param)
		r.fillMapByContext(param)
	}
	client := rpc.UsePaasClient(r.Context.AppServerIP, funcNames[0], r.Context.CookieId)
	if client == nil {
		return nil
	}
	timeOut := 2 * time.Minute
	if len(call.Arguments) == 3 {
		i := call.Argument(2).ToInteger()
		timeOut = time.Duration(i) * time.Second
	}
	result, err := client.DoInvoke(funcNames[1], param, make([]byte, 0), timeOut)
	client.Dispose()
	if err != nil {
		r.Context.Error("invoke", zap.Error(err))
		common.SendMessage("error", r.Context, err.Error())
		panic(r.ToValue(err.Error()))
		return nil
	}
	return r.ToValue(result)
}

func (r *Runtime) builtin_writeServiceLog(call FunctionCall) Value {
	bizType := call.Argument(0)
	level := "info"
	serviceName := "Log"
	funcName := "merkabaUpdate"
	if bizType == nil {
		panic(r.ToValue("日志函数参数错误，请指定bizType"))
	}
	param := make(map[string]any)
	dataObject := call.Argument(1).ToObject(r).self.(*baseObject)
	r.convertToMap(dataObject, param)
	if len(call.Arguments) == 3 {
		level = call.Argument(2).String()
	}
	param["bizType"] = bizType
	param["level"] = level
	r.fillMapByContext(param)
	client := rpc.UsePaasClient(r.Context.AppServerIP, serviceName, r.Context.CookieId)
	if client == nil {
		return nil
	}
	result, err := client.Invoke(funcName, param, make([]byte, 0))
	client.Dispose()
	if err != nil {
		r.Context.Error("invoke", zap.Error(err))
		common.SendMessage("error", r.Context, err.Error())
		panic(r.ToValue("日志RPC异常"))
		return nil
	}
	return r.ToValue(result)
}

func (r *Runtime) builtin_handleData(call FunctionCall) Value {
	var data map[string]any
	dataObject := call.Argument(1).ToObject(r).self.(*baseObject)
	data = make(map[string]any)
	r.convertToMap(dataObject, data)
	format := call.Argument(0).String()
	data["format"] = format
	r.fillMapByContext(data)
	client := rpc.UsePaasClient(r.Context.AppServerIP, "Merkaba", r.Context.CookieId)
	result, err := client.Invoke("onDataReceived", data, make([]byte, 0))
	client.Dispose()
	if err != nil {
		r.Context.Error("Merkaba.onDataReceived", zap.Error(err))
		common.SendMessage("error", r.Context, err.Error())
		return nil
	}
	return r.ToValue(result)
}

func (r *Runtime) builtin_handleFile(call FunctionCall) Value {
	var params map[string]any
	dataObject := call.Argument(0).ToObject(r).self.(*baseObject)
	params = make(map[string]any)
	r.convertToMap(dataObject, params)
	if _, ok := params["format"]; !ok {
		err := zap.Error(errors.New("参数必须包含format"))
		r.Context.Error("handleFile", err)
		common.SendMessage("error", r.Context, err.String)
		return nil
	}
	fileName := call.Argument(1).String()
	file, _ := os.Open(fileName)
	defer file.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		r.Context.Error("handleFile", zap.Error(err))
		common.SendMessage("error", r.Context, err.Error())
		return nil
	}
	client := rpc.UsePaasClient(r.Context.AppServerIP, "Merkaba", r.Context.CookieId)
	result, err := client.Invoke("onFileReceived", params, bytes)
	client.Dispose()
	if err != nil {
		r.Context.Error("Merkaba.onFileReceived", zap.Error(err))
		common.SendMessage("error", r.Context, err.Error())
		return nil
	}
	return r.ToValue(result)
}

func (r *Runtime) builtin_wait(call FunctionCall) Value {
	v := call.Argument(0).ToInteger()
	r.Context.Info("wait", zap.Int64("timeout", v))
	time.Sleep(time.Duration(v) * time.Second)
	return nil
}

func (r *Runtime) builtin_webClient(call FunctionCall) Value {
	domain := call.Argument(0).String()
	webClientMap := WebClients[domain]
	if webClientMap == nil {
		webClientMap = make(map[string]*chromedp.WebClient)
		WebClients[domain] = webClientMap
	}
	webClient := webClientMap[r.Context.TaskName]
	if webClient != nil {
		r.WebClient = webClient
		return r.CreateWebClientObject(webClient)
	}
	r.Context.Info("webClient", zap.String("domain", domain))
	if r.Context.RunMode == 0 {
		common.SendMessage("info", r.Context, fmt.Sprintf("new webClient domain %s", domain))
	}
	option := chromedp.WebClientOption{Domain: domain, Proxy: r.Context.Proxy, Headless: r.Context.Headless}
	webClient = chromedp.NewClient(option, r.Context)
	webClientMap[r.Context.TaskName] = webClient
	r.WebClient = webClient
	webClient.ScriptHandler = r.ScriptHandler
	o := r.CreateWebClientObject(webClient)
	return o
}

func (r *Runtime) initGlobalObject() {
	o := r.globalObject.self
	o._putProp("globalThis", r.globalObject, true, false, true)
	o._putProp("NaN", _NaN, false, false, false)
	o._putProp("undefined", _undefined, false, false, false)
	o._putProp("Infinity", _positiveInf, false, false, false)

	o._putProp("isNaN", r.newNativeFunc(r.builtin_isNaN, nil, "isNaN", nil, 1), true, false, true)
	o._putProp("parseInt", r.newNativeFunc(r.builtin_parseInt, nil, "parseInt", nil, 2), true, false, true)
	o._putProp("parseFloat", r.newNativeFunc(r.builtin_parseFloat, nil, "parseFloat", nil, 1), true, false, true)
	o._putProp("isFinite", r.newNativeFunc(r.builtin_isFinite, nil, "isFinite", nil, 1), true, false, true)
	o._putProp("decodeURI", r.newNativeFunc(r.builtin_decodeURI, nil, "decodeURI", nil, 1), true, false, true)
	o._putProp("decodeURIComponent", r.newNativeFunc(r.builtin_decodeURIComponent, nil, "decodeURIComponent", nil, 1), true, false, true)
	o._putProp("encodeURI", r.newNativeFunc(r.builtin_encodeURI, nil, "encodeURI", nil, 1), true, false, true)
	o._putProp("encodeURIComponent", r.newNativeFunc(r.builtin_encodeURIComponent, nil, "encodeURIComponent", nil, 1), true, false, true)
	o._putProp("escape", r.newNativeFunc(r.builtin_escape, nil, "escape", nil, 1), true, false, true)
	o._putProp("unescape", r.newNativeFunc(r.builtin_unescape, nil, "unescape", nil, 1), true, false, true)
	//注册发送数据到pulsar队列
	o._putProp("success", r.newNativeFunc(r.builtin_success, nil, "success", nil, 1), true, false, true)
	o._putProp("sendData", r.newNativeFunc(r.builtin_sendData, nil, "sendData", nil, 1), true, false, true)
	o._putProp("paramBy", r.newNativeFunc(r.builtin_paramBy, nil, "paramBy", nil, 1), true, false, true)
	o._putProp("remoteCall", r.newNativeFunc(r.builtin_remoteCall, nil, "remoteCall", nil, 1), true, false, true)
	o._putProp("wait", r.newNativeFunc(r.builtin_wait, nil, "wait", nil, 1), true, false, true)
	o._putProp("webClient", r.newNativeFunc(r.builtin_webClient, nil, "webClient", nil, 1), true, false, true)
	o._putProp("breakPoint", r.newNativeFunc(r.builtin_breakPoint, nil, "breakPoint", nil, 1), true, false, true)
	o._putProp("writeFile", r.newNativeFunc(r.builtin_writeFile, nil, "writeFile", nil, 2), true, false, true)
	o._putProp("removeFile", r.newNativeFunc(r.builtin_removeFile, nil, "removeFile", nil, 1), true, false, true)
	o._putProp("downFile", r.newNativeFunc(r.builtin_downFile, nil, "downFile", nil, 1), true, false, true)
	o._putProp("handleData", r.newNativeFunc(r.builtin_handleData, nil, "handleData", nil, 1), true, false, true)
	o._putProp("handleFile", r.newNativeFunc(r.builtin_handleFile, nil, "handleFile", nil, 2), true, false, true)
	o._putProp("writeServiceLog", r.newNativeFunc(r.builtin_writeServiceLog, nil, "writeServiceLog", nil, 4), true, false, true)
	o._putSym(SymToStringTag, valueProp(asciiString(classGlobal), false, false, true))

	// TODO: Annex B

}

func digitVal(d byte) int {
	var v byte
	switch {
	case '0' <= d && d <= '9':
		v = d - '0'
	case 'a' <= d && d <= 'z':
		v = d - 'a' + 10
	case 'A' <= d && d <= 'Z':
		v = d - 'A' + 10
	default:
		return 36
	}
	return int(v)
}

// ECMAScript compatible version of strconv.ParseInt
func parseInt(s string, base int) (Value, error) {
	var n int64
	var err error
	var cutoff, maxVal int64
	var sign bool
	i := 0

	if len(s) < 1 {
		err = strconv.ErrSyntax
		goto Error
	}

	switch s[0] {
	case '-':
		sign = true
		s = s[1:]
	case '+':
		s = s[1:]
	}

	if len(s) < 1 {
		err = strconv.ErrSyntax
		goto Error
	}

	// Look for hex_ prefix.
	if s[0] == '0' && len(s) > 1 && (s[1] == 'x' || s[1] == 'X') {
		if base == 0 || base == 16 {
			base = 16
			s = s[2:]
		}
	}

	switch {
	case len(s) < 1:
		err = strconv.ErrSyntax
		goto Error

	case 2 <= base && base <= 36:
	// valid base; nothing to do

	case base == 0:
		// Look for hex_ prefix.
		switch {
		case s[0] == '0' && len(s) > 1 && (s[1] == 'x' || s[1] == 'X'):
			if len(s) < 3 {
				err = strconv.ErrSyntax
				goto Error
			}
			base = 16
			s = s[2:]
		default:
			base = 10
		}

	default:
		err = errors.New("invalid base " + strconv.Itoa(base))
		goto Error
	}

	// Cutoff is the smallest number such that cutoff*base > maxInt64.
	// Use compile-time constants for common cases.
	switch base {
	case 10:
		cutoff = math.MaxInt64/10 + 1
	case 16:
		cutoff = math.MaxInt64/16 + 1
	default:
		cutoff = math.MaxInt64/int64(base) + 1
	}

	maxVal = math.MaxInt64
	for ; i < len(s); i++ {
		if n >= cutoff {
			// n*base overflows
			return parseLargeInt(float64(n), s[i:], base, sign)
		}
		v := digitVal(s[i])
		if v >= base {
			break
		}
		n *= int64(base)

		n1 := n + int64(v)
		if n1 < n || n1 > maxVal {
			// n+v overflows
			return parseLargeInt(float64(n)+float64(v), s[i+1:], base, sign)
		}
		n = n1
	}

	if i == 0 {
		err = strconv.ErrSyntax
		goto Error
	}

	if sign {
		n = -n
	}
	return intToValue(n), nil

Error:
	return _NaN, err
}

func parseLargeInt(n float64, s string, base int, sign bool) (Value, error) {
	i := 0
	b := float64(base)
	for ; i < len(s); i++ {
		v := digitVal(s[i])
		if v >= base {
			break
		}
		n = n*b + float64(v)
	}
	if sign {
		n = -n
	}
	// We know it can't be represented as int, so use valueFloat instead of floatToValue
	return valueFloat(n), nil
}

var (
	uriUnescaped             [256]bool
	uriReserved              [256]bool
	uriReservedHash          [256]bool
	uriReservedUnescapedHash [256]bool
	emptyEscapeSet           [256]bool
)

func init() {
	for _, c := range "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_.!~*'()" {
		uriUnescaped[c] = true
	}

	for _, c := range ";/?:@&=+$," {
		uriReserved[c] = true
	}

	for i := 0; i < 256; i++ {
		if uriUnescaped[i] || uriReserved[i] {
			uriReservedUnescapedHash[i] = true
		}
		uriReservedHash[i] = uriReserved[i]
	}
	uriReservedUnescapedHash['#'] = true
	uriReservedHash['#'] = true
}

func (r *Runtime) builtin_breakPoint(call FunctionCall) Value {
	v := call.Argument(0).String()
	if strings.HasPrefix(v, "items") {
		r.Context.Info("debug", zap.String("name", v))
	}
	return nil
}

func (r *Runtime) builtin_downFile(call FunctionCall) Value {
	url := call.Argument(0).String()
	fileName := call.Argument(1).String()
	fullFileName := common.Env.Path.Temp + fileName
	err := common.HttpGetFile(url, make(map[string]interface{}), fullFileName)
	if err != nil {
		return nil
	}
	return r.ToValue(fullFileName)
}

func (r *Runtime) builtin_writeFile(call FunctionCall) Value {
	fileName := call.Argument(0).String()
	names := strings.Split(fileName, ".")
	timeNow := time.Now()
	currentDay := timeNow.Format("2006-01-02")
	fileName = names[0] + "_" + currentDay
	var fileExt string
	if len(names) >= 2 {
		fileExt = names[1]
	} else {
		fileExt = ".txt"
	}
	content := call.Argument(1).String()
	fullFileName := common.Env.Path.Temp + fileName + "." + fileExt
	f, err := os.Create(fullFileName)
	if err != nil {
		return nil
	}
	if fileExt == "csv" {
		head, _ := hex.DecodeString("EFBBBF")
		f.Write(head)
	}
	f.WriteString(content)
	f.Close()
	return r.ToValue(fullFileName)
}

func (r *Runtime) builtin_removeFile(call FunctionCall) Value {
	fileName := call.Argument(0).String()
	err := os.Remove(fileName)
	if err != nil {
		return r.ToValue(true)
	} else {
		return r.ToValue(false)
	}
}
