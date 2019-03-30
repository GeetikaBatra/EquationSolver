package main

import (
	"fmt"
	"net/http"
	"log"
	"strconv"
	"strings"
	"regexp"
)
////Struct that contains coefficiants of equation
type coefficiant struct {
	a float64
	b float64
	c float64
    
}
//resulting struct
type result struct{
	x float64
	y float64
}
//Initializes the coeff struct with input coefficiants
// Utilized in tests
func NewCoeff(a,b,c float64) *coefficiant{
	return &coefficiant{a:a,b:b,c:c}
}

//Initializes empty coefficiants
func EmptyNewCoeff() *coefficiant{
	return &coefficiant{}
}

//HTTP Server
func solver(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/equate" {
        http.Error(w, "404 not found.", http.StatusNotFound)
        return
    }
	switch r.Method {
		case "GET":     
		fmt.Fprintf(w, "Welcome")
		case "POST":
			if err := r.ParseForm(); err != nil {
				fmt.Fprintf(w, "ParseForm() err: %v", err)
				return
			}
			first_equation := r.FormValue("first_equation")
			second_equation := r.FormValue("second_equation")
			res:= FindDeterminents(first_equation, second_equation)
			fmt.Fprintf(w, "result_x = %v\n", res.x)
			fmt.Fprintf(w, "result_Y = %v", res.y)
		default:
			fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
		}
}

//Checks equation type using REGEX and extracts coefficaints from it
func (coeff *coefficiant)ValidateEquation(eq_ string){

		res, err := regexp.MatchString("\\d+x[+-]?\\d+y\\=\\d+", eq_)
		if res{

			coeff_num := coeff.ExtractCoeffFromString(eq_)
			coeff.a =  coeff_num[0]
			coeff.b =  coeff_num[1]
			coeff.c = coeff_num[2]
			return
		}
		res, err = regexp.MatchString("x\\+y\\=\\d+", eq_)
		if res{
			coeff_num:= coeff.ExtractCoeffFromString(eq_)
			coeff.a = 1
			coeff.b = 1
			coeff.c = coeff_num[0]
			return
		}
		res, err = regexp.MatchString("x\\-y\\=\\d+", eq_)
		if res{
			coeff_num:= coeff.ExtractCoeffFromString(eq_)
			coeff.a = 1
			coeff.b = -1
			coeff.c = coeff_num[0]
			return
		}
		res, err = regexp.MatchString("[-]?[\\d+]?x\\=\\d+", eq_)
		if res{
			coeff_num:= coeff.ExtractCoeffFromString(eq_)
			if len(coeff_num)==2{
				coeff.a = coeff_num[0]
				coeff.b = 0
				coeff.c = coeff_num[1]
				return
			}
			if(strings.HasPrefix(eq_, "-")){
				coeff.a = -1
			}else{
				coeff.a = 1
			}
			
			coeff.b = 0
			coeff.c = coeff_num[0]
			return
		}
		res, err = regexp.MatchString("[-]?[\\d+]?y\\=\\d+", eq_)
		if res{
			coeff_num:= coeff.ExtractCoeffFromString(eq_)
			if len(coeff_num)==2{
				coeff.a = 0
				coeff.b = coeff_num[0]
				coeff.c = coeff_num[1]
				return
			}
			if(strings.HasPrefix(eq_, "-")){
				coeff.b = -1
			}else{
				coeff.b = 1
			}
			
			coeff.a = 0
			coeff.c = coeff_num[0]
			return
		}
		if err!=nil{
			log.Fatal("Wrong format of equation")
		}
		return 
}

//Extracts all the numbers from an eqaution.
//Uses Regex to find the numbers and converts it to float
func(coeff *coefficiant) ExtractCoeffFromString(eq_ string) []float64{

	re := regexp.MustCompile("[-]?[0-9]+")
	coeff_arr := re.FindAllString(eq_, -1)
	var num_arr []float64

	for _, coeff_ := range coeff_arr{
		if num_, err := strconv.ParseFloat(coeff_, 64); err == nil {
			num_arr = append(num_arr, num_)
		}
	}
	return num_arr
}

//Uses Crammers rule to solve the equation
//Improvement: Use LCM to solve instead of crammer rule.
func FindDeterminents(first_equation, second_equation string) *result{

	var deter float64
	result := result{}
	first_eq := EmptyNewCoeff()
	first_eq.ValidateEquation(first_equation)
	second_eq := EmptyNewCoeff()
	second_eq.ValidateEquation(second_equation)
	deter = (first_eq.a*second_eq.b)-(first_eq.b*second_eq.a)
	if (deter!=0){
		result.x=((first_eq.c*second_eq.b)-(first_eq.b*second_eq.c))/deter
		result.y=((first_eq.a*second_eq.c)-(first_eq.c*second_eq.a))/deter
	}else{
	fmt.Printf("Cramer equations system: determinant is zero there are either no solutions or many solutions exist.")
	
	}
	return &result

}

func main() {
	http.HandleFunc("/equate", solver)
	log.Fatal(http.ListenAndServe(":8080", nil))
	
}
