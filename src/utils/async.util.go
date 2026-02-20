package utils

import (
	"log"
	"runtime/debug"
)

func RunAsync(fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("⚠️  Recovered from async panic: %v\n", r)
				log.Printf("Stack trace:\n%s\n", debug.Stack())
			}
		}()
		fn()
	}()
}

func RunAsyncWithCallback(fn func() (interface{}, error), callback func(interface{}, error)) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("⚠️  Recovered from async panic: %v\n", r)
				log.Printf("Stack trace:\n%s\n", debug.Stack())
			}
		}()
		result, err := fn()
		if callback != nil {
			callback(result, err)
		}
	}()
}
