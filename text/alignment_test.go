package text

import (
	// "fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testCases = []*testCase{}
)

type testCase struct {
	str string
	res string
}

func appendTestCase(str, res string) {
	testCases = append(testCases, &testCase{str, res})
}

func TestAlignment(t *testing.T) {
	// case 1: no '\n'
	appendTestCase(`hello sdada
li sss

kkkkkkkk    ddk`, `>>> hello sdada
>>> li    sss

>>> kkkkkkkk ddk`)
	// case 2: 1个'\n'
	appendTestCase(`hello sdada
li sss

kkkkkkkk    ddk
`, `>>> hello sdada
>>> li    sss

>>> kkkkkkkk ddk
`)
	// case 3: 2个'\n'
	appendTestCase(`hello sdada
li sss

kkkkkkkk    ddk

`, `>>> hello sdada
>>> li    sss

>>> kkkkkkkk ddk

`)
	// case 4: 3个'\n'
	appendTestCase(`hello sdada
li sss

kkkkkkkk    ddk


`, `>>> hello sdada
>>> li    sss

>>> kkkkkkkk ddk


`)
	// case 5: 4个'\n'
	appendTestCase(`hello sdada
li sss

kkkkkkkk    ddk



`, `>>> hello sdada
>>> li    sss

>>> kkkkkkkk ddk



`)
	// case 6
	appendTestCase(`1 22
wwwww ee
ewrewr2d 221`, `>>> 1        22
>>> wwwww    ee
>>> ewrewr2d 221`)

	assert := assert.New(t)

	align := NewAlignment(" ", ">>> ")
	assert.NotNil(align)

	for _, input := range testCases {
		assert.Equal(input.res, align.Format(input.str))
	}
}
