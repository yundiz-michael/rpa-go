package goja

import (
	"log"
)

type Console struct {
	runtime   *Runtime
	formatter *Object
	printer   Printer
}

type Printer interface {
	Log(string)
	Warn(string)
	Debug(string)
	Error(string)
}

type PrinterFunc func(level string, s string)

func (p PrinterFunc) Log(s string) { p("info", s) }

func (p PrinterFunc) Warn(s string) { p("warn", s) }

func (p PrinterFunc) Error(s string) { p("error", s) }

func (p PrinterFunc) Debug(s string) { p("debug", s) }

var defaultPrinter Printer = PrinterFunc(func(level string, s string) { log.Print(level, s) })

func (c *Console) log(p func(string)) func(FunctionCall) Value {
	return func(call FunctionCall) Value {
		if format, ok := AssertFunction(c.formatter.Get("format")); ok {
			ret, err := format(c.formatter, call.Arguments...)
			if err != nil {
				panic(err)
			}
			p(ret.String())
		} else {
			panic(c.runtime.NewTypeError("util.format is not a function"))
		}
		return nil
	}
}

func RequireWithPrinter(printer Printer) ModuleLoader {
	return requireWithPrinter(printer)
}

func requireWithPrinter(printer Printer) ModuleLoader {
	return func(runtime *Runtime, module *Object) {
		c := &Console{
			runtime: runtime,
			printer: printer,
		}
		c.formatter = Require(runtime, "formatter").(*Object)
		o := module.Get("exports").(*Object)
		o.Set("log", c.log(c.printer.Log))
		o.Set("error", c.log(c.printer.Error))
		o.Set("warn", c.log(c.printer.Warn))
		o.Set("debug", c.log(c.printer.Debug))
	}
}

func InitConsole(runtime *Runtime, printer Printer) {
	registry := new(Registry)
	registry.Enable(runtime)
	/*注册格式化插件*/
	registry.RegisterNativeModule("formatter", func(runtime *Runtime, module *Object) {
		u := &Formatter{
			runtime: runtime,
		}
		obj := module.Get("exports").(*Object)
		obj.Set("format", u.js_format)
	})
	runtime.Set("formatter", Require(runtime, "formatter"))
	/*注册console插件*/
	if printer == nil {
		RegisterNativeModule("console", func(runtime *Runtime, module *Object) {
			requireWithPrinter(defaultPrinter)(runtime, module)
		})
	} else {
		registry.RegisterNativeModule("console", RequireWithPrinter(printer))
	}
	runtime.Set("console", Require(runtime, "console"))
}
