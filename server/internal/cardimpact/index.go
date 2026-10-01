// Package cardimpact builds a conservative source index for card changes.
// It never executes card code or infers gameplay from printed descriptions.
package cardimpact

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var cardType = regexp.MustCompile(`^Card([0-9]{7})`)
var cardNumber = regexp.MustCompile(`^[0-9]{7}$`)

type Function struct {
	Name           string   `json:"name"`
	File           string   `json:"file"`
	Line           int      `json:"line"`
	Card           string   `json:"card,omitempty"`
	Receiver       string   `json:"receiver,omitempty"`
	Calls          []string `json:"calls,omitempty"`
	Events         []string `json:"events,omitempty"`
	CardReferences []string `json:"card_references,omitempty"`
	Test           bool     `json:"-"`
}

type Index struct{ Functions []Function }

type RelatedCard struct {
	Number  string   `json:"number"`
	Reasons []string `json:"reasons"`
}

type Report struct {
	Sources      []Function    `json:"sources"`
	Dependencies []Function    `json:"dependencies"`
	RelatedCards []RelatedCard `json:"related_cards"`
	Tests        []Function    `json:"tests"`
	Limits       []string      `json:"limits"`
}

func Build(dir string) (*Index, error) {
	paths, err := filepath.Glob(filepath.Join(dir, "*.go"))
	if err != nil {
		return nil, err
	}
	index := &Index{}
	for _, path := range paths {
		fs := token.NewFileSet()
		file, err := parser.ParseFile(fs, path, nil, 0)
		if err != nil {
			return nil, err
		}
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Body == nil {
				continue
			}
			row := Function{Name: fn.Name.Name, File: filepath.Clean(path), Line: fs.Position(fn.Pos()).Line, Test: strings.HasSuffix(path, "_test.go")}
			if fn.Recv != nil {
				receiver := fn.Recv.List[0].Type
				if ptr, ok := receiver.(*ast.StarExpr); ok {
					receiver = ptr.X
				}
				if name, ok := receiver.(*ast.Ident); ok {
					row.Receiver = name.Name
				}
				if match := cardType.FindStringSubmatch(row.Receiver); match != nil {
					row.Card = match[1]
				}
			}
			calls, refs, events := map[string]bool{}, map[string]bool{}, map[string]bool{}
			ast.Inspect(fn.Type, func(node ast.Node) bool {
				if name, ok := node.(*ast.Ident); ok && strings.HasSuffix(name.Name, "Event") && name.Name != "GameEvent" {
					events[name.Name] = true
				}
				return true
			})
			ast.Inspect(fn.Body, func(node ast.Node) bool {
				if call, ok := node.(*ast.CallExpr); ok {
					switch f := call.Fun.(type) {
					case *ast.Ident:
						calls[f.Name] = true
					case *ast.SelectorExpr:
						calls[f.Sel.Name] = true
					}
				}
				if literal, ok := node.(*ast.BasicLit); ok && literal.Kind == token.STRING {
					value, err := strconv.Unquote(literal.Value)
					if err == nil && cardNumber.MatchString(value) {
						refs[value] = true
					}
				}
				return true
			})
			row.Calls, row.CardReferences = keys(calls), keys(refs)
			row.Events = keys(events)
			index.Functions = append(index.Functions, row)
		}
	}
	return index, nil
}

