## Tiny Go v2
### 목표
- go의 인터프리터 버전을 만들어 보기

### 변경점
- Let-in stmt 제거: 
    go와 더 유사해지기 위함
- Program 대신 Module을 시작점으로 삼음: 
    더 광범위한 호이스팅 표현 위함.

### EBNF

Package -> {Decl}


Decl -> VarDecl | FuncDecl 
VarDecl ->  "var" id {"," id} Type [ "=" Expr {"," Expr }] End
End -> ";"
FuncDecl -> "func" id Params [ReturnTypes] Block


Params -> Omit | "(" Param { "," Param} ")"
Omit -> "()"
Param ->  id Type

Type -> PrimitiveType | "func" ArgTypes [ReturnTypes]
PrimitiveType -> "int" | "bool" | "string" | "error"

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

SimpleStmt -> Assign | Call End  
ShortDecl-> id {"," id } ":=" Expr {"," Expr } End
Assign -> id {"," id} "=" Expr {"," Expr} End
If -> "if" [ShortDecl] Bexp Block ["else" Block ]
For ->  "for" Bexp Block
    |   "for" ShortDecl Bexp End id "=" Expr End Block 
    |   "for" "range" Aexp Block
Return -> "return" [Expr {"," Expr}] End
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
Call ->  BuiltIn Args | (id | Fexp | "("Expr")" ) Args {Args}
// Built-in Function = newError(string) error , scan(id), print(Expr), panic (Lexp)
BuiltIn -> "newError" | "scan" | "print" | "panic" 
Args ->  Omit | "(" Expr {"," Expr} ")" 

ValueForm -> Literal | Fexp
Literal := number | Bool | strlit | "nil"
Bool -> "true" | "false"

id = alpha{alpha| digit | "_"}
number = digit+
strlit = "..." // 부연설명: s = "..." 에 대해 trim(s, "\"") 



