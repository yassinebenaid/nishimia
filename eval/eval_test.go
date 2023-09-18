package eval

import (
	"testing"

	"github.com/yassinebenaid/nishimia/lexer"
	"github.com/yassinebenaid/nishimia/object"
	"github.com/yassinebenaid/nishimia/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"10;", 10},
		{"156;", 156},
		{"785698;", 785698},
		{"-526;", -526},
		{"-11;", -11},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalStringExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"hello world";`, "hello world"},
		{`"hello" + " " + "world";`, "hello world"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testStringObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true;", true},
		{"false;", false},
		{"1 < 5;", true},
		{"4 >= 5;", false},
		{"2*2 <= 50;", true},
		{"1+1*2 <= 85;", true},
		{"false && false;", false},
		{"false || true;", true},
		{"1 == 1;", true},
		{"1 != 1;", false},
		{"true && false;", false},
		{"true && !false;", true},
		{"true || !false;", true},
		{"true || !false;", true},
		{"(true && false) != true;", true},
		{"(true && false) != true;", true},
		{"(false && false) != true;", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestBangOperatorExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true;", false},
		{"!false;", true},
		{"!!false;", false},
		{"!!!true;", false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { +10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
		{"if (1 <= 2) { 10 } else { 20 }", 10},
		{"if (1 >= 2) { 10 } else { 20 }", 20},
		{"if (true && false) { 10 } else { 20 }", 20},
		{"if (false || true) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{"9; return true; 9;", true},
		{"9; return false; 9;", false},
		{"9; return true && false; 9;", false},
		{"9; return false || true; 9;", true},
		{"9; return 10 <= 2; 9;", false},
		{`if 10 > 1 {
			if 10 > 1 {
				if true {
					return 10;
				}
			}
			return 1;
		}`, 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch v := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(v))
		case bool:
			testBooleanObject(t, evaluated, v)
		}
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{"5 + true;", "invalid operation: 5 + true (mismatched types INTEGER and BOOLEAN)"},
		{"5 + true; 5;", "invalid operation: 5 + true (mismatched types INTEGER and BOOLEAN)"},
		{"-true", "invalid operation: -true (operator \"-\" not defined on BOOLEAN)"},
		{"+false", "invalid operation: +false (operator \"+\" not defined on BOOLEAN)"},
		{"!1", "invalid operation: !1 (operator \"!\" not defined on INTEGER)"},
		{"!(1*5)", "invalid operation: !5 (operator \"!\" not defined on INTEGER)"},
		{"if(1+5){5}", "non-boolean value in if-statement , ( got=INTEGER, want=BOOLEAN )"},
		{"if(1){5}", "non-boolean value in if-statement , ( got=INTEGER, want=BOOLEAN )"},
		{"a;", "undefined identifier : a"},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)",
				evaluated, evaluated)
			continue
		}
		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q",
				tt.expectedMessage, errObj.Message)
		}
	}
}

func TestVarStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"var a = 5; a;", 5},
		{"var a = 5 * 5; a;", 25},
		{"var a = 5; var b = a; b;", 5},
		{"var a = 5; var b = a; var c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}

	testNullObject(t, testEval("var num = 10;"))
}

func TestFunction(t *testing.T) {
	input := "func(x) { x + 2; };"

	evaluated := testEval(input)

	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Params) != 1 {
		t.Fatalf("function has wrong params. Params=%+v",
			fn.Params)
	}

	if fn.Params[0].String() != "x" {
		t.Fatalf("param is not 'x'. got=%q", fn.Params[0])
	}

	expectedBody := "{(x + 2)}"

	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}

}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"func(x) { return x; }(5);", 5},
		{"var identity = func(x) { return x; }; identity(5);", 5},
		{"var identity = func(x) { return x; }; identity(5);", 5},
		{"var double = func(x) { return x * 2; }; double(5);", 10},
		{"var add = func(x, y) { return x + y; }; add(5, 5);", 10},
		{"var add = func(x, y) { return x + y; }; add(5 + 5, add(5, 5));", 20},
		{"var num = 4;var add = func(num) { return num + 2; }; add(num);", 6},
		{"var num = 4;var add = func(num) { return num + 2; }; add(num+5);", 11},
		{"var num = 4;var add = func(num) { return num + 2; }; add(num*2);", 10},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestClosures(t *testing.T) {
	input := `
	var i = 0;
	var newAdder = func(x) {
		var it = i;
		return func(y) { return x + y + +it; };
	};
	var addTwo = newAdder(2);
	addTwo(2);
	`
	testIntegerObject(t, testEval(input), 4)
}

func testEval(inp string) object.Object {
	lex := lexer.New(inp)
	par := parser.New(lex)
	program := par.ParseProgram()
	return Eval(program, object.NewEnvirement())
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
		return false
	}

	return true
}

func testStringObject(t *testing.T, obj object.Object, expected string) bool {
	result, ok := obj.(*object.String)
	if !ok {
		t.Errorf("object is not String. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%s, want=%s",
			result.Value, expected)
		return false
	}

	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t",
			result.Value, expected)
		return false
	}

	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func testErrorObject(t *testing.T, obj object.Object) bool {
	if _, ok := obj.(*object.Error); !ok {
		t.Errorf("object is not Error. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}
