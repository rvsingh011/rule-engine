//  Copyright hyperjumptech/grule-rule-engine Authors
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package benchmark

import (
	"fmt"
	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
	"io/ioutil"
	"testing"
)

/**
  Benchmarking `engine.Execute` function by running 100 and 1000 rules with different N values
  Please refer docs/benchmarking_en.md for more info
*/

var knowledgeBase *ast.KnowledgeBase

func Benchmark_Grule_Execution_Engine(b *testing.B) {
	rules := []struct {
		name string
		fun  func(string)
		filename string
	}{
		{"100 rules", load100RulesIntoKnowledgebase, "rules-1.grl"},
		{"100 rules", load100RulesIntoKnowledgebase, "rules-2.grl"},
		{"100 rules", load100RulesIntoKnowledgebase, "rules-3.grl"},
		{"100 rules", load100RulesIntoKnowledgebase, "rules-4.grl"},

	}
	for _, rule := range rules {
		for k := 0; k < 10; k++ {
			b.Run(fmt.Sprintf("%s", rule.name), func(b *testing.B) {
				rule.fun(rule.filename)
				for i := 0; i < b.N; i++ {
					f1 := RideFact{
						Distance: 6000,
						Duration: 121,
					}
					e := engine.NewGruleEngine()
					//Fact1
					dataCtx := ast.NewDataContext()
					err := dataCtx.Add("Fact", &f1)
					if err != nil {
						b.Fail()
					}
					err = e.Execute(dataCtx, knowledgeBase)
					if err != nil {
						fmt.Print(err)
					}
				}
			})
		}
	}
}

func load100RulesIntoKnowledgebase(filename string) {
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
	_ = rb.BuildRuleFromResource("exec_rules_test", "0.1.1", pkg.NewBytesResource([]byte(rules)))
	knowledgeBase = lib.NewKnowledgeBaseInstance("exec_rules_test", "0.1.1")
}


