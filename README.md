# lib-result

[![CI](https://github.com/selfshop-dev/lib-result/actions/workflows/ci.yml/badge.svg)](https://github.com/selfshop-dev/lib-result/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/selfshop-dev/lib-result/branch/main/graph/badge.svg)](https://codecov.io/gh/selfshop-dev/lib-result)
[![Go Report Card](https://goreportcard.com/badge/github.com/selfshop-dev/lib-result)](https://goreportcard.com/report/github.com/selfshop-dev/lib-result)
[![Go version](https://img.shields.io/github/go-mod/go-version/selfshop-dev/lib-result)](go.mod)
[![License](https://img.shields.io/github/license/selfshop-dev/lib-result)](LICENSE)

`Result[T, E]` для Go — явный тип для операций, которые возвращают значение или ошибку. Без внешних зависимостей. Проект организации [selfshop-dev](https://github.com/selfshop-dev).
 
### Installation
 
```bash
go get -u github.com/selfshop-dev/lib-result
```
 
## Overview
 
`lib-result` реализует тип `Result[T, E]`, который делает путь ошибки явным в сигнатуре и позволяет строить цепочки трансформаций без повторяющихся `if err != nil`. Идиоматичная пара `(T, error)` работает хорошо для большинства случаев, но когда нужно трансформировать, собирать или цеплять несколько операций, каждая из которых может упасть, проверки накапливаются быстро.
 
```go
// Без Result
user, err := repo.FindUser(ctx, id)
if err != nil {
    return "", err
}
order, err := repo.LatestOrder(ctx, user.ID)
if err != nil {
    return "", err
}
return order.Reference, nil
 
// С Result
return result.AndThen(
    result.Of(repo.FindUser(ctx, id)),
    func(u User) result.Result[Order, error] {
        return result.Of(repo.LatestOrder(ctx, u.ID))
    },
).Map(func(o Order) string { return o.Reference }).ToGo()
```

Для простого сервисного кода идиоматичная пара `(T, error)` предпочтительнее. Используй `Result` там, где цепочки добавляют реальную читаемость.
 
### Быстрый старт
 
```go
import "github.com/selfshop-dev/lib-result"
 
// Из идиоматичного Go
r := result.Of(repo.FindUser(ctx, id))
 
// Трансформация
name := r.Map(func(u User) string { return u.Name }).UnwrapOr("anonymous")
 
// Обратно в идиоматичный Go
user, err := r.ToGo()
```
 
## Value[T]
 
Для распространённого случая где тип ошибки — `error`, используй алиас [`Value[T]`](result.go):
 
```go
func FindUser(ctx context.Context, id int64) result.Value[User]
 
r := result.Of(repo.FindUser(ctx, id)) // inferred as result.Value[User]
```
 
`Value[T]` — это просто `Result[T, error]`. Когда нужен конкретный тип ошибки — используй полную форму `Result[T, *apperr.Error]`.
 
## Конструкторы
 
Создать `Result` можно явно или через мост из идиоматичного Go:
 
```go
result.Ok[User, error](user)      // явный Ok
result.Err[User, error](err)      // явный Err — паникует если err == nil
result.Of(repo.FindUser(ctx, id)) // из идиоматичной Go пары (value, error)
 
result.OfTyped[User, *apperr.Error](repo.FindUser(ctx, id)) // с конкретным типом ошибки
```
 
`Err` паникует при `nil` — в том числе при типизированном nil (`var e *MyError = nil`). Nil-ошибка это успех; используй `Ok`.
 
## Доступ к значению
 
Безопасные методы никогда не паникуют; небезопасные предназначены только для тестов и инициализации программы.
 
```go
// Безопасные — не паникуют
v, ok := r.Value()          // (T, bool)
e, ok := r.Err()            // (E, bool)
r.UnwrapOr(fallback)        // T или fallback
r.UnwrapOrElse(func(e E) T) // T или результат fn
r.UnwrapOrZero()            // T или zero value
 
// Небезопасные — только для тестов и инициализации
r.Unwrap()            // T или panic
r.Must("load config") // T или panic с контекстом
```
 
## Трансформации
 
Все трансформации — package-level функции, так как Go не поддерживает дополнительные type-параметры в методах.
 
```go
// Трансформировать Ok-значение
result.Map(r, strings.ToUpper)
 
// Трансформировать ошибку
result.MapErr(r, func(e error) *apperr.Error {
    return apperr.Wrap(e, apperr.KindNotFound, "user not found")
})
 
// Цепочка операций — short-circuit на Err
result.AndThen(r, func(u User) result.Result[Order, error] {
    return result.Of(repo.LatestOrder(ctx, u.ID))
})
 
// Fallback на ошибку — short-circuit на Ok
result.OrElse(r, func(e error) result.Result[User, error] {
    if apperr.IsKind(e, apperr.KindNotFound) {
        return result.Ok[User, error](guestUser)
    }
    return result.Err[User, error](e)
})
 
// Комбинирование
result.And(r, other) // other если r Ok, иначе Err r
result.Or(r, other)  // r если Ok, иначе other
```
 
## Коллекции
 
Два варианта для работы со срезом результатов: с остановкой на первой ошибке или со сбором всех сбоев.
 
```go
// Остановиться на первой ошибке
all := result.Collect(results) // Result[[]T, E]
 
// Собрать всё — не останавливаться на ошибках
values, errs := result.CollectAll(results) // ([]T, []E)
```
 
`Collect` подходит когда весь batch бессмысленен при любой ошибке. `CollectAll` — когда нужно обработать все элементы и сообщить обо всех сбоях.
 
## Конвертация обратно в Go
 
Для совместимости со стандартной библиотекой и сторонним кодом `Result` можно вернуть обратно в идиоматичный Go.
 
```go
user, err := r.ToGo()      // (T, error)
user, err := r.ToGoTyped() // (T, E) — без type assertion
value, ok := r.Option()    // (T, bool) — ошибка отбрасывается
```

## Makefile

Основные возможности:

| Цель | Описание |
|---|---|
| `make code-gen` | Запустить `go generate ./...` |
| `make lint` | Запустить golangci-lint |
| `make test` | Генерация кода + тесты с coverage |
| `make prof` | Собрать профили (cpu, mem, block, mutex) |
| `make prof-view` | Открыть профиль в браузере (`FILE=cpu.out` по умолчанию) |

## Лицензия

[`MIT`](LICENSE) © 2026-present [`selfshop-dev`](https://github.com/selfshop-dev)