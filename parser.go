package main

import "strings"

type parsedKudoCommand struct {
	Members []Member
	Text    string
}

func parseKudoCommand(kudoText string) parsedKudoCommand {
	parsed := parsedKudoCommand{}

	if len(kudoText) == 0 {
		return parsed
	}

	lastIndex := 0
	exploded := strings.Split(kudoText, " ")

	for index, target := range exploded {
		targetStr := strings.TrimLeft(target, "@")
		targetStr = strings.Replace(targetStr, ",", "", -1)
		member, err := findMemberByTag(targetStr)
		if err != nil {
			continue
		}

		lastIndex = index
		parsed.Members = append(parsed.Members, member)
	}

	lastMemberFound := exploded[lastIndex]
	startOfMessage := strings.LastIndex(kudoText, lastMemberFound) + len(lastMemberFound) + 1
	endOfMessage := len(kudoText)

	if startOfMessage < endOfMessage {
		parsed.Text = kudoText[startOfMessage:endOfMessage]
	} else {
		parsed.Text = ""
	}

	return parsed
}
