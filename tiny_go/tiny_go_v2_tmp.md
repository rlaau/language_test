## Tiny Go v2
### 목표
- go의 인터프리터 버전을 만들어 보기

### 변경점
- Let-in stmt 제거: 
    go와 더 유사해지기 위함
- Program 대신 Module을 시작점으로 삼음: 
    더 광범위한 호이스팅 표현 위함.

### EBNF

Module -> {Decl}


Decl -> VarDecl | FuncDecl
VarDecl ->  "var" id {"," id} Type [ "=" Expr {"," Expr }] ";"
FuncDecl -> "func" id Params [ReturnTypes] Block


Params -> Omit | "(" Param { "," Param} ")"
Omit -> "()"
Param ->  id Type

Type -> PrimitiveType | "func" ArgTypes [ReturnTypes]
PrimitiveType -> "int" | "bool" | "string" 

ArgTypes -> Omit 
    |   "(" Type {"," Type} ")" 

ReturnTypes -> "(" Type {"," Type} ")" | Type


Stmt -> SimpleStmt
    |   ShortDecl
    |   VarDecl
    |   FuncDecl
    |   Return
    |   If
    |   For
    |   Block

SimpleStmt -> Assign | "scan" id ";" | "print" Expr ";" |  Call ";"
ShortDecl-> id {"," id } ":=" Expr {"," Expr } ";"
Assign -> id {"," id} "=" Expr {"," Expr} ";"
If -> "if" [ShortDecl] Bexp Block ["else" Block ]
For ->  "for" Bexp Block
    |   "for" ShortDecl Bexp ";" id "=" Expr ";" Block 
    |   "for" "range" Aexp Block
Return -> "return" [Expr {"," Expr}] ";"
Block -> "{" {Stmt} "}"



Expr -> Fexp | Lexp  

Fexp -> "func" Params [ReturnTypes] Block

Lexp -> Bexp { ("&&" | "||") Bexp} | "!" Bexp 
Bexp -> Aexp [Relop Aexp]

Relop -> "==" | "!=" | "<" | "<=" | ">" | ">=" 

Aexp -> Term { ("+" | "-") Term } 
Term -> Factor { ("*" | "/") Factor } 

Factor -> ["-"] ("(" Aexp ")" | Atom )

Atom -> id | Call | ValueForm
Call ->  (id | Fexp | "("Expr")" ) Args {Args}
Args ->  Omit | "(" Expr {"," Expr} ")" 

ValueForm -> Literal | Fexp
Literal := number | Bool | strlit
Bool -> "true" | "false"



id = alpha{alpha|digit}
number = digit+
strlit = "..." // 부연설명: s = "..." 에 대해 trim(s, "\"") 



