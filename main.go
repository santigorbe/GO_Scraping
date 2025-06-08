package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

func main() {
	c := colly.NewCollector()
	c.OnHTML(".quote", func(e *colly.HTMLElement) {
		quote := strings.TrimSpace(e.ChildText("span.text"))     // Buscamos un span con la clase "text"
		author := strings.TrimSpace(e.ChildText("small.author")) // Buscamos un small con la clase "author"

		var tags []string
		e.ForEach("div.tags a.tag", func(_ int, el *colly.HTMLElement) {
			tags = append(tags, el.Text)
		})

		fmt.Printf("<<%s>> - %s [%s]\n", quote, author, strings.Join(tags, ", "))
	})

	// Scrapear el resto de paginas
	c.OnHTML("li.next a", func(e *colly.HTMLElement) {
		nextPage := e.Request.AbsoluteURL(e.Attr("href"))
		fmt.Println("Siguiente pagina:", nextPage)
		if err := e.Request.Visit(nextPage); err != nil {
			log.Fatal(err)
		}
	})

	//Limitamos el tiempo en que hace todas las consultas para que no nos bloqueen o saturemos el servidor
	if err := c.Limit(&colly.LimitRule{
		DomainGlob:  "*quotes.toscrape.*",
		Parallelism: 2, // Cuantas consultas paralelas
		RandomDelay: 500 * time.Millisecond,
	}); err != nil {
		log.Fatal(err)
	}

	if err := c.Visit("https://quotes.toscrape.com/page/1"); err != nil {
		log.Fatal(err)
	}
}

// Callbacks mas utilizados:
/*
Para ver que URL visitamos
c.OnRequest(func(r *colly.Request) {
	fmt.Println("Visiting", r.URL)
}

Para chequear el status code
c.OnResponse(func(r *colly.Response) {
	fmt.Println("Status code: ", r.StatusCode)
}

Para capturar errores sin parar todoconst
c.OnError(func(r *colly.Response, err error) {
	log.Println("Error en ", r.Request.URL, ":", err)
})


*/
