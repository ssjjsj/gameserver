package timer

import(
	"gameserver/module"
	"time"
)

type TimerArg struct{
	Duration time.Duration
	ModuleName string
}


var timerMap map[string]*time.Timer


func init(){
	module.StartModule("timer", MessageHandler, OnInit, OnDestroy)
}


func OnInit(){
	//timerMap = make(map[string]*time.Timer)
}


func OnDestroy(){
}


func MessageHandler(data module.CallArg){
	if data.FunctionName == "AddTimer" {
		var arg TimerArg
		arg = data.Args.(TimerArg)
		timer := time.NewTimer(arg.Duration)
		//timerMap[arg.moduleName] = timer

		go onTimer(arg.ModuleName, timer)
	}
}


func onTimer(moduleName string, timer *time.Timer){
	for{
		_ : <- timer.C
		module.ModuleCall(moduleName, "onTimer", nil)
	}
}