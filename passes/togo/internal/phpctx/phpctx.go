package phpctx

import (
	"errors"
	"flag"
	"io"
	"net/http"
	"os/exec"
	"reflect"
	"strings"
)

var (
	// ErrNoStruct ...
	ErrNoStruct     = errors.New("receiver is not a struct")
	// ErrMissingField ...
	ErrMissingField = errors.New("struct is missing field")
	// ErrNotSet ...
	ErrNotSet       = errors.New("dynamic value is not set")
)

var zero = reflect.Value{}
var addr = flag.String("http-address", "localhost:8080", "host and port for serving http")

// PHPVar is a PHPVar
type PHPVar interface{}

// PHPContext is a PHPContext
type PHPContext struct {
	Echo          io.Writer
	dynamicValues map[string]*interface{}
}

func (ctx PHPContext) SetDynamic(name string, value interface{}) {
	ctx.dynamicValues[name] = &value
}

func (ctx PHPContext) GetDynamic(name string) (interface{}, error) {
	v, ok := ctx.dynamicValues[name]
	if !ok {
		return nil, ErrNotSet
	}

	if v == nil {
		return nil, nil
	}

	return *v, nil
}

func (ctx PHPContext) Write(b []byte) (int, error) {
	return ctx.Echo.Write(b)
}

// ListenAndServe (see net/http)
func ListenAndServe() {
	if !flag.Parsed() {
		flag.Parse()
	}
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		panic(err)
	}
}

// GetDynamicProperty returns a dynamic property
func GetDynamicProperty(rcvr interface{}, field string) (interface{}, error) {
	v := reflect.ValueOf(rcvr)
	if v.Kind() == reflect.Struct {
		f := v.FieldByName(field)
		if f == zero {
			return nil, ErrMissingField
		}
		return f.Interface(), nil
	}
	return nil, ErrNoStruct
}

// Shell returns a shell output
func Shell(cmd string) ([]byte, error) {
	cmdParts := strings.SplitN(cmd, " ", 2)
	cmdName := cmdParts[0]
	args := cmdParts[1]
	c := exec.Command(cmdName, args)
	err := c.Run()
	if err != nil {
		return nil, err
	}

	return c.Output()
}

// Server is a server
type Server interface {
	ServePHP(ctx PHPContext)
}
