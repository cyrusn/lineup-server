package main_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"
	"time"
)

func TestMain(t *testing.T) {
	t.Run("add", addStudent)
	for {
		time.Sleep(time.Duration(time.Second * 5))
		t.Run("test", testGet)
	}
}

func testGet(t *testing.T) {
	var wg sync.WaitGroup
	counter := 0

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter++

			resp, err := http.Get("http://localhost:5000/api/arrival")
			if err != nil {
				t.Fatal(err)
			}
			if _, err := ioutil.ReadAll(resp.Body); err != nil {
				t.Fatal(err)
			} else {
				// fmt.Printf("%s", b)
			}
		}()
	}
	wg.Wait()
}

func addStudent(t *testing.T) {
	classes := []string{"a", "b", "c", "d", "e"}
	forms := []int{1, 2, 3, 4, 5}

	for _, i := range forms {
		for _, j := range classes {
			for k := 0; k < 35; k++ {
				url := fmt.Sprintf("http://localhost:5000/api/arrival/%d%s/%d", i, j, k)
				http.Post(url, "", nil)
			}
		}
	}
}
