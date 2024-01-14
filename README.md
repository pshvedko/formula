# formula
calculating the value of an expression using the Dijkstra Shunting Yard algorithm

```go
	f, err := New("1+Sin(2*Pi*x)/2")
	if err != nil {
		log.Fatal(err)
	}
	var b []byte
	b, err = json.Marshal(f)
	if err != nil {
		log.Fatal(err)
	}
	var j Formula
	err = json.Unmarshal(b, &j)
	if err != nil {
		log.Fatal(err)
	}
	var v interface{}
	v, err = j.Evaluate(Bind{"Sin'": math.Sin, "Pi": math.Pi, "x": .25})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(f)
	fmt.Println(j)
	fmt.Println(v)
	// Output:
	// [1 2 Pi * x * Sin' 2 / +]
	// [1 2 Pi * x * Sin' 2 / +]
	// 1.5
  ```
