# 🪐 Orbit

Orbit is a HTTP router for Go. It differs from other frameworks in that `/{params}/in/{the}/url` are converted to the type (any type) you actually want, as part of the framework handling the request. That is to say, instead of your handler recieving strings, it recieves
actual resolved full structs that can include anything you want, like pulling from the database or something.

It also guarantees that if your handler is called, then the data is present, and of the right type.
Orbit's request handling is a bit like the Rocket framework for Rust, hence the whacky space-themed name.

It's unfinished and not close to being usable yet, but the core 'getting params from the request as the right type' bit is working.

## Main points:

### Handlers work like go's `http.Handler`/`http.HandlerFunc`, but the signature is different.

The interface to meet looks like this:

```go
// struct based handler
type Handler interface {
    ServeHTTP(
        http.ResponseWriter, // like net/http
        *http.Request,       // like net/http
        RouteParams,         // map["paramName"]ParamValue
        FromBodyable,        // decoded body, or nil if you don't care
    )
}

// func based handler adapter like http.HandlerFunc
type HandlerFunc func(http.ResponseWriter, *http.Request, RouteParams, FromBodyable)
```

### Routes can contain paramaters

Routes in Orbit look like this:`/foo/bar/{fizz}/baz/{buzz}`.

In the route above:

- `fizz` and `buzz` are parameters and can match any string value
- `/foo/bar/hello/baz/128` matches, with `fizz`=`"hello"`, `buzz`=`"128"`
- `/foo/bar/baz` doesn't match

The Orbit router also matches _methods_ (you can specify a handler only handles GET requests for example).
You specify this when you attach the handler to the router.

Handlers are attached to the router like this:

```go
r := orbit.NewRouter()

r.Handle(
    path,       // The path template to match (/a/b/{c}/d)
    handler,    // your orbit.Handler
    methods,    // []string of HTTP methods to match. Nil means all methods.
    paramTypes, // map[string]FromRequestable - filled and passed to handler on request
    bodyType,   // FromBodyable - type to decode request body to (or nil to skip decoding)
)
```

### Parameter types must implement FromRequestable

To make your types work with Orbit, they need to implement the FromRequestable
interface - this lets Orbit build your type from the param string from an incoming request.

To meet FromRequestable, just add a .FromRequest method that takes a string (which will be
whatever was in the URL for the parameter), and returns the resolved value or err.

Orbit calls FromRequest when trying to derive a value from an incoming request.

```go
// This example lets Orbit build Users from the UID in a request.
// For example, User.FromRequest("125") will return the User with ID 125, or an error.
func (u User)FromRequest(uid string) (any, error) {
    // Our database stores uids as ints, but the request param is always a string.
    uidInt, _ := strconv.Atoi
    user, err = yourAppLogic.getUserByUID(uid)
    return user, err
}
```

### Body types must implement FromBodyable

FromBodyable is just like FromRequestable, except it's used when trying to decode the _body_
of an incoming request to a type.

Instead of being passed a string, instead you get passed an io.ReadCloser of the
request's body that you need to decode.

```go
// This example lets Orbit turn a request body into an Event struct by JSON
// unmarshalling it. This is a contrived example - you should do validation too.
func (e Event)FromBody(body io.ReadCloser) (any, error) {
    var result Event
    _ = json.NewDecoder(body).Decode(&result)
    return result, nil
}
```

## Working Example:

Say you've got an API route to create an event for a given user by posting the
new event as JSON to `/user/{user}/event/{event}`.

An example of a valid request for that endpoint might be `GET /user/5/event/my-event`
with a JSON body representing a new event. The goal is to create an event owner by
user `5` called `my-event`, based on the JSON.

In this example we'll set up a route that matches that endpoint, and calls the handler
with a resolved User, Event name, and Body based on the incoming request.

```go
// ---
/// These our our app's types. We only care about users + events in this example.
/// These types aren't solely for the API, they might mirror what's in our DB.

// Represents a user in our system
type User struct {
    Uid int
    Name string
}

// Represents an event in our system
type Event struct {
    EventID int
    AttendeeNames []string
}

/// ---
/// To make these types work with Orbit, they need to implement the FromRequestable
/// or FromBodyable interfaces. You don't need both, but you can if you want.
//
// Orbit calls FromRequest when trying to derive a value from an incoming request,
// or FromBody when trying to decode from the request's body.

// FromRequest resolves a user based on the UID provided in an API call.
// For example, User.FromRequest("125") will return the User with ID 125, or an error.
func (u User)FromRequest(uid string) (any, error) {
    // Our database stores uids as ints, but the request only has strings.
    uidInt, _ := strconv.Atoi
    user, err = yourAppLogic.getUserByUID(uid)
    return user, err
}

// FromBody resolves an event from the body of a request, by JSON-decoding it.
func (e Event)FromBody(body io.ReadCloser) (any, error) {
    var result Event
    _ = json.NewDecoder(body).Decode(&result)
    return result, nil
}

/// ---
/// Make your handler

// Handles POST /user/{user}/event/{event}
// with a json-encoded User in the body.
func MyNewEventHandler(
    w http.ResponseWriter,
    r *http.Request,
    params RouteParams,
    body FromBodyable,
){
    // The magic of Orbit is that if your handler gets called, the params and body
    // are guaranteed to be present and of the correct type you specify when wiring
    // up the handler to the route.
    //
    // This means you can blindly assert types like this:
    user, _ := params["user"].(User)
    newEventName, _ := params["event"].(orbit.BasicString)
    newEvent, _ := body.(Event)

    // Your app logic
    yourAppStuff.MakeNewEventInDB(user, newEventName, newEvent)

    w.WriteHeader(200)
}

/// ---
/// Make a router and wire the handler up to it.
/// This is where we tell Orbit what types our arguments are.

r := orbit.NewRouter()

r.Handle(
    // The path to match against
    "/user/{user}/event/{event}",
    // your orbit.Handler
    MyNewEventHandler,
    // []string of HTTP methods to match. Nil means all methods.
    []string{"POST"},
    // Param types to decode (types must implement FromRequestable)
    orbit.RouteParams{
        "user": User{},
        "event": orbit.BasicString, // helper: just a string that implements FromRequest.
    },
    // Body type (must implement FromBodyable)
    User{},
)

// Precompiles regexes etc for the router.
// Call this once before starting up.
r.Bake()

log.Fatal(http.ListenAndServe(":8080", r))


```

## Progress

- [x] Matching routes
- [x] Decoding from the url params
- [x] Decoding from the body
- [x] Checking types
- [x] Top level router as a http.Handler
- [ ] Somehow removing reliance on reflection
- [ ] Child routes inheriting params from parent routes.
