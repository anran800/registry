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

//AutoConfigURL ，Delete the registry while exiting
func EditReg(KEYVALUE string){
	key,_,err:=registry.CreateKey(registry.CURRENT_USER,PATH,registry.ALL_ACCESS)

	if err!=nil {
		log.Fatal(err)
	}
	
	defer key.Close()
	err = key.SetStringsValue(KEYNAME,[]string{KEYVALUE})
	if err !=nil {
		log.Fatal(err)
	}
	//refresh
	refreshReg()
	msgChan:=make(chan os.Signal,1)
	signal.Notify(msgChan,syscall.SIGINT,syscall.SIGHUP,syscall.SIGKILL,syscall.SIGQUIT)
	<-msgChan
	func(){
		fmt.Println("Application end")
		err = key.DeleteValue(KEYNAME)
	}()
}


//Refresh
func refreshReg()  {

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
		fmt.Println("Done")
	}
	r1,_,_ :=syscall.Syscall6(uintptr(InternetSetOption),4,
		uintptr(firstP),
		uintptr(internetOptionProxy1),
		uintptr(unsafe.Pointer(&internetInfo)),
		uintptr(sizeOfInfo),0,0)
	if r1 !=0 {
		fmt.Println("Done")
	}

}

