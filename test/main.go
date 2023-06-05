package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/", ServeHTTP)
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err)
	}
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
<html>
<body>  
<table style="width:100%">
  <tr>
    <th>Data</th>
    <th>Content</th> 
  </tr> 
  <tr>
    <td>{{.Data}}</td>
    <td>{{.Content}}</td>
  </tr>
</table> 
</body>
</html>
`
	st := "Hi! How are you?
	Here is the link you wanted:
http://www.google.com"
	data := DataContent{"data", st}

	buf := &bytes.Buffer{}
	t := template.Must(template.New("template1").Parse(html))
	if err := t.Execute(buf, data); err != nil {
		panic(err)
	}
	body := buf.String()
	body = strings.Replace(body, "", "<br>", -1)
	fmt.Fprint(w, body)
}

type DataContent struct {
	Data, Content string
}