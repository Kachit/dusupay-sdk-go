// Command drift compares the DusuPay API documentation (OpenAPI spec, markdown pages,
// Postman collection) against the SDK's Go types and prints a markdown drift report.
//
// Usage (from the repo root):
//
//	go run ./.claude/skills/docs-sync/scripts/drift -docs <dir from fetch_docs.sh>
//
// It only reads files; it never modifies the SDK.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// ---------- mapping ----------

type endpointMap struct {
	Query string `json:"query"`
	Body  string `json:"body"`
	Data  string `json:"data"`
}

type mapping struct {
	Endpoints    map[string]endpointMap `json:"endpoints"`
	Callbacks    map[string]string      `json:"callbacks"`
	CardData     string                 `json:"card_data"`
	PageRequests map[string]string      `json:"page_requests"`
}

// ---------- Go source model ----------

type goField struct {
	kind string // string, number, bool, object, array
	elem string // struct type name for object/array-of-struct fields
}

type goModel struct {
	structs   map[string]map[string]goField
	consts    map[string]bool // string constant values
	headers   map[string]bool // headers set in RequestBuilder
	validated map[string]bool // names in `parameter "x" is empty` messages
}

func loadGo(dir string) *goModel {
	m := &goModel{
		structs:   map[string]map[string]goField{},
		consts:    map[string]bool{},
		headers:   map[string]bool{},
		validated: map[string]bool{},
	}
	underlying := map[string]string{}
	var structTypes = map[string]*ast.StructType{}

	files, _ := filepath.Glob(filepath.Join(dir, "*.go"))
	fset := token.NewFileSet()
	headerRe := regexp.MustCompile(`headers\.Set\("([^"]+)"`)
	validRe := regexp.MustCompile(`parameter "([a-z_.]+)" is empty`)

	for _, path := range files {
		if strings.HasSuffix(path, "_test.go") {
			continue
		}
		src, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		for _, mm := range headerRe.FindAllStringSubmatch(string(src), -1) {
			m.headers[strings.ToLower(mm[1])] = true
		}
		for _, mm := range validRe.FindAllStringSubmatch(string(src), -1) {
			m.validated[mm[1]] = true
		}
		f, err := parser.ParseFile(fset, path, src, 0)
		if err != nil {
			fmt.Fprintf(os.Stderr, "parse %s: %v\n", path, err)
			continue
		}
		ast.Inspect(f, func(n ast.Node) bool {
			switch d := n.(type) {
			case *ast.TypeSpec:
				switch t := d.Type.(type) {
				case *ast.StructType:
					structTypes[d.Name.Name] = t
				case *ast.Ident:
					underlying[d.Name.Name] = t.Name
				}
			case *ast.ValueSpec:
				for _, v := range d.Values {
					if bl, ok := v.(*ast.BasicLit); ok && bl.Kind == token.STRING {
						if s, err := strconv.Unquote(bl.Value); err == nil {
							m.consts[s] = true
						}
					}
				}
			}
			return true
		})
	}

	for name, st := range structTypes {
		fields := map[string]goField{}
		for _, fld := range st.Fields.List {
			if fld.Tag == nil {
				// embedded struct (e.g. ResponseBody) has no json tag of its own
				continue
			}
			tag, _ := strconv.Unquote(fld.Tag.Value)
			jsonName := strings.Split(reflectTag(tag, "json"), ",")[0]
			if jsonName == "" || jsonName == "-" {
				continue
			}
			fields[jsonName] = goKind(fld.Type, underlying, structTypes)
		}
		m.structs[name] = fields
	}

	return m
}

func reflectTag(tag, key string) string {
	re := regexp.MustCompile(key + `:"([^"]*)"`)
	if mm := re.FindStringSubmatch(tag); mm != nil {
		return mm[1]
	}
	return ""
}

