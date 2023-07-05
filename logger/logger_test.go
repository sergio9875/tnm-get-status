package logger

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"malawi-getstatus/utils"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestRedact(t *testing.T) {
	type RedactMe struct {
		Name     string    `redact:"-"`
		Password string    `redact:"complete"`
		First6   string    `redact:"first6"`
		Last4    string    `redact:"last4"`
		Dob      time.Time `redact:"last4"`
		Ptr      *string   `redact:"last4"`
	}
	var ptr = "1234567890"
	dob, _ := time.Parse(
		time.RFC3339,
		"2012-11-01T22:08:41+00:00")
	var obj = RedactMe{
		Name:     "eugenep",
		Password: "password01",
		First6:   "1234567890",
		Last4:    "0987654321",
		Dob:      dob,
		Ptr:      &ptr,
	}

	//val := reflect.ValueOf(&obj).Interface()
	redact(&obj)

	if obj.Name != "eugenep" {
		t.Errorf("Expected: %v, but got %v", "eugenep", obj.Name)
	}
	if obj.Password != strings.Repeat("*", len("password01")) {
		t.Errorf("Expected: %v, but got %v", strings.Repeat("*", len("password01")), obj.Password)
	}
	if obj.First6 != "123456****" {
		t.Errorf("Expected: %v, but got %v", "123456****", obj.First6)
	}
	if obj.Last4 != "******4321" {
		t.Errorf("Expected: %v, but got %v", "******4321", obj.Last4)
	}
	if *obj.Ptr != "******7890" {
		t.Errorf("Expected: %v, but got %v", "******7890", *obj.Ptr)
	}

}

func TestRedact2(t *testing.T) {
	type InnerRedact struct {
		Length   int    `redact:"-"`
		Password string `redact:"complete"`
	}
	type RedactMe struct {
		Name     string      `redact:"-"`
		Password InnerRedact `redact:"complete"`
	}
	var iobj = InnerRedact{
		Length:   10,
		Password: "password01",
	}
	var obj = RedactMe{
		Name:     "eugenep",
		Password: iobj,
	}

	//val := reflect.ValueOf(&obj).Interface()
	redact(&obj)

	if obj.Name != "eugenep" {
		t.Errorf("Expected: %v, but got %v", "eugenep", obj.Name)
	}
	if obj.Password.Length != len("password01") {
		t.Errorf("Expected: %v, but got %v", len("password01"), obj.Password.Length)
	}
	if obj.Password.Password != strings.Repeat("*", len("password01")) {
		t.Errorf("Expected: %v, but got %v", strings.Repeat("*", len("password01")), obj.Password)
	}
}

func TestSanitizePtr(t *testing.T) {
	type InnerRedact struct {
		Length   int    `redact:"-"`
		Password string `redact:"complete"`
	}
	type RedactMe struct {
		Name     string      `redact:"-"`
		Password InnerRedact `redact:"complete"`
	}
	var iobj = InnerRedact{
		Length:   10,
		Password: "password01",
	}
	var obj = RedactMe{
		Name:     "eugenep",
		Password: iobj,
	}

	//val := reflect.ValueOf(&obj).Interface()
	v := Sanitizer(&obj)

	result := *((v).(*RedactMe))

	if result.Name != "eugenep" {
		t.Errorf("Expected: %v, but got %v", "eugenep", result.Name)
	}
	if result.Password.Length != len("password01") {
		t.Errorf("Expected: %v, but got %v", len("password01"), result.Password.Length)
	}
	if result.Password.Password != strings.Repeat("*", len("password01")) {
		t.Errorf("Expected: %v, but got %v", strings.Repeat("*", len("password01")), result.Password)
	}
}

func TestSanitizeObject(t *testing.T) {
	type InnerRedact struct {
		Length   int    `redact:"-"`
		Password string `redact:"complete"`
	}
	type RedactMe struct {
		Name     string      `redact:"-"`
		Password InnerRedact `redact:"complete"`
	}
	var iobj = InnerRedact{
		Length:   10,
		Password: "password01",
	}
	var obj = RedactMe{
		Name:     "eugenep",
		Password: iobj,
	}

	//val := reflect.ValueOf(&obj).Interface()
	v := Sanitizer(obj)

	result := *((v).(*RedactMe))

	if result.Name != "eugenep" {
		t.Errorf("Expected: %v, but got %v", "eugenep", result.Name)
	}
	if result.Password.Length != len("password01") {
		t.Errorf("Expected: %v, but got %v", len("password01"), result.Password.Length)
	}
	if result.Password.Password != strings.Repeat("*", len("password01")) {
		t.Errorf("Expected: %v, but got %v", strings.Repeat("*", len("password01")), result.Password)
	}
}

