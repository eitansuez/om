// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"sync"
)

type DownloadProductService struct {
	CheckProductAvailabilityStub        func(string, string) (bool, error)
	checkProductAvailabilityMutex       sync.RWMutex
	checkProductAvailabilityArgsForCall []struct {
		arg1 string
		arg2 string
	}
	checkProductAvailabilityReturns struct {
		result1 bool
		result2 error
	}
	checkProductAvailabilityReturnsOnCall map[int]struct {
		result1 bool
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *DownloadProductService) CheckProductAvailability(arg1 string, arg2 string) (bool, error) {
	fake.checkProductAvailabilityMutex.Lock()
	ret, specificReturn := fake.checkProductAvailabilityReturnsOnCall[len(fake.checkProductAvailabilityArgsForCall)]
	fake.checkProductAvailabilityArgsForCall = append(fake.checkProductAvailabilityArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("CheckProductAvailability", []interface{}{arg1, arg2})
	fake.checkProductAvailabilityMutex.Unlock()
	if fake.CheckProductAvailabilityStub != nil {
		return fake.CheckProductAvailabilityStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.checkProductAvailabilityReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *DownloadProductService) CheckProductAvailabilityCallCount() int {
	fake.checkProductAvailabilityMutex.RLock()
	defer fake.checkProductAvailabilityMutex.RUnlock()
	return len(fake.checkProductAvailabilityArgsForCall)
}

func (fake *DownloadProductService) CheckProductAvailabilityCalls(stub func(string, string) (bool, error)) {
	fake.checkProductAvailabilityMutex.Lock()
	defer fake.checkProductAvailabilityMutex.Unlock()
	fake.CheckProductAvailabilityStub = stub
}

func (fake *DownloadProductService) CheckProductAvailabilityArgsForCall(i int) (string, string) {
	fake.checkProductAvailabilityMutex.RLock()
	defer fake.checkProductAvailabilityMutex.RUnlock()
	argsForCall := fake.checkProductAvailabilityArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *DownloadProductService) CheckProductAvailabilityReturns(result1 bool, result2 error) {
	fake.checkProductAvailabilityMutex.Lock()
	defer fake.checkProductAvailabilityMutex.Unlock()
	fake.CheckProductAvailabilityStub = nil
	fake.checkProductAvailabilityReturns = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *DownloadProductService) CheckProductAvailabilityReturnsOnCall(i int, result1 bool, result2 error) {
	fake.checkProductAvailabilityMutex.Lock()
	defer fake.checkProductAvailabilityMutex.Unlock()
	fake.CheckProductAvailabilityStub = nil
	if fake.checkProductAvailabilityReturnsOnCall == nil {
		fake.checkProductAvailabilityReturnsOnCall = make(map[int]struct {
			result1 bool
			result2 error
		})
	}
	fake.checkProductAvailabilityReturnsOnCall[i] = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *DownloadProductService) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.checkProductAvailabilityMutex.RLock()
	defer fake.checkProductAvailabilityMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *DownloadProductService) recordInvocation(key string, args []interface{}) {
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