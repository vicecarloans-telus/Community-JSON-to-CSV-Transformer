package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	fmt.Println("Starting JSON Processor...")
	jsonF, err := os.Open("./response.json")

	if err != nil {
		panic("Oops...Cannot open JSON file")
	}

	fmt.Println("Open Reader..")

	defer jsonF.Close()

	b, _ := ioutil.ReadAll(jsonF)

	var r Response

	err = json.Unmarshal(b, &r)

	if err != nil {
		panic("Unable to unmarshal json...Please check the format")
	}

	fmt.Println("Creating CSV...")

	csvF, err := os.Create("data.csv")

	if err != nil {
		panic("Unable to create CSV file...")
	}

	defer csvF.Close()

	w := csv.NewWriter(csvF)

	defer w.Flush()

	messages := r.Data[0].Res.Items

	var content = [][2]string{{"Username", "Email"}}
	for _, message := range messages {
		fmt.Println("Adding...", message.Author.Login)
		err, c := addRow(content, &message.Author)
		if err == nil {
			content = c
		}

		kds := message.Kudos.Items
		if len(kds) > 0 {
			for _, kudo := range kds {
				fmt.Println("Adding...", kudo.User.Login)
				err, c = addRow(content, &kudo.User)
				if err == nil {
					content = c
				}
			}
		}
	}

	fmt.Println("Generating CSV...")

	for _, row := range content {
		err := w.Write(row[:])
		if err != nil {
			fmt.Println("Unable to write to csv...")
		}
	}
}

func addRow(content [][2]string, author *User) (error, [][2]string) {
	fmt.Println("Checking if user ", author.Login, " exists....")
	for _, row := range content {

		if row[0] == author.Login {
			fmt.Println("User ", author.Login, " exists....")
			return fmt.Errorf("User exists"), nil
		}
	}
	var cols [2]string
	cols[0] = author.Login
	cols[1] = author.Email
	content = append(content, cols)
	return nil, content
}

type Response struct {
	Data []ResponseData `json:"data"`
}

type ResponseData struct {
	Res MessagesResponse `json:"response"`
}

type MessagesResponse struct {
	Size  int        `json:"size"`
	Type  string     `json:"type"`
	Items []Messages `json:"items"`
}

type Messages struct {
	Author User          `json:"author"`
	Kudos  KudosResponse `json:"kudos"`
}

type User struct {
	Type  string `json:"user"`
	Email string `json:"email"`
	Login string `json:"login"`
}

type KudosResponse struct {
	Size  int    `json:"size"`
	Type  string `json:"type"`
	Items []Kudo `json:"items"`
}

type Kudo struct {
	Type string `json:"type"`
	User User   `json:"user"`
}
