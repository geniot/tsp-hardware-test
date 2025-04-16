package main

import (
	"log"
	"os"
)

func If[T any](cond bool, vTrue, vFalse T) T {
	if cond {
		return vTrue
	}
	return vFalse
}

func CloseFile(f *os.File) {
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func First[T, U, V any](val T, _ U, _ V) T {
	return val
}

func orPanic(err interface{}) {
	switch v := err.(type) {
	case error:
		if v != nil {
			panic(err)
		}
	case bool:
		if !v {
			panic("condition failed: != true")
		}
	}
}
func orWarn(err interface{}) {
	switch v := err.(type) {
	case error:
		if v != nil {
			log.Println(err)
		}
	case bool:
		if !v {
			log.Println(err)
		}
	}
}

func orPanicRes[T any](res T, err interface{}) T {
	orPanic(err)
	return res
}
