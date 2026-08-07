package narrative

import (
	"context"
	"errors"
	"testing"
)

type fakeProvider struct {
	result ProviderResult
	err    error
	calls  int
}

func (p *fakeProvider) Generate(context.Context, Input) (ProviderResult, error) {
	p.calls++
	return p.result, p.err
}

func baseInput() Input {
	return Input{
		Locale: "ru-RU", Year: 2026, Theme: "city",
		Role:         NamedCode{Code: "findings_seeker", Title: "Искатель находок"},
		Style:        NamedCode{Code: "thoughtful", Title: "Вдумчивый"},
		MainDistrict: NamedCode{Code: "goods", Title: "Товары"},
		SafeFacts:    []SafeFact{{MetricCode: "active_days", Value: 7, Unit: "days"}},
	}
}

func TestServiceAcceptsStrictValidStructuredOutputAndYear(t *testing.T) {
	provider := &fakeProvider{result: ProviderResult{Content: []byte(`{"summary_title":"Городской ритм","summary_text":"Ты был активен 7 дней.","extra":"x"}`), Model: "m"}}
	got := (Service{Provider: provider}).Build(context.Background(), baseInput(), "")
	if got.Source != "template" {
		t.Fatalf("source=%q, extra field must cause fallback", got.Source)
	}

	provider.result.Content = []byte(`{"summary_title":"Итоги 2026","summary_text":"Ты был активен 7 дней."}`)
	got = (Service{Provider: provider}).Build(context.Background(), baseInput(), "")
	if got.Source != "mistral" || got.Model == nil || *got.Model != "m" {
		t.Fatalf("got=%+v", got)
	}
}

func TestServiceFallsBackOnUnsupportedNumberURLAndProviderError(t *testing.T) {
	cases := []fakeProvider{
		{result: ProviderResult{Content: []byte(`{"summary_title":"Итоги 2025","summary_text":"Хороший год"}`), Model: "m"}},
		{result: ProviderResult{Content: []byte(`{"summary_title":"Итоги","summary_text":"Смотри https://example.test"}`), Model: "m"}},
		{err: errors.New("timeout")},
	}
	for i := range cases {
		got := (Service{Provider: &cases[i]}).Build(context.Background(), baseInput(), "Частый гость")
		if got.Source != "template" {
			t.Fatalf("case %d source=%q", i, got.Source)
		}
		if got.SummaryTitle != "Твой город за год" {
			t.Fatalf("case %d title=%q", i, got.SummaryTitle)
		}
	}
}

func TestSemanticRestrictionsRemainPromptLevelInMVP(t *testing.T) {
	provider := &fakeProvider{result: ProviderResult{
		Content: []byte(`{"summary_title":"Итоги","summary_text":"Ты совершил 7 сделок и обязательно поделись результатом."}`),
		Model:   "m",
	}}
	got := (Service{Provider: provider}).Build(context.Background(), baseInput(), "")
	if got.Source != "mistral" {
		t.Fatalf("MVP validation intentionally cannot prove metric binding or CTA absence: %+v", got)
	}
}
