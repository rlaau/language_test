## 추가하려는 기능: 표현으로써의 함수
### 문법
Expr -> 기존 Expr | Func | Func Param
Stmt -> 기존 Stmt | "return" Expr   // 고민 중: 이렇게 되면 "return 문"이 어느 블록에나 들어갈 수 있음. 이걸 허용 후 타입체크를 통해서 함수 블록 안에만 리턴문이 있게 할 지, 아니면 "언어 차원에서" stmt가 함수에만 들어가게 할 지?
Type -> 기존 Type | "func" ((Type Type)| "()" Type | Type)
Param -> "()"| "(" Expr {"," Expr} ")" 
Func -> "func" "(" Args ")" Body
Args -> "()" | "(" Arg { "," Arg} ")"
Arg ->  id Type

## 기존 언어
### EBNF
Program -> {Command}
Command -> Decl | Stmt 

Decl -> "var" id Type [ "=" Expr] ";"| id ":=" Expr ";"
Type -> "int" | "bool" | "string"
Value -> id | number | strlit

id = alpha{alpha|digit}
number = digit+
strlit = "..." // 부연설명: s = "..." 에 대해 trim(s, "\"") 

Stmt -> SimpleStmt ";"
    |   Block
    |   If
    |   For
    |   Let

SimpleStmt -> id "=" Expr | "scan" id | "print" Expr

Block -> "{" Stmts "}"
If -> "if" Bexp "then" Block ["else" Block ]
For ->  "for" Bexp Block
    |   "for" Decl Bexp ";" id "=" Expr ";" Block 
    |   "for" "range" Aexp Block

Let -> "let" Decls "in" "(" Stmts ")"

Stmts -> {Stmt}
Decls -> {Decl}



Expr -> Bexp | Aexp | Sexp

Bexp -> Bterm { "||" Bterm } 
Bterm -> Bfact   { "&&" Bfact   } 

Bfact -> ["!"] Batom 
Batom -> "true" | "false" 
    |   Aexp Relop Aexp  
    |   Sexp Srelop Sexp 
    |   "(" Bexp ")"  
    |   id

Relop -> "==" | "!=" | "<" | "<=" | ">" | ">=" 
Srelop -> "==" | "!="

Aexp -> Aterm { ("+" | "-") Aterm } 
Aterm -> Power { ("*" | "/") Power } 

Power -> Afactor {"^" Afactor} 
Afactor -> ["-"] Aatom 
Aatom -> number | "(" Aexp ")" | id

Sexp -> Satom {"+" Satom}
Satom -> strlit | id