func goKind(e ast.Expr, underlying map[string]string, structs map[string]*ast.StructType) goField {
	switch t := e.(type) {
	case *ast.StarExpr:
		return goKind(t.X, underlying, structs)
	case *ast.ArrayType:
		inner := goKind(t.Elt, underlying, structs)
		return goField{kind: "array", elem: inner.elem}
	case *ast.MapType, *ast.InterfaceType:
		return goField{kind: "object"}
	case *ast.Ident:
		name := t.Name
		if _, ok := structs[name]; ok {
			return goField{kind: "object", elem: name}
		}
		if u, ok := underlying[name]; ok {
			name = u
		}
		switch {
		case name == "string":
			return goField{kind: "string"}
		case name == "bool":
			return goField{kind: "bool"}
		case strings.HasPrefix(name, "int"), strings.HasPrefix(name, "uint"), strings.HasPrefix(name, "float"):
			return goField{kind: "number"}
		}
	}
	return goField{kind: "unknown"}
}

// ---------- findings ----------

type key struct{ goType, field string }

type report struct {
	missing   map[key]map[string]bool            // doc field absent from Go type -> sources
	seen      map[key]bool                       // Go fields confirmed by some source
	kinds     map[key]map[string]map[string]bool // field -> doc kind -> sources
	endpoints map[string]map[string]bool         // endpoint -> sources
	headers   map[string]map[string]bool
	enums     map[string]map[string]bool // enum value -> sources
	required  map[key]map[string]bool
	notes     []string
}

func newReport() *report {
	return &report{
		missing: map[key]map[string]bool{}, seen: map[key]bool{},
		kinds: map[key]map[string]map[string]bool{}, endpoints: map[string]map[string]bool{},
		headers: map[string]map[string]bool{}, enums: map[string]map[string]bool{},
		required: map[key]map[string]bool{},
	}
}

func add(m map[key]map[string]bool, k key, src string) {
	if m[k] == nil {
		m[k] = map[string]bool{}
	}
	m[k][src] = true
}

func addS(m map[string]map[string]bool, k, src string) {
	if m[k] == nil {
		m[k] = map[string]bool{}
	}
	m[k][src] = true
}

func jsonKind(v any) string {
	switch v.(type) {
	case string:
		return "string"
	case float64, json.Number:
		return "number"
	case bool:
		return "bool"
	case map[string]any:
		return "object"
	case []any:
		return "array"
	case nil:
		return "null"
	}
	return "unknown"
}

// compare walks a documented object (decoded example or schema converted to the same
// shape) against a Go struct, recursing into nested structs.
func (r *report) compare(g *goModel, goType string, doc map[string]any, src string) {
	fields, ok := g.structs[goType]
	if !ok {
		r.notes = append(r.notes, fmt.Sprintf("mapping refers to unknown Go type %s (%s)", goType, src))
		return
	}
	for name, val := range doc {
		k := key{goType, name}
		f, ok := fields[name]
		if !ok {
			add(r.missing, k, src)
			continue
		}
		r.seen[k] = true
		kind := jsonKind(val)
		if s, isSchema := val.(schemaLeaf); isSchema {
			kind = string(s)
		}
		if kind != "null" && kind != "unknown" {
			if r.kinds[k] == nil {
				r.kinds[k] = map[string]map[string]bool{}
			}
			if r.kinds[k][kind] == nil {
				r.kinds[k][kind] = map[string]bool{}
			}
			r.kinds[k][kind][src] = true
		}
		if f.elem == "" {
			continue
		}
		switch v := val.(type) {
		case map[string]any:
			r.compare(g, f.elem, v, src)
		case []any:
			for _, item := range v {
				if obj, ok := item.(map[string]any); ok {
					r.compare(g, f.elem, obj, src)
				}
			}
		}
	}
}

// ---------- OpenAPI ----------

// schemaLeaf marks a leaf converted from a schema; its value is the JSON kind.
type schemaLeaf string

type spec struct {
	Paths      map[string]map[string]operation `json:"paths"`
	Components struct {
		Schemas map[string]map[string]any `json:"schemas"`
	} `json:"components"`
}

type operation struct {
	Parameters []struct {
		Name     string         `json:"name"`
		In       string         `json:"in"`
		Required bool           `json:"required"`
		Schema   map[string]any `json:"schema"`
	} `json:"parameters"`
	RequestBody struct {
		Content map[string]struct {
			Schema map[string]any `json:"schema"`
		} `json:"content"`
	} `json:"requestBody"`
	Responses map[string]struct {
		Content map[string]struct {
			Schema map[string]any `json:"schema"`
		} `json:"content"`
	} `json:"responses"`
}

var pathParamRe = regexp.MustCompile(`\{[^}]*\}`)

