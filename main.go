package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/nexidian/gocliselect"
)

type ButtonActions struct {
	Id    string `json:"id"`
	Type  string `json:"type"`
	Url   string `json:"url"`
	Label string `json:"label"`
}
type RequestPayload struct {
	Message       string          `json:"message"`
	Phone         string          `json:"phone"`
	Title         string          `json:"title"`
	Footer        string          `json:"footer"`
	ButtonActions []ButtonActions `json:"buttonActions"`
	Document      string          `json:"document"`
	FileName      string          `json:"fileName"`
}

func SendMessage(payload RequestPayload, sendPdf bool) *http.Response {
	SUA_INSTANCIA := os.Getenv("SUA_INSTANCIA")
	SEU_TOKEN := os.Getenv("SEU_TOKEN")
	ACTION := "/send-text"
	if sendPdf {
		ACTION = "/send-document/pdf"
	}
	// ACTION := "/send-button-actions"
	fmt.Println("ACTION", ACTION)
	url := "https://api.z-api.io/instances/" + SUA_INSTANCIA + "/token/" + SEU_TOKEN + ACTION
	fmt.Println("url:", url)
	fmt.Println("payload.Message", payload.Message)
	json_data, _ := json.Marshal(payload)
	r, _ := http.Post(
		url,
		"application/json",
		bytes.NewBuffer(json_data),
	)
	return r
}

func ParseResponse(r *http.Response) []byte {
	readedBody, _ := io.ReadAll(r.Body)
	var res interface{}
	json.Unmarshal(readedBody, &res)
	body, _ := json.Marshal(&res)
	return body
}

func ParseMessage(msg string, element CandidatesStruct) string {
	messageChangedEmail := strings.Replace(msg, "@email", element.Email, 1)
	messageChangedName := strings.Replace(messageChangedEmail, "@nome", element.Name, 1)
	messageChangedPhone := strings.Replace(messageChangedName, "@phone", element.Phone, 1)
	messageChangedPass := strings.Replace(messageChangedPhone, "@pass", element.Password, 1)
	newMessage := strings.Replace(messageChangedPass, "@nome", element.Name, 1)
	return newMessage
}

func SelectFieldInPath(menuQuestion string) string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	path := filepath.Dir(ex)
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	files, err := f.Readdir(0)
	if err != nil {
		panic(err)
	}

	menu := gocliselect.NewMenu(menuQuestion)
	for _, v := range files {
		// v.Name()
		isXlsx := strings.Contains(v.Name(), ".xlsx")
		if !v.IsDir() && isXlsx {
			menu.AddItem(v.Name(), v.Name())
		}
	}
	choice := menu.Display()
	return choice
}

func GetInputFromTerminal(inputQuestion string) string {
	fmt.Println(inputQuestion)
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	return line
}

func main() {
	excelPath := SelectFieldInPath("Qual o Excel Quer Ler ?")
	fileToSend := GetInputFromTerminal("Quer enviar algum arquivo ? Cole o link a baixo")
	line := ""
	if len(fileToSend) > 0 {
		// line = GetInputFromTerminal("Digite o texto que deseja enviar (use @email, @nome, @phone e @pass para personalizar a mensagem):")
	}

	fmt.Printf("\nEsse é um exemplo de mensagem:\n")
	fmt.Printf("-----\n")
	SendMessagesToWhatsapp(excelPath, fileToSend, line)
}

func SendMessagesToWhatsapp(excelPath string, fileToSend string, message string) {
	candidates := ExelParser(excelPath)
	for _, element := range candidates {
		parsedMessage := ParseMessage(message, element)
		payload := RequestPayload{}
		payload.Phone = element.Phone
		payload.Message = parsedMessage
		n := rand.Intn(1000)
		time.Sleep(time.Duration(n) * time.Millisecond)

		response := SendMessage(payload, false)
		defer response.Body.Close()
		body := ParseResponse(response)
		var buf bytes.Buffer
		json.HTMLEscape(&buf, body)
		fmt.Println("buf.String()", buf.String())

		if len(fileToSend) > 0 {
			payload.Document = fileToSend
			payload.FileName = "Convite Oficial"
		}
		responsePdfMesage := SendMessage(payload, true)
		defer responsePdfMesage.Body.Close()
		bodyPdf := ParseResponse(responsePdfMesage)
		var bufPdf bytes.Buffer
		json.HTMLEscape(&bufPdf, bodyPdf)
		fmt.Println("buf.String()", bufPdf.String())
	}
}
