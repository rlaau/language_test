## 추가하려는 기능: 표현으로써의 함수
### 문법
Expr -> 기존 Expr | FuncExpr
FuncExpr -> "func" Params [RetrunTypes] FuncBody

Atom -> 기존 Atom | Call
Call ->  (id | FuncExpr | "("Expr")" ) Args {Args}


Decl -> 기존 Decl | FuncDecl
FuncDecl -> "func" id Params [ReturnTypes] FuncBody

FuncBody -> "{" FuncStmts "}"
FuncStmts-> {FuncStmt}
FuncStmt -> Stmt | "return" [Expr {"," Expr}] ";"

Type -> 기존 Type | "func" ArgTypes [ReturnTypes]
ArgTypes -> "(" [Type {"," Type}] ")" 
ReturnTypes -> ArgTypes | Type

Params -> "("")" | "(" Param { "," Param} ")"
Param ->  id Type
Args ->  "("")"| "(" Expr {"," Expr} ")" 

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
    |   atom

Relop -> "==" | "!=" | "<" | "<=" | ">" | ">=" 
Srelop -> "==" | "!="

Aexp -> Aterm { ("+" | "-") Aterm } 
Aterm -> Afactor { ("*" | "/") Afactor } 

Afactor -> ["-"] Aatom 
Aatom -> number | "(" Aexp ")" | atom

Sexp -> Satom {"+" Satom}
Satom -> strlit | atom

Atom -> id