func endpointKey(method, path string) string {
	path = strings.SplitN(path, "?", 2)[0]
	path = "/" + strings.Trim(strings.TrimSpace(path), "/")
	return strings.ToUpper(method) + " " + pathParamRe.ReplaceAllString(path, "{}")
}

func (s *spec) resolve(sch map[string]any, depth int) map[string]any {
	if depth > 10 || sch == nil {
		return sch
	}
	if ref, ok := sch["$ref"].(string); ok {
		return s.resolve(s.Components.Schemas[strings.TrimPrefix(ref, "#/components/schemas/")], depth+1)
	}
	if all, ok := sch["allOf"].([]any); ok {
		merged := map[string]any{"type": "object", "properties": map[string]any{}}
		props := merged["properties"].(map[string]any)
		var required []any
		for _, part := range all {
			if pm, ok := part.(map[string]any); ok {
				r := s.resolve(pm, depth+1)
				if p, ok := r["properties"].(map[string]any); ok {
					for k, v := range p {
						props[k] = v
					}
				}
				if req, ok := r["required"].([]any); ok {
					required = append(required, req...)
				}
			}
		}
		if p, ok := sch["properties"].(map[string]any); ok {
			for k, v := range p {
				if _, exists := props[k]; !exists {
					props[k] = v
				}
			}
		}
		merged["required"] = required
		return merged
	}
	return sch
}

// toShape converts a schema to the same map/[]any/leaf shape as a decoded JSON example.
func (s *spec) toShape(sch map[string]any, depth int, enums func(string)) any {
	sch = s.resolve(sch, depth)
	if sch == nil {
		return schemaLeaf("unknown")
	}
	if e, ok := sch["enum"].([]any); ok {
		for _, v := range e {
			if str, ok := v.(string); ok {
				enums(str)
			}
		}
	}
	if props, ok := sch["properties"].(map[string]any); ok {
		out := map[string]any{}
		for k, v := range props {
			if vm, ok := v.(map[string]any); ok {
				out[k] = s.toShape(vm, depth+1, enums)
			}
		}
		return out
	}
	switch sch["type"] {
	case "array":
		if items, ok := sch["items"].(map[string]any); ok {
			return []any{s.toShape(items, depth+1, enums)}
		}
		return []any{}
	case "string":
		return schemaLeaf("string")
	case "number", "integer":
		return schemaLeaf("number")
	case "boolean":
		return schemaLeaf("bool")
	case "object":
		if ex, ok := sch["example"].(map[string]any); ok {
			return ex
		}
		return map[string]any{}
	}
	return schemaLeaf("unknown")
}

func (r *report) checkSpec(g *goModel, mp *mapping, path string) {
	raw, err := os.ReadFile(path)
	if err != nil {
		r.notes = append(r.notes, "openapi spec not found: "+path)
		return
	}
	var s spec
	if err := json.Unmarshal(raw, &s); err != nil {
		r.notes = append(r.notes, "openapi spec unreadable: "+err.Error())
		return
	}
	const src = "openapi"
	for p, ops := range s.Paths {
		for method, op := range ops {
			ek := endpointKey(method, p)
			addS(r.endpoints, ek, src)
			em := mp.Endpoints[ek]
			query := map[string]any{}
			for _, prm := range op.Parameters {
				switch prm.In {
				case "header":
					addS(r.headers, strings.ToLower(prm.Name), src+" "+ek)
				case "query":
					query[prm.Name] = s.toShape(prm.Schema, 0, func(string) {})
					if prm.Required && em.Query != "" && !g.validated[prm.Name] {
						add(r.required, key{em.Query, prm.Name}, src)
					}
				}
			}
			if len(query) > 0 && em.Query != "" {
				r.compare(g, em.Query, query, src+" "+ek)
			}
			for _, c := range op.RequestBody.Content {
				enumSink := func(v string) { addS(r.enums, v, src+" "+ek) }
				resolved := s.resolve(c.Schema, 0)
				if shape, ok := s.toShape(c.Schema, 0, enumSink).(map[string]any); ok && em.Body != "" {
					r.compare(g, em.Body, shape, src+" "+ek)
				}
				if req, ok := resolved["required"].([]any); ok && em.Body != "" {
					for _, name := range req {
						if n, ok := name.(string); ok && !g.validated[n] {
							add(r.required, key{em.Body, n}, src)
						}
					}
				}
			}
			for code, resp := range op.Responses {
				if !strings.HasPrefix(code, "2") || em.Data == "" {
					continue
				}
				for _, c := range resp.Content {
					shape, _ := s.toShape(c.Schema, 0, func(string) {}).(map[string]any)
					if data, ok := shape["data"].(map[string]any); ok {
						r.compare(g, em.Data, data, src+" "+ek+" "+code)
					}
				}
			}
		}
	}
}

