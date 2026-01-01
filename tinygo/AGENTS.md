# AGENTS.md for tinygo

## Project context

- This repository contains a Tiny Go language implementation.
- The language is built as a pipeline: lexer -> parser -> resolver -> evaluator.
- The language spec is documented in `tinygo/tiny_go_v1.md`.

## Working assumptions

- Keep changes consistent with the current pipeline structure.
- Prefer reading the spec and matching resolver behavior when adding evaluator logic.
- Expect the language to evolve with new types, features, and built-ins; avoid hard-coding one-off rules unless the spec says so.

## Where to look first

- Specifiaction : `tinygo/tiny_go_v1.md`
- Lexer: `tinygo/lexer`
- Parser: `tinygo/parser`
- Resolver: `tinygo/resolver`
- Evaluator: `tinygo/evaluator`

## Guidance for new work

- Follow the spec in `tinygo/tiny_go_v1.md` for current behavior.
- When changing scoping or name resolution, align with resolver rules and tests.
- Keep interfaces between stages simple and explicit (AST, ResolveTable, HoistInfo, InitOrder).
