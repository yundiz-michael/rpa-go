package goja

import (
	"bytes"
)

type Formatter struct {
	runtime *Runtime
}

func (u *Formatter) format(f rune, val Value, w *bytes.Buffer) bool {
	switch f {
	case 's':
		w.WriteString(val.String())
	case 'd':
		w.WriteString(val.ToNumber().String())
	case 'j':
		if json, ok := u.runtime.Get("JSON").(*Object); ok {
			if stringify, ok := AssertFunction(json.Get("stringify")); ok {
				res, err := stringify(json, val)
				if err != nil {
					w.WriteString(err.Error())
				} else {
					w.WriteString(res.String())
				}
			}
		}
	case '%':
		w.WriteByte('%')
		return false
	default:
		w.WriteByte('%')
		w.WriteRune(f)
		return false
	}
	return true
}

func (u *Formatter) Format(b *bytes.Buffer, f string, args ...Value) {
	pct := false
	argNum := 0
	for _, chr := range f {
		if pct {
			if argNum < len(args) {
				if u.format(chr, args[argNum], b) {
					argNum++
				}
			} else {
				b.WriteByte('%')
				b.WriteRune(chr)
			}
			pct = false
		} else {
			if chr == '%' {
				pct = true
			} else {
				b.WriteRune(chr)
			}
		}
	}

	for _, arg := range args[argNum:] {
		b.WriteByte(' ')
		b.WriteString(arg.String())
	}
}

func (u *Formatter) js_format(call FunctionCall) Value {
	var b bytes.Buffer
	var fmt string
	var args []Value
	if len(call.Arguments) == 1 {
		fmt = "%j"
		args = call.Arguments
	} else if arg := call.Argument(0); !IsUndefined(arg) {
		fmt = arg.String()
		args = call.Arguments[1:]
	}
	u.Format(&b, fmt, args...)
	return u.runtime.ToValue(b.String())
}

func NewFormatter(runtime *Runtime) *Formatter {
	return &Formatter{
		runtime: runtime,
	}
}
