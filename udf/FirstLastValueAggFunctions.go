package udf

/*
  add the following udfInfo array into job definition, you can use the following 6 aggregate functions in your SQL:
  1, FirstValueLong(aLongColumn)
  2, LastValueLong(aLongColumn)
  3, FirstValueDouble(aDoubleColumn)
  4, LastValueDouble(aDoubleColumn)
  5, FirstValueString(aStringColumn)
  6, LastValueString(aStringColumn)

"udfInfo": [{
	"name": "FirstLastValueAggFunctions.go",
	"content": "H4sIAAAAAAAA/81XXW+bMBR9Dr/ijocJpo6k0rSHSp2Ubtpeuk3rpL1UfXCIYXQER8Z0H1X/+2yDA7GJkzRxkrwgLvf6Hl+f46PMUfwLpRiqaeJ52WxOKIPAG/i4iMk0K9JhSia+eKeU0NL3Qs9LqiKGrMhYEMKjN+AJ0Q1Os5JhGrwkk/tvFa7w41PoPXke+zvHoGJQMlrFTBRNEUNwe5cVvChBMU8X2XLlgMArVRFCitnHjJaiVydZLJElkOMiIJFYK4R3MBLRAcWsogXU4dvRnTfgC6tokeWr+1yjvjZl9g/DxWW3l+wt4/1NxafX55t2RnFczaocMRw8oLzbXo63XhMuAc3nuJg2GM6A54YrVuQtKYpZ73LDIdxXJeM5M/KAgf3EkIj5AinwhptdIGq2e34ht7oKS4kboizqVs9ihmmKg0k3ZBnBRD6jKJKDGO7y49UgefYD5RW+JkVaU3c51hJ4AT/RUkL4yjEGMfujj95SMm4pwNmQdSvlQXPp6OfI88RBifQoaMclT6uueCEHLU+MZ0XLNJOpTzZMNw2JdgYkmncYKWKWtp9w/Wb0DSHgb2/fnIG8iuw9eb/mQ9ReIHI24ktnMo06+SOqlw/PJD35bNS3UdOxjL7g34HvW+F/lgQWgMbLQxOhK3NoCv54Gf9Exa+0WUaNQqwg3lPMj3lBKkLNW63ZW+e6tnOBi7i7nnk46xkg7wEr7JyUOGiksrOaxW2uiXkp1KvlpYyNpKxVnIKSNUgHErLW1YWOa4N2IWMN/FFUrGHYWcQGDVxoWAe9Vwm3F8QHUk1yrFtyHV1jynXSlrasik5Bziaqg1uzamwTdZIT5NSemwbPNmi1iSNbtIKxR5NueeHWphfQ3Rh1V+Ra0G7WW0jcqDkFhRugDm3ZbvW9mW0/T97GFo5r3XsTdw8lnNq3G2m3V8d3RjPzP3UdXWPgddKWBq6KTkHeJqqDG7hqbBN4KXPc+Xe9vqlv39/Mv9UejuzfCsYe/bulhVv/XkB3499djWtBu39voXCj5hQEboA6tH87lfdm9v08dRs7OK59703bPYxwat99yv4PTrFQWtoZAAA="
}]

  Note: content = base64(gzip(ThisSourceFile)), this compression can also be done at http://www.txtwizard.net/compression
*/

import (
	"encoding/gob"
	"errors"
)

func init() {
	gob.Register(&objQueue{})
}

type objQueue struct {
	data []interface{}
}

func (o *objQueue) getFirst() interface{} {
	if len(o.data) > 0 {
		return o.data[0]
	}

	return nil
}

func (o *objQueue) getLast() interface{} {
	size := len(o.data)
	if size > 0 {
		return o.data[size-1]
	}

	return nil
}

func (o *objQueue) accumulate(val interface{}) {
	o.data = append(o.data, val)
}

func (o *objQueue) retract(val interface{}) {
	// just remove the first one
	size := len(o.data)
	if size > 0 {
		o.data = o.data[1:]
	}
}

func (o *objQueue) reset() {
	o.data = nil
}

func (o *objQueue) merge(b *objQueue) {
	o.data = append(o.data, b.data...)
}

/////////////////////////////////////////////////////////////
// FirstValueLong
type FirstValueLong struct {
}

func (f FirstValueLong) Open(ctx interface{}) {
}

func (f FirstValueLong) Accumulate(acci interface{}, value interface{}) {
	acc := acci.(*objQueue)
	if value != nil {
		acc.accumulate(value)
	}
}

func (f FirstValueLong) Retract(acci interface{}, value interface{}) {
	acc := acci.(*objQueue)
	acc.retract(value)
}

func (f FirstValueLong) GetValue(acci interface{}) (int64, error) {
	acc := acci.(*objQueue)
	ret := acc.getFirst()
	if ret != nil {
		return ret.(int64), nil
	}
	return 0, errors.New("")
}

func (f FirstValueLong) Merge(acciA interface{}, acciB interface{}) {
	a := acciA.(*objQueue)
	b := acciB.(*objQueue)
	a.merge(b)
}

func (f FirstValueLong) CreateAccumulator() interface{} {
	return &objQueue{}
}

func (f FirstValueLong) ResetAccumulator(acci interface{}) {
	acc := acci.(*objQueue)
	acc.reset()
}

func (f FirstValueLong) Close() {
}

/////////////////////////////////////////////////////////////
// LastValueLong
type LastValueLong struct {
}

func (f LastValueLong) Open(ctx interface{}) {
}

func (f LastValueLong) Accumulate(acci interface{}, value interface{}) {
	acc := acci.(*objQueue)
	if value != nil {
		acc.accumulate(value)
	}
}

