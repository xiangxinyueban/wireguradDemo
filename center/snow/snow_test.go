package snow

import (
	"fmt"
	"log"
	"sync"
	"testing"
)

func TestSnow(t *testing.T) {
	sn, err := NewSnow(1)
	if err != nil {
		log.Fatalln(err)
	}
	var wg sync.WaitGroup
	var results sync.Map
	for i := 1; i < 1; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			id := sn.GetID()
			_, exist := results.LoadOrStore(id, true)
			if exist {
				log.Fatalln("duplicated", id)
			}
			fmt.Println(id)
		}()
	}
	wg.Wait()
}
