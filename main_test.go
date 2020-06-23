package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"testing"
)

func TestPart(t *testing.T) {
	length := 167895
	partLen := 400

	// 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10
	for start, end := 0, 0; start < length; start += partLen {
		if length-start >= partLen {
			end = start + partLen - 1
		} else {
			end = length - 1
		}

		fmt.Printf("start: %d, end: %d\n", start, end)
	}
}

func TestD(t *testing.T) {
	fmt.Printf("%d\n", 30/2/5)
}

func TestWriteFile(t *testing.T) {
	file, err := os.Create("test.txt")
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	defer file.Close()
	defer file.Sync()

	newOffset, err := file.Seek(3, 0)
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	fmt.Printf("offset: %d\n", newOffset)
	file.Write([]byte("3"))
}

func TestRead(t *testing.T) {
	file, err := os.Open("test.txt")
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	defer file.Close()

	data, _ := ioutil.ReadAll(file)
	fmt.Printf("%v\n", data)
}

func TestCurrent(t *testing.T) {
	file, err := os.OpenFile("test.txt", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	defer file.Close()
	defer file.Sync()

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		for j := 0; j < 100; j++ {
			for i := 65; i <= 90; i++ {
				n, err := file.Write([]byte{byte(i)})
				fmt.Printf("go1, n: %d, err: %s\n", n, err)
			}
		}
		wg.Done()
	}()

	go func() {
		for j := 0; j < 100; j++ {
			for i := 97; i <= 122; i++ {
				n, err := file.Write([]byte{byte(i)})
				fmt.Printf("go2, n: %d, err: %s\n", n, err)
			}
		}
		wg.Done()
	}()
	wg.Wait()
}

func TestAPPend(t *testing.T) {
	file, err := os.OpenFile("test.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	defer file.Close()
	defer file.Sync()

}

func TestContentLen(t *testing.T) {
	url := "https://nchc.dl.sourceforge.net/project/evolution-x/raphael/EvolutionX_4.4_raphael-10.0-20200602-1022-OFFICIAL.zip"
	length := ContentLen(url)
	fmt.Printf("len: %d\n", length)
}
