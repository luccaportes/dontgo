# DontGo

[![GoDoc](https://godoc.org/github.com/lucp2/dontgo?status.svg)](https://godoc.org/github.com/lucp2/dontgo)

An unofficial client for using dontpad.com with go. Create temporary online persistence.

## dontpad.com

Dontpad is an online and easily shareable notepad. You can create and edit pages using any identifier simply by accessing dontpad.com/{identifier}

## Installation
You can install with a simple go get  
` go get github.com/lucp2/dontgo`

## Usage

There are basically four methods you can call. All are exemplified here and documented in godoc.

```
package main

import(
    "github.com/lucp2/dontgo"
    "fmt"
)

func main() {
	// first lets write something in a page with the identifier "dontgo_test"
	// you can pass variables with any type to be written. With we can not simply convert it
	// to string, we will transform it into json.
	err := dontgo.Write("dontgo_test", "hello world")

	// now with access dontpad.com/dontgo_test we will see our "hello world" there

	// if there is anything wrong, we will tell you. 
	// What can go wrong are basically connection problems
	// or we couldn't convert what you passed to string (what did you pass?)
	if err != nil {
		fmt.Println(err)
		return
	}

	// we can also retrieve what we've just written
	a := dontgo.Read("dontgo_test")
	fmt.Println(a) //will print "hello world"

	// if there is anything wrong with the read function (such as connection), it will panic.
	// this is proposital, allowing you to use it directly as a value.
	// but if you don't this to happen, use ReadNoPanic instead.

	a, err = dontgo.ReadNoPanic("dontgo_test")
	if err != nil {
		// now you can treat it
		fmt.Println(err)
		return
	}
	fmt.Println(a) //will print "hello world"

	// the write function will overwrite anything that is already in the identifier's page.
	// if that is not your intention, use Append.
	err = dontgo.Append("dontgo_test", ", how you doing?")
	if err != nil {
		fmt.Println(err)
		return
	}

	// and if we read it again now
	a = dontgo.Read("dontgo_test")
	fmt.Println(a) //will print "hello world, how you doing?"

	// we can also clear a page simply by calling Clear
	err = dontgo.Clear("dontgo_test")
	if err != nil {
		fmt.Println(err)
		return
	}

	// and if we read it again now
	a = dontgo.Read("dontgo_test")
	fmt.Println(a) //will print nothing
}

```

## Warning

Do not ever use this in production. First because this is unofficial and there is no guarantee that dontpad will not change the way they store the contents. And second because anyone is able to overwrite the identifiers without any authorization.

