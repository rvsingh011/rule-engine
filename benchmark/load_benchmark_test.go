package benchmark

import (
	"fmt"
	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
	"io/ioutil"
	"testing"
)

/**
  Benchmarking `ast.KnowledgeBase` by loading 100 and 1000 rules into knowledgeBase
  Please refer docs/benchmarking_en.md for more info
*/
type RideFact struct {
	Distance           int32
	Duration           int32
	RideType           string
	IsFrequentCustomer bool
	Result             bool
	NetAmount          float32
}

func Benchmark_Grule_Load_Rules(b *testing.B) {
	rules := []struct {
		name string
		fun  func(string)
		filename string
	}{
		{"100 rules - 1", load100RulesIntoKnowledgeBase, "rule-1.grl"},
		{"100 rules - 2", load100RulesIntoKnowledgeBase, "rule-2.grl"},
		{"100 rules - 3", load100RulesIntoKnowledgeBase, "rule-3.grl"},
		{"100 rules - 4", load100RulesIntoKnowledgeBase, "rule-4.grl"},
	}
	for _, rule := range rules {
		for k := 0; k < 10; k++ {
			b.Run(fmt.Sprintf("%s", rule.name), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					rule.fun(rule.filename)
				}
			})
		}
	}
}

func load100RulesIntoKnowledgeBase(filename string) {
	input, _ := ioutil.ReadFile(filename)
	rules := string(input)
	fact := &RideFact{
		Distance: 6000,
		Duration: 121,
	}
	dctx := ast.NewDataContext()
	_ = dctx.Add("Fact", fact)

	lib := ast.NewKnowledgeLibrary()
	rb := builder.NewRuleBuilder(lib)
	_ = rb.BuildRuleFromResource("load_rules_test", "0.1.1", pkg.NewBytesResource([]byte(rules)))
	_ = lib.NewKnowledgeBaseInstance("load_rules_test", "0.1.1")
}
