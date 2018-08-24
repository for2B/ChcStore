package htmltool

import (
	"fmt"
	"net/http"
)

func ResponHtml(w http.ResponseWriter, title string, body string) {
	fmt.Fprintln(w, `<title>`+title+`</title><h1 style="text-align:center;margin-top:50px;">`+body+`</h1>`)
}

func ResponDialog(w http.ResponseWriter, title string, body string) {
	fmt.Fprintln(w, `
		<title>`+title+`</title>
		<h1 style="text-align:center;margin-top:50px;">
		
		
		
		
		
		`+body+

		`
		</h1>
		`)
}
