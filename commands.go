package main

import (
	"fmt"
	"net/http"
	"strings"
)

func handleNewKudoCommand(w http.ResponseWriter, memberFrom *Member, target string, extra string) {
	member, err := findMemberByTag(strings.TrimLeft(target, "@"))
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	kudo := Kudo{0, extra, member, memberFrom, 0}
	kudos = append(kudos, kudo)
	dbSaveKudo(&kudo)

	notifyUser("New kudo from <@"+memberFrom.ID+">!\n"+extra, *member)
	notifyChannel("New kudo from <@" + memberFrom.ID + "> was given to <@" + member.ID + ">!\n" + extra)

	fmt.Fprint(w, "Kudo has been registered!")
}

func handleLikeCommand(w http.ResponseWriter, memberFrom *Member, target string) {
	member, err := findMemberByTag(strings.TrimLeft(target, "@"))
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	for i := len(kudos) - 1; i >= 0; i-- {
		if kudos[i].MemberTo.ID == member.ID {
			if dbKudoLiked(kudos[i], *memberFrom) {
				fmt.Fprint(w, "You have already liked this.")
				return
			}

			kudos[i].LikeCount++
			dbUpdateKudoLikes(kudos[i])
			addMemberToLike(kudos[i], memberFrom)

			fmt.Fprint(w, "Thank you!")
			notifyUser(fmt.Sprint("<@", memberFrom.ID, "> likes your kudo! Total likes: ", kudos[i].LikeCount, "\n", kudos[i].Kudo), *member)

			return
		}
	}

	fmt.Fprint(w, "Found nothing to like...")
}
