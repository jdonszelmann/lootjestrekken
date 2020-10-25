package main

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"os"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestIntegration(t *testing.T) {
	log.Info("Running in memory test")
	IntegrationHelper(t, "inmemory", "")

	log.Info("Running db test")
	d := os.TempDir()
	p := d + "/LootjesTrekkenTestIntegration"
	err := os.RemoveAll(p)
	assert.NoError(t, err)
	err = os.Mkdir(p, os.ModePerm)
	assert.NoError(t, err)
	IntegrationHelper(t, "db", p)
	err = os.RemoveAll(p)
	assert.NoError(t, err)
}

func IntegrationHelper(t *testing.T, storetype, dbloc string) {
	ctx, cancel := context.WithCancel(context.Background())
	go runServer(ctx, "0.0.0.0", 12345, storetype, dbloc)
	defer cancel()

	time.Sleep(500 * time.Millisecond)

	res, err := http.Get("http://localhost:12345/")
	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)

	res, err = http.Get("http://localhost:12345/t")
	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)

	arr := make([]byte, 1024)
	n, err := res.Body.Read(arr)
	assert.Equal(t, err, io.EOF)
	assert.Equal(t, n, 0)
	assert.Equal(t, arr[:n], []byte(""))

	res, err = http.Get("http://localhost:12345/t/test/add")
	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)

	res, err = http.Get("http://localhost:12345/t")
	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)

	arr = make([]byte, 1024)
	n, err = res.Body.Read(arr)
	assert.Equal(t, err, io.EOF)
	assert.Greater(t, n, 0)
	assert.Equal(t, arr[:n], []byte("test"))

	res, err = http.Get("http://localhost:12345/t/test/people")
	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)

	arr = make([]byte, 1024)
	n, err = res.Body.Read(arr)
	assert.Equal(t, err, io.EOF)
	assert.Equal(t, n, 0)
	assert.Equal(t, arr[:n], []byte(""))

	res, err = http.Get("http://localhost:12345/t/test/people/jonathaan/add")
	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)

	res, err = http.Get("http://localhost:12345/t/test/people/jonathan/add")
	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)

	res, err = http.Get("http://localhost:12345/t/test/people")
	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)

	arr = make([]byte, 1024)
	n, err = res.Body.Read(arr)
	assert.Equal(t, err, io.EOF)
	assert.Greater(t, n, 0)
	assert.Equal(t, arr[:n], []byte("jonathaan\njonathan"))

	res, err = http.Get("http://localhost:12345/t/test/people/jonathaan/remove")
	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)

	res, err = http.Get("http://localhost:12345/t/test/people")
	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)

	arr = make([]byte, 1024)
	n, err = res.Body.Read(arr)
	assert.Equal(t, err, io.EOF)
	assert.Greater(t, n, 0)
	assert.Equal(t, arr[:n], []byte("jonathan"))
}


func TestSame(t *testing.T) {
	num := 100

	var wg sync.WaitGroup
	wg.Add(num)

	for i := 0; i < num; i++ {
		i := i
		go func() {
			port := 12340 + i

			ctx, cancel := context.WithCancel(context.Background())
			go runServer(ctx, "0.0.0.0", port, "inmemory", "")
			time.Sleep(1 * time.Second)

			res, err := http.Get(fmt.Sprintf("http://localhost:%d/t/test/add", port))
			assert.NoError(t, err)
			assert.Equal(t, res.StatusCode, http.StatusOK)

			res, err = http.Get(fmt.Sprintf("http://localhost:%d/t/test/people/a/add", port))
			assert.NoError(t, err)
			assert.Equal(t, res.StatusCode, http.StatusOK)

			res, err = http.Get(fmt.Sprintf("http://localhost:%d/t/test/people/b/add", port))
			assert.NoError(t, err)
			assert.Equal(t, res.StatusCode, http.StatusOK)

			res, err = http.Get(fmt.Sprintf("http://localhost:%d/t/test/people/c/add", port))
			assert.NoError(t, err)
			assert.Equal(t, res.StatusCode, http.StatusOK)

			res, err = http.Get(fmt.Sprintf("http://localhost:%d/t/test/people/d/add", port))
			assert.NoError(t, err)
			assert.Equal(t, res.StatusCode, http.StatusOK)

			res, err = http.Get(fmt.Sprintf("http://localhost:%d/t/test/trek", port))
			assert.NoError(t, err)
			assert.Equal(t, res.StatusCode, http.StatusOK)

			res, err = http.Get(fmt.Sprintf("http://localhost:%d/t/test/people/a/getrokken", port))
			assert.NoError(t, err)
			assert.Equal(t, res.StatusCode, http.StatusOK)
			arr := make([]byte, 1024)
			n, err := res.Body.Read(arr)
			assert.Equal(t, err, io.EOF)
			assert.Greater(t, n, 0)
			assert.NotEqual(t, arr[n-1:n], []byte("a"))
			assert.True(t, reflect.DeepEqual(arr[n-1:n], []byte("b")) || reflect.DeepEqual(arr[n-1:n], []byte("c")) || reflect.DeepEqual(arr[n-1:n], []byte("d")))

			res, err = http.Get(fmt.Sprintf("http://localhost:%d/t/test/people/b/getrokken", port))
			assert.NoError(t, err)
			assert.Equal(t, res.StatusCode, http.StatusOK)
			arr = make([]byte, 1024)
			n, err = res.Body.Read(arr)
			assert.Equal(t, err, io.EOF)
			assert.Greater(t, n, 0)
			assert.NotEqual(t, arr[n-1:n], []byte("b"))
			assert.True(t, reflect.DeepEqual(arr[n-1:n], []byte("a")) || reflect.DeepEqual(arr[n-1:n], []byte("c")) || reflect.DeepEqual(arr[n-1:n], []byte("d")))

			res, err = http.Get(fmt.Sprintf("http://localhost:%d/t/test/people/c/getrokken", port))
			assert.NoError(t, err)
			assert.Equal(t, res.StatusCode, http.StatusOK)
			arr = make([]byte, 1024)
			n, err = res.Body.Read(arr)
			assert.Equal(t, err, io.EOF)
			assert.Greater(t, n, 0)
			assert.NotEqual(t, arr[n-1:n], []byte("c"))
			assert.True(t, reflect.DeepEqual(arr[n-1:n], []byte("a")) || reflect.DeepEqual(arr[n-1:n], []byte("b")) || reflect.DeepEqual(arr[n-1:n], []byte("d")))

			res, err = http.Get(fmt.Sprintf("http://localhost:%d/t/test/people/d/getrokken", port))
			assert.NoError(t, err)
			assert.Equal(t, res.StatusCode, http.StatusOK)
			arr = make([]byte, 1024)
			n, err = res.Body.Read(arr)
			assert.Equal(t, err, io.EOF)
			assert.Greater(t, n, 0)
			assert.NotEqual(t, arr[n-1:n], []byte("d"))
			assert.True(t, reflect.DeepEqual(arr[n-1:n], []byte("a")) || reflect.DeepEqual(arr[n-1:n], []byte("b")) || reflect.DeepEqual(arr[n-1:n], []byte("c")))

			cancel()
			wg.Done()
		}()
	}

	wg.Wait()
}
