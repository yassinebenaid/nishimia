# nishimia

A fully fnctional interpreter built on top of [go programming languae](https://go.dev), This is an educational project, and so its not that fancy , but though its still fully functional with support for :

- variables and bindings
- functions
- if-conditions
- Closures
- functions are first-class sitizens , this means you can pass them as arguments or return them as values,
- error handling out of the box

nothing of the above features is perfect , but as said, its still work perfectly

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
