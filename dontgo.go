// Package dontgo is an unofficial client for using dontpad.com with go.
// It allows you to write and read in the dontpad website.
// You might want to use this package as a temporary way of applying persistence to your variables.
// Do not ever use this in production. First because this is unofficial and there is no guarantee that dontpad
// will not change the way they store the contents. And second because anyone is able to overwrite the identifiers
// without any authorization.
// You can also access what you wrote on the website dontpad.com/{identifier}
package dontgo

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"strconv"
	"encoding/json"

	"github.com/PuerkitoBio/goquery"
	"gopkg.in/headzoo/surf.v1"
)

const (
	baseURL = "http://dontpad.com/"
)

// Write function receives a identifier string and an empty interface content, which contains what you want to write
// in the identifier page. You can pass pretty much anything, if it is any type of int, float or bool we will simply
// convert it into a string. If it is anything else, we will convert it to json with the builtin json package.
// If it can not convert, we will return its error so you can treat it. This function overwrites any previous content
// using the same identifier.
func Write(identifier string, content interface{}) error {
	if identifier[len(identifier)-4:] == ".zip"{
		fmt.Println("Warning: Using .zip at the end of the identifier causes dontpad.com to generate a zip file and not allowing to write.")
		fmt.Println("We have removed the \".zip\" for you")
		identifier = identifier[:len(identifier)-4]
	}
	v := url.Values{}
	text, err := anyToString(content)
	if err != nil{
		return err
	}
	v.Set("text", text)
	s := v.Encode()
	req, err := http.NewRequest("POST", baseURL+identifier, strings.NewReader(s))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// Read function gets an identifier and returns a string with its contents. You can also access it on
// dontpad.com/{identifier}. We do not return an error because we want you to be able to access it directly as value,
// if there are any errors, keep in mind that it will panic. If you don't want this to happen, use ReadNoPanic instead.
func Read(identifier string) (r string) {
	if identifier[len(identifier)-4:] == ".zip"{
		fmt.Println("Warning: Using .zip at the end of the identifier causes dontpad.com to generate a zip file and not returning the actual content.")
		fmt.Println("We have removed the \".zip\" for you")
		identifier = identifier[:len(identifier)-4]
	}
	bow := surf.NewBrowser()
	err := bow.Open("http://dontpad.com/" + identifier)
	if err != nil{
		panic(err)
	}
	bow.Find("textarea").Each(func(_ int, s *goquery.Selection) {
		r = s.Text()
	})
	return r
}

// ReadNoPanic function is essentially the same as Read function excepts it returns an error together with the value,
// allowing you to treat it.
func ReadNoPanic(identifier string) (r string, err error) {
	if identifier[len(identifier)-4:] == ".zip"{
		fmt.Println("Warning: Using .zip at the end of the identifier causes dontpad.com to generate a zip file and not returning the actual content.")
		fmt.Println("We have removed the \".zip\" for you")
		identifier = identifier[:len(identifier)-4]
	}
	bow := surf.NewBrowser()
	err = bow.Open("http://dontpad.com/" + identifier)
	if err != nil{
		return "", err
	}
	bow.Find("textarea").Each(func(_ int, s *goquery.Selection) {
		r = s.Text()
	})
	return r, nil
}

// Append function is essentially the same as the Write function, excepts it does not overwrite the contents, if there
// is any, using the same identifier. It appends to the end of the previous content.
func Append(identifier string, content interface{}) error {
	currentText := Read(identifier)
	text, err := anyToString(content)
	if err != nil{
		return err
	}
	return Write(identifier, currentText+text)
}

// Clear function deletes any content previously set using a identifier.
func Clear(identifier string) error {
	return Write(identifier, "")
}

// We need generics.
func anyToString(i interface{}) (string, error) {
	switch i.(type) {
	case string:
		return i.(string), nil
	case int:
		return strconv.FormatInt(int64(i.(int)), 10), nil
	case uint:
		return strconv.FormatInt(int64(i.(uint)), 10), nil
	case int8:
		return strconv.FormatInt(int64(i.(int8)), 10), nil
	case uint8:
		return strconv.FormatInt(int64(i.(uint8)), 10), nil
	case int16:
		return strconv.FormatInt(int64(i.(int16)), 10), nil
	case uint16:
		return strconv.FormatInt(int64(i.(uint16)), 10), nil
	case int32: //or rune
		return strconv.FormatInt(int64(i.(int32)), 10), nil
	case uint32:
		return strconv.FormatInt(int64(i.(uint32)), 10), nil
	case int64:
		return strconv.FormatInt(i.(int64), 10), nil
	case uint64:
		return strconv.FormatUint(i.(uint64), 10), nil
	case float32:
		return fmt.Sprintf("%f", i), nil
	case float64:
		return fmt.Sprintf("%f", i), nil
	case bool:
		return strconv.FormatBool(i.(bool)), nil
	default:
		r, err := json.Marshal(i)
		return string(r), err
	}
}