func srcTest() string {
	// Determine caller func
	pc, file, lineno, ok := runtime.Caller(1)
	src := ""
	if ok {
		slice := strings.Split(runtime.FuncForPC(pc).Name(), "/")
		src = slice[len(slice)-1]
		slice = strings.Split(file, "/")
		file := slice[len(slice)-1]
		src = fmt.Sprintf("%s at %s:%d", src, file, lineno-1)
	}
	return src
}

func TestNonFormatLoggers(t *testing.T) {
	isOkay := func(fn interface{}, reqId string, level string, logLine string) bool {
		var buf bytes.Buffer
		log.SetOutput(&buf)
		defer func() {
			log.SetOutput(os.Stderr)
		}()

		fn.(func(string, ...interface{}))(reqId, logLine)
		src := srcTest()
		result := false
		if level == "INFO" {
			result = strings.HasSuffix(buf.String(), fmt.Sprintf("[%s] %s\n",
				level, logLine))
		} else {
			result = strings.HasSuffix(buf.String(), fmt.Sprintf("[%s] (%s) %s\n",
				level, src, logLine))
		}
		buf.Reset()
		return result
	}

	type args struct {
		fn      interface{}
		reqId   string
		level   string
		logLine string
	}

	tests := []struct {
		name      string
		arg       args
		assertion assert.BoolAssertionFunc
	}{
		{"Trace good", args{Trace, "reqId", TRACE, "hello info"}, assert.True},
		{"Println good", args{Println, "", INFO, "hello info"}, assert.True},
		{"Debug good", args{Debug, "reqId", DEBUG, "hello debug"}, assert.True},
		{"Info good", args{Info, "reqId", INFO, "hello info"}, assert.True},
		{"Warn good", args{Warn, "reqId", WARN, "hello warn"}, assert.True},
		{"Error good", args{Error, "reqId", ERROR, "hello error"}, assert.True},
	}

	os.Setenv("LOG_LEVEL", "TRACE")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assertion(t, isOkay(tt.arg.fn, tt.arg.reqId, tt.arg.level, tt.arg.logLine))
		})
	}
	os.Setenv("LOG_LEVEL", "")
}

func TestFormatLoggers(t *testing.T) {
	isOkay := func(fn interface{}, level string, format string, a ...interface{}) bool {
		var buf bytes.Buffer
		log.SetOutput(&buf)
		defer func() {
			log.SetOutput(os.Stderr)
		}()

		fn.(func(string, string, ...interface{}))("reqId", format, a...)
		src := srcTest()
		result := false
		msg := fmt.Sprintf(format, a...)
		if level == "INFO" {
			result = strings.HasSuffix(buf.String(), fmt.Sprintf("[%s] %s\n",
				level, msg))
		} else {
			result = strings.HasSuffix(buf.String(), fmt.Sprintf("[%s] (%s) %s\n",
				level, src, msg))
		}
		buf.Reset()
		return result
	}

	type args struct {
		fn     interface{}
		level  string
		format string
		a      interface{}
	}

	tests := []struct {
		name      string
		arg       args
		assertion assert.BoolAssertionFunc
	}{
		{"Printf good", args{Printf, INFO, "hello %s", "info"}, assert.True},
		{"Debugf good", args{Debugf, DEBUG, "hello %s", "debug"}, assert.True},
		{"Infof good", args{Infof, INFO, "hello %s", "info"}, assert.True},
		{"Warnf good", args{Warnf, WARN, "hello %s", "warn"}, assert.True},
		{"Errorf good", args{Errorf, ERROR, "hello %s", "error"}, assert.True},
	}

	os.Setenv("LOG_LEVEL", "TRACE")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assertion(t, isOkay(tt.arg.fn, tt.arg.level, tt.arg.format, tt.arg.a))
		})
	}
	os.Setenv("LOG_LEVEL", "")
}

