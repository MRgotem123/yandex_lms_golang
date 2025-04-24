package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/MRgotem123/yandex_lms_golang/HttpCalculator/WorkWithSQL"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"unicode"
)

var UserID string

type Values struct {
	Expression string
	Status     string
	Result     string
}

type SendValues struct {
	Id             string `json:"id"`
	Arg1           string `json:"arg1"`
	Arg2           string `json:"arg2"`
	Operation      string `json:"operation"`
	Operation_time int    `json:"operation_time"`
}

type resultData struct {
	ID     string `json:"id"`
	Result string `json:"result"`
}

var ExpressionsMap = make(map[string]Values)
var QueueTask = make(map[string]SendValues)

var TIME_ADDITION_MS = 100
var TIME_SUBTRACTION_MS = 100
var TIME_MULTIPLICATION_MS = 500
var TIME_DIVISIONS_MS = 500

var Slises_easyExpr [][]string

var ExpressionToRPN = make(map[string][]string)

var mu sync.Mutex

func IsOperator(r rune) bool {
	return r == '+' || r == '-' || r == '*' || r == '/'
}
func GetExpression(w http.ResponseWriter, r *http.Request) (string, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {

		return "", fmt.Errorf("ошибка чтения тела запроса: %w", err)
	}
	expression := string(body)

	err = NormalExpression(expression)
	if err != nil {
		return "", err
	}

	return expression, nil
}

func GenerateRandomID(length int, large string) (string, error) {
	bytes := make([]byte, length/2)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	randomID := hex.EncodeToString(bytes)

	if large == "small" {
		return "_id" + randomID, nil
	}
	return "id" + randomID, nil
}

/*func AddToMap(expression, id, status, result string) {
	if id == "" {
		log.Fatal("Ошибка: ID не может быть пустым")
		return
	}

	val := ExpressionsMap[id]

	if expression != "" {
		val.Expression = expression
	}
	if status != "" {
		val.Status = status
	}
	if result != "" {
		val.Result = result
	}

	ExpressionsMap[id] = val
}*/

func AddToQueueTaskMap(id, arg1, arg2, operation string, operation_time int) {
	if id == "" {
		log.Fatal("Ошибка: ID не может быть пустым")
		return
	}

	sval := QueueTask[id]

	sval.Id = id

	if arg1 != "" {
		sval.Arg1 = arg1
	}
	if arg2 != "" {
		sval.Arg2 = arg2
	}
	if operation != "" {
		sval.Operation = operation
	}
	if operation_time != 0 {
		sval.Operation_time = operation_time
	}

	QueueTask[id] = sval
}

func Expressions(w http.ResponseWriter, r *http.Request) {
	if UserID == "" {
		http.Error(w, `{"Войдите в акаунт или харегестрируйтесь."}`, http.StatusUnauthorized)
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/v1/expressions")
	if path == "" {
		lenExpressions, err := WorkWithSQL.LenExpresions()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		if lenExpressions == 0 {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintln("выражения отсутствуют!")))
			return
		}
		w.WriteHeader(http.StatusOK)
		for key, i := range ExpressionsMap {
			w.Write([]byte(fmt.Sprintln(key, i)))
			fmt.Println(key, i)
		}
		return
	}

	id := strings.TrimPrefix(path, "/")

	OutExpression, err := WorkWithSQL.GetExpression(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintln("id не найден!")))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintln(OutExpression)))
	fmt.Println(OutExpression)
}

func NormalExpression(calculate string) error {
	if calculate == "" {
		return errors.New("параметр 'calculate' отсутствует или пуст")
	}

	openParentheses := 0
	lastCharWasOperator := false
	lastCharWasDot := false

	for i, r := range calculate {
		if unicode.IsDigit(r) {
			lastCharWasOperator = false
			lastCharWasDot = false
			continue
		}

		if r == '.' {
			if i == 0 || i == len(calculate)-1 || lastCharWasOperator || calculate[i-1] == '(' {
				return errors.New("некорректное использование десятичной точки")
			}

			if lastCharWasDot {
				return errors.New("число не может содержать две точки подряд")
			}

			lastCharWasDot = true
			continue
		}

		if unicode.IsSpace(r) {
			continue
		}

		if IsOperator(r) {
			// Оператор не может быть первым, если это не унарный минус
			if i == 0 && r != '-' {
				return errors.New("выражение не может начинаться с оператора (кроме унарного минуса)")
			}

			// Оператор не может быть последним
			if i == len(calculate)-1 {
				return errors.New("выражение не может заканчиваться оператором")
			}

			// Оператор * или / не может идти сразу после (
			if i > 0 && calculate[i-1] == '(' && (r == '*' || r == '/') {
				return errors.New("'*' или '/' не могут идти сразу после '('")
			}

			// Два оператора подряд запрещены (кроме унарного минуса)
			if lastCharWasOperator && r != '-' {
				return errors.New("два оператора подряд недопустимы")
			}

			lastCharWasOperator = true
			lastCharWasDot = false
			continue
		}

		if r == '(' {
			openParentheses++
			lastCharWasOperator = false
			lastCharWasDot = false
			continue
		}

		if r == ')' {
			if openParentheses == 0 {
				return errors.New("неверное количество закрывающих скобок")
			}
			openParentheses--
			lastCharWasOperator = false
			lastCharWasDot = false
			continue
		}

		return errors.New("в параметре 'calculate' присутствуют невалидные знаки")
	}

	if openParentheses != 0 {
		return errors.New("неверное количество открывающих скобок")
	}

	if lastCharWasDot {
		return errors.New("число не может заканчиваться точкой")
	}

	return nil
}

