package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

// Go Concurrency Patterns: Pipelines and cancellation
//
//Informally, a pipeline is a series of stages connected by channels, where each stage is a group of goroutines running the same function. In each stage, the goroutines
//receive values from upstream via inbound channels
//perform some function on that data, usually producing new values
//send values downstream via outbound channels
//
//Multiple functions can read from the same channel until that channel is closed; this is called fan-out. This provides a way to distribute work amongst a group of workers to parallelize CPU use and I/O.
//A function can read from multiple inputs and proceed until all are closed by multiplexing the input channels onto a single channel that's closed when all the inputs are closed. This is called fan-in.
//
//Sends on a closed channel panic, so it's important to ensure all sends are done before calling close.
//
//There is a pattern to our pipeline functions:
//stages close their outbound channels when all the send operations are done.
//stages keep receiving values from inbound channels until those channels are closed.
//
//Instead, we need to provide a way for downstream stages to indicate to the senders that they will stop accepting input.
//
//we can do this by closing a channel, because a receive operation on a closed channel can always proceed immediately, yielding the element type's zero value.
//
//This close is effectively a broadcast signal to the senders.
//
//Here are the guidelines for pipeline construction:
//stages close their outbound channels when all the send operations are done.
//stages keep receiving values from inbound channels until those channels are closed or the senders are unblocked.

// A result is the product of reading and summing a file using MD5.
type result struct {
	path string
	sum  [md5.Size]byte
	err  error
}

// sumFiles starts goroutines to walk the directory tree at root and digest each
// regular file.  These goroutines send the results of the digests on the result
// channel and send the result of the walk on the error channel.  If done is
// closed, sumFiles abandons its work.
func sumFiles(done <-chan struct{}, root string) (<-chan result, <-chan error) {
	// For each regular file, start a goroutine that sums the file and sends
	// the result on c.  Send the result of the walk on errc.
	c := make(chan result)
	errc := make(chan error, 1)
	go func() { // HL
		var wg sync.WaitGroup
		// 每个子文件开启新goroutine，add+进行操作+Done
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.Mode().IsRegular() {
				return nil
			}
			wg.Add(1)
			go func() { // HL
				data, err := ioutil.ReadFile(path)
				select {
				case c <- result{path, md5.Sum(data), err}: // HL
				case <-done: // HL
				}
				wg.Done()
			}()
			// Abort the walk if done is closed.
			select {
			case <-done: // HL
				return errors.New("walk canceled")
			default:
				return nil
			}
		})
		// Walk has returned, so all calls to wg.Add are done.  Start a
		// goroutine to close c once all the sends are done.
		go func() { // HL
			wg.Wait()
			close(c) // HL
		}()
		// No select needed here, since errc is buffered.
		errc <- err // HL
	}()
	return c, errc
}

// MD5All reads all the files in the file tree rooted at root and returns a map
// from file path to the MD5 sum of the file's contents.  If the directory walk
// fails or any read operation fails, MD5All returns an error.  In that case,
// MD5All does not wait for inflight read operations to complete.
func MD5All(root string) (map[string][md5.Size]byte, error) {
	// MD5All closes the done channel when it returns; it may do so before
	// receiving all the values from c and errc.
	done := make(chan struct{}) // HLdone
	defer close(done)           // HLdone

	c, errc := sumFiles(done, root) // HLdone

	m := make(map[string][md5.Size]byte)
	// range进行持续观察chan中所有信息直到close
	for r := range c { // HLrange
		if r.err != nil { // 有一个错误则return
			return nil, r.err
		}
		m[r.path] = r.sum
	}
	if err := <-errc; err != nil {
		return nil, err
	}
	return m, nil
}

func main() {
	m, err := MD5All(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	var paths []string
	for path := range m {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		fmt.Printf("%x  %s\n", m[path], path)
	}
}
