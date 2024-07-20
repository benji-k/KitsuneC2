//This package uses the observer design pattern to let packages cross communicate about their state changes. Each struct represents the type
//of notification another package can choose to observe.

package notifications

type NotificationType string

const (
	FAIL    = "FAIL"
	INFO    = "INFO"
	SUCCESS = "SUCCESS"
)

type Notification struct {
	Message string
	NType   NotificationType
}

//------------------Begin notification types-----------------------

// Gets dispatched when a new implant registered
type implantRegister struct {
	callbackfuncs []func(Notification)
}

var ImplantRegisterNotification implantRegister = implantRegister{}

func (i *implantRegister) Subscribe(callbackFunc func(Notification)) {
	i.callbackfuncs = append(i.callbackfuncs, callbackFunc)
}

func (i *implantRegister) Dispatch(n Notification) {
	for _, cbFunc := range i.callbackfuncs {
		cbFunc(n)
	}
}
