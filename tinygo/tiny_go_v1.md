## Tiny Go v1
### 목표
- go로 인터프리팅 되는 go 언어를 만들어 보자

### 예시
```go
func main(){
    a, b := 4, 2
    divided, err := divide(a,b);
    if err != ok {
        print(errString(err))
        panic(err)
    }
    print("4 나누기 2는" + intToString(divided))

}

func divide(a int, b int) (int, error) {
    if b == 0 {
        return 0, newError("can't divide by zero");
    }
    return a/b, ok
}

func intToString(i int) string {
    if i == 0 {
        return digitToString(0)
    }

    lastDigit := i- 10*(i/10)
    reduced := i/10
    return intToString(reduced)+digitToString(lastDigit) 
}

func digitToString(i int) string {
    if i > 9 || i < 0 {
        panic("err")
    }
    if i == 0 {
        return "0"
    }
    //...반복은 생략...//

    if i == 9 {
        return "9"
    }
    return "0"
}

```

### 타입과 값
Built-in Type과 값 <br>
- int, number
- bool, true | false
- string, strlit
- error, strlit // tiny go에서는 error를 타입으로 다룬다.
- funcion 타입 

타입 간 연산<br>
- 서로 다른 타입 간 연산은 불가함

지원하는 연산<br>
- int : 이항연산, 단항연산, 비교연산
- bool: 일치연산, 논리연산
- string : 일치연산, +
- error : 일치연산
- function 타입 : 연산 제공하지 않음


- 이항연산 : +, -, *, /
- 단항연산 : -
- 일치연산 : ==, !=
- 비교연산 : 일치연산, <, <=, >, >=
- 논리연산 : &&, ||, !

### 실행모델과 스코핑
- 인터프리터
- main()부터 실행
- 정적 스코프
- 패키지 레벨에서 호이스팅 존재, 로컬 블록에선 호이스팅 없음.
- 호이스팅은 정확히 말하자면, go의 init order임.

### 에러 모델
- 에러를 표현, 타입으로 취급함
- newError를 통해 에러 표현 생성 가능
- errString을 통해 에러의 문자열 값 가져오기 가능

### 선언, 할당, 바인딩
선언
- 다중 선언의 LHS, RHS는 반드시 일치해야 함. (go-like하진 않지만, tiny go는 개수 일치만 accept함)
- 초기화가 없는 선언은 zero value로 초기화 된 선언으로 여김.
- 변수 셰도잉 허용함
- 단, 빌트인의 셰도잉은 허용하지 않음.
- 함수 매개변수의 중복을 허용하지 않음. (func (a int, a string) 불가)
할당
- 할당 역시 LHS, RHS 개수가 반드시 같아야 함
바인딩
- var로 변수 선언 시엔, 좌변의 변수들은 반드시 새 변수여야 함.
- ":=" 로 변수 선언 시엔, 좌변에 반드시 하나 이상의 새 변수 필요.
- ":=" 는 로컬 블록 내에서만 사용 가능.
- a, b = 1, 2 식의 동시 할당 및 선언 가능.
### 표준 환경
Built in function <br>
- 모든 빌트인 함수는 사용 시 반드시 호출되어야 하며
- 호출 시 결과값으로 다른 함수를 리턴하지 않음. (연쇄 호출 불가)
- 구문법 상으롣, 빌트인 함수는 바로 앞 하나의 괄호에 묶인 인자 리스트만을 가져감.
```go
    func newError(s string) error   // string 표현을 strlit으로 변환 후 error value로 리턴
    func errString(e error) string  // error의 strlit value를 string으로 리턴
    func len(s string) int
    func scan(id)       // id에 stdin의 값을 문자열로 받음
    func print(Expr)    // stdout에 string 타입의 Expr 출력
    func panic(Lexp)    // 프로그램 전체에 panic 전파
```
Predefine Operator <br>
- +, -, *, /
- -
- == , != , <, <=, >, >=
- &&, ||, !

### 구문법 (EBNF)
```ocaml
Package -> {Decl}

Decl -> VarDecl | FuncDecl  
VarDecl ->  "var" id {"," id} Type [ "=" Expr {"," Expr }] End
End -> ";"
FuncDecl -> "func" id Params [ReturnTypes] Block

Params -> Omit | "(" Param { "," Param} ")"
Omit -> "()"
Param ->  id Type

Type -> PrimitiveType | FuncType
FuncType ->  "func" ArgTypes [ReturnTypes]
PrimitiveType -> "int" | "bool" | "string" | "error"
ArgTypes -> Omit 
    |   "(" Type {"," Type} ")" 

ReturnTypes ->  Type
| "(" Type {"," Type} ")" 

Stmt -> Assign
    |   CallStmt
    |   ShortDecl
    |   VarDecl
    |   FuncDecl
    |   Return
    |   If
    |   For
    |   Block
Assign -> id {"," id} "=" Expr {"," Expr} End
CallStmt-> Call End
Call -> Primary Args {Args} (*| BuiltInCall*)
ShortDecl-> id {"," id } ":=" Expr {"," Expr } End
Return -> "return" [Expr {"," Expr}] End
If -> "if" [ShortDecl] Bexp Block ["else" Block ]

For ->  "for" Bexp Block
    |   "for" ShortDecl Bexp End id "=" Expr End Block 
    |   "for" "range" Aexp Block
Block -> "{" {Stmt} "}"


Expr ->  Bexp { "||" Bexp} 
Bexp -> Bterm {"&&" Bterm }
Bterm -> Aexp [Relop Aexp] | "!" Bterm 

Relop -> "==" | "!=" | "<" | "<=" | ">" | ">=" 
Aexp -> Term { ("+" | "-") Term } 
Term -> Factor { ("*" | "/") Factor } 
Factor -> ["-"]  Atom

Atom -> Primary {Args} (*| BuiltInCall*) //(* Atom = Primary | Call {call이 builtInCall 포함}*)
Primary -> "(" Expr ")" | id  |  ValueForm

BuiltInCall -> ("newError" | "errString" | "scan" | "print" | "panic" | "len") Args

ValueForm -> Literal | Fexp
Literal := number | "true" | "false" | strlit | "ok"
Fexp -> "func" Params [ReturnTypes] Block

Args ->  Omit | "(" Expr {"," Expr} ")" 

id = alpha{alpha| digit | "_"}
number = digit+
strlit = "..." // 부연설명: s = "..." 에 대해 trim(s, "\"") 
```
