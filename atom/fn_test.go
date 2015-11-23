package atom

func S(s string) Text {
	return Text{Content: s}
}

func URLs(urls ...string) []Link {
	a := make([]Link, len(urls))
	for i, url := range urls {
		a[i] = Link{URL: url}
	}
	return a
}

func Persons(names ...string) []Person {
	a := make([]Person, len(names))
	for i, name := range names {
		a[i] = Person{Name: name}
	}
	return a
}
