// Code generated by counterfeiter. DO NOT EDIT.
package internal

import (
	"com/alexander/debendency/pkg/commands"
	"sync"
)

type FakeCommand struct {
	CommandStub        func(string, ...string) (string, int, error)
	commandMutex       sync.RWMutex
	commandArgsForCall []struct {
		arg1 string
		arg2 []string
	}
	commandReturns struct {
		result1 string
		result2 int
		result3 error
	}
	commandReturnsOnCall map[int]struct {
		result1 string
		result2 int
		result3 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeCommand) Command(arg1 string, arg2 ...string) (string, int, error) {
	fake.commandMutex.Lock()
	ret, specificReturn := fake.commandReturnsOnCall[len(fake.commandArgsForCall)]
	fake.commandArgsForCall = append(fake.commandArgsForCall, struct {
		arg1 string
		arg2 []string
	}{arg1, arg2})
	stub := fake.CommandStub
	fakeReturns := fake.commandReturns
	fake.recordInvocation("Command", []interface{}{arg1, arg2})
	fake.commandMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2...)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	return fakeReturns.result1, fakeReturns.result2, fakeReturns.result3
}

func (fake *FakeCommand) CommandCallCount() int {
	fake.commandMutex.RLock()
	defer fake.commandMutex.RUnlock()
	return len(fake.commandArgsForCall)
}

func (fake *FakeCommand) CommandCalls(stub func(string, ...string) (string, int, error)) {
	fake.commandMutex.Lock()
	defer fake.commandMutex.Unlock()
	fake.CommandStub = stub
}

func (fake *FakeCommand) CommandArgsForCall(i int) (string, []string) {
	fake.commandMutex.RLock()
	defer fake.commandMutex.RUnlock()
	argsForCall := fake.commandArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeCommand) CommandReturns(result1 string, result2 int, result3 error) {
	fake.commandMutex.Lock()
	defer fake.commandMutex.Unlock()
	fake.CommandStub = nil
	fake.commandReturns = struct {
		result1 string
		result2 int
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeCommand) CommandReturnsOnCall(i int, result1 string, result2 int, result3 error) {
	fake.commandMutex.Lock()
	defer fake.commandMutex.Unlock()
	fake.CommandStub = nil
	if fake.commandReturnsOnCall == nil {
		fake.commandReturnsOnCall = make(map[int]struct {
			result1 string
			result2 int
			result3 error
		})
	}
	fake.commandReturnsOnCall[i] = struct {
		result1 string
		result2 int
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeCommand) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.commandMutex.RLock()
	defer fake.commandMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeCommand) recordInvocation(key string, args []interface{}) {
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

var _ commands.Command = new(FakeCommand)