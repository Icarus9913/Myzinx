package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Person struct {
	ID   uint32
	Age  int32
}

func main() {
	var p Person = Person{
		ID: 9527,
		Age: 18,
	}

	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.LittleEndian, p.ID)
	if nil!=err{
		panic("编码失败1"+err.Error())
	}

	err = binary.Write(buf, binary.LittleEndian, p.Age)
	if nil!=err{
		panic("编码失败2"+err.Error())
	}

	///////////////////////////////////////////////////////////////////////////

	dataBuf := new(bytes.Buffer)
	_, err = dataBuf.ReadFrom(buf)
	if nil!=err{
		panic(err)
	}

	var ID1 uint32
	var Age1 int32

	err = binary.Read(dataBuf, binary.LittleEndian, &ID1)
	if nil!=err{
		panic(err)
	}

	err = binary.Read(dataBuf, binary.LittleEndian, &Age1)
	if nil!=err{
		panic(err)
	}

	fmt.Println(ID1)
	fmt.Println(Age1)
}
