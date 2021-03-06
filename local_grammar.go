package fmr

import (
	"fmt"

	"github.com/liuzl/ling"
	"github.com/mitchellh/hashstructure"
)

func localGrammar(d *ling.Document) (*Grammar, error) {
	if d == nil {
		return nil, fmt.Errorf("document is empty")
	}
	if len(d.Spans) == 0 {
		return nil, nil
	}
	g := &Grammar{Name: "local", Rules: make(map[string]*Rule)}
	for _, span := range d.Spans {
		if span.Annotations["value"] == nil {
			continue
		}
		m, ok := span.Annotations["value"].(map[string]interface{})
		if !ok {
			continue
		}
		terms := []*Term{&Term{Value: span.String(), Type: Terminal}}
		for k, values := range m {
			rb := &RuleBody{terms, nil}
			switch values.(type) {
			case []string:
				args := []*Arg{}
				for _, v := range values.([]string) {
					args = append(args, &Arg{"string", v})
				}
				list := &FMR{"nf.list", args}
				rb.F = &FMR{
					"nf.entity",
					[]*Arg{&Arg{"string", k}, &Arg{"func", list}},
				}
			}
			hash, err := hashstructure.Hash(rb, nil)
			if err != nil {
				return nil, err
			}
			if _, has := g.Rules[k]; has {
				g.Rules[k].Body[hash] = rb
			} else {
				g.Rules[k] = &Rule{k, map[uint64]*RuleBody{hash: rb}}
			}
		}
	}
	if len(g.Rules) == 0 {
		return nil, nil
	}
	if err := g.refine("l"); err != nil {
		return nil, err
	}
	return g, nil
}
