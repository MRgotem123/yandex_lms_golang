# MRgotem - Http Calculator

Это многопользовательский Http калькулятор, который принимает выражения через HTTP-запрос и вычисляет его при помощи Обратной польской нотации (RPN), и GOrutins.
Http калькулятор состоит из Агента и Оркестратора.

**Мой тг:** https://t.me/G_O_T_E_M123
(если возникнут вопросы, пишите)

**Архетектура кода:** https://unidraw.io/app/board/80203b3c44f12d115650

## Возможности:
  - Обработка математических выражений с операциями `+`, `-`, `*`, `/`.
  - Поддержка скобок для указания порядка операций.
  - Проверка выражения на наличие синтаксических ошибок.
  - Поддерживает десятичные числа.
  - Поддерживает отрицательные числа.
  - Разбивает сложные выражения на простые, для ускорения вычислений.
  - Сохраняет историю вычислений.
  - Поиско выражений по их ID.
  - Регулировка количества GOrutins.
  - Регулировка времени вычисления каждой операции.

## Установка:
  Требуется:
   - Go 1.23.1, или выше.
   - git

  1. Клонируйте репозиторий:
     ```bash
     git clone --no-checkout https://github.com/MRgotem123/yandex_lms_golang/HttpCalculator.git
     cd HttpCalculator
     ```

## Запуск:

  1. Запуск Оркестратора:
     ```bash
     go run Orchestrator/
     ```
     
  2. Запуск Агента:
     (в новом окне терминала)
     ```bash
     go run Agent/
     ```

**Адрес:** http://localhost:9090/api/v1

## Примеры ввода:
  (в новом окне терминала)
  - пример регестрации, для терменала:
    ```bash
      curl -X POST "http://localhost:9090/api/v1/resister" \
     -H "Content-Type: application/json" \
     -d '{"login":"matvey","password":"123"}'
    ```

  - пример регестрации, для PowerShell:
    ```bash
      $body = @{
        login    = "matvey"
        password = "123"
      } | ConvertTo-Json -Depth 3

      Invoke-RestMethod `
        -Uri "http://localhost:9090/api/v1/resister" `
        -Method Post `
        -Body $body `
        -ContentType "application/json"
    ```

  - пример входа, для терменала:
    ```bash
      curl -X POST "http://localhost:9090/api/v1/login" \
     -H "Content-Type: application/json" \
     -d '{"login":"matvey","password":"123"}'
    ```

  - пример входа, для PowerShell:
    ```bash
    $body = @{
        login    = "matvey"
        password = "123"
      } | ConvertTo-Json -Depth 3

      Invoke-RestMethod `
        -Uri "http://localhost:9090/api/v1/login" `
        -Method Post `
        -Body $body `
        -ContentType "application/json"
    ```
  
  - пример выражения, для терминала:
     ```bash
     curl -X POST http://localhost:9090/api/v1/calculate -d "(2+2)*(4-8)"
     ```
     
  - пример выражения, для PowerShell:
      ```bash
      Invoke-RestMethod -Uri "http://localhost:9090/api/v1/calculate" -Method Post -Body "(2+2)*(4-8)"
      ```
    
  - посмотреть все задачи, для терминала:
      ```bash
      curl -X POST http://localhost:9090/api/v1/expressions
      ```
      
  - посмотреть все задачи, для PowerShell:
      ```bash
      Invoke-RestMethod -Uri "http://localhost:9090/api/v1/expressions" -Method Post
      ```

  - посмотреть задачу по определённому id, для терминала:
      ```bash
      curl -X POST http://localhost:9090/api/v1/expressions/ -d "id36f8aa562f"
      ```
      
  - посмотреть задачу по определённому id, для PowerShell:
      ```bash
      Invoke-RestMethod -Uri "http://localhost:9090/api/v1/expressions/" -Method Post -Body "id36f8aa562f"
      ```

## Коды ответа:
  1. Содержит невалидные символы, код ответа: 422
     ```bash
     curl -X POST http://localhost:9090/api/v1/calculate -d "!(2+2)*(4-8)"
     ```

  2. ID успешно записан, код ответа: 201
     ```bash
     curl -X POST http://localhost:9090/api/v1/calculate -d "(2+2)*(4-8)"
     ```

  3. Список всех выражений, 500 если выражений не было, 200 если есть хотя бы одно выражение.
     ```bash
     curl -X POST http://localhost:9090/api/v1/expressions
     ```

  4. Поиск выражения по ID, 500 если ошибка, 404 если нет такого ID, 200 если ID найдено.
     ```bash
     curl -X POST http://localhost:9090/api/v1/expressions/ -d "id36f8aa562f"
     ```

  5.  Агент запрашивает: получение задачи на выполнение, 500 что-то пошло не так, 404 нет задачи, 200 задача успешно получена.

## Запуск тестов:

  - запуск теста для Orchestrator:
    ```bash
    go test Orchestrator/
    ```

  - запуск теста для Agent:
    ```bash
    go test Agent/
    ```
