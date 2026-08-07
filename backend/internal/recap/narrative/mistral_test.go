package narrative

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMistralHTTPProviderBuildsStrictSingleRequest(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if r.Method != http.MethodPost {
			t.Errorf("method=%s", r.Method)
		}
		if r.URL.Path != "/v1/chat/completions" {
			t.Errorf("path=%s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer secret" {
			t.Errorf("authorization=%q", r.Header.Get("Authorization"))
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body["temperature"] != float64(0) {
			t.Errorf("temperature=%v", body["temperature"])
		}
		format, ok := body["response_format"].(map[string]any)
		if !ok || format["type"] != "json_schema" {
			t.Fatalf("response_format=%#v", body["response_format"])
		}
		schema, ok := format["json_schema"].(map[string]any)
		if !ok {
			t.Fatalf("json_schema=%#v", format["json_schema"])
		}
		if schema["strict"] != true || schema["schema_definition"] == nil {
			t.Fatalf("schema=%#v", schema)
		}
		if _, exists := schema["schema"]; exists {
			t.Fatalf("legacy schema key present: %#v", schema)
		}
		w.Header().Set("Content-Type", "application/json")
		response := map[string]any{"model": "actual-model", "choices": []any{map[string]any{"message": map[string]any{"content": `{"summary_title":"Итоги","summary_text":"Готово"}`}}}}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	provider := MistralHTTPProvider{APIKey: "secret", Endpoint: server.URL + "/v1/chat/completions", Client: server.Client()}
	got, err := provider.Generate(context.Background(), baseInput())
	if err != nil {
		t.Fatal(err)
	}
	if calls != 1 {
		t.Fatalf("calls=%d", calls)
	}
	if got.Model != "actual-model" {
		t.Fatalf("model=%q", got.Model)
	}
	if string(got.Content) != `{"summary_title":"Итоги","summary_text":"Готово"}` {
		t.Fatalf("content=%s", got.Content)
	}
}

func TestMistralHTTPProviderRejectsNon2xx(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { http.Error(w, "busy", http.StatusTooManyRequests) }))
	defer server.Close()
	provider := MistralHTTPProvider{APIKey: "secret", Endpoint: server.URL, Client: server.Client()}
	if _, err := provider.Generate(context.Background(), baseInput()); err == nil {
		t.Fatal("expected error")
	}
}