// ---------- markdown pages ----------

var (
	codeBlockRe  = regexp.MustCompile("(?s)```(?:json|javascript)\\s*\n(.*?)```")
	commentRe    = regexp.MustCompile(`(?s)/\*.*?\*/`)
	tableRowRe   = regexp.MustCompile(`<tr><td>([a-z_][a-z0-9_.]*)(?:<mark[^>]*>\*</mark>)?</td>`)
	tableBlockRe = regexp.MustCompile(`(?s)<table>.*?</table>`)
)

func decodeAll(s string) []any {
	dec := json.NewDecoder(strings.NewReader(commentRe.ReplaceAllString(s, "")))
	var out []any
	for {
		var v any
		if err := dec.Decode(&v); err != nil {
			return out
		}
		out = append(out, v)
	}
}

func pageRequestType(mp *mapping, page string) string {
	for prefix, t := range mp.PageRequests {
		if strings.HasPrefix(page, prefix) {
			return t
		}
	}
	return ""
}

func (r *report) checkPages(g *goModel, mp *mapping, dir string) {
	files, _ := filepath.Glob(filepath.Join(dir, "*.md"))
	if len(files) == 0 {
		r.notes = append(r.notes, "no markdown pages found in "+dir)
	}
	for _, path := range files {
		raw, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		page := strings.ReplaceAll(strings.TrimSuffix(filepath.Base(path), ".md"), "__", "/")
		src := "page " + page
		text := string(raw)
		reqType := pageRequestType(mp, page)

		for _, block := range codeBlockRe.FindAllStringSubmatch(text, -1) {
			for _, v := range decodeAll(block[1]) {
				obj, ok := v.(map[string]any)
				if !ok {
					continue
				}
				switch {
				case obj["event"] != nil && obj["payload"] != nil:
					payload, _ := obj["payload"].(map[string]any)
					tt, _ := payload["transaction_type"].(string)
					if t := mp.Callbacks[strings.ToUpper(tt)]; t != "" {
						r.compare(g, t, payload, src+" (callback "+fmt.Sprint(obj["event"])+")")
					}
				case obj["card_no"] != nil && mp.CardData != "":
					r.compare(g, mp.CardData, obj, src+" (card data)")
				case obj["transaction_method"] != nil && obj["merchant_reference"] != nil && obj["code"] == nil && reqType != "":
					r.compare(g, reqType, obj, src+" (request example)")
				}
			}
		}

		// Parameter tables on request pages (header row starts with "Parameter"). Header tables
		// ("Header Name") and response tables on other pages are skipped via reqType/target.
		for _, table := range tableBlockRe.FindAllString(text, -1) {
			if !strings.Contains(table, ">Parameter</th>") {
				continue
			}
			target := reqType
			if strings.Contains(table, "<td>card_no</td>") {
				target = mp.CardData
			}
			if target == "" {
				continue
			}
			doc := map[string]any{}
			for _, mm := range tableRowRe.FindAllStringSubmatch(table, -1) {
				setPath(doc, mm[1])
			}
			r.compare(g, target, doc, src+" (parameter table)")
		}
	}
}

// setPath turns "extra_params.bank_code" into nested maps so compare recurses.
func setPath(doc map[string]any, dotted string) {
	parts := strings.Split(dotted, ".")
	cur := doc
	for i, p := range parts {
		if i == len(parts)-1 {
			if _, exists := cur[p]; !exists {
				cur[p] = schemaLeaf("unknown")
			}
			return
		}
		next, ok := cur[p].(map[string]any)
		if !ok {
			next = map[string]any{}
			cur[p] = next
		}
		cur = next
	}
}

// ---------- Postman ----------

