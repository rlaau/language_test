# Language Test

## GOAL

MyDSL을 만들어보자.

MyDSL은 반응형 · 비동기 · 함수형 패러다임을 반영한
DSL이다.

---

## Language Spec

1. 인터프리터 언어다.
2. 반응형 아키텍쳐에 용이해야 한다.
3. 비동기 계산이 가능해야 한다.
4. 함수형 패러다임을 반영한다.

---

## Host Runtime

- Host Language: Go
- 이유:
  - 런타임 제어가 명확함
  - 고루틴/채널을 통한 비동기 모델 실험이 용이함
  - 내가 익숙함

