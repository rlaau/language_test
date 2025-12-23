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
)

func main() {
	// Ctrl+C (SIGINT), 종료 시그널 처리
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// 시그널 오면 우아하게 종료
	go func() {
		<-sigCh
		fmt.Println("\n(CTRL+C) 종료합니다. bye 👋")
		os.Exit(0)
	}()

	in := bufio.NewReader(os.Stdin)

	fmt.Println("------------------------")
	fmt.Println("|  Welcome tiny go REPL |")
	fmt.Println("|  - multi-line: 빈 줄로 실행")
	fmt.Println("|  - exit/quit 로 종료")
	fmt.Println("------------------------")

	for {
		code, ok := readMultiline(in)
		if !ok {
			// EOF (Ctrl+D) 등
			fmt.Println("\n입력 종료. bye 👋")
			return
		}

		trim := strings.TrimSpace(code)
		if trim == "" {
			continue
		}
		if trim == "exit" || trim == "quit" {
			fmt.Println("bye 👋")
			return
		}

		lx := lexer.NewLexer()
		lx.Set(code)

		ps := parser.NewParser(lx)
		parsed, err := ps.ParsePackage()
		if err != nil {
			// println 대신, 보기 좋게
			fmt.Printf("error: %v\n", err)
			continue
		}

		fmt.Println(parsed.String())
	}
}

// 여러 줄 입력을 받아 하나의 string으로 합쳐 반환.
// 규칙: 첫 프롬프트 >>>, 이후 ... , 빈 줄이면 종료(실행)
func readMultiline(r *bufio.Reader) (string, bool) {
	var b strings.Builder

	// 첫 줄
	fmt.Print(">>> ")
	line, err := r.ReadString('\n')
	if err != nil {
		// EOF면 false
		return "", false
	}
	line = strings.TrimRight(line, "\r\n")

	// 빈 줄이면 그냥 빈 입력
	if strings.TrimSpace(line) == "" {
		return "", true
	}
	b.WriteString(line)
	b.WriteByte('\n')

	// 다음 줄들
	for {
		fmt.Print("... ")
		line, err := r.ReadString('\n')
		if err != nil {
			return b.String(), false
		}
		line = strings.TrimRight(line, "\r\n")

		// 빈 줄이면 입력 종료
		if strings.TrimSpace(line) == "" {
			break
		}
		b.WriteString(line)
		b.WriteByte('\n')
	}

	return b.String(), true
}
