package main

import(
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-lookup/providers"
)

func main(){

     root := flag.String("root", "https://whosonfirst.mapzen.com/data/", "")
     wofid := flag.Int("wofid", -1, "")     

     flag.Parse()
     
     p, e := providers.NewWOFFSProvider(*root)

     if e != nil {
     	panic(e)
     }

     f, e := p.GetFeatureById(*wofid)

     if e != nil {
     	panic(e)
     }

     fmt.Println(f)
}
