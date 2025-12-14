### Goal
let's make tiny Go
### EBNF
Program -> {Command}
Command -> Decl | Stmt 

Decl -> "var" id Type [ "=" Expr] ";"| id ":=" Expr ";"
Type -> "int" | "bool" | "string"
Value -> id | number | strlit

id = alpha{alpha|digit}
number = digit+
strlit = s = "..." 에 대해 trim(s, "\"") 

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