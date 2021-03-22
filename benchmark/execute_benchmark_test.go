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
    IsReversal       bool
    AlreadyReversed  bool
    IsCurrentBatch bool
    IsPreviousBatch bool
    Action        string
}

func (mf *Facts) GetAction(typeString string) string {
    return fmt.Sprintf("%s", typeString)
}

var reversalFact = []Facts{
	{
		IsReversal: true,
		AlreadyReversed: false,
		IsPreviousBatch: true,
	},{
		IsReversal: false,
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
		IsAuth : true ,
		Target: "Hello",
		IsNetworking: true,
	},
	{
		IsAuth : true ,
		Target: "Hfnwe",
		IsNetworking: true,
	},
}



func TestGrule_Execution_Engine(b *testing.T) {
	// build Knowledge base 
	FactKB = CreateKnowledgeBase(FactFile, "fact")
	AuthKB = CreateKnowledgeBase(authFactFile, "authFact")


	e:= engine.NewGruleEngine()
	dataCtx := ast.NewDataContext()
	err := dataCtx.Add("facts", &reversalFact[0])
	if err != nil {
		fmt.Println("There was a error", err.Error())
	}
	err = e.Execute(dataCtx, FactKB)
	if err != nil {
		fmt.Print(err)
	}


	fmt.Println(reversalFact[0])

	// // check for multiple facts 
	// for k := 0; k < 10; k++ {
	// 	b.Run(fmt.Sprintf("%s", rule.name), func(b *testing.B) {
	// 		rule.fun(rule.filename)
	// 		for i := 0; i < b.N; i++ {
				
	// 			e := engine.NewGruleEngine()
	// 			//Fact1
	// 			dataCtx := ast.NewDataContext()
	// 			err := dataCtx.Add("Fact", &f1)
	// 			if err != nil {
	// 				b.Fail()
	// 			}
	// 			err = e.Execute(dataCtx, knowledgeBase)
	// 			if err != nil {
	// 				fmt.Print(err)
	// 			}
	// 		}
	// 	})
	// }


}

func CreateKnowledgeBase(filename string, FactType string) *ast.KnowledgeBase{
	input, _ := ioutil.ReadFile(filename)
	rules := string(input)
	lib := ast.NewKnowledgeLibrary()
	rb := builder.NewRuleBuilder(lib)

	_ = rb.BuildRuleFromResource(FactType, "0.1.1", pkg.NewBytesResource([]byte(rules)))
	return  lib.NewKnowledgeBaseInstance(FactType, "0.1.1")
}


