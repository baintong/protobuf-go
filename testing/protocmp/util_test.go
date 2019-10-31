// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protocmp

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"google.golang.org/protobuf/internal/encoding/pack"
	testpb "google.golang.org/protobuf/internal/testprotos/test"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

func TestEqual(t *testing.T) {
	type test struct {
		x, y interface{}
		opts cmp.Options
		want bool
	}
	var tests []test

	allTypesDesc := (*testpb.TestAllTypes)(nil).ProtoReflect().Descriptor()

	// Test nil and empty messages of differing types.
	tests = append(tests, []test{{
		x:    (*testpb.TestAllTypes)(nil),
		y:    (*testpb.TestAllTypes)(nil),
		opts: cmp.Options{Transform()},
		want: true,
	}, {
		x:    (*testpb.TestAllTypes)(nil),
		y:    new(testpb.TestAllTypes),
		opts: cmp.Options{Transform()},
		want: false,
	}, {
		x:    (*testpb.TestAllTypes)(nil),
		y:    dynamicpb.NewMessage(allTypesDesc),
		opts: cmp.Options{Transform()},
		want: false,
	}, {
		x:    (*testpb.TestAllTypes)(nil),
		y:    new(testpb.TestAllTypes),
		opts: cmp.Options{Transform(), IgnoreEmptyMessages()},
		want: true,
	}, {
		x:    (*testpb.TestAllTypes)(nil),
		y:    dynamicpb.NewMessage(allTypesDesc),
		opts: cmp.Options{Transform(), IgnoreEmptyMessages()},
		want: true,
	}, {
		x:    new(testpb.TestAllTypes),
		y:    new(testpb.TestAllTypes),
		opts: cmp.Options{Transform()},
		want: true,
	}, {
		x:    new(testpb.TestAllTypes),
		y:    dynamicpb.NewMessage(allTypesDesc),
		opts: cmp.Options{Transform()},
		want: true,
	}, {
		x:    new(testpb.TestAllTypes),
		y:    new(testpb.TestAllExtensions),
		opts: cmp.Options{Transform()},
		want: false,
	}, {
		x:    struct{ I interface{} }{(*testpb.TestAllTypes)(nil)},
		y:    struct{ I interface{} }{(*testpb.TestAllTypes)(nil)},
		opts: cmp.Options{Transform()},
		want: true,
	}, {
		x:    struct{ I interface{} }{(*testpb.TestAllTypes)(nil)},
		y:    struct{ I interface{} }{new(testpb.TestAllTypes)},
		opts: cmp.Options{Transform()},
		want: false,
	}, {
		x:    struct{ I interface{} }{(*testpb.TestAllTypes)(nil)},
		y:    struct{ I interface{} }{dynamicpb.NewMessage(allTypesDesc)},
		opts: cmp.Options{Transform()},
		want: false,
	}, {
		x:    struct{ I interface{} }{(*testpb.TestAllTypes)(nil)},
		y:    struct{ I interface{} }{new(testpb.TestAllTypes)},
		opts: cmp.Options{Transform(), IgnoreEmptyMessages()},
		want: true,
	}, {
		x:    struct{ I interface{} }{(*testpb.TestAllTypes)(nil)},
		y:    struct{ I interface{} }{dynamicpb.NewMessage(allTypesDesc)},
		opts: cmp.Options{Transform(), IgnoreEmptyMessages()},
		want: true,
	}, {
		x:    struct{ I interface{} }{new(testpb.TestAllTypes)},
		y:    struct{ I interface{} }{new(testpb.TestAllTypes)},
		opts: cmp.Options{Transform()},
		want: true,
	}, {
		x:    struct{ I interface{} }{new(testpb.TestAllTypes)},
		y:    struct{ I interface{} }{dynamicpb.NewMessage(allTypesDesc)},
		opts: cmp.Options{Transform()},
		want: true,
	}, {
		x:    struct{ M proto.Message }{(*testpb.TestAllTypes)(nil)},
		y:    struct{ M proto.Message }{(*testpb.TestAllTypes)(nil)},
		opts: cmp.Options{Transform()},
		want: true,
	}, {
		x:    struct{ M proto.Message }{(*testpb.TestAllTypes)(nil)},
		y:    struct{ M proto.Message }{new(testpb.TestAllTypes)},
		opts: cmp.Options{Transform()},
		want: false,
	}, {
		x:    struct{ M proto.Message }{(*testpb.TestAllTypes)(nil)},
		y:    struct{ M proto.Message }{dynamicpb.NewMessage(allTypesDesc)},
		opts: cmp.Options{Transform()},
		want: false,
	}, {
		x:    struct{ M proto.Message }{(*testpb.TestAllTypes)(nil)},
		y:    struct{ M proto.Message }{new(testpb.TestAllTypes)},
		opts: cmp.Options{Transform(), IgnoreEmptyMessages()},
		want: true,
	}, {
		x:    struct{ M proto.Message }{(*testpb.TestAllTypes)(nil)},
		y:    struct{ M proto.Message }{dynamicpb.NewMessage(allTypesDesc)},
		opts: cmp.Options{Transform(), IgnoreEmptyMessages()},
		want: true,
	}, {
		x:    struct{ M proto.Message }{new(testpb.TestAllTypes)},
		y:    struct{ M proto.Message }{new(testpb.TestAllTypes)},
		opts: cmp.Options{Transform()},
		want: true,
	}, {
		x:    struct{ M proto.Message }{new(testpb.TestAllTypes)},
		y:    struct{ M proto.Message }{dynamicpb.NewMessage(allTypesDesc)},
		opts: cmp.Options{Transform()},
		want: true,
	}}...)

	// Test IgnoreUnknown.
	raw := pack.Message{
		pack.Tag{1, pack.BytesType}, pack.String("Hello, goodbye!"),
	}.Marshal()
	tests = append(tests, []test{{
		x:    apply(&testpb.TestAllTypes{OptionalSint64: proto.Int64(5)}, setUnknown{raw}),
		y:    &testpb.TestAllTypes{OptionalSint64: proto.Int64(5)},
		opts: cmp.Options{Transform()},
		want: false,
	}, {
		x:    apply(&testpb.TestAllTypes{OptionalSint64: proto.Int64(5)}, setUnknown{raw}),
		y:    &testpb.TestAllTypes{OptionalSint64: proto.Int64(5)},
		opts: cmp.Options{Transform(), IgnoreUnknown()},
		want: true,
	}, {
		x:    apply(&testpb.TestAllTypes{OptionalSint64: proto.Int64(5)}, setUnknown{raw}),
		y:    &testpb.TestAllTypes{OptionalSint64: proto.Int64(6)},
		opts: cmp.Options{Transform(), IgnoreUnknown()},
		want: false,
	}, {
		x:    apply(&testpb.TestAllTypes{OptionalSint64: proto.Int64(5)}, setUnknown{raw}),
		y:    apply(dynamicpb.NewMessage(allTypesDesc), setField{6, int64(5)}),
		opts: cmp.Options{Transform()},
		want: false,
	}, {
		x:    apply(&testpb.TestAllTypes{OptionalSint64: proto.Int64(5)}, setUnknown{raw}),
		y:    apply(dynamicpb.NewMessage(allTypesDesc), setField{6, int64(5)}),
		opts: cmp.Options{Transform(), IgnoreUnknown()},
		want: true,
	}}...)

	// Test IgnoreDefaultScalars.
	tests = append(tests, []test{{
		x: &testpb.TestAllTypes{
			DefaultInt32:  proto.Int32(81),
			DefaultUint32: proto.Uint32(83),
			DefaultFloat:  proto.Float32(91.5),
			DefaultBool:   proto.Bool(true),
			DefaultBytes:  []byte("world"),
		},
		y: &testpb.TestAllTypes{
			DefaultInt64:       proto.Int64(82),
			DefaultUint64:      proto.Uint64(84),
			DefaultDouble:      proto.Float64(92e3),
			DefaultString:      proto.String("hello"),
			DefaultForeignEnum: testpb.ForeignEnum_FOREIGN_BAR.Enum(),
		},
		opts: cmp.Options{Transform()},
		want: false,
	}, {
		x: &testpb.TestAllTypes{
			DefaultInt32:  proto.Int32(81),
			DefaultUint32: proto.Uint32(83),
			DefaultFloat:  proto.Float32(91.5),
			DefaultBool:   proto.Bool(true),
			DefaultBytes:  []byte("world"),
		},
		y: &testpb.TestAllTypes{
			DefaultInt64:       proto.Int64(82),
			DefaultUint64:      proto.Uint64(84),
			DefaultDouble:      proto.Float64(92e3),
			DefaultString:      proto.String("hello"),
			DefaultForeignEnum: testpb.ForeignEnum_FOREIGN_BAR.Enum(),
		},
		opts: cmp.Options{Transform(), IgnoreDefaultScalars()},
		want: true,
	}, {
		x: &testpb.TestAllTypes{
			OptionalInt32:  proto.Int32(81),
			OptionalUint32: proto.Uint32(83),
			OptionalFloat:  proto.Float32(91.5),
			OptionalBool:   proto.Bool(true),
			OptionalBytes:  []byte("world"),
		},
		y: &testpb.TestAllTypes{
			OptionalInt64:       proto.Int64(82),
			OptionalUint64:      proto.Uint64(84),
			OptionalDouble:      proto.Float64(92e3),
			OptionalString:      proto.String("hello"),
			OptionalForeignEnum: testpb.ForeignEnum_FOREIGN_BAR.Enum(),
		},
		opts: cmp.Options{Transform(), IgnoreDefaultScalars()},
		want: false,
	}, {
		x: &testpb.TestAllTypes{
			OptionalInt32:  proto.Int32(0),
			OptionalUint32: proto.Uint32(0),
			OptionalFloat:  proto.Float32(0),
			OptionalBool:   proto.Bool(false),
			OptionalBytes:  []byte(""),
		},
		y: &testpb.TestAllTypes{
			OptionalInt64:       proto.Int64(0),
			OptionalUint64:      proto.Uint64(0),
			OptionalDouble:      proto.Float64(0),
			OptionalString:      proto.String(""),
			OptionalForeignEnum: testpb.ForeignEnum_FOREIGN_FOO.Enum(),
		},
		opts: cmp.Options{Transform(), IgnoreDefaultScalars()},
		want: true,
	}, {
		x: apply(new(testpb.TestAllExtensions),
			setExtension{testpb.E_DefaultInt32Extension, int32(81)},
			setExtension{testpb.E_DefaultUint32Extension, uint32(83)},
			setExtension{testpb.E_DefaultFloatExtension, float32(91.5)},
			setExtension{testpb.E_DefaultBoolExtension, bool(true)},
			setExtension{testpb.E_DefaultBytesExtension, []byte("world")}),
		y: apply(new(testpb.TestAllExtensions),
			setExtension{testpb.E_DefaultInt64Extension, int64(82)},
			setExtension{testpb.E_DefaultUint64Extension, uint64(84)},
			setExtension{testpb.E_DefaultDoubleExtension, float64(92e3)},
			setExtension{testpb.E_DefaultStringExtension, string("hello")}),
		opts: cmp.Options{Transform()},
		want: false,
	}, {
		x: apply(new(testpb.TestAllExtensions),
			setExtension{testpb.E_DefaultInt32Extension, int32(81)},
			setExtension{testpb.E_DefaultUint32Extension, uint32(83)},
			setExtension{testpb.E_DefaultFloatExtension, float32(91.5)},
			setExtension{testpb.E_DefaultBoolExtension, bool(true)},
			setExtension{testpb.E_DefaultBytesExtension, []byte("world")}),
		y: apply(new(testpb.TestAllExtensions),
			setExtension{testpb.E_DefaultInt64Extension, int64(82)},
			setExtension{testpb.E_DefaultUint64Extension, uint64(84)},
			setExtension{testpb.E_DefaultDoubleExtension, float64(92e3)},
			setExtension{testpb.E_DefaultStringExtension, string("hello")}),
		opts: cmp.Options{Transform(), IgnoreDefaultScalars()},
		want: true,
	}, {
		x: apply(new(testpb.TestAllExtensions),
			setExtension{testpb.E_OptionalInt32Extension, int32(0)},
			setExtension{testpb.E_OptionalUint32Extension, uint32(0)},
			setExtension{testpb.E_OptionalFloatExtension, float32(0)},
			setExtension{testpb.E_OptionalBoolExtension, bool(false)},
			setExtension{testpb.E_OptionalBytesExtension, []byte("")}),
		y: apply(new(testpb.TestAllExtensions),
			setExtension{testpb.E_OptionalInt64Extension, int64(0)},
			setExtension{testpb.E_OptionalUint64Extension, uint64(0)},
			setExtension{testpb.E_OptionalDoubleExtension, float64(0)},
			setExtension{testpb.E_OptionalStringExtension, string("")}),
		opts: cmp.Options{Transform()},
		want: false,
	}, {
		x: apply(new(testpb.TestAllExtensions),
			setExtension{testpb.E_OptionalInt32Extension, int32(0)},
			setExtension{testpb.E_OptionalUint32Extension, uint32(0)},
			setExtension{testpb.E_OptionalFloatExtension, float32(0)},
			setExtension{testpb.E_OptionalBoolExtension, bool(false)},
			setExtension{testpb.E_OptionalBytesExtension, []byte("")}),
		y: apply(new(testpb.TestAllExtensions),
			setExtension{testpb.E_OptionalInt64Extension, int64(0)},
			setExtension{testpb.E_OptionalUint64Extension, uint64(0)},
			setExtension{testpb.E_OptionalDoubleExtension, float64(0)},
			setExtension{testpb.E_OptionalStringExtension, string("")}),
		opts: cmp.Options{Transform(), IgnoreDefaultScalars()},
		want: true,
	}, {
		x: &testpb.TestAllTypes{
			DefaultFloat: proto.Float32(91.6),
		},
		y:    &testpb.TestAllTypes{},
		opts: cmp.Options{Transform(), IgnoreDefaultScalars()},
		want: false,
	}, {
		x: &testpb.TestAllTypes{
			OptionalForeignMessage: &testpb.ForeignMessage{},
		},
		y:    &testpb.TestAllTypes{},
		opts: cmp.Options{Transform(), IgnoreDefaultScalars()},
		want: false,
	}}...)

	// Test IgnoreEmptyMessages.
	tests = append(tests, []test{{
		x:    []*testpb.TestAllTypes{nil, {}, {OptionalInt32: proto.Int32(5)}},
		y:    []*testpb.TestAllTypes{nil, {}, {OptionalInt32: proto.Int32(5)}},
		opts: cmp.Options{Transform()},
		want: true,
	}, {
		x:    []*testpb.TestAllTypes{nil, {}, {OptionalInt32: proto.Int32(5)}},
		y:    []*testpb.TestAllTypes{{OptionalInt32: proto.Int32(5)}},
		opts: cmp.Options{Transform(), IgnoreEmptyMessages()},
		want: false,
	}, {
		x:    &testpb.TestAllTypes{OptionalForeignMessage: &testpb.ForeignMessage{}},
		y:    &testpb.TestAllTypes{OptionalForeignMessage: nil},
		opts: cmp.Options{Transform()},
		want: false,
	}, {
		x:    &testpb.TestAllTypes{OptionalForeignMessage: &testpb.ForeignMessage{}},
		y:    &testpb.TestAllTypes{OptionalForeignMessage: nil},
		opts: cmp.Options{Transform(), IgnoreEmptyMessages()},
		want: true,
	}, {
		x:    &testpb.TestAllTypes{OptionalForeignMessage: &testpb.ForeignMessage{C: proto.Int32(5)}},
		y:    &testpb.TestAllTypes{OptionalForeignMessage: nil},
		opts: cmp.Options{Transform(), IgnoreEmptyMessages()},
		want: false,
	}, {
		x:    &testpb.TestAllTypes{RepeatedForeignMessage: []*testpb.ForeignMessage{}},
		y:    &testpb.TestAllTypes{RepeatedForeignMessage: nil},
		opts: cmp.Options{Transform()},
		want: true,
	}, {
		x:    &testpb.TestAllTypes{RepeatedForeignMessage: []*testpb.ForeignMessage{nil, {}}},
		y:    &testpb.TestAllTypes{RepeatedForeignMessage: nil},
		opts: cmp.Options{Transform()},
		want: false,
	}, {
		x:    &testpb.TestAllTypes{RepeatedForeignMessage: []*testpb.ForeignMessage{nil, {}}},
		y:    &testpb.TestAllTypes{RepeatedForeignMessage: nil},
		opts: cmp.Options{Transform(), IgnoreEmptyMessages()},
		want: true,
	}, {
		x:    &testpb.TestAllTypes{RepeatedForeignMessage: []*testpb.ForeignMessage{nil, {C: proto.Int32(5)}, {}}},
		y:    &testpb.TestAllTypes{RepeatedForeignMessage: nil},
		opts: cmp.Options{Transform(), IgnoreEmptyMessages()},
		want: false,
	}, {
		x:    &testpb.TestAllTypes{RepeatedForeignMessage: []*testpb.ForeignMessage{nil, {C: proto.Int32(5)}, {}}},
		y:    &testpb.TestAllTypes{RepeatedForeignMessage: []*testpb.ForeignMessage{{}, {}, nil, {}, {C: proto.Int32(5)}, {}}},
		opts: cmp.Options{Transform()},
		want: false,
	}, {
		x:    &testpb.TestAllTypes{RepeatedForeignMessage: []*testpb.ForeignMessage{nil, {C: proto.Int32(5)}, {}}},
		y:    &testpb.TestAllTypes{RepeatedForeignMessage: []*testpb.ForeignMessage{{}, {}, nil, {}, {C: proto.Int32(5)}, {}}},
		opts: cmp.Options{Transform(), IgnoreEmptyMessages()},
		want: true,

		// TODO
	}, {
		x:    &testpb.TestAllTypes{MapStringNestedMessage: map[string]*testpb.TestAllTypes_NestedMessage{}},
		y:    &testpb.TestAllTypes{MapStringNestedMessage: nil},
		opts: cmp.Options{Transform()},
		want: true,
	}, {
		x:    &testpb.TestAllTypes{MapStringNestedMessage: map[string]*testpb.TestAllTypes_NestedMessage{"1": nil, "2": {}}},
		y:    &testpb.TestAllTypes{MapStringNestedMessage: nil},
		opts: cmp.Options{Transform()},
		want: false,
	}, {
		x:    &testpb.TestAllTypes{MapStringNestedMessage: map[string]*testpb.TestAllTypes_NestedMessage{"1": nil, "2": {}}},
		y:    &testpb.TestAllTypes{MapStringNestedMessage: nil},
		opts: cmp.Options{Transform(), IgnoreEmptyMessages()},
		want: true,
	}, {
		x:    &testpb.TestAllTypes{MapStringNestedMessage: map[string]*testpb.TestAllTypes_NestedMessage{"1": nil, "2": {A: proto.Int32(5)}, "3": {}}},
		y:    &testpb.TestAllTypes{MapStringNestedMessage: nil},
		opts: cmp.Options{Transform(), IgnoreEmptyMessages()},
		want: false,
	}, {
		x:    &testpb.TestAllTypes{MapStringNestedMessage: map[string]*testpb.TestAllTypes_NestedMessage{"1": nil, "2": {A: proto.Int32(5)}, "3": {}}},
		y:    &testpb.TestAllTypes{MapStringNestedMessage: map[string]*testpb.TestAllTypes_NestedMessage{"1": {}, "1a": {}, "1b": nil, "2": {A: proto.Int32(5)}, "4": {}}},
		opts: cmp.Options{Transform()},
		want: false,
	}, {
		x:    &testpb.TestAllTypes{MapStringNestedMessage: map[string]*testpb.TestAllTypes_NestedMessage{"1": nil, "2": {A: proto.Int32(5)}, "3": {}}},
		y:    &testpb.TestAllTypes{MapStringNestedMessage: map[string]*testpb.TestAllTypes_NestedMessage{"1": {}, "1a": {}, "1b": nil, "2": {A: proto.Int32(5)}, "4": {}}},
		opts: cmp.Options{Transform(), IgnoreEmptyMessages()},
		want: true,
	}}...)

	// Test IgnoreEnums and IgnoreMessages.
	tests = append(tests, []test{{
		x: &testpb.TestAllTypes{
			OptionalNestedMessage:  &testpb.TestAllTypes_NestedMessage{A: proto.Int32(1)},
			RepeatedNestedMessage:  []*testpb.TestAllTypes_NestedMessage{{A: proto.Int32(2)}},
			MapStringNestedMessage: map[string]*testpb.TestAllTypes_NestedMessage{"3": {A: proto.Int32(3)}},
		},
		y:    &testpb.TestAllTypes{},
		opts: cmp.Options{Transform()},
		want: false,
	}, {
		x: &testpb.TestAllTypes{
			OptionalNestedMessage:  &testpb.TestAllTypes_NestedMessage{A: proto.Int32(1)},
			RepeatedNestedMessage:  []*testpb.TestAllTypes_NestedMessage{{A: proto.Int32(2)}},
			MapStringNestedMessage: map[string]*testpb.TestAllTypes_NestedMessage{"3": {A: proto.Int32(3)}},
		},
		y:    &testpb.TestAllTypes{},
		opts: cmp.Options{Transform(), IgnoreMessages(&testpb.TestAllTypes{})},
		want: true,
	}, {
		x: &testpb.TestAllTypes{
			OptionalNestedEnum:  testpb.TestAllTypes_FOO.Enum(),
			RepeatedNestedEnum:  []testpb.TestAllTypes_NestedEnum{testpb.TestAllTypes_BAR},
			MapStringNestedEnum: map[string]testpb.TestAllTypes_NestedEnum{"baz": testpb.TestAllTypes_BAZ},
		},
		y:    &testpb.TestAllTypes{},
		opts: cmp.Options{Transform()},
		want: false,
	}, {
		x: &testpb.TestAllTypes{
			OptionalNestedEnum:  testpb.TestAllTypes_FOO.Enum(),
			RepeatedNestedEnum:  []testpb.TestAllTypes_NestedEnum{testpb.TestAllTypes_BAR},
			MapStringNestedEnum: map[string]testpb.TestAllTypes_NestedEnum{"baz": testpb.TestAllTypes_BAZ},
		},
		y:    &testpb.TestAllTypes{},
		opts: cmp.Options{Transform(), IgnoreEnums(testpb.TestAllTypes_NestedEnum(0))},
		want: true,
	}, {
		x: &testpb.TestAllTypes{
			OptionalNestedEnum:  testpb.TestAllTypes_FOO.Enum(),
			RepeatedNestedEnum:  []testpb.TestAllTypes_NestedEnum{testpb.TestAllTypes_BAR},
			MapStringNestedEnum: map[string]testpb.TestAllTypes_NestedEnum{"baz": testpb.TestAllTypes_BAZ},

			OptionalNestedMessage:  &testpb.TestAllTypes_NestedMessage{A: proto.Int32(1)},
			RepeatedNestedMessage:  []*testpb.TestAllTypes_NestedMessage{{A: proto.Int32(2)}},
			MapStringNestedMessage: map[string]*testpb.TestAllTypes_NestedMessage{"3": {A: proto.Int32(3)}},
		},
		y: &testpb.TestAllTypes{},
		opts: cmp.Options{Transform(),
			IgnoreMessages(&testpb.TestAllExtensions{}),
			IgnoreEnums(testpb.ForeignEnum(0)),
		},
		want: false,
	}}...)

	// Test IgnoreFields and IgnoreOneofs.
	tests = append(tests, []test{{
		x:    &testpb.TestAllTypes{OptionalInt32: proto.Int32(5)},
		y:    &testpb.TestAllTypes{OptionalInt32: proto.Int32(6)},
		opts: cmp.Options{Transform()},
		want: false,
	}, {
		x: &testpb.TestAllTypes{OptionalInt32: proto.Int32(5)},
		y: &testpb.TestAllTypes{},
		opts: cmp.Options{Transform(),
			IgnoreFields(&testpb.TestAllTypes{}, "optional_int32")},
		want: true,
	}, {
		x: &testpb.TestAllTypes{OptionalInt32: proto.Int32(5)},
		y: &testpb.TestAllTypes{OptionalInt32: proto.Int32(6)},
		opts: cmp.Options{Transform(),
			IgnoreFields(&testpb.TestAllTypes{}, "optional_int32")},
		want: true,
	}, {
		x: &testpb.TestAllTypes{OptionalInt32: proto.Int32(5)},
		y: &testpb.TestAllTypes{OptionalInt32: proto.Int32(6)},
		opts: cmp.Options{Transform(),
			IgnoreFields(&testpb.TestAllTypes{}, "optional_int64")},
		want: false,
	}, {
		x:    &testpb.TestAllTypes{OneofField: &testpb.TestAllTypes_OneofUint32{5}},
		y:    &testpb.TestAllTypes{OneofField: &testpb.TestAllTypes_OneofString{"5"}},
		opts: cmp.Options{Transform()},
		want: false,
	}, {
		x: &testpb.TestAllTypes{OneofField: &testpb.TestAllTypes_OneofUint32{5}},
		y: &testpb.TestAllTypes{OneofField: &testpb.TestAllTypes_OneofString{"5"}},
		opts: cmp.Options{Transform(),
			IgnoreFields(&testpb.TestAllTypes{}, "oneof_uint32"),
			IgnoreFields(&testpb.TestAllTypes{}, "oneof_string")},
		want: true,
	}, {
		x: &testpb.TestAllTypes{OneofField: &testpb.TestAllTypes_OneofUint32{5}},
		y: &testpb.TestAllTypes{OneofField: &testpb.TestAllTypes_OneofString{"5"}},
		opts: cmp.Options{Transform(),
			IgnoreOneofs(&testpb.TestAllTypes{}, "oneof_field")},
		want: true,
	}, {
		x: apply(new(testpb.TestAllExtensions),
			setExtension{testpb.E_OptionalStringExtension, "hello"}),
		y: apply(new(testpb.TestAllExtensions),
			setExtension{testpb.E_OptionalStringExtension, "goodbye"}),
		opts: cmp.Options{Transform()},
		want: false,
	}, {
		x: apply(new(testpb.TestAllExtensions),
			setExtension{testpb.E_OptionalStringExtension, "hello"}),
		y: new(testpb.TestAllExtensions),
		opts: cmp.Options{Transform(),
			IgnoreDescriptors(testpb.E_OptionalStringExtension.TypeDescriptor())},
		want: true,
	}, {
		x: apply(new(testpb.TestAllExtensions),
			setExtension{testpb.E_OptionalStringExtension, "hello"}),
		y: apply(new(testpb.TestAllExtensions),
			setExtension{testpb.E_OptionalStringExtension, "goodbye"}),
		opts: cmp.Options{Transform(),
			IgnoreDescriptors(testpb.E_OptionalStringExtension.TypeDescriptor())},
		want: true,
	}, {
		x: apply(new(testpb.TestAllExtensions),
			setExtension{testpb.E_OptionalStringExtension, "hello"}),
		y: apply(new(testpb.TestAllExtensions),
			setExtension{testpb.E_OptionalStringExtension, "goodbye"}),
		opts: cmp.Options{Transform(),
			IgnoreDescriptors(testpb.E_OptionalInt32Extension.TypeDescriptor())},
		want: false,
	}}...)

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := cmp.Equal(tt.x, tt.y, tt.opts)
			if got != tt.want {
				if !got {
					t.Errorf("cmp.Equal = false, want true; diff:\n%v", cmp.Diff(tt.x, tt.y, tt.opts))
				} else {
					t.Errorf("cmp.Equal = true, want false")
				}
			}
		})
	}
}

type setField struct {
	num protoreflect.FieldNumber
	val interface{}
}
type setUnknown struct {
	raw protoreflect.RawFields
}
type setExtension struct {
	typ protoreflect.ExtensionType
	val interface{}
}

// apply applies a sequence of mutating operations to m.
func apply(m proto.Message, ops ...interface{}) proto.Message {
	mr := m.ProtoReflect()
	md := mr.Descriptor()
	for _, op := range ops {
		switch op := op.(type) {
		case setField:
			fd := md.Fields().ByNumber(op.num)
			mr.Set(fd, protoreflect.ValueOf(op.val))
		case setUnknown:
			mr.SetUnknown(op.raw)
		case setExtension:
			mr.Set(op.typ.TypeDescriptor(), protoreflect.ValueOf(op.val))
		}
	}
	return m
}