func Analys(rpnExpression []string) ([]int, []int) {
	count := 0
	itemToDaleyt := []int{}
	itemToID := []int{}
	for item, _ := range rpnExpression {
		_, err := strconv.ParseFloat(rpnExpression[item], 64)
		if err != nil {
			fmt.Errorf("Ошибка определения числа strconv.ParseFloat(): %v", err)
		}
		if err == nil {
			count++
		} else {
			count--
		}
		if count == 2 {
			_, err = strconv.ParseFloat(rpnExpression[item+1], 64)
			if err == nil {
				count--
			} else {
				Slises_easyExpr = append(Slises_easyExpr, []string{rpnExpression[item-1], rpnExpression[item], rpnExpression[item+1]})

				itemToDaleyt = append(itemToDaleyt, item-1, item)
				itemToID = append(itemToID, item+1)

				count = 0
			}
		}
	}
	return itemToDaleyt, itemToID
}

func ToRPN(expression string) ([]string, error) {
	var output []string
	var operators []rune

	precedence := map[rune]int{
		'+': 1,
		'-': 1,
		'*': 2,
		'/': 2,
	}

	for i := 0; i < len(expression); i++ {
		char := rune(expression[i])

		// Обрабатываем унарный минус (если он стоит в начале или после оператора/скобки)
		if char == '-' && (i == 0 || (!unicode.IsDigit(rune(expression[i-1])) && expression[i-1] != ')')) {
			number := "-"
			i++

			for i < len(expression) && (unicode.IsDigit(rune(expression[i])) || expression[i] == '.') {
				number += string(expression[i])
				i++
			}
			i--
			output = append(output, number)
		} else if unicode.IsDigit(char) || char == '.' {
			// Собираем обычное число
			number := string(char)
			for i+1 < len(expression) && (unicode.IsDigit(rune(expression[i+1])) || expression[i+1] == '.') {
				i++
				number += string(expression[i])
			}
			output = append(output, number)
		} else if char == '(' {
			operators = append(operators, char)
		} else if char == ')' {

			for len(operators) > 0 && operators[len(operators)-1] != '(' {
				output = append(output, string(operators[len(operators)-1]))
				operators = operators[:len(operators)-1]
			}

			if len(operators) > 0 && operators[len(operators)-1] == '(' {
				operators = operators[:len(operators)-1]
			}
		} else if _, exists := precedence[char]; exists {

			for len(operators) > 0 && precedence[operators[len(operators)-1]] >= precedence[char] {
				output = append(output, string(operators[len(operators)-1]))
				operators = operators[:len(operators)-1]
			}
			operators = append(operators, char)
		}
	}

	for len(operators) > 0 {
		output = append(output, string(operators[len(operators)-1]))
		operators = operators[:len(operators)-1]
	}

	return output, nil
}

func DivideID(id string, not string) (string, string) {
	for i := range id {
		if id[i] == '_' {
			ID := id[:i]
			if not == "Not" {
				return ID, ""
			}
			newid := id[:0] + id[i:]
			return ID, newid
		}
	}
	return "", id
}

func IDLocation(expressionID []string) []int {
	idIndex := []int{}

	for index, str := range expressionID {
		hasLetter := false
		hasDigit := false

		for _, char := range str {
			if unicode.IsLetter(char) {
				hasLetter = true
			}
			if unicode.IsDigit(char) {
				hasDigit = true
			}
			if hasLetter && hasDigit {
				idIndex = append(idIndex, index)
				break
			}
		}
	}

	return idIndex
}

