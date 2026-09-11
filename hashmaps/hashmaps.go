package main

import (
	"fmt"
)

type user struct {
	username      string
	password      string
	occupied      bool
	nextUserindex int
}
type hashmap struct {
	initial   []user
	secondary []user
}

func searchmap(username string, hash func(string) uint32, hmap hashmap) (user, error) {
	value := int(hash(username)) % len(hmap.initial)
	fmt.Println("Value is: ", value)
	secondaryindexlayer := 0
	var u string = hmap.initial[value].username
	var index int = hmap.initial[value].nextUserindex
	if u != username && index != -1 && u != "" {
		fmt.Println("New second layer")

		for u != username {
			secondaryindexlayer++
			u = hmap.secondary[index].username
			if u != username && hmap.secondary[index].nextUserindex != -1 {
				index = hmap.secondary[index].nextUserindex
			} else if u == username {
				fmt.Println(username, "found", secondaryindexlayer, "layer(s) deep")
				return hmap.secondary[index], nil
			} else if hmap.secondary[index].nextUserindex == -1 {
				fmt.Println("2nd Layer search triggered, but user not found")
				break
			}
		}
		fmt.Println("Error")
		return hmap.secondary[index], fmt.Errorf("User Not Found")
	} else if (index == -1 || u == "") && u != username {
		fmt.Println("No 2nd level chain found, and user not found.", username)
		return hmap.initial[value], fmt.Errorf("User Not Found")
	} else {
		fmt.Println("Found in Initial")
		return hmap.initial[value], nil
	}
}

func addtomap(u user, hash func(string) uint32, hmap hashmap, collisioncount int) (int, hashmap) {

	value := hash(u.username)
	u.occupied = true
	u.nextUserindex = -1
	index := int(value) % len(hmap.initial)
	if !hmap.initial[index].occupied {
		hmap.initial[index] = u

	} else {
		collisioncount++

		hmap.secondary = append(hmap.secondary, u)

		last := len(hmap.secondary) - 1
		var currentUindex = hmap.initial[index].nextUserindex
		if currentUindex != -1 {
			var nextUindex = hmap.secondary[currentUindex].nextUserindex

			for nextUindex != -1 {
				currentUindex = nextUindex
				nextUindex = hmap.secondary[nextUindex].nextUserindex
			}
			hmap.secondary[currentUindex].nextUserindex = last
		} else {
			hmap.initial[index].nextUserindex = last
		}
	}
	return collisioncount, hmap

}

func getloadpercentage(hashmap []user, loadcount int) float64 {
	if len(hashmap) == 0 {
		fmt.Println("Load factor is at: 0 %")
		return 0
	}
	loadfactor := float64(loadcount) / float64(len(hashmap))
	return loadfactor
}

func resize(hmap hashmap, scalingfactor int) (int, hashmap) {
	nhmap := hashmap{make([]user, len(hmap.initial)*scalingfactor), make([]user, 0, 1000)}
	collisioncount := 0

	for i := 0; i < len(hmap.initial); i++ {
		if hmap.initial[i].occupied {

			collisioncount, nhmap = addtomap(hmap.initial[i], djb2, nhmap, collisioncount)
		}
	}
	for i := 0; i < len(hmap.secondary); i++ {
		if hmap.secondary[i].occupied {
			collisioncount, nhmap = addtomap(hmap.secondary[i], djb2, nhmap, collisioncount)
		}
	}

	return collisioncount, nhmap
}
