package test

import (
	"bufio"
	"context"
	"fmt"
	"strings"
	"testing"
	"testrand/reader/eval"
)

func BenchmarkRead(b *testing.B) {
	sample := strings.NewReader(`
	(begin
		(define local-word-count (new-hashmap))
		(define i 0)
		(loop (not (eq? i 1000)) (begin
		(foreach-array (string-split (read-file "sample1.txt") " ")
			(lambda (word)
				(put-hashmap local-word-count word (+ (get-hashmap local-word-count word 0) 1))
			)
		)
		(set i (+ i 1))
		))
	)
	`)
	read := eval.New(bufio.NewReader(sample))
	result, err := read.Read()
	if err != nil {
		panic(err)
	}
	env, err := eval.NewGlobalEnvironment()
	ctx := context.Background()
	if err != nil {
		panic(err)
	}
	b.StartTimer()
	result, err = eval.Eval(ctx, result, env)
	b.StopTimer()
	fmt.Println(result)
	if err != nil {
		panic(err)
	}
}

func BenchmarkLoop(b *testing.B) {
	sample := strings.NewReader(`
(begin
(define a 0)
(loop (< a 10000) (begin
(set a (+ a 1))
))
)
	`)
	read := eval.New(bufio.NewReader(sample))
	result, err := read.Read()
	if err != nil {
		panic(err)
	}
	env, err := eval.NewGlobalEnvironment()
	ctx := context.Background()
	if err != nil {
		panic(err)
	}
	b.StartTimer()
	result, err = eval.Eval(ctx, result, env)
	b.StopTimer()
	fmt.Println(result)
	if err != nil {
		panic(err)
	}
}