func TestLogIt(t *testing.T) {
	_ = os.Setenv("LOG_LEVEL", "TRACE")

	Init()
	isOkay := func(fn interface{}, reqId string, level string, format string, a ...interface{}) bool {
		var buf bytes.Buffer
		log.SetOutput(&buf)
		defer func() {
			log.SetOutput(os.Stderr)
		}()
		cl := os.Getenv("LOG_LEVEL")
		if level == DEBUG {
			_ = os.Setenv("LOG_LEVEL", "INFO")
		}
		fn.(func(string, string, ...interface{}))(reqId, format, a...)
		src := srcTest()
		result := false
		a[0] = utils.JsonIt(Sanitizer(a[0]))
		msg := fmt.Sprintf(format, a[0])
		if level == "INFO" {
			msg2 := buf.String()
			result = strings.HasSuffix(msg2, fmt.Sprintf("[%s] %s",
				level, msg))
		} else {
			result = strings.HasSuffix(buf.String(), fmt.Sprintf("[%s] (%s) %s",
				level, src, msg))
		}
		buf.Reset()
		if level == DEBUG {
			_ = os.Setenv("LOG_LEVEL", cl)
		}
		return result
	}

	type args struct {
		fn     interface{}
		reqId  string
		level  string
		format string
		a      interface{}
	}

	type RedactMe struct {
		Name     string `redact:"-"`
		Password string `redact:"complete"`
		First6   string `redact:"first6"`
		Last4    string `redact:"last4"`
	}

	var obj = RedactMe{
		Name:     "eugenep",
		Password: "password01",
		First6:   "1234567890",
		Last4:    "0987654321",
	}

	var i = make([]interface{}, 1)
	i[0] = obj

	tests := []struct {
		name      string
		arg       args
		assertion assert.BoolAssertionFunc
	}{
		{"Tracef good", args{Tracef, "reqId", TRACE, "hello %+v\n", obj}, assert.True},
		{"Printf good", args{Printf, "", INFO, "hello %+v\n", obj}, assert.True},
		{"Debugf good", args{Debugf, "reqId", DEBUG, "hello %+v\n", obj}, assert.False}, // Note this checks the logging level environment variable
		{"Infof good", args{Infof, "reqId", INFO, "hello %+v\n", obj}, assert.True},
		{"Warnf good", args{Warnf, "reqId", WARN, "hello %+v\n", obj}, assert.True},
		{"Errorf good", args{Errorf, "reqId", ERROR, "hello %+v\n", obj}, assert.True},
	}

	os.Setenv("LOG_LEVEL", "TRACE")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assertion(t, isOkay(tt.arg.fn, tt.arg.reqId, tt.arg.level, tt.arg.format, tt.arg.a))
		})
	}
	os.Setenv("LOG_LEVEL", "")
}

func TestFatalf(t *testing.T) {
	type args struct {
		requestId string
		format    string
		a         []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test",
			args: args{
				requestId: "123",
				format:    "message %v",
				a:         []interface{}{"word"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					errString := err.(string)
					if errString != fmt.Sprintf(tt.args.format, tt.args.a...) {
						t.Errorf("Expected[%s] got [%s]", fmt.Sprintf(tt.args.format, tt.args.a...), errString)
					}
				}
			}()
			Fatalf(tt.args.requestId, tt.args.format, tt.args.a...)
		})
	}
}

func TestFatal(t *testing.T) {
	type args struct {
		requestId string
		message   string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test",
			args: args{
				requestId: "123",
				message:   "message %v",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					errString := err.(string)
					if errString != tt.args.message {
						t.Errorf("Expected[%s] got [%s]", tt.args.message, errString)
					}
				}
			}()
			Fatal(tt.args.requestId, tt.args.message)
		})
	}
}

func Test_first6Last4MaskFunc(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "Good",
			args: args{
				str: "12345600001234",
			},
			want: "123456****1234",
		},
		{
			name: "All",
			args: args{
				str: "1234560",
			},
			want: "*******",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := first6Last4MaskFunc(tt.args.str); got != tt.want {
				t.Errorf("first6Last4MaskFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_last4MaskFunc(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "Good",
			args: args{
				str: "12345600001234",
			},
			want: "**********1234",
		},
		{
			name: "All",
			args: args{
				str: "123",
			},
			want: "***",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := last4MaskFunc(tt.args.str); got != tt.want {
				t.Errorf("last4MaskFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_first6MaskFunc(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "Good",
			args: args{
				str: "12345600001234",
			},
			want: "123456********",
		},
		{
			name: "All",
			args: args{
				str: "12345",
			},
			want: "*****",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := first6MaskFunc(tt.args.str); got != tt.want {
				t.Errorf("first6MaskFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}
