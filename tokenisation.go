package orbit

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Tokenise takes a precompiled route regex and an actual request path, and
// returns a map of params + their values extracted from the request path.
func tokenise(regex *regexp.Regexp, expectedParams []string, path string) (map[string]string, error) {

	// Use the regex on the path
	matches := regex.FindAllStringSubmatch(path, -1)

	// If it didn't match, return now
	if len(matches) != 1 {
		return nil, errRouteDoesNotMatch(path)
	}
	params := matches[0][1:]

	// Check the number of matches is equal to the number of params we were expecting.
	// It shouldn't happen, since a parent of this will have already matched the regex to the path.
	if len(params) != len(expectedParams) {
		return nil, errMisconfigured(fmt.Sprintf("wrong param count (want %d got %v)", len(expectedParams), params))
	}

	// Push the extracted params into the result map.
	results := make(map[string]string, len(params))
	for idx := 0; idx < len(params); idx++ {
		results[expectedParams[idx]] = params[idx]
	}

	return results, nil

}

// Takes a route template (e.g. /a/b/{c}/d/{e}) and builds a regex that will
// extract parameters from a real request path (e.g. /a/b/foo/d/bar).
//
// Returns the regex, an ordered slice of the tokens that will be matched, and
// an err if the regexp fails to build or if the request is invalid.
//
// For example, buildMatcherRegex("/a/b/{param1}/d/e/{param2}/f") will return
// (<some regex>, ["param1", "param2"], nil).
func buildMatcherRegex(path string) (*regexp.Regexp, []string, error) {

	// Grab the positions of braces in the path.
	positions, err := getPositionsOfSquirlies(path)
	if err != nil {
		return nil, nil, err
	}

	// Preallocate memory for the names slice. Because braces come in pairs,
	// positions will always have an even number of elements.
	// We'll fill this as we go.
	names := make([]string, 0, len(positions)/2)

	// Wheat index did our last param squirly brace end at?
	lastEnd := 0

	// Start off the regex pattern. We'll add to this as we go.
	var rxp strings.Builder
	rxp.WriteString("^")

	// Loop through every pair of positions and extract the name within them,
	// and then add to the regex. Note the +2 - we're jumping through pairs.
	for idx := 0; idx < len(positions); idx += 2 {

		tStart := positions[idx] // Where does this {param} start?
		tEnd := positions[idx+1] // Where does this {param} end?

		// Add everything since the last closing brace to the regex
		// If this is the first time through the loop, then we're just adding
		// everything since the start.
		rxp.WriteString(regexp.QuoteMeta(path[lastEnd:tStart]))
		lastEnd = tEnd

		// Add this param's name to the names slice
		// +1/-1 here to trim the {}'s (e.g. {foo} -> foo)
		names = append(names, path[tStart+1:tEnd-1])

		// Add the match group to the regex.
		rxp.WriteString("([a-zA-Z0-9_-]+)")

	}

	// Add the remaining path to the end of the regex
	rxp.WriteString(regexp.QuoteMeta(path[lastEnd:]))
	rxp.WriteString("/?$")

	fmt.Printf("\nCompiled regex for %s: %s\n\n", path, rxp.String())

	// Try compiling the regex
	compiledRxp, err := regexp.Compile(rxp.String())
	if err != nil {
		return nil, nil, fmt.Errorf("regex compilation failed: %s", err.Error())
	}

	return compiledRxp, names, nil

}

// Scans through the path checking that all the {squirlies} are balanced and
// go no deeper than exactly one level, and then returns a slice of the
// positions of those squirlies within the path.
func getPositionsOfSquirlies(path string) ([]int, error) {

	var positions []int = []int{} // positions of squirlies
	var inside bool = false       // are we currently inside a set of squirlies?

	for idx := 0; idx < len(path); idx++ {

		// Found an opening squirly
		if path[idx] == '{' {

			// If we're already inside a set of squirlies then we can't open new ones. Return an error.
			if inside {
				return nil, errors.New("found nested / unbalanced curly braces in route path")
			}

			// Otherwise add the start brace's index to the slice.
			positions = append(positions, idx)
			inside = true
			continue

		}

		// Found a closing squirly
		if path[idx] == '}' {

			// If we're NOT inside a set of squirlies then there's nothing to close so return an error
			if !inside {
				return nil, errors.New("found nested / unbalanced curly braces in route path")
			}

			// Otherwise add the end brace's index to the slice (+1 so we include the brace itself!)
			inside = false
			positions = append(positions, idx+1)
			continue

		}

	}

	// If were still inside, we're missing a closing squirly.
	if inside {
		return nil, errors.New("found nested / unbalanced curly braces in route path")
	}

	return positions, nil

}
