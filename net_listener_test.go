package shmipc

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net"
	"sync"
	"testing"
	"time"
)

func TestNewListenerAndClose(t *testing.T) {
	rawListener, err := net.Listen("unix", "/tmp/test.sock")
	assert.NoError(t, err)

	listener := newListener(rawListener, defaultBacklog)
	assert.NotNil(t, listener)

	err = listener.Close()
	assert.NoError(t, err)
}

func TestListenerAcceptAndClose(t *testing.T) {
	listener, err := ListenWithBacklog("/tmp/test2.sock", defaultBacklog)
	assert.NoError(t, err)
	defer listener.Close()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		conn, err := listener.Accept()
		assert.NoError(t, err)
		assert.NotNil(t, conn)
		conn.Close()
	}()

	conn, err := net.Dial("unix", "/tmp/test2.sock")
	assert.NoError(t, err)
	conn.Close()

	wg.Wait()
}

func TestStreamWrapperReadWrite(t *testing.T) {
	listener, err := ListenWithBacklog("/tmp/test3.sock", defaultBacklog)
	assert.NoError(t, err)
	defer listener.Close()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		conn, err := listener.Accept()
		assert.NoError(t, err)
		defer conn.Close()

		data := make([]byte, 5)
		n, err := conn.Read(data)
		assert.NoError(t, err)
		assert.Equal(t, 5, n)
		assert.Equal(t, []byte("hello"), data)
	}()

	conn, err := net.Dial("unix", "/tmp/test3.sock")
	assert.NoError(t, err)
	defer conn.Close()

	n, err := conn.Write([]byte("hello"))
	assert.NoError(t, err)
	assert.Equal(t, 5, n)

	wg.Wait()
}

func TestStreamWrapperDeadline(t *testing.T) {
	listener, err := ListenWithBacklog("/tmp/test4.sock", defaultBacklog)
	assert.NoError(t, err)
	defer listener.Close()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		conn, err := listener.Accept()
		assert.NoError(t, err)
		defer conn.Close()

		conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))

		data, err := ioutil.ReadAll(conn)
		assert.Error(t, err)
		assert.True(t, len(data) > 0)
	}()

	conn, err := net.Dial("unix", "/tmp/test4.sock")
	assert.NoError(t, err)
	defer conn.Close()

	time.Sleep(200 * time.Millisecond)
	n, err := conn.Write(bytes.Repeat([]byte("a"), 4096))
	assert.NoError(t, err)
	assert.Equal(t, 4096, n)

	wg.Wait()
}
