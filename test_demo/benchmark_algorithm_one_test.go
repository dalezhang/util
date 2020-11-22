package test_demo

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)
func TestReadFull(t *testing.T) {
	data := []byte("sacasstringeqw")
	input := bytes.NewBuffer(data)
	find := []byte("string")
	// The number of bytes we are looking for.
	size := len(find)
	// Declare the buffers we need to process the stream.
	buf := make([]byte, size)
	end := size - 1
	n, _ := io.ReadFull(input, buf[:end])
	fmt.Println("n = ", n)
	fmt.Println("buf = ", string(buf))
	n, _ = io.ReadFull(input, buf[end:])
	fmt.Println("buf = ", string(buf))
	fmt.Println("buf[0]", string(buf[0]))
	// Slice that front byte out.
	copy(buf, buf[1:])
	fmt.Println("buf = ", string(buf))

}

// go test -run none -bench AlgorithmOne -benchtime 3s -benchmem

func BenchmarkAlgorithmOne(b *testing.B) {
	var output bytes.Buffer
	in := []byte("his name is elvis, super star.")
	find := []byte("elvis")
	repl := []byte("Elvis")

	b.ResetTimer()
	fmt.Println("N = ",b.N)

	for i := 0; i < b.N; i++ {
		output.Reset()
		algOne(in, find, repl, &output)
	}
}