// Analyze follows concrete helper calls transitively and reverse callers for a
// changed file. Selector names are deliberately over-approximated: ambiguous
// receiver types expand to every matching implementation instead of silently
// excluding a possible dependency. Card methods form review boundaries.
func (index *Index) Analyze(number, changedFile string) (Report, error) {
	result := Report{Limits: []string{
		"Static review aid, not a proof of independence: closures, registry dispatch, attached behaviors and shared state require combination tests.",
		"Matching method names over-approximate receiver types. Card methods are dependency boundaries; related hooks and explicit card references are listed separately.",
		"Tests are source references, not measured coverage. Run go test ./... and go vet ./... even when the focused list is empty.",
	}}
	if (number == "") == (changedFile == "") {
		return result, fmt.Errorf("specify exactly one card number or changed file")
	}
	if number != "" && !cardNumber.MatchString(number) {
		return result, fmt.Errorf("invalid card number %q", number)
	}
	sources, dependencies := map[int]bool{}, map[int]bool{}
	reasons := map[string]map[string]bool{}
	addReason := func(card, reason string) {
		if card == "" || card == number {
			return
		}
		if reasons[card] == nil {
			reasons[card] = map[string]bool{}
		}
		reasons[card][reason] = true
	}
	for id, fn := range index.Functions {
		if fn.Test {
			continue
		}
		match := number != "" && (fn.Card == number || contains(fn.CardReferences, number))
		if changedFile != "" {
			match = filepath.Clean(fn.File) == filepath.Clean(changedFile)
		}
		if match {
			sources[id] = true
		}
	}
	if len(sources) == 0 {
		return result, fmt.Errorf("no source functions found for query; generated definitions and trait-only cards may have no behavior methods")
	}
	for id := range sources {
		dependencies[id] = true
	}
	for changed := true; changed; {
		changed = false
		for id, fn := range index.Functions {
			if !dependencies[id] {
				continue
			}
			for other, candidate := range index.Functions {
				if candidate.Test || candidate.Card != "" || dependencies[other] {
					continue
				}
				if contains(fn.Calls, candidate.Name) {
					dependencies[other] = true
					changed = true
				}
			}
		}
	}
	// For changed shared helpers include every transitive caller, including cards.
	callers := map[int]bool{}
	for id := range sources {
		callers[id] = true
	}
	if changedFile != "" {
		for changed := true; changed; {
			changed = false
			for id, fn := range index.Functions {
				if fn.Test || callers[id] {
					continue
				}
				for target := range callers {
					if contains(fn.Calls, index.Functions[target].Name) {
						callers[id] = true
						changed = true
						break
					}
				}
			}
		}
	}
	hooks := map[string]bool{}
	events := map[string]bool{}
	for id := range sources {
		fn := index.Functions[id]
		for _, event := range fn.Events {
			events[event] = true
		}
		if fn.Card != "" && fn.Name != "ID" && fn.Name != "Name" && !strings.HasPrefix(fn.Name, "HasActive") {
			hooks[fn.Name] = true
		}
	}
	for id, fn := range index.Functions {
		if fn.Test {
			continue
		}
		if sources[id] {
			result.Sources = append(result.Sources, fn)
		}
		if dependencies[id] && !sources[id] {
			result.Dependencies = append(result.Dependencies, fn)
		}
		if fn.Card != "" && hooks[fn.Name] {
			addReason(fn.Card, "shared hook: "+fn.Name)
		}
		for _, event := range fn.Events {
			if events[event] {
				addReason(fn.Card, "shared rule event: "+event)
			}
		}
		if callers[id] && changedFile != "" {
			addReason(fn.Card, "calls changed code: "+fn.Name)
		}
		if number != "" && contains(fn.CardReferences, number) {
			addReason(fn.Card, "explicit reference to "+number)
		}
	}
	for _, card := range sortedReasonKeys(reasons) {
		result.RelatedCards = append(result.RelatedCards, RelatedCard{card, keys(reasons[card])})
	}
	selectedNames := map[string]bool{}
	for id := range sources {
		selectedNames[index.Functions[id].Name] = true
	}
	for id := range dependencies {
		if index.Functions[id].Card == "" {
			selectedNames[index.Functions[id].Name] = true
		}
	}
	for _, fn := range index.Functions {
		if !fn.Test || !strings.HasPrefix(fn.Name, "Test") {
			continue
		}
		match := number != "" && contains(fn.CardReferences, number)
		for _, call := range fn.Calls {
			if selectedNames[call] {
				match = true
			}
		}
		if match {
			result.Tests = append(result.Tests, fn)
		}
	}
	return result, nil
}

func contains(values []string, wanted string) bool {
	for _, value := range values {
		if value == wanted {
			return true
		}
	}
	return false
}
func keys(values map[string]bool) []string {
	result := make([]string, 0, len(values))
	for value := range values {
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}
func sortedReasonKeys(values map[string]map[string]bool) []string {
	result := make([]string, 0, len(values))
	for value := range values {
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}
