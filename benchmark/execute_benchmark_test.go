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

var FactKB *ast.KnowledgeBase
var AuthKB *ast.KnowledgeBase

var FactFile = "./fact-file.grl"
var authFactFile = "./auth-file.grl"

type Facts struct {
	IsReversal      bool `json:"IsReversal"`
	AlreadyReversed bool
	IsCurrentBatch  bool
	IsPreviousBatch bool
	Action          string
}

func (mf *Facts) GetAction(typeString string) string {
	return fmt.Sprintf("%s", typeString)
}

var reversalFact = []Facts{
	{
		IsReversal:      true,
		AlreadyReversed: false,
		IsPreviousBatch: true,
	}, {
		IsReversal:      false,
		AlreadyReversed: false,
		IsPreviousBatch: true,
	},
}

type AuthFacts struct {
	IsAuth       bool
	Target       string
	IsNetworking bool
	IsRequestMsg bool
	UnknownEvent bool
	Types        string
	Name         string
	Handlers     []string
}

func (mf *AuthFacts) GetType(typeString string) string {
	return fmt.Sprintf("%s", typeString)
}
func (mf *AuthFacts) GetName(nameString string) string {
	return fmt.Sprintf("%s", nameString)
}
func (mf *AuthFacts) GetHandler(aString []string, subString string, subString2 string) []string {
	val := fmt.Sprintf("%s", subString)
	val2 := fmt.Sprintf("%s", subString2)
	aString = append(aString, val, val2)
	return aString
}

var authFacts = []AuthFacts{
	{
		IsAuth:       true,
		Target:       "Hello",
		IsNetworking: true,
	},
	{
		IsAuth:       true,
		Target:       "Hfnwe",
		IsNetworking: true,
	},
}

func Benchmark_Grule_Execution_Engine(b *testing.B) {
	// build Knowledge base
	FactKB = CreateKnowledgeBase(FactFile, "fact")
	AuthKB = CreateKnowledgeBase(authFactFile, "authFact")

	var foundAuthFact, foundReversalFact bool 
	var authFactParsed AuthFacts
	var reversalFactParsed Facts

	for j := 0; j < 10 ; j++ {

		b.Run(fmt.Sprintf("%d", j), func(b *testing.B) {

		for i := 0; i < b.N; i++ {
			var fact interface{}
			if i%2 == 0 {
				fact = authFacts[0]
			} else {
				fact = reversalFact[0]
			}
			if authFactParsed, foundAuthFact = fact.(AuthFacts); foundAuthFact {
				// fmt.Println("This fact is Auth Fact")
			} else if reversalFactParsed, foundReversalFact = fact.(Facts); foundReversalFact {
				// fmt.Println("This fact is reversal Fact")
			}
	
			e := engine.NewGruleEngine()
			//Fact1
			dataCtx := ast.NewDataContext()
	
			if foundAuthFact {
				err := dataCtx.Add("facts", &authFactParsed)
				if err != nil {
					b.Fail()
				}
				err = e.Execute(dataCtx, AuthKB)
				if err != nil {
					fmt.Print(err)
				}
			} else if foundReversalFact {
				err := dataCtx.Add("facts", &reversalFactParsed)
				if err != nil {
					b.Fail()
				}
	
				err = e.Execute(dataCtx, FactKB)
				if err != nil {
					fmt.Print(err)
				}
			}
	
		}
	})
	}

}

func CreateKnowledgeBase(filename string, FactType string) *ast.KnowledgeBase {
	input, _ := ioutil.ReadFile(filename)
	rules := string(input)
	lib := ast.NewKnowledgeLibrary()
	rb := builder.NewRuleBuilder(lib)

	_ = rb.BuildRuleFromResource(FactType, "0.1.1", pkg.NewBytesResource([]byte(rules)))
	return lib.NewKnowledgeBaseInstance(FactType, "0.1.1")
}
