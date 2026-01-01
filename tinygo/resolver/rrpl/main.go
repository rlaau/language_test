package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/rlaaudgjs5638/langTest/tinygo/lexer"
	"github.com/rlaaudgjs5638/langTest/tinygo/parser"
	"github.com/rlaaudgjs5638/langTest/tinygo/resolver"
)

func main() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println("\n(CTRL+C) 종료합니다. bye")
		os.Exit(0)
	}()
	in := bufio.NewReader(os.Stdin)
	fmt.Println("---------------------------")
	fmt.Println("|  Welcome tiny go RRPL   |")
	fmt.Println("|  - multi-line: 빈 줄로 실행")
	fmt.Println("|  - exit/quit 로 종료")
	fmt.Println("---------------------------")

	for {
		code, ok := readMultiline(in)
		if !ok {
			fmt.Println("\n입력 종료. bye")
			return
		}

		trim := strings.TrimSpace(code)
		if trim == "" {
			continue
		}
		if trim == "exit" || trim == "quit" {
			fmt.Println("bye")
			return
		}

		lx := lexer.NewLexer()
		lx.Set(code)

		ps := parser.NewParser(lx)
		parsed, err := ps.ParsePackage()
		if err != nil {
			fmt.Printf("parse error: %v\n", err)
			printParseErrorLocation(code, ps.CurrentToken().Pos)
			continue
		}
		fmt.Println(parsed.String())
		table, hoist, initOrder, _, rerr := resolver.Resolve(parsed)
		if rerr != nil {
			fmt.Printf("resolve error: %v\n", rerr)
			continue
		}
		fmt.Println(table.Print())
		fmt.Println(hoist.Print())
		fmt.Println(initOrder.Print(hoist))
	}
}

func readMultiline(r *bufio.Reader) (string, bool) {
	var b strings.Builder

	fmt.Print(">>> ")
	line, err := r.ReadString('\n')
	if err != nil {
		return "", false
	}
	line = strings.TrimRight(line, "\r\n")

	if strings.TrimSpace(line) == "" {
		return "", true
	}
	b.WriteString(line)
	b.WriteByte('\n')

	for {
		fmt.Print("... ")
		line, err := r.ReadString('\n')
		if err != nil {
			return b.String(), false
		}
		line = strings.TrimRight(line, "\r\n")

		if strings.TrimSpace(line) == "" {
			break
		}
		b.WriteString(line)
		b.WriteByte('\n')
	}

	return b.String(), true
}

func printParseErrorLocation(input string, pos int) {
	if len(input) == 0 {
		return
	}
	caretPos := pos - 1
	if caretPos < 0 {
		caretPos = 0
	}
	if caretPos > len(input) {
		caretPos = len(input)
	}

	lineStart := strings.LastIndex(input[:caretPos], "\n")
	if lineStart == -1 {
		lineStart = 0
	} else {
		lineStart++
	}
	lineEnd := strings.Index(input[caretPos:], "\n")
	if lineEnd == -1 {
		lineEnd = len(input)
	} else {
		lineEnd = caretPos + lineEnd
	}

	line := input[lineStart:lineEnd]
	column := caretPos - lineStart
	fmt.Println(line)
	fmt.Println(strings.Repeat(" ", column) + "^ parse error here")
}
