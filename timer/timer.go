package timer

import(
	"gameserver/module"
	"time"
	"fmt"
)

type TimerArg struct{
	Duration time.Duration
	ModuleName string
}


var moduleList []string


func InitTimerModule(){
	module.RegistModule("timer", MessageHandler, OnInit, OnDestroy)
}


func OnInit(){
	moduleList = make([]string, 0)
	timer := time.Tick(time.Millisecond*100)

	go onTimer(timer)
}


func OnDestroy(){
}


func MessageHandler(data module.CallArg){
	fmt.Println("timer on message handle:"+data.FunctionName)
	if data.FunctionName == "AddTimer"{
		var timerArg = data.Args.(TimerArg)
		moduleList = append(moduleList, timerArg.ModuleName)
		fmt.Println(timerArg.ModuleName+" add timer")
	}
}


func onTimer(timer <-chan time.Time){
	for{
		<- timer
		//fmt.Println("timer on time")
		for i:=0; i<len(moduleList); i++{
			moduleName := moduleList[i]
			//fmt.Println("on timer:"+moduleName)
			module.ModuleCall(moduleName, "onTimer", nil)
		}
	}
}