# nishimia
a fully fnctional interpreter built on top of golang, this is an educational project, and so its not that fancy , but though , its still fully functional, with support for
variables and bindings, functions , and if-conditions , also functions are first-class sitizens , this means you can pass them as arguments or return them as values, and ofcourse error handling out of the box

here is a sinppet of all available features :

```go
var five = 5;
var ten = 10;

var add = func(x, y) {
	return x + y;
};

var multiply = func(x, y) {
	return x * y;
};

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
	} else {
		return false;
	}
};

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

var getAdditionClosure = func(x) {
	return func(i) { return x + i;};
};

var addition = add(five, ten);
var multiplication = multiply(five, ten);
var devision = devide(five, ten);
var maximum = max(five, ten);
var fiveIsPositive = isPositive(five);
var fiveIsZero = isZero(five);
var fiveIsNotZero = isNotZero(five);
var tenIsNegativeOrZero = isNegativeOrZero(five);
var name = "yassine benaid";
```
