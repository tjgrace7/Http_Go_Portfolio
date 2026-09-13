package main

import (
	"fmt"
	"math/rand"
	"testing"
)

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

// Collision Chain Testing
func TestCollisionChain(t *testing.T) {
	//Change this variable to one of the variables in the slice to see where it lands in the hashmap. If you use the djb2 algorithm, eve should iterate through the 2nd slice in the hashmap twice before returning
	users := []user{
		{username: "tjgrace", password: "12345"},
		{username: "tyler", password: "12345"},
		{username: "eve", password: "12345"},
		{username: "john", password: "12354"},
		{username: "elvis", password: "12345"},
		{username: "willis", password: "12345"},
		{username: "August", password: "12345"},
		{username: "Hannah", password: "12345"},
		{username: "george", password: "12345"},
		{username: "will", password: "12345"},
	}
	collisioncount := 0
	var hmap hashmap = hashmap{make([]user, 10), make([]user, 0, 100)}
	for _, u := range users {
		collisioncount, hmap = addtomap(u, djb2, hmap, collisioncount)
	}
	testusernams := []string{
		//A layer is the number of times a nextUserIndex has been called from the user struct
		//Eva will return a user who is two layers deep in the hash map
		//Henry will return a user who is not in the Hashmap, but will match a Hash code with eve, tyler, and tjgrace. This will iterate through the 2nd slice several times
		//tjgrace will be found directly in the initial map
		"eve", "Henry", "tjgrace", "mark",
	}
	for _, u := range testusernams {
		t.Run(u, func(t *testing.T) {
			username, err := searchmap(u, djb2, hmap)

			if u == "Henry" || u == "mark" {
				if err == nil {
					t.Errorf("expected error for %q, git user %v", u, username)
				}
			} else if err != nil {
				t.Errorf("Error Searching for User %q", u)
			}
			fmt.Println("Username", username.username)
		})
	}

}
func TestResizeSearch(t *testing.T) {
	var hmap hashmap = hashmap{make([]user, 10), make([]user, 0, 100)}
	collisioncount := 0
	users := []user{
		{username: "tjgrace", password: "12345"},
		{username: "tyler", password: "12345"},
		{username: "eve", password: "12345"},
		{username: "john", password: "12354"},
		{username: "elvis", password: "12345"},
		{username: "willis", password: "12345"},
		{username: "August", password: "12345"},
		{username: "Hannah", password: "12345"},
		{username: "george", password: "12345"},
		{username: "will", password: "12345"},
	}
	//This will resize the hashmap to fit users in a different spot and search for a new user
	for i := 0; i < len(users); i++ {
		collisioncount, hmap = addtomap(users[i], djb2, hmap, collisioncount)
		loadfactor := getloadpercentage(hmap.initial, i)
		//change loadfactor to .90 to avoid firing test
		if loadfactor > 0.90 {
			fmt.Println("Hash Map Resized")
			collisioncount, hmap = resize(hmap, 10)
		}

	}
	us, err := searchmap("eve", djb2, hmap)
	if err != nil {
		t.Error("Error finding user")
	}
	fmt.Println("User found", us)
}

func BenchmarkXxx(b *testing.B) {
	var hmap hashmap = hashmap{make([]user, 100000), make([]user, 0, 100)}
	var loopiterations int = 0
	collisioncount := 0
	var testuser user
	for b.Loop() {
		loopiterations++
		var u user
		i := rand.Intn(10000)
		u.username = randomString(15)
		u.password = "12345"
		collisioncount, hmap = addtomap(u, djb2, hmap, collisioncount)
		loadfactor := getloadpercentage(hmap.initial, loopiterations)
		if loadfactor > 0.75 {
			fmt.Println("Hash Map Resized")
			collisioncount, hmap = resize(hmap, 2)
		}
		//Takes user as test user for searching where the most recent iteration of i == 219 happens.
		//This would obviously not be used in a loading settings, but is useful for finding a random test user that is generated in this loop.
		if i == 219 {
			testuser = u
		}

	}
	if testuser.username == "" {
		fmt.Println("Error: No user entered into test. Random integer didn't hit i")
		return
	}
	fmt.Println("Test User is: ", testuser.username)
	fmt.Println("Loop Iterations", loopiterations)
	fmt.Println("collisions", collisioncount, "hashmap2 length: ", len(hmap.secondary))
	founduser, err := searchmap(testuser.username, djb2, hmap)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Found User: ", founduser.username, " Password: ", founduser.password)
}
