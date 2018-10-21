// Package dontgo is an unofficial client for using dontpad.com with go.
// It allows you to write and read in the dontpad website.
// You might want to use this package as a temporary way of applying persistence to your variables.
// Do not ever use this in production. First because this is unofficial and there is no guarantee that dontpad
// will not change the way they store the contents. And second because anyone is able to overwrite the identifiers
// without any authorization.
package dontgo

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"gopkg.in/headzoo/surf.v1"
	"strings"
	"strconv"
	"encoding/json"
)

const (
	baseURL = "http://dontpad.com/"
)

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

func Read(identifier string) (r string) {
	if identifier[len(identifier)-4:] == ".zip"{
		fmt.Println("Warning: Using .zip at the end of the identifier causes dontpad.com to generate a zip file and not returning the actual content.")
		fmt.Println("We have removed the \".zip\" for you")
		identifier = identifier[:len(identifier)-4]
	}
	bow := surf.NewBrowser()
	bow.Open("http://dontpad.com/" + identifier)
	bow.Find("textarea").Each(func(_ int, s *goquery.Selection) {
		r = s.Text()
	})
	return r
}

func Append(identifier string, content interface{}) error {
	currentText := Read(identifier)
	text, err := anyToString(content)
	if err != nil{
		return err
	}
	return Write(identifier, currentText+text)
}

func Clear(identifier string) error {
	return Write(identifier, "")
}

// We need generics
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
