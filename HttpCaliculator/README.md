**Клонирование репозитория:** git clone https://github.com/MRgotem123/yandex_lms_golang/HttpCaliculator.git

Запуск Оркестратора: go run Orchestrator.go

Запуск Агента: go run Agent.go

пример выражения для терменала: curl -X POST http://localhost:9090/api/v1/calculate -d "(2+2)*(4-8)"

пример выражения для PowerShell: Invoke-RestMethod -Uri "http://localhost:9090/api/v1/calculate" -Method Post -Body "(2+2)*(4-8)"

посмотреть все задачи: curl -X POST http://localhost:9090/api/v1/expressions

посмотреть задачу по определённому id: curl -X POST http://localhost:9090/api/v1/expressions/ -d "id36f8aa562f"