type pmItem struct {
	Name    string   `json:"name"`
	Item    []pmItem `json:"item"`
	Request *struct {
		Method string `json:"method"`
		URL    any    `json:"url"`
		Header []struct {
			Key string `json:"key"`
		} `json:"header"`
		Body *struct {
			Raw string `json:"raw"`
		} `json:"body"`
	} `json:"request"`
}

func (r *report) checkPostman(g *goModel, mp *mapping, path string) {
	raw, err := os.ReadFile(path)
	if err != nil {
		r.notes = append(r.notes, "postman collection not found: "+path)
		return
	}
	var root pmItem
	if err := json.Unmarshal(raw, &root); err != nil {
		r.notes = append(r.notes, "postman collection unreadable: "+err.Error())
		return
	}
	var walk func(it pmItem)
	walk = func(it pmItem) {
		if rq := it.Request; rq != nil {
			rawURL := ""
			switch u := rq.URL.(type) {
			case string:
				rawURL = u
			case map[string]any:
				rawURL, _ = u["raw"].(string)
			}
			rawURL = strings.ReplaceAll(rawURL, "{{baseUrl}}", "")
			ek := endpointKey(rq.Method, rawURL)
			src := "postman \"" + it.Name + "\""
			addS(r.endpoints, ek, "postman")
			em := mp.Endpoints[ek]
			for _, h := range rq.Header {
				addS(r.headers, strings.ToLower(h.Key), src)
			}
			if parts := strings.SplitN(rawURL, "?", 2); len(parts) == 2 && em.Query != "" {
				query := map[string]any{}
				for _, kv := range strings.Split(parts[1], "&") {
					name := strings.SplitN(kv, "=", 2)[0]
					query[name] = "x"
				}
				r.compare(g, em.Query, query, src)
			}
			if rq.Body != nil && strings.TrimSpace(rq.Body.Raw) != "" && em.Body != "" {
				for _, v := range decodeAll(rq.Body.Raw) {
					if obj, ok := v.(map[string]any); ok {
						r.compare(g, em.Body, obj, src)
					}
				}
			}
		}
		for _, c := range it.Item {
			walk(c)
		}
	}
	walk(root)
}

// ---------- output ----------

func srcList(m map[string]bool) string {
	var out []string
	for s := range m {
		out = append(out, s)
	}
	sort.Strings(out)
	if len(out) > 4 {
		return strings.Join(out[:4], "; ") + fmt.Sprintf("; +%d more", len(out)-4)
	}
	return strings.Join(out, "; ")
}

func sortedKeys(m map[key]map[string]bool) []key {
	var ks []key
	for k := range m {
		ks = append(ks, k)
	}
	sort.Slice(ks, func(i, j int) bool {
		if ks[i].goType != ks[j].goType {
			return ks[i].goType < ks[j].goType
		}
		return ks[i].field < ks[j].field
	})
	return ks
}

