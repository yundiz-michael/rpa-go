package common

/*
   #cgo LDFLAGS: -L/workspace/xpa/go/merkaba/libs/ -lpthread -lvncserver
   typedef void (*MouseCallbackFn)(char *taskName,int button,int x,int y);
   typedef void (*KeyboardCallbackFn)(char *taskName,int down,int key);
   typedef void (*ConnectionCallbackFn)(char* taskName,int clientCount);
   typedef struct {
       MouseCallbackFn mouse;
       KeyboardCallbackFn keyboard;
       ConnectionCallbackFn connection;
   } Callbacks;

   typedef unsigned long uintptr_t;

   uintptr_t CreateVNC(int width, int height, int port, char *title, char *webRoot, Callbacks callbacks);
   void ProcessVNC(uintptr_t screen);
   void StopVNC(uintptr_t screen);
   void RefreshVNC(uintptr_t screen, unsigned char *buffer);
   extern void keyboardCallbackEvent(char* taskName,int button,int x,int y);
   extern void mouseCallbackEvent(char* taskName,int down,int key);
   extern void connectionCallbackEvent(char* taskName,int clientCount);

*/
import "C"
import (
	"github.com/nfnt/resize"
	"go.uber.org/zap"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"sync/atomic"
	"unsafe"
)

var VncPortFrom uint64 = 4320
var VncWidth int = 900
var VncHeight int = 600
var VncBPP int = 4

type VncInstance struct {
	TaskName string
	Port     int
	State    chan string
	Event    uintptr
	Screen   uintptr
}

var VncInstances = make(map[string]*VncInstance)

//export onMouseCallback
func onMouseCallback(taskName *C.char, button int, x int, y int) {
	//name := C.GoString(taskName)
	//fmt.Printf("todo 转发鼠标事件给chrome或其他应用,taskName=%s,button=%d x=%d\n", name, button, x)
}

//export onKeyboardCallback
func onKeyboardCallback(taskName *C.char, down int, key int) {
	//fmt.Printf("todo 转发键盘事件给chrome或其他应用, key=%d\n", key)
}

//export onConnectionCallback
func onConnectionCallback(taskName *C.char, clientCount int) {
	if clientCount != 0 {
		return
	}
	if !Env.Vnc.AutoClose {
		return
	}
	name := C.GoString(taskName)
	StopVNC(name, "NoConnection")
}

func HasVNC(taskName string) bool {
	if _, ok := VncInstances[taskName]; ok {
		return true
	} else {
		return false
	}
}

func StartVNC(width int, height int, taskName string) *VncInstance {
	if i, ok := VncInstances[taskName]; ok {
		return i
	}
	atomic.AddUint64(&VncPortFrom, 1)
	s := &VncInstance{
		TaskName: taskName,
		Port:     (int)(VncPortFrom),
		State:    make(chan string, 1),
	}
	vncRoot := RootPath + "web/"
	VncInstances[taskName] = s
	/*处理回调*/
	cCallbacks := C.Callbacks{}
	cCallbacks.mouse = C.MouseCallbackFn(C.mouseCallbackEvent)
	cCallbacks.keyboard = C.KeyboardCallbackFn(C.keyboardCallbackEvent)
	cCallbacks.connection = C.ConnectionCallbackFn(C.connectionCallbackEvent)
	screen := C.CreateVNC(C.int(width), C.int(height), C.int(s.Port),
		C.CString(taskName),
		C.CString(vncRoot),
		cCallbacks)
	s.Screen = uintptr(screen)
	go func() {
		for {
			select {
			case state := <-s.State:
				if state == "WebClose" || state == "NoConnection" || state == "ScriptStop" {
					C.StopVNC((C.ulong)(s.Screen))
					delete(VncInstances, taskName)
					LoggerStd.Info("VNC stop ", zap.String("taskName", taskName), zap.String("state", state))
					return
				}
			default:
				C.ProcessVNC((C.ulong)(s.Screen))
			}
		}
	}()
	return s
}

func RefreshVNC(taskName string, img image.Image) {
	w := VncWidth
	h := VncHeight
	newImage := resize.Resize(uint(w), 0, img, resize.Lanczos3)
	if instance, ok := VncInstances[taskName]; ok {
		buffer := make([]uint8, w*h*VncBPP)
		var i, j int
		for j = 0; j < h; j++ {
			for i = 0; i < w; i++ {
				r, g, b, _ := color.NRGBAModel.Convert(newImage.At(i, j)).RGBA()
				buffer[(j*w+i)*VncBPP+0] = uint8(r) /* red */
				buffer[(j*w+i)*VncBPP+1] = uint8(g) /* green */
				buffer[(j*w+i)*VncBPP+2] = uint8(b) /* blue */
			}
		}
		C.RefreshVNC((C.ulong)(instance.Screen), (*C.uchar)(unsafe.Pointer(&buffer[0])))
	}
}

func StopVNC(taskName string, state string) {
	if instance, ok := VncInstances[taskName]; ok {
		instance.State <- state
	}
}

func ForTestVNCInstance() {
	width := 800
	height := 600
	StartVNC(width, height, "demo")
	file, _ := os.Open(Env.Path.Config + "demo.jpg")
	img, _ := jpeg.Decode(file)
	newImage := resize.Resize(uint(width), 0, img, resize.Lanczos3)
	file.Close()
	RefreshVNC("demo", newImage)
}
