package main

import (
	"flag"
	"fmt"

	//"io/ioutil"

	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

type Data_Env struct {
	App_Struct           string
	App_Struct_CallBacks []string
	FuncCallBacks        []string
	IntFuncCallBacks     []string
	JSO_list             []string
	CreateConstructor    string

	TypeVar_list map[string]string
	ChildsList   []string

	ObjectList []string
}

func New_Data_Env() *Data_Env {
	de := Data_Env{}
	de.TypeVar_list = make(map[string]string)
	return &de
}

var Template_File string = `
package main

import (
    "fmt"
    "syscall/js"
    
	object "github.com/wanderer69/js_object"
)

type FuncCallBack func(this js.Value, args []js.Value, jso *object.JSObject) interface{}

%v

%v

%v
/*
%v
*/
`

func (de *Data_Env) CreateFile() string {
	ss0 := strings.Join(de.FuncCallBacks, " ")
	ss1 := strings.Join(de.ObjectList, " ")

	res := fmt.Sprintf(Template_File, de.App_Struct, ss0, de.CreateConstructor, ss1)
	return res
}

var Template_App_Struct string = `
type AppT struct {
	Jsoa            []*object.JSObject
	Jsod            map[string]*object.JSObject
	%v
	Context         map[string]interface{}	
}
`

func (de *Data_Env) Create_App_Struct() {
	ll := []string{}
	for i := range de.App_Struct_CallBacks {
		ll = append(ll, fmt.Sprintf("%v\r\n", de.App_Struct_CallBacks[i]))
	}
	de.App_Struct = fmt.Sprintf(Template_App_Struct, strings.Join(ll, "	"))
}

var Template_CallBack string = `	%v_%v          func( js.Value, []js.Value, *object.JSObject, *AppT) interface{}
`

func (de *Data_Env) Add_CallBack(ItemName string, MethodName string) {
	res := fmt.Sprintf(Template_CallBack, ItemName, MethodName)
	de.App_Struct_CallBacks = append(de.App_Struct_CallBacks, res)
}

var Template_ObjectListItem string = `	%v := jsod["%v"]
`

func (de *Data_Env) Add_ObjectList(ItemName string, ClassObjectName string) {
	res := fmt.Sprintf(Template_ObjectListItem, ItemName, ClassObjectName)
	de.ObjectList = append(de.ObjectList, res)
}

var Template_FuncCallBack string = `
func (at AppT) %v_%v_CallBack(this js.Value, args []js.Value, jso *object.JSObject, atp *AppT) interface{} {
	fmt.Printf("%v_%v_CallBack\r\n")
	if atp.%v_%v != nil {
		atp.%v_%v(this, args, jso, atp)
	}
	return nil
}
`

func (de *Data_Env) Add_FuncCallBack(ItemName string, MethodName string) {
	res := fmt.Sprintf(Template_FuncCallBack, ItemName, MethodName, ItemName, MethodName, ItemName, MethodName, ItemName, MethodName)
	de.FuncCallBacks = append(de.FuncCallBacks, res)
	de.Add_CallBack(ItemName, MethodName)
}

var Template_CreateConstructor string = `
func CreateConstructor() ([]*object.JSObject, map[string]*object.JSObject) {
	doc := object.NewDocument()
	doc.Type = "figma"
	jsoa := []*object.JSObject{}
	jso_dict := make(map[string]*object.JSObject)
	jso := &object.JSObject{}
	%v
	%v
	%v
	return jsoa, jso_dict
}
`

func (de *Data_Env) Add_CreateConstructor() {
	sl := []string{}
	for _, v := range de.TypeVar_list {
		sl = append(sl, fmt.Sprintf("%v\r\n", v))
	}
	//	fmt.Printf(sl "%v\r\n", sl)
	ss := strings.Join(sl, " ")
	scl := strings.Join(de.ChildsList, " ")
	de.CreateConstructor = fmt.Sprintf(Template_CreateConstructor, ss, strings.Join(de.JSO_list, " "), scl)
}

func (de *Data_Env) Add_ChildsList(ClassGroupName string, ChildName string, Anchor string) {
	sf := `	jso_g, ok = jso_dict["%v"]
	if !ok {
		fmt.Printf("Error! no group object %v!\r\n")
	}
	jso_c, ok = jso_dict["%v"]
	if !ok {
		fmt.Printf("Error! no child object %v!\r\n")
	}
	jso_g.Childs = append(jso_g.Childs, jso_c)
	jso_c.Parent = jso_g
`
	/*
	   	sfa := `	jso_g.CorrectText("%v")

	   `
	*/
	_, ok := de.TypeVar_list["jso_g"]
	if !ok {
		de.TypeVar_list["jso_g"] = "var jso_g *object.JSObject"
	}

	_, ok = de.TypeVar_list["jso_c"]
	if !ok {
		de.TypeVar_list["jso_c"] = "var jso_c *object.JSObject"
	}

	_, ok = de.TypeVar_list["ok"]
	if !ok {
		de.TypeVar_list["ok"] = "var ok bool"
	}
	ss := fmt.Sprintf(sf, ClassGroupName, ClassGroupName, ChildName, ChildName)
	de.ChildsList = append(de.ChildsList, ss)
	/*
		if len(Anchor) > 0 {
			ss := fmt.Sprintf(sfa, Anchor)
			de.ChildsList = append(de.ChildsList, ss)
		}
	*/
}

func (de *Data_Env) Add_CorrectText(ClassGroupName string, Anchor string) {
	sfa := `	jso_g, ok = jso_dict["%v"]
	if !ok {
		fmt.Printf("Error! no group object %v!\r\n")
	}
	jso_g.CorrectText("%v")	
`
	_, ok := de.TypeVar_list["jso_g"]
	if !ok {
		de.TypeVar_list["jso_g"] = "var jso_g *object.JSObject"
	}
	_, ok = de.TypeVar_list["ok"]
	if !ok {
		de.TypeVar_list["ok"] = "var ok bool"
	}
	if len(Anchor) > 0 {
		ss := fmt.Sprintf(sfa, ClassGroupName, ClassGroupName, Anchor)
		de.ChildsList = append(de.ChildsList, ss)
	}
}

func (de *Data_Env) Add_JSObject(TypeName string, Prefix string, ItemName string, ObjectName string) string {
	sl := []string{}
	result := ""
	switch TypeName {
	case "group":
		sf := `jso = doc.NewBlock("%v")
	bl = object.Block{}
	bl.ClickCBName = "%v_%v_click_CallBack"
	jso.ObjectExtender = bl
	jsoa = append(jsoa, jso)
	jso_dict["group_%v"] = jso

`
		_, ok := de.TypeVar_list["bl"]
		if !ok {
			de.TypeVar_list["bl"] = "var bl object.Block"
		}
		l2 := fmt.Sprintf(sf, ItemName, Prefix, ObjectName, ObjectName)
		sl = append(sl, l2)
	case "ordinary":
		sf := `jso = doc.NewOrdinary("%v")
	or = object.Ordinary{}
	jso.ObjectExtender = or
	jsoa = append(jsoa, jso)
	jso_dict["ordinary_%v"] = jso

`
		_, ok := de.TypeVar_list["or"]
		if !ok {
			de.TypeVar_list["or"] = "var or object.Ordinary"
		}
		l2 := fmt.Sprintf(sf, ItemName, ObjectName)
		sl = append(sl, l2)
	case "label":
		sf := `jso = doc.NewLabel("%v", "")
	lb = object.Label{}
	lb.ClickCBName = "%v_%v_click_CallBack"
	jso.ObjectExtender = lb
	jsoa = append(jsoa, jso)
	jso_dict["label_%v"] = jso

`
		_, ok := de.TypeVar_list["lb"]
		if !ok {
			de.TypeVar_list["lb"] = "var lb object.Label"
		}
		l2 := fmt.Sprintf(sf, ItemName, Prefix, ObjectName, ObjectName)
		sl = append(sl, l2)
	case "image":
		sf := `jso = doc.NewImage("%v")
	im = object.Image{}
	im.ClickCBName = "%v_%v_click_CallBack"
	jso.ObjectExtender = im
	jsoa = append(jsoa, jso)
	jso_dict["image_%v"] = jso

`
		_, ok := de.TypeVar_list["im"]
		if !ok {
			de.TypeVar_list["im"] = "var im object.Image"
		}
		l2 := fmt.Sprintf(sf, ItemName, Prefix, ObjectName, ObjectName)
		sl = append(sl, l2)
	case "text":
		sf := `jso = doc.NewText("%v", "")
	tx = object.Text{}
	tx.ChangeCBName = "%v_%v_change_CallBack"
	jso.ObjectExtender = tx
	jsoa = append(jsoa, jso)
	jso_dict["text_%v"] = jso
`
		_, ok := de.TypeVar_list["tx"]
		if !ok {
			de.TypeVar_list["tx"] = "var tx object.Text"
		}
		l2 := fmt.Sprintf(sf, ItemName, Prefix, ObjectName, ObjectName)
		sl = append(sl, l2)
	case "button":
		sf := `jso = doc.NewButton("%v", "")
	bt = object.Button{}
	bt.ClickCBName = "%v_%v_click_CallBack"
	jso.ObjectExtender = bt
	jsoa = append(jsoa, jso)
	jso_dict["button_%v"] = jso
`
		_, ok := de.TypeVar_list["bt"]
		if !ok {
			de.TypeVar_list["bt"] = "var bt object.Button"
		}
		l2 := fmt.Sprintf(sf, ItemName, Prefix, ObjectName, ObjectName)
		sl = append(sl, l2)
	case "list":
		sf := `	list_id := []string{}
	list_data := []string{}
	jso = doc.NewList("%v", list_id, list_data)
	ls = object.List{}
	l.ChangeCBName = "%v_%v_change_CallBack"
	jso.ObjectExtender = l
	jsoa = append(jsoa, jso)
	jso_dict["list_%v"] = jso
`
		_, ok := de.TypeVar_list["ls"]
		if !ok {
			de.TypeVar_list["ls"] = "var ls object.List"
		}
		l2 := fmt.Sprintf(sf, ItemName, Prefix, ObjectName, ObjectName)
		sl = append(sl, l2)
	case "selector":
		sf := `	list_id := []string{}
	list_data := []string{}
	jso = doc.NewSelector("%v", list_id, list_data)
	se = object.Selector{}
	se.ChangeCBName = "%v_%v_change_CallBack"
	jso.ObjectExtender = se
	jsoa = append(jsoa, jso)
	jso_dict["selector_%v"] = jso
`
		_, ok := de.TypeVar_list["se"]
		if !ok {
			de.TypeVar_list["se"] = "var se object.Selector"
		}
		l2 := fmt.Sprintf(sf, ItemName, Prefix, ObjectName, ObjectName)
		sl = append(sl, l2)
	case "table":
		/*
			for i, _ := range lo.HeaderItems {
				list_id = append(list_id, fmt.Sprintf("header_%v", i))
			}
			for i, _ := range lo.TableRowItems {
				ld := []string{}
				for j, _ := range lo.TableRowItems {
					ld = append(ld, fmt.Simple("cell_%v_%v", i, j)
				}
				list_data = append(list_data, ld)
			}

		*/
		sf := `	list_id := []string{}
	list_data := [][]string{}
	jso = doc.NewTable("%v", list_id, list_data)
	tb := object.Table{}
	tb.ClickCBName = "%v_%v_click_CallBack"
	tb.ChangeCBName = "%v_%v_change_CallBack"
	jso.ObjectExtender = tb
	jsoa = append(jsoa, jso)
	jso_dict["table_%v"] = jso
`
		_, ok := de.TypeVar_list["tb"]
		if !ok {
			de.TypeVar_list["tb"] = "var tb object.Table"
		}
		l2 := fmt.Sprintf(sf, ItemName, Prefix, ObjectName, Prefix, ObjectName, ObjectName)
		sl = append(sl, l2)
		/*
			case "tree":
				lo := TreeO{}
				transcode(dci.Object.ObjectExtender, &lo)

				type ConvertFunc func(tree_data []TreeItemO) []TreeData
				var cf ConvertFunc
				//g_num := 1
				cf = func(tree_data []TreeItemO) []TreeData {
					tia := []TreeData{}
					for i, _ := range tree_data {
						ti := TreeData{}
						ti.Data = tree_data[i].Data
						ti.TreeDatas = cf(tree_data[i].TreeItems)
						tia = append(tia, ti)
					}
					return tia
				}
				tree_data := cf(lo.TreeItems)

				list_id := "tree"
				jso = doc.NewTree(dci.Object.ObjectID, list_id, tree_data)
				l := Tree{}
				l.ClickCBName = lo.ClickCB
				l.ChangeCBName = lo.ChangeCB
				jso.ObjectExtender = l
		*/
	}
	de.JSO_list = append(de.JSO_list, sl...)
	return result
}

func main() {

	var file_input string
	flag.StringVar(&file_input, "file_input", "", "input file name with path")
	var path_output string
	flag.StringVar(&path_output, "path_output", "", "path to output file app_data.go")

	flag.Parse()

	if len(file_input) == 0 {
		fmt.Printf("Must be path and name input file\r\n")
		flag.PrintDefaults()
		return
	}
	if len(path_output) == 0 {
		fmt.Printf("Must be path to output file\r\n")
		flag.PrintDefaults()
		return
	}

	stati, err := os.Stat(file_input)
	if os.IsNotExist(err) {
		fmt.Printf("File does not exist.")
		return
	}
	if stati.IsDir() {
		fmt.Printf("Exist directory. No file.")
		return
	}

	stato, err := os.Stat(path_output)
	if os.IsNotExist(err) {
		fmt.Printf("Path does not exist.")
	}
	if !stato.IsDir() {
		fmt.Printf("Exist file. No directory.")
		return
	}

	bs, err := os.ReadFile(file_input)
	if err != nil {
		log.Fatal(err)
	}

	text := string(bs)

	de := New_Data_Env() // &Data_Env{}

	doc, err := html.Parse(strings.NewReader(text))
	if err != nil {
		log.Fatal(err)
	}
	type ElemStack struct {
		Block  string
		Name   string
		Childs []string
		Anchor string
	}
	esd := make(map[string]ElemStack)
	ord_cont := 1
	var f func(*html.Node, []ElemStack)
	f = func(n *html.Node, es []ElemStack) {
		var esc ElemStack
		if len(es) > 0 {
			esc = es[len(es)-1]
		}
		if n.Type == html.ElementNode && n.Data == "div" {
			for _, a := range n.Attr {
				// fmt.Printf("attr %v\r\n", a)
				if a.Key == "class" {
					// fmt.Println(a.Val)
					class_name := a.Val
					cnl := strings.Split(class_name, "_")
					switch cnl[0] {
					case "desktop":
					case "rectangle":
					case "button":
						// fmt.Printf("%v button %v a.Key %v, class_name %v, cnl[1] %v\r\n", es, cnl[0], a.Key, class_name, cnl[1])
						if len(es) > 0 {
							esc_, ok := esd[esc.Name]
							if ok {
								de.Add_JSObject(cnl[0], "Button", class_name, cnl[1])
								de.Add_FuncCallBack("Button_"+cnl[1], "click")
								de.Add_ObjectList("button_"+cnl[1], "button_"+cnl[1])
								esc_.Childs = append(esc_.Childs, "button_"+cnl[1])
								esd[esc.Name] = esc_
							} else {
								// error!!!
								log.Printf("error: len stack element == 0\r\n")
							}
						}
					case "label":
						// fmt.Printf("%v label %v a.Key %v, class_name %v, cnl[1] %v\r\n", es, cnl[0], a.Key, class_name, cnl[1])
						if len(es) > 0 {
							esc_, ok := esd[esc.Name]
							if ok {
								de.Add_JSObject(cnl[0], "Label", class_name, cnl[1])
								de.Add_FuncCallBack("Label_"+cnl[1], "click")
								de.Add_ObjectList("label_"+cnl[1], "label_"+cnl[1])
								esc_.Childs = append(esc_.Childs, "label_"+cnl[1])
								esd[esc.Name] = esc_
								escc := ElemStack{cnl[0], cnl[1], []string{}, ""}
								es = append(es, escc)
								esd[escc.Name] = escc
							} else {
								// error!!!
								log.Printf("error: len stack element == 0\r\n")
							}
						}
					case "group":
						// fmt.Printf("group %v\r\n", cnl[1])
						de.Add_JSObject(cnl[0], "Group", class_name, cnl[1])
						de.Add_FuncCallBack("Group_"+cnl[1], "click")
						de.Add_ObjectList("group_"+cnl[1], "group_"+cnl[1])
						escc := ElemStack{cnl[0], cnl[1], []string{}, ""}
						es = append(es, escc)
						esd[escc.Name] = escc
					case "text":
						// fmt.Printf("%v text %v a.Key %v, class_name %v, cnl[1] %v\r\n", es, cnl[0], a.Key, class_name, cnl[1])
						if len(es) > 0 {
							esc_, ok := esd[esc.Name]
							if ok {
								de.Add_JSObject(cnl[0], "Text", class_name, cnl[1])
								de.Add_FuncCallBack("Text_"+cnl[1], "change")
								de.Add_ObjectList("text_"+cnl[1], "text_"+cnl[1])
								esc_.Childs = append(esc_.Childs, "text_"+cnl[1])
								esd[esc.Name] = esc_
								escc := ElemStack{cnl[0], cnl[1], []string{}, ""}
								es = append(es, escc)
								esd[escc.Name] = escc
							} else {
								// error!!!
								log.Printf("error: len stack element == 0\r\n")
							}
						}
					case "image":
						// fmt.Printf("%v len(es) %v image %v a.Key %v, class_name %v, cnl[1] %v\r\n", es, len(es), cnl[0], a.Key, class_name, cnl[1])
						if len(es) > 0 {
							esc_, ok := esd[esc.Name]
							if ok {
								de.Add_JSObject(cnl[0], "Image", class_name, cnl[1])
								de.Add_FuncCallBack("Image_"+cnl[1], "click")
								de.Add_ObjectList("image_"+cnl[1], "image_"+cnl[1])
								esc_.Childs = append(esc_.Childs, "image_"+cnl[1])
								esd[esc.Name] = esc_
							} else {
								// error!!!
								log.Printf("error: len stack element == 0\r\n")
							}
						}
					default:
						// fmt.Printf("%v ordinary %v a.Key %v, class_name %v, cnl[1] %v\r\n", es, cnl[0], a.Key, class_name, cnl[1])
						if len(es) > 0 {
							esc_, ok := esd[esc.Name]
							if ok {
								ss := fmt.Sprintf("%v_%v", cnl[1], ord_cont)
								ord_cont = ord_cont + 1
								de.Add_JSObject("ordinary", "", class_name, ss)
								//de.Add_FuncCallBack("Ordinary_" + ss, "click")
								de.Add_ObjectList("ordinary_"+cnl[1], "ordinary_"+cnl[1])
								esc_.Childs = append(esc_.Childs, "ordinary_"+cnl[1])
								esd[esc.Name] = esc_
							} else {
								// error!!!
								log.Printf("error: len stack element == 0\r\n")
							}
						}
					}
					//break
				}

			}
		} else {
			// fmt.Printf("n.Type %v, n.Data %v\r\n", n.Type, n.Data)
			if n.Type == 3 && n.Data == "span" { // html.TextNode
				for _, a := range n.Attr {
					// fmt.Printf("attr %v\r\n", a)
					if a.Key == "class" {
						// fmt.Println(a.Val)
						class_name := a.Val
						cnl := strings.Split(class_name, "_")
						switch cnl[0] {
						case "example":
							// fmt.Printf("%v example %v a.Key %v, class_name %v, cnl[1] %v\r\n", es, cnl[0], a.Key, class_name, cnl[1])
							if len(es) > 0 {
								esc_, ok := esd[esc.Name]
								if ok {
									ss := fmt.Sprintf("%v_%v", cnl[1], ord_cont)
									ord_cont = ord_cont + 1
									de.Add_JSObject("ordinary", "", class_name, ss)
									//de.Add_FuncCallBack("Ordinary_" + ss, "click")
									//esc_.Childs = append(esc_.Childs, cnl[1])
									// fmt.Printf("esc.Block %v class_name %v\r\n", esc.Block, class_name)
									if esc.Block == "text" || esc.Block == "label" {
										esc_.Anchor = class_name
									}
									esd[esc.Name] = esc_
								} else {
									// error!!!
									log.Printf("error: len stack element == 0\r\n")
								}
							}
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c, es)
		}
	}
	f(doc, []ElemStack{})
	for _, v := range esd {
		// fmt.Printf("-> %#v\r\n", v)
		de.Add_CorrectText(v.Block+"_"+v.Name, v.Anchor)
		for i := range v.Childs {
			de.Add_ChildsList(v.Block+"_"+v.Name, v.Childs[i], v.Anchor)
		}
	}

	de.Create_App_Struct()

	de.Add_CreateConstructor()

	ss := de.CreateFile()

	//fmt.Printf("%v", ss)
	fo := filepath.Join(path_output, "app_data.go")
	_ = os.WriteFile(fo, []byte(ss), 0644)
}
