package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/deeptech-kmitl/Cicero-Backend/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type IRiLogger interface {
	Print() IRiLogger
	Save()
	setResponse(res any)
}

type RiLogger struct {
	Time       string `json:"time"`
	Ip         string `json:"ip"`
	Method     string `json:"method"`
	StatusCode int    `json:"status_code"`
	Path       string `json:"path"`
	Response   any    `json:"response"`
}

func InitRiLogger(c *fiber.Ctx, res any) IRiLogger {
	log := &RiLogger{
		Time:       time.Now().Local().Format("2006-01-02 15:04:05"),
		Ip:         c.IP(),
		Method:     c.Method(),
		Path:       c.Path(),
		StatusCode: c.Response().StatusCode(),
	}
	log.setResponse(res)
	return log
}

// เพื่อ print log ออกทาง console
func (l *RiLogger) Print() IRiLogger {
	utils.Debug(l)
	return l

}

// เพื่อบันทึก log ลงในไฟล์
func (l *RiLogger) Save() {
	data := utils.Output(l)

	fileName := fmt.Sprintf("./assets/logs/rilogger_%v.txt", strings.ReplaceAll(time.Now().Format("2006-01-02"), "-", ""))
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()
	file.WriteString(string(data) + "\n")
}

// เพื่อเก็บ log ของ response ที่ส่งไปให้ client
func (l *RiLogger) setResponse(res any) {
	l.Response = res
}
