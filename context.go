package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/labstack/gommon/color"
	"github.com/mattn/go-colorable"
	"github.com/mkideal/pkg/debug"
)

type (
	// Context provide running context
	Context struct {
		router     []string
		path       string
		argv       interface{}
		nativeArgs []string
		flagSet    *flagSet
		command    *Command
		writer     io.Writer
		color      color.Color

		HTTPRequest  *http.Request
		HTTPResponse http.ResponseWriter
	}

	// Validator validate flag before running command
	Validator interface {
		Validate(*Context) error
	}

	// AutoHelper represents interface for showing help information automaticly
	AutoHelper interface {
		AutoHelp() bool
	}
)

func newContext(path string, router, args []string, argv interface{}, clr color.Color) (*Context, error) {
	ctx := &Context{
		path:       path,
		router:     router,
		argv:       argv,
		nativeArgs: args,
		color:      clr,
		flagSet:    newFlagSet(),
	}
	if argv != nil {
		ctx.flagSet = parseArgv(args, argv, ctx.color)
		if ctx.flagSet.err != nil {
			return nil, ctx.flagSet.err
		}
	}
	return ctx, nil
}

// Path returns full command name
// `./app hello world -a --xyz=1` will returns "hello world"
func (ctx *Context) Path() string {
	return ctx.path
}

// Router returns full command name with string array
// `./app hello world -a --xyz=1` will returns ["hello" "world"]
func (ctx *Context) Router() []string {
	return ctx.router
}

// NativeArgs returns native args
// `./app hello world -a --xyz=1` will return ["-a" "--xyz=1"]
func (ctx *Context) NativeArgs() []string {
	return ctx.nativeArgs
}

// Args returns free args
// `./app hello world -a=1 abc xyz` will return ["abc" "xyz"]
func (ctx *Context) Args() []string {
	return ctx.flagSet.args
}

// NArg returns length of Args
func (ctx *Context) NArg() int {
	return len(ctx.flagSet.args)
}

// Argv returns parsed args object
func (ctx *Context) Argv() interface{} {
	return ctx.argv
}

// IsSet determins wether `flag` be set
func (ctx *Context) IsSet(flag string) bool {
	fl, ok := ctx.flagSet.flagMap[flag]
	if !ok {
		return false
	}
	return fl.actual != ""
}

// FormValues returns parsed args as url.Values
func (ctx *Context) FormValues() url.Values {
	if ctx.flagSet == nil {
		debug.Panicf("ctx.flagSet == nil")
	}
	return ctx.flagSet.values
}

// Command returns current command instance
func (ctx *Context) Command() *Command {
	return ctx.command
}

// Usage returns current command's usage with current context
func (ctx *Context) Usage() string {
	return ctx.command.Usage(ctx)
}

// WriteUsage writes usage to writer
func (ctx *Context) WriteUsage() {
	ctx.String(ctx.Usage())
}

// Writer returns writer
func (ctx *Context) Writer() io.Writer {
	if ctx.writer == nil {
		ctx.writer = colorable.NewColorableStdout()
	}
	return ctx.writer
}

// Write implements io.Writer
func (ctx *Context) Write(data []byte) (n int, err error) {
	return ctx.Writer().Write(data)
}

// Color returns color instance
func (ctx *Context) Color() *color.Color {
	return &ctx.color
}

// String writes formatted string to writer
func (ctx *Context) String(format string, args ...interface{}) *Context {
	fmt.Fprintf(ctx.Writer(), format, args...)
	return ctx
}

// JSON writes json string of obj to writer
func (ctx *Context) JSON(obj interface{}) *Context {
	data, err := json.Marshal(obj)
	if err == nil {
		fmt.Fprintf(ctx.Writer(), string(data))
	}
	return ctx
}

// JSONln writes json string of obj end with "\n" to writer
func (ctx *Context) JSONln(obj interface{}) *Context {
	return ctx.JSON(obj).String("\n")
}

// JSONIndent writes pretty json string of obj to writer
func (ctx *Context) JSONIndent(obj interface{}, prefix, indent string) *Context {
	data, err := json.MarshalIndent(obj, prefix, indent)
	if err == nil {
		fmt.Fprintf(ctx.Writer(), string(data))
	}
	return ctx
}

// JSONIndentln writes pretty json string of obj end with "\n" to writer
func (ctx *Context) JSONIndentln(obj interface{}, prefix, indent string) *Context {
	return ctx.JSONIndent(obj, prefix, indent).String("\n")
}