func Orchestrator(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	if UserID == "" {
		http.Error(w, `{"Войдите в акаунт или харегестрируйтесь."}`, http.StatusUnauthorized)
	}

	expression, err := GetExpression(w, r)
	if err != nil {
		log.Printf("ошибка парсинга выражения: %v", err)

		http.Error(w, `{"error": "Некорректное выражение"}`, http.StatusUnprocessableEntity)
		return
	}
	log.Println("340. ПОЛУЧИЛ ВЫРАЖЕНИЕ:", expression)
	ExpressionID, err := GenerateRandomID(10, "")
	if err != nil {
		log.Printf("Ошибка при генерации ID: %v", err)

		http.Error(w, `{"error": "Ошибка генерации ID"}`, http.StatusInternalServerError)
		return
	}

	lenExpr, err := WorkWithSQL.LenExpresions()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if lenExpr > 0 {
		findExprID, err := WorkWithSQL.FindSaymExpression(expression)
		if err == nil {
			w.Write([]byte(fmt.Sprint("Это выражение уже было посчитано! Вот его ID:", findExprID)))
			return
		}
	}

	// Добавляем в карту
	err = WorkWithSQL.InsertExpresions(UserID, ExpressionID, expression, 1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Отправляем ответ
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Уникальный id на ваше выражение: %s\n", ExpressionID)))

	// Преобразуем в RPN
	ExpressionToRPN[ExpressionID], err = ToRPN(expression)
	if err != nil {
		log.Printf("Ошабка при переводе в RPN: %v\n", err)

		http.Error(w, `{"error": "Ошабка при переводе в RPN"}`, http.StatusInternalServerError)
		return
	}
	log.Println("363. ПЕРЕВЁЛ В RPN:", ExpressionToRPN)

	itemsToDalayte, itemToID := Analys(ExpressionToRPN[ExpressionID])
	log.Println("366. РАЗБИЛ НА ПРОСТЫЕ ЗАДАЧИ:", Slises_easyExpr)

	// Создаем задачи
	for i := range Slises_easyExpr {
		id, err := GenerateRandomID(6, "small")
		if err != nil {
			http.Error(w, fmt.Sprintf("Ошибка при генерации id для простого выражения: %v\n", err), http.StatusInternalServerError)
			return
		}
		switch Slises_easyExpr[i][2] {
		case "+":
			AddToQueueTaskMap(ExpressionID+id, Slises_easyExpr[i][0], Slises_easyExpr[i][1], Slises_easyExpr[i][2], TIME_ADDITION_MS)
		case "-":
			AddToQueueTaskMap(ExpressionID+id, Slises_easyExpr[i][0], Slises_easyExpr[i][1], Slises_easyExpr[i][2], TIME_SUBTRACTION_MS)
		case "*":
			AddToQueueTaskMap(ExpressionID+id, Slises_easyExpr[i][0], Slises_easyExpr[i][1], Slises_easyExpr[i][2], TIME_MULTIPLICATION_MS)
		case "/":
			if Slises_easyExpr[i][1] == "0" {
				log.Println("деление на ноль!")

				http.Error(w, `{"error": "нелзя делить на ноль!"}`, http.StatusUnprocessableEntity)
				return
			}
			AddToQueueTaskMap(ExpressionID+id, Slises_easyExpr[i][0], Slises_easyExpr[i][1], Slises_easyExpr[i][2], TIME_DIVISIONS_MS)
		}
		ExpressionToRPN[ExpressionID][itemToID[i]] = id
	}

	Slises_easyExpr = nil
	for i := len(itemsToDalayte) - 1; i >= 0; i-- {
		index := itemsToDalayte[i]
		ExpressionToRPN[ExpressionID] = append(ExpressionToRPN[ExpressionID][:index], ExpressionToRPN[ExpressionID][index+1:]...)
	}
}

func OrchestratorReturn(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	log.Println("405. Вызван OrchestratorReturn", r.Method)

	if r.Method == http.MethodGet {
		if len(QueueTask) == 0 {
			log.Println("Нет доступных задач для агента")

			http.Error(w, `{"error": "Нет доступных задач для агента"}`, http.StatusNotFound)
			return
		}

		var selectedTask SendValues
		var taskID string
		for id, task := range QueueTask {
			selectedTask = task
			taskID = id
			break
		}
		delete(QueueTask, taskID)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(selectedTask)

	}
	if r.Method == http.MethodPost {
		// Если POST, получаем результат и сохраняем его
		log.Println("493. Получен POST-запрос с результатом")

		var task resultData
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			log.Printf("Ошибка декодирования JSON: %v\n", err)

			http.Error(w, `{"error": "Ошибка декодирования"}`, http.StatusInternalServerError)
			return
		}

		fmt.Println("ПОЛУЧИЛ РЕЗУЛЬТАТ:", task.ID, task.Result)

		// Обновление результата
		ID, _ := DivideID(task.ID, "")
		idLocation := IDLocation(ExpressionToRPN[ID])
		log.Println("место ID и выражение:", idLocation, ExpressionToRPN[ID])
		updated := false
		for _, i := range idLocation {
			ID2, smlid := DivideID(task.ID, "")
			log.Println("сравниваем ExpressionToRPN[ID2][i]:", ExpressionToRPN[ID2][i], "и smlid:", smlid)
			if ExpressionToRPN[ID2][i] == smlid {
				log.Printf("Обновляем ID=%s на значение %s\n", task.ID, task.Result)
				ExpressionToRPN[ID2][i] = task.Result
				updated = true
			}
		}
		if !updated {
			log.Panicf("459. Не найдено совпадение ID в ExpressionToRPN[ID]")
		}

		log.Println("ПОДСТАВИЛ РЕЗУЛЬТАТ В ВЫРАЖЕНИЕ:", ExpressionToRPN[ID])

		idLocation = IDLocation(ExpressionToRPN[ID])
		if len(idLocation) == 0 {
			if len(ExpressionToRPN[ID]) == 1 {
				//записываем ответ в основную мапу!
				ID2, _ := DivideID(task.ID, "")
				err := WorkWithSQL.UpdateExpressionResult(ID2, ExpressionToRPN[ID2][0], 3)
				if err != nil {
					log.Println("Ошибка записи ответа в SQL:", err)
					http.Error(w, `{"Ошибка записи ответа."}`, http.StatusInternalServerError)
				}
				log.Println("записываем ответ в основную мапу!", ID2)
				return
			}
			itemsToDalayte, itemToID := Analys(ExpressionToRPN[ID])
			log.Println("Повторно разбиваем на простые задачи:", Slises_easyExpr)
			log.Printf("ИНДЕКСЫ НА УДАЛЕНИЕ: %v. ИНДЕКСЫ ПОД ID: %v\n", itemsToDalayte, itemToID)

			ID2, _ := DivideID(task.ID, "")
			// Создаем задачи
			for i := range Slises_easyExpr {
				id, err := GenerateRandomID(6, "small")
				if err != nil {
					log.Printf("Ошибка при генерации id для простого выражения: %v\n", err)

					http.Error(w, `{"error": "Ошибка генерации ID"}`, http.StatusInternalServerError)
					return
				}
				switch Slises_easyExpr[i][2] {
				case "+":
					AddToQueueTaskMap(ID2+id, Slises_easyExpr[i][0], Slises_easyExpr[i][1], Slises_easyExpr[i][2], TIME_ADDITION_MS)
				case "-":
					AddToQueueTaskMap(ID2+id, Slises_easyExpr[i][0], Slises_easyExpr[i][1], Slises_easyExpr[i][2], TIME_SUBTRACTION_MS)
				case "*":
					AddToQueueTaskMap(ID2+id, Slises_easyExpr[i][0], Slises_easyExpr[i][1], Slises_easyExpr[i][2], TIME_MULTIPLICATION_MS)
				case "/":
					if Slises_easyExpr[i][1] == "0" {
						log.Println("деление на ноль!")

						http.Error(w, `{"error": "нелзя делить на ноль!"}`, http.StatusUnprocessableEntity)
						return
					}
					AddToQueueTaskMap(ID2+id, Slises_easyExpr[i][0], Slises_easyExpr[i][1], Slises_easyExpr[i][2], TIME_DIVISIONS_MS)
				}
				Slises_easyExpr[i] = []string{id}
				ExpressionToRPN[ID2][itemToID[i]] = id
			}
			for i, item := range QueueTask {
				log.Printf("РАЗБИТЫЕ ЗАДАЧИ: %v, ПО ID: %v\n", item, i)
			}

			Slises_easyExpr = nil
			for i := len(itemsToDalayte) - 1; i >= 0; i-- {
				index := itemsToDalayte[i]
				ExpressionToRPN[ID2] = append(ExpressionToRPN[ID2][:index], ExpressionToRPN[ID2][index+1:]...)
			}
			log.Println("ВЫРАЖЕНИЕ ПОСЛЕ УДАЛЕНИЯ И ПОДСТАНОВКИ ID:", ExpressionToRPN[ID2])
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Результат принят"))
	}
}

func Registrate(w http.ResponseWriter, r *http.Request) {
	userID, err, statusKode := WorkWithSQL.Registrate(w, r)
	if err != nil {
		http.Error(w, `{"Ошибка регестрации": err}`, statusKode)
	}
	UserID = userID
	w.Write([]byte("Вы успешно зарегестрировались! :)"))
}

func Login(w http.ResponseWriter, r *http.Request) {
	userID, err, statusKode := WorkWithSQL.Login(w, r)
	if err != nil {
		http.Error(w, `{"Ошибка входа": err}`, statusKode)
	}
	UserID = userID
	w.Write([]byte("Вы успешно вошли в аккаунт! :)"))
}
