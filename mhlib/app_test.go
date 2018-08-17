package mhlib

import (
	"testing"
)


func TestSelfRender(t *testing.T) {

	// With datasources
	templateString := `
gomplate:
  datasources:
  - "http_obj=https://httpbin.org/get"
  datasourceheaders: []
values:
  a: foo
  b: '[[ .values.a ]]bar' # foobar
  c: '[[ .values.b ]]baz' # foobarbaz
  f: 'Func Test: [[ net.LookupIP "example.com" ]]'
`
	expected := `
gomplate:
  datasources:
  - "http_obj=https://httpbin.org/get"
  datasourceheaders: []
values:
  a: foo
  b: 'foobar' # foobar
  c: 'foobarbaz' # foobarbaz
  f: 'Func Test: 10.114.236.103'
`
	out, err := selfRender(templateString, true)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	if out != expected {
		t.Logf("Actual: %s\nExpected: %s\n", out, expected)
		t.Fail()
	}

	// Without datasources
	templateString = `
values:
  a: foo
  b: '[[ .values.a ]]bar' # foobar
  c: '[[ .values.b ]]baz' # foobarbaz
`
	expected = `
values:
  a: foo
  b: 'foobar' # foobar
  c: 'foobarbaz' # foobarbaz
`
	out, err = selfRender(templateString, true)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	if out != expected {
		t.Logf("Actual: %s\nExpected: %s\n", out, expected)
		t.Fail()
	}


}
