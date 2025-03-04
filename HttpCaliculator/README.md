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
    git clone https://github.com/MRgotem123/yandex_lms_golang/HttpCaliculator.git
    ```
  2. Перейдите в папку с кодом:
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

## Примеры ввода:
  - пример выражения, для терминала: curl -X POST http://localhost:9090/api/v1/calculate -d "(2+2)*(4-8)"
  - пример выражения, для PowerShell: Invoke-RestMethod -Uri "http://localhost:9090/api/v1/calculate" -Method Post -Body "(2+2)*(4-8)"
    
  - посмотреть все задачи, для терминала: curl -X POST http://localhost:9090/api/v1/expressions
  - посмотреть все задачи, для PowerShell: Invoke-RestMethod -Uri "http://localhost:9090/api/v1/expressions" -Method Post

  - посмотреть задачу по определённому id, для терминала: curl -X POST http://localhost:9090/api/v1/expressions/ -d "id36f8aa562f"
  - посмотреть задачу по определённому id, для PowerShell: Invoke-RestMethod -Uri "http://localhost:9090/api/v1/expressions/" -Method Post -Body "id36f8aa562f"

