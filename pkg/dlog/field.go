package dlog

import (
	"fmt"
	"go.uber.org/zap"
	"strings"
	"time"
)

func FieldAid(value string) Field {
	return String("aid", value)
}

func FieldMod(value string) Field {
	value = strings.Replace(value, " ", ".", -1)
	return String("mod", value)
}

func FieldAddr(value string) Field {
	return String("addr", value)
}

func FieldAddrAny(value interface{}) Field {
	return Any("addr", value)
}

func FieldName(value string) Field {
	return String("name", value)
}

func FieldType(value string) Field {
	return String("type", value)
}

func FieldCode(value int32) Field {
	return Int32("code", value)
}

func FieldCost(value time.Duration) Field {
	return String("cost", fmt.Sprintf("%.3f", float64(value.Round(time.Microsecond))/float64(time.Millisecond)))
}

func FieldKey(value string) Field {
	return String("key", value)
}

func FieldKeyAny(value interface{}) Field {
	return Any("key", value)
}

func FieldValue(value string) Field {
	return String("value", value)
}

func FieldValueAny(value interface{}) Field {
	return Any("value", value)
}

func FieldErrKind(value string) Field {
	return String("errKind", value)
}

func FieldErr(err error) Field {
	return zap.Error(err)
}

func FieldStringErr(err string) Field {
	return String("err", err)
}

func FieldExtMessage(vals ...interface{}) Field {
	return zap.Any("ext", vals)
}

func FieldStack(value []byte) Field {
	return ByteString("stack", value)
}

func FieldMethod(value string) Field {
	return String("method", value)
}

func FieldEvent(value string) Field {
	return String("event", value)
}
