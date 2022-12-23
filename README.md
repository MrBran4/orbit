# 🪐 Orbit

Orbit is a HTTP router for Go. It differs from other frameworks in that `/{params}/in/{the}/url` are converted to the type (any type) you actually want, as part of the framework handling the request. That is to say, instead of your handler recieving strings, it recieves
actual resolved full structs that can include anything you want, like pulling from the database or something.

It also guarantees that if your handler is called, then the data is present and valid. Orbit's request handling is a bit like the Rocket framework for Rust, hence the whacky space-themed name.

It's unfinished and not close to being usable yet, but the core 'getting params from the request as the right type' bit is working.

## Quick Example:

Say you've got an API route to get attendees of a specific event, owned by a specific user.
Your route template would look like `/user/{user}/event/{event}`, and an example of an
incoming request intended for that endpoint might be `GET /user/5/event/25`

You also have existing User and Event types in your app:

```go
type User struct {
    Uid int
    Name string
}

type Event struct {
    EventID int
    AttendeeNames []string
}
```

First, make your types meet the `orbit.FromRequestable` interface by implementing a `FromRequest` method. It's super simple, you're given the string from the param, and you return a valid _whatever_. For example:

```go
// FromRequest resolves a user based on the UID provided in an API call.
// For example, FromRequest(125) will return the User with ID 125, or an error.
func (u User)FromRequest(uid string) (any, error) {
    // - Cast UID to an int
    // - Look user up in database
    // - Populate a user struct with the data from the db
    return user, nil
}

// Same for Events...
func (e Event)FromRequest(eventID string) (any, error) {
    // ...
```

Ignore the use of Any here for now - Orbit uses reflection to make sure the returned type matches what you asked for.
See [golang/go#30602](https://github.com/golang/go/issues/30602) for why that would help and why it's unavoidable right now.

Next, just make your handler take an orbit.RouteParams argument. orbit.RequestParams is just an alias for `map[string]FromRequestable` - it'll contain the types from your request. You also need to provide the types when you wire up the handler.

Your handler might look like:

```go
func MyHandleFunc(w *http.ResponseWriter, r http.Request, params orbit.RouteParams) {
    // Do something.
}

r := orbit.NewRouter()
r.HandleGet(
    "/user/{user_id}/event/{event_id}",
    orbit.RouteParams{
        "user": User,
        "event": Event
    },
    MyHandleFunc,
)

```

When a request comes in matching `/user/{user_id}/event/{event_id}`, orbit will extract the _string_ parameters
from the path (the user_id and event_id bits), then it'll pass those to the corresponding FromRequest(...) funcs
you wrote earlier, storing the result in a new RouteParams struct which gets passed to your handler.

For completeness, that means for the request `/user/5/event/25`, orbit eventually winds up making this call to your code:

```go
MyHandleFunc(
    w /* http.ResponseWriter for the request */,
    r /* The incoming request, for if you need the body or headers in your handler */.
    orbit.RouteParams{
        "user": User{Uid: 5, Name: "Joe Bloggs"}, /* The result of User.FromRequest("5") */
        "event": Event{Uid: 5, Name: "Joe Bloggs"}, /* The result of Event.FromRequest("25") */

    }
)
```

Orbit _won't_ call your handler if it fails to build the RouteParams.
This means that if your handler gets called you can _safely_ just do params["user"] and _expect_ a user to be there.

## Soon but not yet

**Child routes** inherit data from parent routes. For example:

```go
router := orbit.NewRouter()

// Subrouter (can contain many child routes):
pr := router.NewGroup("/user/{user_id}", orbit.RouteParams{"user": User{}})

// Route hung off the subrouter
pr.HandleGet("/event/{event_id}", orbit.RouteParams{"event": Event}, MyHandleFunc)
```

In this example, when MyHandleFunc is called:

- The Group extracts its information first (grabs the user from the db).
  If it fails to do this, then nothing further happens.
- The child handler extracts its information second (grabs the event from the db).
  Again if this fails, nothing further happens.
- The params from **both** levels are **combined**.

MyHandleFunc is then called with the params of both itself and its parent routers:

```go
MyHandleFunc(
    w /* http.ResponseWriter for the request */,
    r /* The incoming request, for if you need the body or headers in your handler */.
    orbit.RouteParams{
        "user": User{Uid: 5, Name: "Joe Bloggs"}, /* The result of User.FromRequest("5") */
        "event": Event{Uid: 5, Name: "Joe Bloggs"}, /* The result of Event.FromRequest("25") */

    }
)
```

Looks pointless in this example, but it's easy to imagine you have more than one child route of /user which
all want to know about the user data.
