package p2

import "p1"

func Used_in_another_package() {
	p1.Global_1 = "modified"
	p1.Global_2 = "modified"
	p1.Global_3 = "modified"
	p1.Global_5 = "modified"

	p1.Global_Const_1 = "modified" // want "assignment to global variable marked with const"
	p1.Global_Const_2 = "modified" // want "assignment to global variable marked with const"
	p1.Global_Const_3 = "modified" // want "assignment to global variable marked with const"
	p1.Global_Const_5 = "modified" // want "assignment to global variable marked with const"
	p1.Global_Const_6 = "modified" // want "assignment to global variable marked with const"
}

func Package_is_hidden() {
	type p struct {
		Global_1       string
		Global_Const_1 string
	}
	var p1 p
	p1.Global_1 = "modified"
	p1.Global_Const_1 = "modified"
}
