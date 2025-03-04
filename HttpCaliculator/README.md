# MRgotem - Http Calculator

Это Http калькулятор, который принимает выражения через HTTP-запрос и вычисляет его при помощи Обратной польской нотации (RPN), и GOrutins.

## Возможности:
  - Обработка математических выражений с операциями `+`, `-`, `*`, `/`.
  - Поддержка скобок для указания порядка операций.
  - Проверка выражения на наличие синтаксических ошибок.
  - Поддерживает десятичные цифры.
  - Разбивает сложные выражения на простые, для ускорения вычеслений.
  - Сохраняет историю вычислений.
  - Поиско выражений по их ID.
  - Регулировка количества GOrutins.
  - Регулировка времени вычисления каждой опирации.

## Установка:
  (Требуется версия Go 1.23.1, или выше.)

  1. Клонируйте репозиторий:
     ```bash
     git clone --no-checkout https://github.com/MRgotem123/yandex_lms_golang.git
     cd yandex_lms_golang
     git sparse-checkout init --cone
     git sparse-checkout set HttpCaliculator
     git checkout
     ```

  3. Перейдите в папку с кодом:
     ```bash
     cd HttpCaliculator
     ```

## Запуск:

  1. Запуск Оркестратора:
     ```bash
     go run Orchestrator.go
     ```
     
  2. Запуск Агента:
     ```bash
     go run Agent.go
     ```

**Адрес:** http://localhost:9090/api/v1/calculate

## Примеры ввода:
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

  3. Список всех выражений, 500 если выражений небыло, 200 если есть хотя-бы одно выражение.
     ```bash
     curl -X POST http://localhost:9090/api/v1/expressions
     ```

  4. Поиск выражения по ID, 500 если ошибка, 404 если нет такого ID, 200 если ID найдено.
     ```bash
     curl -X POST http://localhost:9090/api/v1/expressions/ -d "id36f8aa562f"
     ```

  5.  Агент запрашивает: получение задачи на выполнение, 500 что-то пошло не так, 404 нет задачи, 200 задача успешно получина.
