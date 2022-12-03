// package main

// import (
// 	"fmt"
// 	"time"
// )

// func g(v int) {
// 	fmt.Println(v*2, v*v, " ")
// }

// func SpawnGoroutines(n int) {
// 	for i := 0; i < n; i++ {
// 		// fmt.Println(i*2, i*i, " ")
// 		go g(i)
// 	}
// }

// func main() {

// 	go SpawnGoroutines(10)

// 	// sleep main goroutine
// 	time.Sleep(1 * time.Second)
// }
