package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/rlaaudgjs5638/langTest/tinygo/evaluator"
	"github.com/rlaaudgjs5638/langTest/tinygo/lexer"
	"github.com/rlaaudgjs5638/langTest/tinygo/parser"
	"github.com/rlaaudgjs5638/langTest/tinygo/resolver"
)

// * 예시 코드
// func main(){a,b:=4,2; divided,err:=divide(a,b); if err!=ok{print(errString(err));panic(err);}print("4 divide 2 is"+intToString(divided));} func divide(a int,b int)(int,error){if b==0{return 0,newError("can't divide by zero"); }return a/b,ok;} func intToString(i int)string{if i==0{return digitToString(0); } lastDigit:=i-10*(i/10);reduced:=i/10;return intToString(reduced)+digitToString(lastDigit);} func digitToString(i int)string{if i>9||i<0{panic("out of digit range"); } if i==0{return "0" ;} if i==1{return "1";} if i==2{return "2";} if i==3{return "3";}if i==4{return "4";}if i==5{return "5";}if i==6{return "6";}if i==7{return "7";}if i==8{return "8";}if i==9{return "9";}return "0";}
// 출력 결과: 4 divide 2 is02
func main() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println("\n(CTRL+C) exit")
		os.Exit(0)
	}()

	in := bufio.NewReader(os.Stdin)

	fmt.Println("------------------------")
	fmt.Println("| Welcome tiny go REPL |")
	fmt.Println("| - multi-line: empty line to run")
	fmt.Println("| - exit/quit to exit")
	fmt.Println("------------------------")

	for {
		code, ok := readMultiline(in)
		if !ok {
			fmt.Println("\ninput ended")
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
		pkg, err := ps.ParsePackage()
		if err != nil {
			fmt.Printf("parse error: %v\n", err)
			printParseErrorLocation(code, ps.CurrentToken().Pos)
			continue
		}

		table, hoist, order, builtins, err := resolver.Resolve(pkg)
		if err != nil {
			fmt.Printf("resolve error: %v\n", err)
			continue
		}

		_, err = evaluator.Evaluate(*pkg, hoist, order, table, builtins)
		if err != nil {
			fmt.Printf("eval error: %v\n", err)
			continue
		}
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
