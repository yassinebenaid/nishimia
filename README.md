# nishimia

A fully functional interpreter built on top of [go programming language](https://go.dev), This is an educational project, and so its not that fancy , but though its still fully functional with support for :

- variables and bindings
- functions
- built in functions
- if-conditions
- Closures
- functions are first-class citizens , this means you can pass them as arguments or return them as values,
- error handling out of the box

none of the above features perfect , but as said, its still work perfectly .

here is a sinppet of all available features :

```go
var five = 5;
var ten = 10;

var add = func(x, y) {
	return x + y;
};

var added = add(five,ten);


var multiply = func(x, y) {
	return x * y;
};

var multiplied = multiply(five, add(ten,10));

var devide = func(x, y) {
	if y > 0 {
		return x / y;
	} else {
		return 0;
	}
};

var isPositive = func(x) {
	if x >= 0 {
		return true;
	}

	return false;
};

var positive = isPositive(-10);

var isZero = func(x) {
	return x == 0;
};

var isNotZero = func(x) {
	return x != 0;
};

var isNegativeOrZero = func(x) {
	return x <= 0;
};

var max = func(x, y) {
	if x > y {
		return x ;
	}

	if x < y {
		return y ;
	}

	return x;
};


var name = "yassine benaid";
len(name); // len is built in here , and this comment is not supported by the way

var getAdditionClosure = func(x) {
	return func(i) { return x + i;};
};


var additionClosure = getAdditionClosure(2);
additionClosure(5);

var AcceptClosure = func(closure,value){
	return closure(value);
};

AcceptClosure(func(v){
	return v * 15;
},10)
```

## Installation && Testing
To get started , clone this repository , then in the project directory run : `go build -o nishimia`  , this will built the executable for you, 

The interpreter comes with a `repl` , you can use to run the code in the fly, run `./nishimia` with no extra arguments to get started

![image](https://github.com/yassinebenaid/nishimia/assets/101285507/c4902ca9-e6e0-4a4d-b3b3-5886bdd2a018)


To run a source code from a file , run `./nishimia path/to/file.ns`


# Warning
please if you like this language and decided to rewrite all your SAAS projects and services using it, please don't , just don't  , thanks.
