package articles

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/januarsyah901/ascard/api/internal/httpx"
)

func fields(t *testing.T, body string) httpx.Fields {
	t.Helper()
	var f httpx.Fields
	if err := json.Unmarshal([]byte(body), &f); err != nil {
		t.Fatalf("bad test body %s: %v", body, err)
	}
	return f
}

func ptr(s string) *string { return &s }

func TestValidate(t *testing.T) {
	q := func(s string) string { b, _ := json.Marshal(s); return string(b) }
	emoji := "😀" // 2 UTF-16 code units

	cases := []struct {
		name string
		body string
		ok   bool
		want Input
	}{
		// title
		{"title trimmed", `{"title":"  Hi  ","content":"x"}`, true, Input{Title: "Hi", Content: "x"}},
		{"title missing", `{"content":"x"}`, false, Input{}},
		{"title null", `{"title":null,"content":"x"}`, false, Input{}},
		{"title number", `{"title":5,"content":"x"}`, false, Input{}},
		{"title array", `{"title":["a"],"content":"x"}`, false, Input{}},
		{"title blank", `{"title":" \t\n ","content":"x"}`, false, Input{}},
		{"title only BOM and nbsp", `{"title":"\ufeff\u00a0\u3000","content":"x"}`, false, Input{}},
		{"title NEL is not JS whitespace", `{"title":"\u0085","content":"x"}`, true, Input{Title: "\u0085", Content: "x"}},
		{"title 200", `{"title":` + q(strings.Repeat("t", 200)) + `,"content":"x"}`, true, Input{Title: strings.Repeat("t", 200), Content: "x"}},
		{"title 201", `{"title":` + q(strings.Repeat("t", 201)) + `,"content":"x"}`, false, Input{}},
		{"title 199 + 2 spaces counts untrimmed", `{"title":` + q(" "+strings.Repeat("t", 199)+" ") + `,"content":"x"}`, false, Input{}},
		{"title 198 + 2 spaces", `{"title":` + q(" "+strings.Repeat("t", 198)+" ") + `,"content":"x"}`, true, Input{Title: strings.Repeat("t", 198), Content: "x"}},
		{"title 100 emoji = 200 units", `{"title":` + q(strings.Repeat(emoji, 100)) + `,"content":"x"}`, true, Input{Title: strings.Repeat(emoji, 100), Content: "x"}},
		{"title 101 emoji = 202 units", `{"title":` + q(strings.Repeat(emoji, 101)) + `,"content":"x"}`, false, Input{}},
		// content
		{"content kept untrimmed", `{"title":"t","content":"  body \n"}`, true, Input{Title: "t", Content: "  body \n"}},
		{"content missing", `{"title":"t"}`, false, Input{}},
		{"content blank", `{"title":"t","content":" \n "}`, false, Input{}},
		{"content number", `{"title":"t","content":1}`, false, Input{}},
		{"content object", `{"title":"t","content":{}}`, false, Input{}},
		{"content 100000", `{"title":"t","content":` + q(strings.Repeat("c", 100000)) + `}`, true, Input{Title: "t", Content: strings.Repeat("c", 100000)}},
		{"content 100001", `{"title":"t","content":` + q(strings.Repeat("c", 100001)) + `}`, false, Input{}},
		{"content 99999 + space", `{"title":"t","content":` + q(strings.Repeat("c", 100000)+" ") + `}`, false, Input{}},
		// category
		{"category trimmed", `{"title":"t","content":"x","category":"  Mat  "}`, true, Input{Title: "t", Content: "x", Category: ptr("Mat")}},
		{"category cut to 50", `{"title":"t","content":"x","category":` + q(strings.Repeat("k", 60)) + `}`, true, Input{Title: "t", Content: "x", Category: ptr(strings.Repeat("k", 50))}},
		{"category cut after trim", `{"title":"t","content":"x","category":` + q("   "+strings.Repeat("k", 50)+"z") + `}`, true, Input{Title: "t", Content: "x", Category: ptr(strings.Repeat("k", 50))}},
		{"category cut never splits pair", `{"title":"t","content":"x","category":` + q(strings.Repeat("k", 49)+emoji) + `}`, true, Input{Title: "t", Content: "x", Category: ptr(strings.Repeat("k", 49))}},
		{"category empty", `{"title":"t","content":"x","category":""}`, true, Input{Title: "t", Content: "x"}},
		{"category blank", `{"title":"t","content":"x","category":"   "}`, true, Input{Title: "t", Content: "x"}},
		{"category null", `{"title":"t","content":"x","category":null}`, true, Input{Title: "t", Content: "x"}},
		{"category number", `{"title":"t","content":"x","category":7}`, true, Input{Title: "t", Content: "x"}},
		{"category bool", `{"title":"t","content":"x","category":true}`, true, Input{Title: "t", Content: "x"}},
		// published
		{"published true", `{"title":"t","content":"x","published":true}`, true, Input{Title: "t", Content: "x", Published: true}},
		{"published false", `{"title":"t","content":"x","published":false}`, true, Input{Title: "t", Content: "x"}},
		{"published string true", `{"title":"t","content":"x","published":"true"}`, true, Input{Title: "t", Content: "x"}},
		{"published 1", `{"title":"t","content":"x","published":1}`, true, Input{Title: "t", Content: "x"}},
		{"published null", `{"title":"t","content":"x","published":null}`, true, Input{Title: "t", Content: "x"}},
		{"published missing", `{"title":"t","content":"x"}`, true, Input{Title: "t", Content: "x"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := validate(fields(t, tc.body))
			if ok != tc.ok {
				t.Fatalf("ok = %v, want %v", ok, tc.ok)
			}
			if !ok {
				return
			}
			if got.Title != tc.want.Title || got.Content != tc.want.Content || got.Published != tc.want.Published {
				t.Fatalf("got %+v, want %+v", got, tc.want)
			}
			if (got.Category == nil) != (tc.want.Category == nil) ||
				(got.Category != nil && *got.Category != *tc.want.Category) {
				t.Fatalf("category = %v, want %v", deref(got.Category), deref(tc.want.Category))
			}
		})
	}
}

func TestValidateEmptyFields(t *testing.T) {
	if _, ok := validate(httpx.Fields{}); ok {
		t.Fatal("empty body must fail")
	}
}

func deref(s *string) any {
	if s == nil {
		return nil
	}
	return *s
}
