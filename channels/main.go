package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type user struct {
	Name string `json:"first_name"`
	Age  int
}

func main() {
	processCh := make(chan user, 20)
	go func() {
		for i := 1; i <= 20; i++ {
			processCh <- user{
				Name: fmt.Sprintf("name_%d", i),
				Age:  i,
			}
		}

		close(processCh)
	}()

	var wg sync.WaitGroup
	wg.Add(2)

	printCh := make(chan user, 10)
	go func() {
		processUsers(processCh, printCh, &wg)
	}()
	go func() {
		printUsers(printCh, &wg)
	}()

	wg.Wait()
}

func processUsers(ch chan user, printCh chan user, wg *sync.WaitGroup) {
	for u := range ch {
		time.Sleep(time.Second * 2)
		u.Age = u.Age * 10
		u.Name = fmt.Sprintf("%s_processed", u.Name)
		printCh <- u
	}

	close(printCh)
	wg.Done()
}

func printUsers(ch chan user, wg *sync.WaitGroup) {
	for u := range ch {
		time.Sleep(time.Second)
		// fmt.Printf("User: name=%s, age=%d", u.name, u.age)
		raw, err := json.Marshal(u)
		if err != nil {
			fmt.Printf("error=%s\n", err.Error())
		} else {
			fmt.Printf("User json: %s\n", string(raw))
		}

		// us := User{}
		// err := json.Unmarshal(`{"first_name":"artem", "Age": 10}`, &us)
		// if err != nil {

		// }
	}
	wg.Done()
}
