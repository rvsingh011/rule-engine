package main

import (
	"fmt"
	"time"

	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
)

type MyFact struct {
	IntAttribute     int64
	StringAttribute  string
	BooleanAttribute bool
	FloatAttribute   float64
	TimeAttribute    time.Time
	WhatToSay        string
}

func (mf *MyFact) GetWhatToSay(sentence string) string {
	return fmt.Sprintf("Let say \"%s\"", sentence)
}

func (mf *MyFact) SetIntAttribute(value int64) {
	mf.IntAttribute = value
}

func (mf *MyFact) Execute(value string) {
	fmt.Println("Executed this rule", value)

}

func main() {

	// lets prepare a rule definition
	drls := `
rule CheckValues "Check the default values" salience 10 {
    when 
        MF.IntAttribute == 123 && MF.StringAttribute == "Some string value"
    then
		MF.Execute("CheckValues");
        MF.WhatToSay = MF.GetWhatToSay("Hello Grule");
		MF.SetIntAttribute(567);
		Changed("MF.IntAttribute");
        Retract("CheckValues");

}

rule CheckValues1 "Check the default values" salience 9 {
    when 
		MF.IntAttribute == 567 && MF.StringAttribute == "Some string value"
    then
		MF.Execute("CheckValues1");
		MF.SetIntAttribute(123);
        MF.WhatToSay = MF.GetWhatToSay("PQRS");
		Changed("MF.IntAttribute");
        Retract("CheckValues1");
}

rule CheckValues2 "Check the default values" salience 8 {
    when 
		MF.IntAttribute == 123 && MF.StringAttribute == "Some string value"
    then
		MF.Execute("CheckValues2");
        MF.WhatToSay = MF.GetWhatToSay("QWERTY");
        Retract("CheckValues2");
}
`
	myFact := &MyFact{
		IntAttribute:     123,
		StringAttribute:  "Some string value",
		BooleanAttribute: true,
		FloatAttribute:   1.234,
		TimeAttribute:    time.Now(),
	}

	dataCtx := ast.NewDataContext()
	err := dataCtx.Add("MF", myFact)
	if err != nil {
		panic(err)
	}

	knowledgeLibrary := ast.NewKnowledgeLibrary()
	ruleBuilder := builder.NewRuleBuilder(knowledgeLibrary)
	// Add the rule definition above into the library and name it 'TutorialRules'  version '0.0.1'
	bs := pkg.NewBytesResource([]byte(drls))
	err = ruleBuilder.BuildRuleFromResource("TutorialRules", "0.0.1", bs)
	if err != nil {
		panic(err)
	}

	knowledgeBase := knowledgeLibrary.NewKnowledgeBaseInstance("TutorialRules", "0.0.1")

	engine := engine.NewGruleEngine()
	err = engine.Execute(dataCtx, knowledgeBase)
	if err != nil {
		panic(err)
	}

	fmt.Println("Final Print", myFact.WhatToSay)
	// this should prints
	// Lets Say "Hello Grule"
}
