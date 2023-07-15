package gohttp

import (
	"fmt"
	"testing"
)

var (
	uris                   = []string{"/", "/test/{string}", "/test", "/test/{float}/{int}"}
	uriShouldBeParametized = []bool{false, true, false, true}
	workingTestUris        = []string{"/", "/test/aaa", "/test", "/test/109.111/123"}
	badTestUris            = []string{"/abc", "/test/", "/test/sdf", "/test/test/aaaaa"}
)

func TestCompileUri(t *testing.T) {
	for index, _ := range uris {
		_uri := compileUri(uris[index])

		fmt.Printf("\turi compiles correctly %s (templated? %t)...", _uri.fullUri, _uri.hasTemplatedParams)

		if _uri.hasTemplatedParams != uriShouldBeParametized[index] {
			t.Error(t.Name(), "_uri should have a templated param but doesn't")
			t.FailNow()
		} else {
			fmt.Println("It compiled!")
		}
	}
}

func TestUriMatches(t *testing.T) {
	for index, _ := range uris {
		_uri := compileUri(uris[index])
		_targetUri := compileUri(workingTestUris[index])

		fmt.Printf("\turi match %s and %s...", _uri.fullUri, _targetUri.fullUri)

		if !_uri.uriMatches(&_targetUri) {
			t.Error(t.Name(), "_targetUri doesn't match _uri (but it should)")
			t.Error(_uri.fullUri, _targetUri.fullUri)
			t.FailNow()
		} else {
			fmt.Println("They match!")
		}
	}
}

func TestUriDontMatches(t *testing.T) {
	for index, _ := range uris {
		_uri := compileUri(uris[index])
		_targetUri := compileUri(badTestUris[index])

		fmt.Printf("\turi don't match %s and %s...", _uri.fullUri, _targetUri.fullUri)

		if _uri.uriMatches(&_targetUri) {
			t.Error(t.Name(), "_targetUri matches _uri (but it shouldn't)")
			t.Error()
			t.FailNow()
		} else {
			fmt.Println("They don't match!")
		}
	}
}
