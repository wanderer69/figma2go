package main

import (
	"flag"
	"fmt"

	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

type DataEnv struct {
	AppStruct          string
	AppStructCallbacks []string
	FuncCallBacks      []string
	IntFuncCallBacks   []string
	JSOList            []string
	CreateConstructor  string

	TypeVarList map[string]string
	ChildsList  []string

	ObjectList []string
}

func NewDataEnv() *DataEnv {
	de := DataEnv{}
	de.TypeVarList = make(map[string]string)
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

func (de *DataEnv) CreateFile() string {
	ss0 := strings.Join(de.FuncCallBacks, " ")
	ss1 := strings.Join(de.ObjectList, " ")

	res := fmt.Sprintf(Template_File, de.AppStruct, ss0, de.CreateConstructor, ss1)
	return res
}

var TemplateAppStruct string = `
type AppT struct {
	Jsoa            []*object.JSObject
	Jsod            map[string]*object.JSObject
	%v
	Context         map[string]interface{}	
}
`

func (de *DataEnv) CreateAppStruct() {
	ll := []string{}
	for i := range de.AppStructCallbacks {
		ll = append(ll, fmt.Sprintf("%v\r\n", de.AppStructCallbacks[i]))
	}
	de.AppStruct = fmt.Sprintf(TemplateAppStruct, strings.Join(ll, "	"))
}

var TemplateCallBack string = `	%v_%v          func( js.Value, []js.Value, *object.JSObject, *AppT) interface{}
`

func (de *DataEnv) AddCallback(ItemName string, MethodName string) {
	res := fmt.Sprintf(TemplateCallBack, ItemName, MethodName)
	de.AppStructCallbacks = append(de.AppStructCallbacks, res)
}

var TemplateObjectListItem string = `	%v := jsod["%v"]
`

func (de *DataEnv) AddObjectList(ItemName string, ClassObjectName string) {
	res := fmt.Sprintf(TemplateObjectListItem, ItemName, ClassObjectName)
	de.ObjectList = append(de.ObjectList, res)
}

var TemplateFuncCallback string = `
func (at AppT) %v_%v_CallBack(this js.Value, args []js.Value, jso *object.JSObject, atp *AppT) interface{} {
	fmt.Printf("%v_%v_CallBack\r\n")
	if atp.%v_%v != nil {
		atp.%v_%v(this, args, jso, atp)
	}
	return nil
}
`

func (de *DataEnv) AddFuncCallback(ItemName string, MethodName string) {
	res := fmt.Sprintf(TemplateFuncCallback, ItemName, MethodName, ItemName, MethodName, ItemName, MethodName, ItemName, MethodName)
	de.FuncCallBacks = append(de.FuncCallBacks, res)
	de.AddCallback(ItemName, MethodName)
}

var TemplateCreateConstructor string = `
func CreateConstructor() ([]*object.JSObject, map[string]*object.JSObject) {
	doc := object.NewDocument()
	doc.Type = "figma"
	jsoArray := []*object.JSObject{}
	jsoDict := make(map[string]*object.JSObject)
	jso := &object.JSObject{}
	%v
	%v
	%v
	return jsoArray, jsoDict
}
`

func (de *DataEnv) AddCreateConstructor() {
	sl := []string{}
	for _, v := range de.TypeVarList {
		sl = append(sl, fmt.Sprintf("%v\r\n", v))
	}
	//	fmt.Printf(sl "%v\r\n", sl)
	ss := strings.Join(sl, " ")
	scl := strings.Join(de.ChildsList, " ")
	de.CreateConstructor = fmt.Sprintf(TemplateCreateConstructor, ss, strings.Join(de.JSOList, " "), scl)
}

func (de *DataEnv) AddChildsList(ClassGroupName string, ChildName string, Anchor string) {
	sf := `	jsoG, ok = jsoDict["%v"]
	if !ok {
		fmt.Printf("Error! no group object %v!\r\n")
	}
	jsoC, ok = jsoDict["%v"]
	if !ok {
		fmt.Printf("Error! no child object %v!\r\n")
	}
	jsoG.Childs = append(jsoG.Childs, jsoC)
	jsoC.Parent = jsoG
`
	/*
	   	sfa := `	jsoG.CorrectText("%v")

	   `
	*/
	_, ok := de.TypeVarList["jsoG"]
	if !ok {
		de.TypeVarList["jsoG"] = "var jsoG *object.JSObject"
	}

	_, ok = de.TypeVarList["jsoC"]
	if !ok {
		de.TypeVarList["jsoC"] = "var jsoC *object.JSObject"
	}

	_, ok = de.TypeVarList["ok"]
	if !ok {
		de.TypeVarList["ok"] = "var ok bool"
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

func (de *DataEnv) AddCorrectText(ClassGroupName string, Anchor string) {
	sfa := `	jsoG, ok = jsoDict["%v"]
	if !ok {
		fmt.Printf("Error! no group object %v!\r\n")
	}
	jsoG.CorrectText("%v")	
`
	_, ok := de.TypeVarList["jsoG"]
	if !ok {
		de.TypeVarList["jsoG"] = "var jsoG *object.JSObject"
	}
	_, ok = de.TypeVarList["ok"]
	if !ok {
		de.TypeVarList["ok"] = "var ok bool"
	}
	if len(Anchor) > 0 {
		ss := fmt.Sprintf(sfa, ClassGroupName, ClassGroupName, Anchor)
		de.ChildsList = append(de.ChildsList, ss)
	}
}

func (de *DataEnv) AddJSObject(TypeName string, Prefix string, ItemName string, ObjectName string) string {
	sl := []string{}
	result := ""
	switch TypeName {
	case "group":
		sf := `	jso = doc.NewBlock("%v")
	bl = object.Block{}
	bl.ClickCBName = "%v_%v_click_CallBack"
	jso.ObjectExtender = bl
	jsoArray = append(jsoArray, jso)
	jsoDict["group_%v"] = jso

`
		_, ok := de.TypeVarList["bl"]
		if !ok {
			de.TypeVarList["bl"] = "var bl object.Block"
		}
		l2 := fmt.Sprintf(sf, ItemName, Prefix, ObjectName, ObjectName)
		sl = append(sl, l2)
	case "ordinary":
		sf := `	jso = doc.NewOrdinary("%v")
	or = object.Ordinary{}
	jso.ObjectExtender = or
	jsoArray = append(jsoArray, jso)
	jsoDict["ordinary_%v"] = jso

`
		_, ok := de.TypeVarList["or"]
		if !ok {
			de.TypeVarList["or"] = "var or object.Ordinary"
		}
		l2 := fmt.Sprintf(sf, ItemName, ObjectName)
		sl = append(sl, l2)
	case "label":
		sf := `	jso = doc.NewLabel("%v", "")
	lb = object.Label{}
	lb.ClickCBName = "%v_%v_click_CallBack"
	jso.ObjectExtender = lb
	jsoArray = append(jsoArray, jso)
	jsoDict["label_%v"] = jso

`
		_, ok := de.TypeVarList["lb"]
		if !ok {
			de.TypeVarList["lb"] = "var lb object.Label"
		}
		l2 := fmt.Sprintf(sf, ItemName, Prefix, ObjectName, ObjectName)
		sl = append(sl, l2)
	case "image":
		sf := `	jso = doc.NewImage("%v")
	im = object.Image{}
	im.ClickCBName = "%v_%v_click_CallBack"
	jso.ObjectExtender = im
	jsoArray = append(jsoArray, jso)
	jsoDict["image_%v"] = jso

`
		_, ok := de.TypeVarList["im"]
		if !ok {
			de.TypeVarList["im"] = "var im object.Image"
		}
		l2 := fmt.Sprintf(sf, ItemName, Prefix, ObjectName, ObjectName)
		sl = append(sl, l2)
	case "text":
		sf := `	jso = doc.NewText("%v", "")
	tx = object.Text{}
	tx.ChangeCBName = "%v_%v_change_CallBack"
	jso.ObjectExtender = tx
	jsoArray = append(jsoArray, jso)
	jsoDict["text_%v"] = jso
`
		_, ok := de.TypeVarList["tx"]
		if !ok {
			de.TypeVarList["tx"] = "var tx object.Text"
		}
		l2 := fmt.Sprintf(sf, ItemName, Prefix, ObjectName, ObjectName)
		sl = append(sl, l2)
	case "button":
		sf := `	jso = doc.NewButton("%v", "")
	bt = object.Button{}
	bt.ClickCBName = "%v_%v_click_CallBack"
	jso.ObjectExtender = bt
	jsoArray = append(jsoArray, jso)
	jsoDict["button_%v"] = jso
`
		_, ok := de.TypeVarList["bt"]
		if !ok {
			de.TypeVarList["bt"] = "var bt object.Button"
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
	jsoArray = append(jsoArray, jso)
	jsoDict["list_%v"] = jso
`
		_, ok := de.TypeVarList["ls"]
		if !ok {
			de.TypeVarList["ls"] = "var ls object.List"
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
	jsoArray = append(jsoArray, jso)
	jsoDict["selector_%v"] = jso
`
		_, ok := de.TypeVarList["se"]
		if !ok {
			de.TypeVarList["se"] = "var se object.Selector"
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
	jsoArray = append(jsoArray, jso)
	jsoDict["table_%v"] = jso
`
		_, ok := de.TypeVarList["tb"]
		if !ok {
			de.TypeVarList["tb"] = "var tb object.Table"
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
	de.JSOList = append(de.JSOList, sl...)
	return result
}

func main() {
	var fileInput string
	flag.StringVar(&fileInput, "file_input", "", "input file name with path")
	var pathOutput string
	flag.StringVar(&pathOutput, "path_output", "", "path to output file app_data.go")

	flag.Parse()

	if len(fileInput) == 0 {
		fmt.Printf("Must be path and name input file\r\n")
		flag.PrintDefaults()
		return
	}
	if len(pathOutput) == 0 {
		fmt.Printf("Must be path to output file\r\n")
		flag.PrintDefaults()
		return
	}

	stati, err := os.Stat(fileInput)
	if os.IsNotExist(err) {
		fmt.Printf("File does not exist.")
		return
	}
	if stati.IsDir() {
		fmt.Printf("Exist directory. No file.")
		return
	}

	statOutput, err := os.Stat(pathOutput)
	if os.IsNotExist(err) {
		fmt.Printf("Path does not exist.")
	}
	if !statOutput.IsDir() {
		fmt.Printf("Exist file. No directory.")
		return
	}

	bs, err := os.ReadFile(fileInput)
	if err != nil {
		log.Fatal(err)
	}

	text := string(bs)

	de := NewDataEnv()

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
	ordCont := 1
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
					className := a.Val
					cnl := strings.Split(className, "_")
					switch cnl[0] {
					case "desktop":
					case "rectangle":
					case "button":
						// fmt.Printf("%v button %v a.Key %v, class_name %v, cnl[1] %v\r\n", es, cnl[0], a.Key, class_name, cnl[1])
						if len(es) > 0 {
							escLocal, ok := esd[esc.Name]
							if ok {
								de.AddJSObject(cnl[0], "Button", className, cnl[1])
								de.AddFuncCallback("Button_"+cnl[1], "click")
								de.AddObjectList("button_"+cnl[1], "button_"+cnl[1])
								escLocal.Childs = append(escLocal.Childs, "button_"+cnl[1])
								esd[esc.Name] = escLocal
							} else {
								// error!!!
								log.Printf("error: len stack element == 0\r\n")
							}
						}
					case "label":
						// fmt.Printf("%v label %v a.Key %v, class_name %v, cnl[1] %v\r\n", es, cnl[0], a.Key, class_name, cnl[1])
						if len(es) > 0 {
							escLocal, ok := esd[esc.Name]
							if ok {
								de.AddJSObject(cnl[0], "Label", className, cnl[1])
								de.AddFuncCallback("Label_"+cnl[1], "click")
								de.AddObjectList("label_"+cnl[1], "label_"+cnl[1])
								escLocal.Childs = append(escLocal.Childs, "label_"+cnl[1])
								esd[esc.Name] = escLocal
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
						de.AddJSObject(cnl[0], "Group", className, cnl[1])
						de.AddFuncCallback("Group_"+cnl[1], "click")
						de.AddObjectList("group_"+cnl[1], "group_"+cnl[1])
						escc := ElemStack{cnl[0], cnl[1], []string{}, ""}
						es = append(es, escc)
						esd[escc.Name] = escc
					case "text":
						// fmt.Printf("%v text %v a.Key %v, class_name %v, cnl[1] %v\r\n", es, cnl[0], a.Key, class_name, cnl[1])
						if len(es) > 0 {
							escLocal, ok := esd[esc.Name]
							if ok {
								de.AddJSObject(cnl[0], "Text", className, cnl[1])
								de.AddFuncCallback("Text_"+cnl[1], "change")
								de.AddObjectList("text_"+cnl[1], "text_"+cnl[1])
								escLocal.Childs = append(escLocal.Childs, "text_"+cnl[1])
								esd[esc.Name] = escLocal
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
							escLocal, ok := esd[esc.Name]
							if ok {
								de.AddJSObject(cnl[0], "Image", className, cnl[1])
								de.AddFuncCallback("Image_"+cnl[1], "click")
								de.AddObjectList("image_"+cnl[1], "image_"+cnl[1])
								escLocal.Childs = append(escLocal.Childs, "image_"+cnl[1])
								esd[esc.Name] = escLocal
							} else {
								// error!!!
								log.Printf("error: len stack element == 0\r\n")
							}
						}
					default:
						// fmt.Printf("%v ordinary %v a.Key %v, class_name %v, cnl[1] %v\r\n", es, cnl[0], a.Key, class_name, cnl[1])
						if len(es) > 0 {
							escLocal, ok := esd[esc.Name]
							if ok {
								ss := fmt.Sprintf("%v_%v", cnl[1], ordCont)
								ordCont = ordCont + 1
								de.AddJSObject("ordinary", "", className, ss)
								//de.Add_FuncCallBack("Ordinary_" + ss, "click")
								de.AddObjectList("ordinary_"+cnl[1], "ordinary_"+cnl[1])
								escLocal.Childs = append(escLocal.Childs, "ordinary_"+cnl[1])
								esd[esc.Name] = escLocal
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
						className := a.Val
						cnl := strings.Split(className, "_")
						switch cnl[0] {
						case "example":
							// fmt.Printf("%v example %v a.Key %v, class_name %v, cnl[1] %v\r\n", es, cnl[0], a.Key, class_name, cnl[1])
							if len(es) > 0 {
								escLocal, ok := esd[esc.Name]
								if ok {
									ss := fmt.Sprintf("%v_%v", cnl[1], ordCont)
									ordCont = ordCont + 1
									de.AddJSObject("ordinary", "", className, ss)
									//de.Add_FuncCallBack("Ordinary_" + ss, "click")
									//esc_.Childs = append(esc_.Childs, cnl[1])
									// fmt.Printf("esc.Block %v class_name %v\r\n", esc.Block, class_name)
									if esc.Block == "text" || esc.Block == "label" {
										escLocal.Anchor = className
									}
									esd[esc.Name] = escLocal
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
		de.AddCorrectText(v.Block+"_"+v.Name, v.Anchor)
		for i := range v.Childs {
			de.AddChildsList(v.Block+"_"+v.Name, v.Childs[i], v.Anchor)
		}
	}

	de.CreateAppStruct()

	de.AddCreateConstructor()

	ss := de.CreateFile()

	//fmt.Printf("%v", ss)
	fo := filepath.Join(pathOutput, "app_data.go")
	_ = os.WriteFile(fo, []byte(ss), 0644)
}
