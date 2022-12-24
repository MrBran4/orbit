package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/mrbran4/orbit"
)

// ---
/// These our our app's types. We only care about users + events in this example.
/// These types aren't solely for the API, they might mirror what's in our DB.

// Represents a user in our system
type user struct {
	Uid  int
	Name string
}

// FromRequest resolves a user based on the UID provided in an API call.
// For example, User.FromRequest("125") will return the User with ID 125, or an error.
func (u user) FromRequest(uid string) (any, error) {
	// In a real app you'd use the database here
	uidInt, _ := strconv.Atoi(uid)
	return user{
		Uid:  uidInt,
		Name: "Joe Bloggs",
	}, nil
}

// Represents an event in our system
type event struct {
	EventID       int      `json:"event_id"`
	AttendeeNames []string `json:"attendees"`
}

// FromBody resolves an event from the body of a request, by JSON-decoding it.
func (e event) FromBody(body io.ReadCloser) (any, error) {
	var result event
	_ = json.NewDecoder(body).Decode(&result)
	return result, nil
}

/// ---
/// Make your handler

// Handles POST /user/{user}/event/{event}
// with a json-encoded User in the body.
func myNewEventHandler(
	w http.ResponseWriter,
	r *http.Request,
	params orbit.RouteParams,
	body orbit.FromBodyable,
) {
	// The magic of Orbit is that if your handler gets called, the params and body
	// are guaranteed to be present and of the correct type you specify when wiring
	// up the handler to the route.
	//
	// This means you can blindly assert types like this:
	usr, _ := params["user"].(user)
	ename, _ := params["event"].(orbit.BasicString)
	evt, _ := body.(event)

	// Your app logic
	fmt.Printf("User: %+v\nEvent Key (from url): %s\nEvent (from body): %+v\n", usr, ename, evt)

	w.WriteHeader(200)
}

func main() {

	// Make a router and wire the handler up to it.
	// This is where we tell Orbit what types our arguments are.
	r := orbit.NewRouter()

	r.Handle(
		// The path to match against
		"/user/{user}/event/{event}",
		// your orbit.Handler
		orbit.HandlerFunc(myNewEventHandler),
		// []string of HTTP methods to match. Nil means all methods.
		[]string{"POST"},
		// Param types to decode (types must implement FromRequestable)
		orbit.RouteParams{
			"user":  user{},
			"event": orbit.BasicString(""), // helper: just a string that implements FromRequest.
		},
		// Body type (must implement FromBodyable)
		event{},
	)

	// Precompiles regexes etc for the router.
	// Call this once before starting up.
	r.Bake()

	log.Fatal(http.ListenAndServe(":8080", r))

}

// Now test with:
//   curl -XPOST http://localhost:8080/user/123/event/my-new-event --data '{"event_id": 123, "attendees": ["Amy", "Betty", "Cressida"]}'
//
// Output (in terminal):
//   User: {Uid:123 Name:Joe Bloggs}
//   Event Key (from url): my-new-event
//   Event (from body): {EventID:123 AttendeeNames:[Amy Betty Cressida]}