func (r *report) print(g *goModel, mp *mapping) {
	fmt.Println("# DusuPay docs drift report")

	fmt.Println("\n## Endpoints")
	var eps []string
	for ep := range r.endpoints {
		eps = append(eps, ep)
	}
	sort.Strings(eps)
	newEP := 0
	for _, ep := range eps {
		if _, ok := mp.Endpoints[ep]; !ok {
			fmt.Printf("- NEW (not in SDK mapping): `%s` (%s)\n", ep, srcList(r.endpoints[ep]))
			newEP++
		}
	}
	for ep := range mp.Endpoints {
		if r.endpoints[ep] == nil {
			fmt.Printf("- NOT IN DOCS (removed?): `%s`\n", ep)
			newEP++
		}
	}
	if newEP == 0 {
		fmt.Println("- all documented endpoints are mapped to SDK types")
	}

	fmt.Println("\n## Fields documented but missing from SDK types")
	if len(r.missing) == 0 {
		fmt.Println("- none")
	}
	for _, k := range sortedKeys(r.missing) {
		fmt.Printf("- `%s.%s` (%s)\n", k.goType, k.field, srcList(r.missing[k]))
	}

	fmt.Println("\n## Type conflicts (doc kind vs SDK kind, or sources disagreeing)")
	conflicts := 0
	var kks []key
	for k := range r.kinds {
		kks = append(kks, k)
	}
	sort.Slice(kks, func(i, j int) bool { return kks[i].goType+kks[i].field < kks[j].goType+kks[j].field })
	for _, k := range kks {
		goK := g.structs[k.goType][k.field].kind
		kinds := r.kinds[k]
		_, matches := kinds[goK]
		if len(kinds) == 1 && matches {
			continue
		}
		conflicts++
		var parts []string
		for kind, srcs := range kinds {
			parts = append(parts, fmt.Sprintf("%s in [%s]", kind, srcList(srcs)))
		}
		sort.Strings(parts)
		fmt.Printf("- `%s.%s`: SDK=%s; docs: %s\n", k.goType, k.field, goK, strings.Join(parts, " | "))
	}
	if conflicts == 0 {
		fmt.Println("- none")
	}

	fmt.Println("\n## SDK fields not confirmed by any source (informational; examples are often partial)")
	types := map[string]bool{}
	for _, em := range mp.Endpoints {
		for _, t := range []string{em.Query, em.Body, em.Data} {
			if t != "" {
				types[t] = true
			}
		}
	}
	for _, t := range mp.Callbacks {
		types[t] = true
	}
	types[mp.CardData] = true
	var tnames []string
	for t := range types {
		tnames = append(tnames, t)
	}
	sort.Strings(tnames)
	unconfirmed := 0
	for _, t := range tnames {
		var fs []string
		for f := range g.structs[t] {
			if !r.seen[key{t, f}] {
				fs = append(fs, f)
			}
		}
		if len(fs) > 0 {
			sort.Strings(fs)
			fmt.Printf("- `%s`: %s\n", t, strings.Join(fs, ", "))
			unconfirmed++
		}
	}
	if unconfirmed == 0 {
		fmt.Println("- none")
	}

	fmt.Println("\n## Headers in docs not set by the SDK RequestBuilder")
	var hs []string
	for h := range r.headers {
		if !g.headers[h] && h != "accept" {
			hs = append(hs, h)
		}
	}
	sort.Strings(hs)
	if len(hs) == 0 {
		fmt.Println("- none")
	}
	for _, h := range hs {
		fmt.Printf("- `%s` (%s)\n", h, srcList(r.headers[h]))
	}

	fmt.Println("\n## Enum values without a matching SDK constant")
	var evs []string
	for v := range r.enums {
		if !g.consts[v] {
			evs = append(evs, v)
		}
	}
	sort.Strings(evs)
	if len(evs) == 0 {
		fmt.Println("- none")
	}
	for _, v := range evs {
		fmt.Printf("- `%s` (%s)\n", v, srcList(r.enums[v]))
	}

	fmt.Println("\n## Required by the spec but no `parameter \"x\" is empty` check in the SDK")
	if len(r.required) == 0 {
		fmt.Println("- none")
	}
	for _, k := range sortedKeys(r.required) {
		fmt.Printf("- `%s.%s` (%s)\n", k.goType, k.field, srcList(r.required[k]))
	}

	if len(r.notes) > 0 {
		fmt.Println("\n## Tool notes")
		for _, n := range r.notes {
			fmt.Println("- " + n)
		}
	}
}

func main() {
	docs := flag.String("docs", "", "directory produced by fetch_docs.sh (pages/, openapi.json, postman.json)")
	src := flag.String("src", ".", "SDK source directory (repo root)")
	mappingPath := flag.String("mapping", ".claude/skills/docs-sync/scripts/drift/mapping.json", "type mapping file")
	flag.Parse()
	if *docs == "" {
		fmt.Fprintln(os.Stderr, "usage: drift -docs <dir> [-src .] [-mapping mapping.json]")
		os.Exit(2)
	}

	rawMap, err := os.ReadFile(*mappingPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "read mapping:", err)
		os.Exit(2)
	}
	var mp mapping
	if err := json.Unmarshal(rawMap, &mp); err != nil {
		fmt.Fprintln(os.Stderr, "parse mapping:", err)
		os.Exit(2)
	}

	g := loadGo(*src)
	r := newReport()
	r.checkSpec(g, &mp, filepath.Join(*docs, "openapi.json"))
	r.checkPages(g, &mp, filepath.Join(*docs, "pages"))
	r.checkPostman(g, &mp, filepath.Join(*docs, "postman.json"))
	r.print(g, &mp)
}
