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
	Code             string
}

func (mf *MyFact) GetWhatToSay(sentence string) string {
	return fmt.Sprintf("Let say \"%s\"", sentence)
}

func (mf *MyFact) SetCode(code string) {
	mf.Code = code
	fmt.Println(mf.Code)
}

func (mf *MyFact) GetCode() string {
	return mf.Code
}

func (mf *MyFact) Print(value string) {
	fmt.Println(value)
	fmt.Println("Code", mf.Code)
}

func main() {
	// lets prepare a rule definition
	drls := `
rule checkValue1 "Check the default values" salience 98 {
	when 
		MF.GetCode() == "100"
	then
		MF.SetCode("200");
		MF.Print("checkValue1");
		Retract("checkValue1");
}

rule CheckCode "Check the default values" salience 99 {
	when 
		MF.Code == "000"
	then
		MF.SetCode("100");
		MF.Print("CheckCode");
		Retract("CheckCode");
}


rule CheckValues "Check the default values" salience 100 {
    when 
        MF.IntAttribute == 123 && MF.StringAttribute == "Some string value"
    then
		MF.Print("ChekValue");
		MF.SetCode("000");
		MF.Print("ChekValue");
        Retract("CheckValues");
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

	fmt.Println("Final Print", myFact.Code)
	// this should prints
	// Lets Say "Hello Grule"
}
