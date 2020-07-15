// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"sync"

	"github.com/pivotal-cf/om/download_clients"
)

type FileArtifacter struct {
	NameStub        func() string
	nameMutex       sync.RWMutex
	nameArgsForCall []struct {
	}
	nameReturns struct {
		result1 string
	}
	nameReturnsOnCall map[int]struct {
		result1 string
	}
	SHA256Stub        func() string
	sHA256Mutex       sync.RWMutex
	sHA256ArgsForCall []struct {
	}
	sHA256Returns struct {
		result1 string
	}
	sHA256ReturnsOnCall map[int]struct {
		result1 string
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FileArtifacter) Name() string {
	fake.nameMutex.Lock()
	ret, specificReturn := fake.nameReturnsOnCall[len(fake.nameArgsForCall)]
	fake.nameArgsForCall = append(fake.nameArgsForCall, struct {
	}{})
	fake.recordInvocation("Name", []interface{}{})
	fake.nameMutex.Unlock()
	if fake.NameStub != nil {
		return fake.NameStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.nameReturns
	return fakeReturns.result1
}

func (fake *FileArtifacter) NameCallCount() int {
	fake.nameMutex.RLock()
	defer fake.nameMutex.RUnlock()
	return len(fake.nameArgsForCall)
}

func (fake *FileArtifacter) NameCalls(stub func() string) {
	fake.nameMutex.Lock()
	defer fake.nameMutex.Unlock()
	fake.NameStub = stub
}

func (fake *FileArtifacter) NameReturns(result1 string) {
	fake.nameMutex.Lock()
	defer fake.nameMutex.Unlock()
	fake.NameStub = nil
	fake.nameReturns = struct {
		result1 string
	}{result1}
}

func (fake *FileArtifacter) NameReturnsOnCall(i int, result1 string) {
	fake.nameMutex.Lock()
	defer fake.nameMutex.Unlock()
	fake.NameStub = nil
	if fake.nameReturnsOnCall == nil {
		fake.nameReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.nameReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *FileArtifacter) SHA256() string {
	fake.sHA256Mutex.Lock()
	ret, specificReturn := fake.sHA256ReturnsOnCall[len(fake.sHA256ArgsForCall)]
	fake.sHA256ArgsForCall = append(fake.sHA256ArgsForCall, struct {
	}{})
	fake.recordInvocation("SHA256", []interface{}{})
	fake.sHA256Mutex.Unlock()
	if fake.SHA256Stub != nil {
		return fake.SHA256Stub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.sHA256Returns
	return fakeReturns.result1
}

func (fake *FileArtifacter) SHA256CallCount() int {
	fake.sHA256Mutex.RLock()
	defer fake.sHA256Mutex.RUnlock()
	return len(fake.sHA256ArgsForCall)
}

func (fake *FileArtifacter) SHA256Calls(stub func() string) {
	fake.sHA256Mutex.Lock()
	defer fake.sHA256Mutex.Unlock()
	fake.SHA256Stub = stub
}

func (fake *FileArtifacter) SHA256Returns(result1 string) {
	fake.sHA256Mutex.Lock()
	defer fake.sHA256Mutex.Unlock()
	fake.SHA256Stub = nil
	fake.sHA256Returns = struct {
		result1 string
	}{result1}
}

func (fake *FileArtifacter) SHA256ReturnsOnCall(i int, result1 string) {
	fake.sHA256Mutex.Lock()
	defer fake.sHA256Mutex.Unlock()
	fake.SHA256Stub = nil
	if fake.sHA256ReturnsOnCall == nil {
		fake.sHA256ReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.sHA256ReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *FileArtifacter) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.nameMutex.RLock()
	defer fake.nameMutex.RUnlock()
	fake.sHA256Mutex.RLock()
	defer fake.sHA256Mutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FileArtifacter) recordInvocation(key string, args []interface{}) {
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

var _ download_clients.FileArtifacter = new(FileArtifacter)