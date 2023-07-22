package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}
var todoList []string

func getCmd(input string) string {
    inputArr := strings.Split(input, " ")
    return inputArr[0]
}

func getMessage(input string) string {
    inputArr := strings.Split(input, " ")
    var result string
    for i := 1; i < len(inputArr); i++ {
        result += inputArr[i]
    }
    return result
}

func updateTodoList(input string) {
    tmpList := todoList
    todoList = []string{}
    for _, val := range tmpList {
        if val == input {
            continue
        }
        todoList = append(todoList, val)
    }
}



func main() {
	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w,r,"websockets.html")
	})
	http.HandleFunc("/todo", func (w http.ResponseWriter, r *http.Request) {
		fmt.Println("setting up http server")
		conn, err := upgrader.Upgrade(w,r,nil)

		if err != nil {
			log.Fatal(err)
			return;
		}

		defer conn.Close()

		for {
			messageType, messageBytes, err := conn.ReadMessage()

			if err != nil {
				log.Fatal(err)
				return;
			}

			// if err := conn.WriteMessage(messageType, messageBytes); err != nil {
			// 	log.Println(err)
			// 	return
			// }

			message := string(messageBytes)
			cmd := getCmd(message)
			msg := getMessage(message)

			if cmd == "add" {
                todoList = append(todoList, msg)
            } else if cmd == "remove" {
                updateTodoList(msg)
            }

			output := "Current Todos: \n"
            for _, todo := range todoList {
                output += "\n - " + todo + "\n"
            }
            output += "\n----------------------------------------"
			messageBytes = []byte(output)

			err = conn.WriteMessage(messageType, messageBytes)
            if err != nil {
                log.Println("write failed:", err)
                break;
            }


		}
	
	})

	http.ListenAndServe(":8000", nil)
}

func initialiseHandShake() {

}