func (f LastValueLong) Retract(acci interface{}, value interface{}) {
	acc := acci.(*objQueue)
	acc.retract(value)
}

func (f LastValueLong) GetValue(acci interface{}) (int64, error) {
	acc := acci.(*objQueue)
	ret := acc.getLast()
	if ret != nil {
		return ret.(int64), nil
	}
	return 0, errors.New("")
}

func (f LastValueLong) Merge(acciA interface{}, acciB interface{}) {
	a := acciA.(*objQueue)
	b := acciB.(*objQueue)
	a.merge(b)
}

func (f LastValueLong) CreateAccumulator() interface{} {
	return &objQueue{}
}

func (f LastValueLong) ResetAccumulator(acci interface{}) {
	acc := acci.(*objQueue)
	acc.reset()
}

func (f LastValueLong) Close() {
}

/////////////////////////////////////////////////////////////
// FirstValueDouble
type FirstValueDouble struct {
}

func (f FirstValueDouble) Open(ctx interface{}) {
}

func (f FirstValueDouble) Accumulate(acci interface{}, value interface{}) {
	acc := acci.(*objQueue)
	if value != nil {
		acc.accumulate(value)
	}
}

func (f FirstValueDouble) Retract(acci interface{}, value interface{}) {
	acc := acci.(*objQueue)
	acc.retract(value)
}

func (f FirstValueDouble) GetValue(acci interface{}) (float64, error) {
	acc := acci.(*objQueue)
	ret := acc.getFirst()
	if ret != nil {
		return ret.(float64), nil
	}
	return 0, errors.New("")
}

func (f FirstValueDouble) Merge(acciA interface{}, acciB interface{}) {
	a := acciA.(*objQueue)
	b := acciB.(*objQueue)
	a.merge(b)
}

func (f FirstValueDouble) CreateAccumulator() interface{} {
	return &objQueue{}
}

func (f FirstValueDouble) ResetAccumulator(acci interface{}) {
	acc := acci.(*objQueue)
	acc.reset()
}

func (f FirstValueDouble) Close() {
}

/////////////////////////////////////////////////////////////
// LastValueDouble
type LastValueDouble struct {
}

func (f LastValueDouble) Open(ctx interface{}) {
}

func (f LastValueDouble) Accumulate(acci interface{}, value interface{}) {
	acc := acci.(*objQueue)
	if value != nil {
		acc.accumulate(value)
	}
}

func (f LastValueDouble) Retract(acci interface{}, value interface{}) {
	acc := acci.(*objQueue)
	acc.retract(value)
}

func (f LastValueDouble) GetValue(acci interface{}) (float64, error) {
	acc := acci.(*objQueue)
	ret := acc.getLast()
	if ret != nil {
		return ret.(float64), nil
	}
	return 0, errors.New("")
}

func (f LastValueDouble) Merge(acciA interface{}, acciB interface{}) {
	a := acciA.(*objQueue)
	b := acciB.(*objQueue)
	a.merge(b)
}

func (f LastValueDouble) CreateAccumulator() interface{} {
	return &objQueue{}
}

func (f LastValueDouble) ResetAccumulator(acci interface{}) {
	acc := acci.(*objQueue)
	acc.reset()
}

func (f LastValueDouble) Close() {
}

/////////////////////////////////////////////////////////////
// FirstValueString
type FirstValueString struct {
}

func (f FirstValueString) Open(ctx interface{}) {
}

func (f FirstValueString) Accumulate(acci interface{}, value interface{}) {
	acc := acci.(*objQueue)
	if value != nil {
		acc.accumulate(value)
	}
}

func (f FirstValueString) Retract(acci interface{}, value interface{}) {
	acc := acci.(*objQueue)
	acc.retract(value)
}

func (f FirstValueString) GetValue(acci interface{}) (string, error) {
	acc := acci.(*objQueue)
	ret := acc.getFirst()
	if ret != nil {
		return ret.(string), nil
	}
	return "", errors.New("")
}

func (f FirstValueString) Merge(acciA interface{}, acciB interface{}) {
	a := acciA.(*objQueue)
	b := acciB.(*objQueue)
	a.merge(b)
}

func (f FirstValueString) CreateAccumulator() interface{} {
	return &objQueue{}
}

func (f FirstValueString) ResetAccumulator(acci interface{}) {
	acc := acci.(*objQueue)
	acc.reset()
}

func (f FirstValueString) Close() {
}

/////////////////////////////////////////////////////////////
// LastValueString
type LastValueString struct {
}

func (f LastValueString) Open(ctx interface{}) {
}

func (f LastValueString) Accumulate(acci interface{}, value interface{}) {
	acc := acci.(*objQueue)
	if value != nil {
		acc.accumulate(value)
	}
}

func (f LastValueString) Retract(acci interface{}, value interface{}) {
	acc := acci.(*objQueue)
	acc.retract(value)
}

func (f LastValueString) GetValue(acci interface{}) (string, error) {
	acc := acci.(*objQueue)
	ret := acc.getLast()
	if ret != nil {
		return ret.(string), nil
	}
	return "", errors.New("")
}

func (f LastValueString) Merge(acciA interface{}, acciB interface{}) {
	a := acciA.(*objQueue)
	b := acciB.(*objQueue)
	a.merge(b)
}

func (f LastValueString) CreateAccumulator() interface{} {
	return &objQueue{}
}

func (f LastValueString) ResetAccumulator(acci interface{}) {
	acc := acci.(*objQueue)
	acc.reset()
}

func (f LastValueString) Close() {
}
