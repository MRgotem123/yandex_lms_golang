import (
"fmt"
"math/rand"
"time"
)

type User struct {
ID    int
Name  string
Email string
Age   int
}

type Report struct {
User     User
ReportID []int
Date     time.Time
}

func CreateReport(user User, reportDate string) (Report, error) {
date, err := time.Parse("2006-01-02", reportDate)
if err != nil {
return Report{}, fmt.Errorf("неправильный формат даты: %w", err)
}

digits := rand.Perm(10)[:10] // Уникальные числа от 0 до 9

report := Report{
User:     user,
ReportID: digits,
Date:     date,
}

return report, nil
}

func PrintReport(report Report) {
fmt.Printf(
"ID: %d \nName: %s \nEmail: %s \nAge: %d\nReportID: %v\nDate: %s\n",
report.User.ID, report.User.Name, report.User.Email, report.User.Age,
report.ReportID, report.Date.Format("2006-01-02"),
)
}

func GenerateUserReports(users []User, reportDate string) []Report {
var reports []Report
for _, user := range users {
report, err := CreateReport(user, reportDate)
if err != nil {
fmt.Printf("Ошибка при создании отчета для пользователя ID %d: %v\n", user.ID, err)
continue
}
reports = append(reports, report)
}
return reports
}