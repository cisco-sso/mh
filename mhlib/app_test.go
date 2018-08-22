package mhlib

import (
	"fmt"
	"net"
	"testing"
)

func TestSelfRender(t *testing.T) {

	// Without datasources
	templateString := `
values:
  a: foo
  b: '[[ .values.a ]]bar' # foobar
  c: '[[ .values.b ]]baz' # foobarbaz
`
	expected := `
values:
  a: foo
  b: 'foobar' # foobar
  c: 'foobarbaz' # foobarbaz
`
	if out, err := selfRender(templateString); out != expected {
		if err != nil {
			t.Log(err)
			t.FailNow()
		}
		t.Logf("\nActual: %s\nExpected: %s\n", out, expected)
		t.FailNow()
	}

}

func TestSelfRenderDataSources(t *testing.T) {

	// With datasources
	templateString := `
gomplate:
  datasources:
  - "http_obj=https://httpbin.org/get"
  datasourceheaders: []
values:
  a: 'Func Test: [[ net.LookupIP "example.com" ]]'
`

	ips, err := net.LookupIP("example.com")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	expected := fmt.Sprintf(`
gomplate:
  datasources:
  - "http_obj=https://httpbin.org/get"
  datasourceheaders: []
values:
  a: 'Func Test: %s'
`, ips[0].String())
	if out, err := selfRender(templateString); out != expected {
		if err != nil {
			t.Log(err)
			t.FailNow()
		}
		t.Logf("\nActual: %s\nExpected: %s\n", out, expected)
		t.FailNow()
	}
}
