// k8s
package k8s

import (
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
	_init_ "github.com/mss-boot-io/cd-template/imports/k8s/jsii"
)

type Quantity interface {
	Value() interface{}
}

// The jsii proxy struct for Quantity
type jsiiProxy_Quantity struct {
	_ byte // padding
}

func (j *jsiiProxy_Quantity) Value() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"value",
		&returns,
	)
	return returns
}

func Quantity_FromNumber(value *float64) Quantity {
	_init_.Initialize()

	if err := validateQuantity_FromNumberParameters(value); err != nil {
		panic(err)
	}
	var returns Quantity

	_jsii_.StaticInvoke(
		"k8s.Quantity",
		"fromNumber",
		[]interface{}{value},
		&returns,
	)

	return returns
}

func Quantity_FromString(value *string) Quantity {
	_init_.Initialize()

	if err := validateQuantity_FromStringParameters(value); err != nil {
		panic(err)
	}
	var returns Quantity

	_jsii_.StaticInvoke(
		"k8s.Quantity",
		"fromString",
		[]interface{}{value},
		&returns,
	)

	return returns
}
