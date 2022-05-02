package yrpc

import (
	"errors"
	"fmt"
	rpc2 "github.com/AnyISalIn/yrpc/rpc"
	"reflect"
)

func suitableStreamMethods(typ reflect.Type, revr reflect.Value) (map[string]rpc2.StreamHandler, error) {
	methods := make(map[string]rpc2.StreamHandler)

	sname := reflect.Indirect(revr).Type().Name()

	if sname == "" {
		s := "rpc.Register: no service name for type " + typ.String()
		return nil, errors.New(s)
	}

	sh := reflect.TypeOf(rpc2.StreamHandler(nil))

	for m := 0; m < typ.NumMethod(); m++ {
		tmethod := typ.Method(m)
		vmethod := revr.Method(m)
		mname := tmethod.Name
		// Method must be exported.
		if !tmethod.IsExported() {
			continue
		}

		if vmethod.Type().AssignableTo(sh) {
			methods[fmt.Sprintf("%s.%s", sname, mname)] = vmethod.Convert(sh).Interface().(rpc2.StreamHandler)
		}
	}

	return methods, nil
}
