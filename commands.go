package main

import (
	"fmt"
	"net/http"
	"time"
)

func handleNewKudoCommand(w http.ResponseWriter, memberFrom Member, command string, value int) {

	kudo := Kudo{ID: 0, MemberFrom: memberFrom, Value: value, Color: randomColor(), Original: command, Date: time.Now()}

	parsed := parseKudoCommand(command)

	kudo.Text = parsed.Text
	kudo.Recipients = parsed.Members

	if len(kudo.Recipients) == 0 {
		printKudoUsage(w)
		return
	}

	if len(kudo.Text) == 0 {
		fmt.Fprint(w, "Please enter kudo message!")
		return
	}

	dbSaveKudo(&kudo)
	kudos = append(kudos, kudo)

	text := "New kudo from <@" + memberFrom.ID + ">!\n>" + kudo.Text

	if value == -1 {
		text = "New boo from <@" + memberFrom.ID + ">!\n>" + kudo.Text
	}

	for _, recip := range kudo.Recipients {
		notifyUser(text, recip)
	}

	if value == 1 {
		fmt.Fprint(w, "Kudo has been registered!")
	} else {
		fmt.Fprint(w, "Boo has been registered!")
	}
}
