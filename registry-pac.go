package registry
import(
	"fmt"
	"log"
	"syscall"
	"unsafe"
	"os"
	"os/signal"
	"golang.org/x/sys/windows/registry"
)

const (
	PATH = "Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings"
	KEYNAME = "AutoConfigURL"
)

//增加AutoConfigURL ，设置并启动自动脚本.退出时删除注册表
func EditReg(KEYVALUE string){
	key,_,err:=registry.CreateKey(registry.CURRENT_USER,PATH,registry.ALL_ACCESS)

	if err!=nil {
		log.Fatal(err)
	}
	//写入string地址
	defer key.Close()
	err = key.SetStringsValue(KEYNAME,[]string{KEYVALUE})
	if err !=nil {
		log.Fatal(err)
	}
	//刷新代理
	refreshReg()
	//捕获程序退出信号
	msgChan:=make(chan os.Signal,1)
	signal.Notify(msgChan,syscall.SIGINT,syscall.SIGHUP,syscall.SIGKILL,syscall.SIGQUIT)
	<-msgChan
	//关闭的时候删除注册表
	func(){
		fmt.Println("程序结束")
		err = key.DeleteValue(KEYNAME)
	}()
}


//刷新注册表文件，使脚本立即生效
func refreshReg()  {
	//TODO:错误处理
	//参数准备
	firstP := int64(0)
	internetOptionProxy := int64(39)
	internetOptionProxy1 := int64(37)
	internetInfo := int64(0)
	sizeOfInfo := int64(0)
	wininet,_ := syscall.LoadLibrary("wininet.dll")
	defer syscall.FreeLibrary(wininet)
	InternetSetOption,_ :=syscall.GetProcAddress(syscall.Handle(wininet),"InternetSetOptionA")
	r,_,_ :=syscall.Syscall6(uintptr(InternetSetOption),4,
		uintptr(firstP),
		uintptr(internetOptionProxy),
		uintptr(unsafe.Pointer(&internetInfo)),
		uintptr(sizeOfInfo),0,0)
	if r !=0 {
		fmt.Println("成功")
	}
	r1,_,_ :=syscall.Syscall6(uintptr(InternetSetOption),4,
		uintptr(firstP),
		uintptr(internetOptionProxy1),
		uintptr(unsafe.Pointer(&internetInfo)),
		uintptr(sizeOfInfo),0,0)
	if r1 !=0 {
		fmt.Println("成功")
	}

}

