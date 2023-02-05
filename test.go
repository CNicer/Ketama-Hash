package main

import (
	"crypto/md5"
	"fmt"
)

func main(){
	key := []byte("hello world")
	digest := md5.Sum(key)
	fmt.Println(len(digest))
}