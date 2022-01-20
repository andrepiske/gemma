package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
)

type Reader struct {
	src *bufio.Reader
}

type RValue interface{}

func NewReader(s io.Reader) *Reader {
	return &Reader{src: bufio.NewReader(s)}
}

func (r *Reader) ReadValue() RValue {
	p, _ := r.src.Peek(1)

	if p[0] == '*' {
		fmt.Println("will read an Array")
		s, _ := r.src.ReadString('\n')

		arrayLen, _ := strconv.Atoi(s[1 : len(s)-2])

		fmt.Printf("array has %d elements\n", arrayLen)

		result := make([]RValue, arrayLen)

		for i := 0; i < arrayLen; i++ {
			// read element of array
			v := r.ReadValue()
			result[i] = v
		}

		return result

	} else if p[0] == '+' {
		fmt.Println("will read a Simple String")

		s, _ := r.src.ReadString('\n')
		v := s[:len(s)-1]
		fmt.Printf("Read simple string '%s'\n", v)

		return v

	} else if p[0] == '-' {
		fmt.Println("will read an Error")
		fmt.Println("not implemented")
		panic(0)
	} else if p[0] == ':' {
		fmt.Println("will read an Integer")
		fmt.Println("not implemented")
		panic(0)
	} else if p[0] == '$' {
		fmt.Println("will read a complex String")
		s, _ := r.src.ReadString('\n')
		strLen, _ := strconv.Atoi(s[1 : len(s)-2])
		buf := make([]byte, strLen)
		r.src.Read(buf)
		r.src.ReadByte() // \r
		r.src.ReadByte() // \n
		v := string(buf)
		fmt.Printf("string has length %d and it is '%s'\n", strLen, v)
		return v
	} else {
		fmt.Printf("wuuuut is '%v'?\n", p)
		panic(p)
	}

	return nil
}

func handleClient(conn net.Conn) {
	r := NewReader(conn)
	r.ReadValue()

	conn.Close()
}

func main() {
	ln, _ := net.Listen("tcp", ":7111")

	for {
		fmt.Println("waiting new connection...")
		conn, _ := ln.Accept()
		go handleClient(conn)
	}
}
