## 함수 표현식에 대해 확장된 언어
### EBNF
Program -> {Command}
Command -> Decl | Stmt 

Decl -> "var" id Type [ "=" Expr] ";"
    |   id ":=" Expr ";" 
    |   FuncDecl

FuncDecl -> "func" id Params [ReturnTypes] FuncBody
FuncBody -> "{" FuncStmts "}"
FuncStmts-> {FuncStmt}
FuncStmt -> Stmt | "return" [Expr {"," Expr}] ";"

Params -> "("")" | "(" Param { "," Param} ")"
Param ->  id Type

Type -> "int" 
    |   "bool" 
    |   "string"
    |   "func" ArgTypes [ReturnTypes]

ArgTypes -> "(" ")" 
    |   "(" Type {"," Type} ")" 
    |   "(" Param, {"," Param} ")"

ReturnTypes -> "(" Type {"," Type} ")" | Type


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



Expr -> Fexp | Bexp | Aexp | Sexp

Fexp -> -> "func" Params [RetrunTypes] FuncBody

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


Atom -> id | Call
Call ->  (id | Fexp | "("Expr")" ) Args {Args}
Args ->  "("")"| "(" Expr {"," Expr} ")" 

id = alpha{alpha|digit}
number = digit+
strlit = "..." // 부연설명: s = "..." 에 대해 trim(s, "\"") 
