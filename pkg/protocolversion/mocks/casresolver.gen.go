// Code generated by counterfeiter. DO NOT EDIT.
package mocks

import (
	"net/url"
	"sync"

	"github.com/trustbloc/orb/pkg/context/common"
)

type CASResolver struct {
	ResolveStub        func(webCASURL *url.URL, cid string, data []byte) ([]byte, string, error)
	resolveMutex       sync.RWMutex
	resolveArgsForCall []struct {
		webCASURL *url.URL
		cid       string
		data      []byte
	}
	resolveReturns struct {
		result1 []byte
		result2 string
		result3 error
	}
	resolveReturnsOnCall map[int]struct {
		result1 []byte
		result2 string
		result3 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *CASResolver) Resolve(webCASURL *url.URL, cid string, data []byte) ([]byte, string, error) {
	var dataCopy []byte
	if data != nil {
		dataCopy = make([]byte, len(data))
		copy(dataCopy, data)
	}
	fake.resolveMutex.Lock()
	ret, specificReturn := fake.resolveReturnsOnCall[len(fake.resolveArgsForCall)]
	fake.resolveArgsForCall = append(fake.resolveArgsForCall, struct {
		webCASURL *url.URL
		cid       string
		data      []byte
	}{webCASURL, cid, dataCopy})
	fake.recordInvocation("Resolve", []interface{}{webCASURL, cid, dataCopy})
	fake.resolveMutex.Unlock()
	if fake.ResolveStub != nil {
		return fake.ResolveStub(webCASURL, cid, data)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	return fake.resolveReturns.result1, fake.resolveReturns.result2, fake.resolveReturns.result3
}

func (fake *CASResolver) ResolveCallCount() int {
	fake.resolveMutex.RLock()
	defer fake.resolveMutex.RUnlock()
	return len(fake.resolveArgsForCall)
}

func (fake *CASResolver) ResolveArgsForCall(i int) (*url.URL, string, []byte) {
	fake.resolveMutex.RLock()
	defer fake.resolveMutex.RUnlock()
	return fake.resolveArgsForCall[i].webCASURL, fake.resolveArgsForCall[i].cid, fake.resolveArgsForCall[i].data
}

func (fake *CASResolver) ResolveReturns(result1 []byte, result2 string, result3 error) {
	fake.ResolveStub = nil
	fake.resolveReturns = struct {
		result1 []byte
		result2 string
		result3 error
	}{result1, result2, result3}
}

func (fake *CASResolver) ResolveReturnsOnCall(i int, result1 []byte, result2 string, result3 error) {
	fake.ResolveStub = nil
	if fake.resolveReturnsOnCall == nil {
		fake.resolveReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 string
			result3 error
		})
	}
	fake.resolveReturnsOnCall[i] = struct {
		result1 []byte
		result2 string
		result3 error
	}{result1, result2, result3}
}

func (fake *CASResolver) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.resolveMutex.RLock()
	defer fake.resolveMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *CASResolver) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ common.CASResolver = new(CASResolver)
