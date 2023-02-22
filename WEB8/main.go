package main

import (
	"os"
	"text/template"
)

type User struct {
	Name  string
	Email string
	Age   int
}

func (u User) IsOld() bool {
	return u.Age > 30
}

func main() {

	user := User{Name: "jjm", Email: "jjm@hanmail.net", Age: 30}
	user2 := User{Name: "jjm2", Email: "jjm2@hanmail.net", Age: 32}
	users := []User{user, user2}
	//template 생성
	//tmpl, err := template.New("Tmp11").Parse("Name: {{.Name}}\nEmail: {{.Email}}\nAge: {{.Age}}")
	tmpl, err := template.New("Tmp11").ParseFiles("../templates/tmpl1.tmpl", "../templates/tmpl2.tmpl")
	if err != nil {
		panic(err)
	}

	// tmpl.Execute(os.Stdout, user)
	// tmpl.Execute(os.Stdout, user2)

	// tmpl.ExecuteTemplate(os.Stdout, "tmpl2.tmpl", user)
	// tmpl.ExecuteTemplate(os.Stdout, "tmpl2.tmpl", user2)

	tmpl.ExecuteTemplate(os.Stdout, "tmpl2.tmpl", users)

}
