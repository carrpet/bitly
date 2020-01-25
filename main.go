package main

import (
	"bitly/client"
	"fmt"
)

// The problem:

// To get all the bitlinks for a user's default group then for each bitlink
// get the number of user clicks by country, add all the clicks up by country
// then we divide those numbers by 30 and return a list of countries and average
//clicks

//Pseudocode and naive implementation:
/*
countrytoclickshashtable(ctcht) = {}
groupId = GET(user)
bitlinks = GET(groups/{groupId}/bitlinks)
for each link in bitlinks:
     clicksbycountry = GET(bitlinks/{link}/countries)
     ctcht{clicksbycountry[country]} += clicksbycountry[clicks]
return map(divideclicksby30,ctcht)

*/
func main() {
	token := "5ad8274a49bcd964f23d4b685c272c37de718711"
	bc := client.BitlyClientInfo{Token: token}

	userInfo, err := client.GetUserInfo(bc)
	if err != nil {
		panic(err)
	}
	fmt.Printf("UserInfo has name: %s and GroupGuid: %s", userInfo.Name, userInfo.GroupGuid)

	// retrieve all the links and stick into a hashtable
	groupslinks, err := client.GetBitlinksForGroup(bc, userInfo.GroupGuid)
	if err != nil {
		panic(err)
	}
	fmt.Printf("GroupsBitlink has link: %s and id: %s", groupslinks.Links[0].Link, groupslinks.Links[0].ID)

}
