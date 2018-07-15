package protojsonUtils
import(
	//"os"
	//"io/ioutil"
	//"github.com/bitly/go-simplejson"
	"encoding/json"
)

// var jsonsMap map[int][]byte
// var jsonIndexNameMap map[string]int


func JsonToData(jsonData []byte, output interface{}){
	json.Unmarshal(jsonData, output)
}


func DataToJson(data interface{})([]byte){
	bytes, err := json.Marshal(data)
	if err == nil {
		return bytes
	}

	return nil
}


// func ProtoNameToProtoId(protoName string){

// }


// func initJson(){
// 	fi, err := os.Open("./json/proto.json")
//     if err != nil {
//         panic(err)
//     }
//     defer fi.Close()

// 	fd, err := ioutil.ReadAll(fi)
// 	if err == nil {
// 		js, err := simplejson.NewJson(fd)
// 		if err == nil {
// 			arrayJson := js.GetPath("protos")
// 			eleIndex := 0
// 			for {
// 				element := arrayJson.Get("array").GetIndex(eleIndex)
// 				bytes, err := element.Encode()
// 				if err == nil {
// 					jsonsMap[eleIndex] = bytes
// 					protoName, err := jsonIndexNameMap[element.Get("protoname").String()]
// 					if err == nil {
// 						jsonIndexNameMap[protoName] = eleIndex
// 					}
// 				}
// 			}
// 		}
// 	}
// }


// func init(){
// 	jsonIndexNameMap = make(map[string]int)
// 	jsonsMap = make(map[int][]byte)
// 	initJson()
// }