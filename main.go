package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {
	var guess bool
	flag.BoolVar(&guess, "guess", false, "Insert best guess about attendace as if it were factual.")
	flag.Parse()

	file := flag.Args()[0]
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	csvReader := csv.NewReader(f)
	people, parties, err := parseCSV(csvReader, guess)
	if err != nil {
		panic(err)
	}

	host := strings.TrimRight(flag.Args()[1], "/")
	peopleBody, err := json.Marshal(peopleReq{People: people})
	if err != nil {
		panic(err)
	}
	b := bytes.NewBuffer(peopleBody)
	r, err := http.NewRequest("PUT", host+"/people", b)
	if err != nil {
		panic(err)
	}
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		panic(resp.Status + string(body))
	}

	partyBody, err := json.Marshal(partiesReq{Parties: parties})
	if err != nil {
		panic(err)
	}
	b = bytes.NewBuffer(partyBody)
	r, err = http.NewRequest("PUT", host+"/parties", b)
	if err != nil {
		panic(err)
	}
	resp, err = http.DefaultClient.Do(r)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		panic(resp.Status + string(body))
	}
}

type peopleReq struct {
	People []Person `json:"people"`
}

type partiesReq struct {
	Parties []Party `json:"parties"`
}

func parseCSV(r *csv.Reader, guess bool) ([]Person, []Party, error) {
	parties := map[string]Party{}
	people := map[string]Person{}
	_, err := r.Read() // ignore the first line
	if err != nil {
		return nil, nil, err
	}
	var plusOnes []Person
	partyPeople := map[string][]Person{}
	for {
		col, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, err
		}
		firstName, lastName := col[0], col[1]
		partyFName, partyLName := col[2], col[3]
		partyID := col[5]
		id := col[11]
		attendingGuess := col[10]
		person := Person{
			ID:      id,
			PartyID: partyID,
			Name:    firstName,
		}
		if guess && attendingGuess != "Maybe" {
			person.Replied = true
			person.Reply = (attendingGuess == "Definitely" || attendingGuess == "Probably")
		}
		if lastName == "" && strings.HasSuffix(firstName, " +1") {
			plusOnes = append(plusOnes, person)
			continue
		}
		person.Name = person.Name + " " + lastName
		people[id] = person
		if partyFName == firstName && partyLName == lastName {
			var codeRunes []rune
			for _, char := range strings.ToLower(col[7]) {
				if strings.ContainsRune("abcdefghijklmnopqrstuvwxyz0123456789", char) {
					codeRunes = append(codeRunes, char)
				}
			}
			codeWord := string(codeRunes)
			address := col[6]
			partySortValue := col[4]
			numPeople := parties[partyID].NumPeople + 1
			party := Party{
				ID:        partyID,
				LeadID:    id,
				Name:      firstName + " " + lastName + "’s Party",
				NumPeople: numPeople,
				Address:   address,
				MagicWord: codeWord,
				SortValue: partySortValue,
			}
			if person.Replied && person.Reply {
				party.NumComing++
			}
			parties[partyID] = party
		} else {
			party := parties[partyID]
			party.NumPeople++
			parties[partyID] = party
		}
		partyPeople[partyID] = append(partyPeople[partyID], person)
	}
	for _, plusOne := range plusOnes {
		party, ok := parties[plusOne.PartyID]
		if !ok {
			fmt.Println("Plus one", plusOne.ID, "has an invalid party ID")
			continue
		}
		party.NumPeople++
		if plusOne.Replied && plusOne.Reply {
			party.NumComing++
		}
		parties[party.ID] = party
		personName := strings.TrimSuffix(plusOne.Name, "'s +1")
		var found bool
		for _, person := range partyPeople[party.ID] {
			if strings.Split(person.Name, " ")[0] == personName {
				found = true
				person.GetsPlusOne = true
				people[person.ID] = person
				break
			}
		}
		if !found {
			fmt.Println("Plus one", plusOne.ID, "doesn't belong to anyone...")
		}
	}
	resPeople := make([]Person, 0, len(people))
	for _, person := range people {
		resPeople = append(resPeople, person)
	}
	resParties := make([]Party, 0, len(parties))
	for _, party := range parties {
		resParties = append(resParties, party)
	}
	return resPeople, resParties, nil
}
