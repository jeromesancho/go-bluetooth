package gatt



import (
   "sync"
   "github.com/jeromesancho/go-bluetooth/bluez"
   "github.com/jeromesancho/go-bluetooth/util"
   "github.com/jeromesancho/go-bluetooth/props"
   "github.com/godbus/dbus"
)

var GattProfile1Interface = "org.bluez.GattProfile1"


// NewGattProfile1 create a new instance of GattProfile1
//
// Args:
// - servicePath: <application dependent>
// - objectPath: <application dependent>
func NewGattProfile1(servicePath string, objectPath dbus.ObjectPath) (*GattProfile1, error) {
	a := new(GattProfile1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  servicePath,
			Iface: GattProfile1Interface,
			Path:  dbus.ObjectPath(objectPath),
			Bus:   bluez.SystemBus,
		},
	)
	
	a.Properties = new(GattProfile1Properties)

	_, err := a.GetProperties()
	if err != nil {
		return nil, err
	}
	
	return a, nil
}


/*
GattProfile1 GATT Profile hierarchy

Local profile (GATT client) instance. By registering this type of object
an application effectively indicates support for a specific GATT profile
and requests automatic connections to be established to devices
supporting it.

*/
type GattProfile1 struct {
	client     				*bluez.Client
	propertiesSignal 	chan *dbus.Signal
	objectManagerSignal chan *dbus.Signal
	objectManager       *bluez.ObjectManager
	Properties 				*GattProfile1Properties
	watchPropertiesChannel chan *dbus.Signal
}

// GattProfile1Properties contains the exposed properties of an interface
type GattProfile1Properties struct {
	lock sync.RWMutex `dbus:"ignore"`

	/*
	UUIDs 128-bit GATT service UUIDs to auto connect.
	*/
	UUIDs []string

}

//Lock access to properties
func (p *GattProfile1Properties) Lock() {
	p.lock.Lock()
}

//Unlock access to properties
func (p *GattProfile1Properties) Unlock() {
	p.lock.Unlock()
}




// SetUUIDs set UUIDs value
func (a *GattProfile1) SetUUIDs(v []string) error {
	return a.SetProperty("UUIDs", v)
}



// GetUUIDs get UUIDs value
func (a *GattProfile1) GetUUIDs() ([]string, error) {
	v, err := a.GetProperty("UUIDs")
	if err != nil {
		return []string{}, err
	}
	return v.Value().([]string), nil
}



// Close the connection
func (a *GattProfile1) Close() {
	
	a.unregisterPropertiesSignal()
	
	a.client.Disconnect()
}

// Path return GattProfile1 object path
func (a *GattProfile1) Path() dbus.ObjectPath {
	return a.client.Config.Path
}

// Client return GattProfile1 dbus client
func (a *GattProfile1) Client() *bluez.Client {
	return a.client
}

// Interface return GattProfile1 interface
func (a *GattProfile1) Interface() string {
	return a.client.Config.Iface
}

// GetObjectManagerSignal return a channel for receiving updates from the ObjectManager
func (a *GattProfile1) GetObjectManagerSignal() (chan *dbus.Signal, func(), error) {

	if a.objectManagerSignal == nil {
		if a.objectManager == nil {
			om, err := bluez.GetObjectManager()
			if err != nil {
				return nil, nil, err
			}
			a.objectManager = om
		}

		s, err := a.objectManager.Register()
		if err != nil {
			return nil, nil, err
		}
		a.objectManagerSignal = s
	}

	cancel := func() {
		if a.objectManagerSignal == nil {
			return
		}
		a.objectManagerSignal <- nil
		a.objectManager.Unregister(a.objectManagerSignal)
		a.objectManagerSignal = nil
	}

	return a.objectManagerSignal, cancel, nil
}


// ToMap convert a GattProfile1Properties to map
func (a *GattProfile1Properties) ToMap() (map[string]interface{}, error) {
	return props.ToMap(a), nil
}

// FromMap convert a map to an GattProfile1Properties
func (a *GattProfile1Properties) FromMap(props map[string]interface{}) (*GattProfile1Properties, error) {
	props1 := map[string]dbus.Variant{}
	for k, val := range props {
		props1[k] = dbus.MakeVariant(val)
	}
	return a.FromDBusMap(props1)
}

// FromDBusMap convert a map to an GattProfile1Properties
func (a *GattProfile1Properties) FromDBusMap(props map[string]dbus.Variant) (*GattProfile1Properties, error) {
	s := new(GattProfile1Properties)
	err := util.MapToStruct(s, props)
	return s, err
}

// ToProps return the properties interface
func (a *GattProfile1) ToProps() bluez.Properties {
	return a.Properties
}

// GetWatchPropertiesChannel return the dbus channel to receive properties interface
func (a *GattProfile1) GetWatchPropertiesChannel() chan *dbus.Signal {
	return a.watchPropertiesChannel
}

// SetWatchPropertiesChannel set the dbus channel to receive properties interface
func (a *GattProfile1) SetWatchPropertiesChannel(c chan *dbus.Signal) {
	a.watchPropertiesChannel = c
}

// GetProperties load all available properties
func (a *GattProfile1) GetProperties() (*GattProfile1Properties, error) {
	a.Properties.Lock()
	err := a.client.GetProperties(a.Properties)
	a.Properties.Unlock()
	return a.Properties, err
}

// SetProperty set a property
func (a *GattProfile1) SetProperty(name string, value interface{}) error {
	return a.client.SetProperty(name, value)
}

// GetProperty get a property
func (a *GattProfile1) GetProperty(name string) (dbus.Variant, error) {
	return a.client.GetProperty(name)
}

// GetPropertiesSignal return a channel for receiving udpdates on property changes
func (a *GattProfile1) GetPropertiesSignal() (chan *dbus.Signal, error) {

	if a.propertiesSignal == nil {
		s, err := a.client.Register(a.client.Config.Path, bluez.PropertiesInterface)
		if err != nil {
			return nil, err
		}
		a.propertiesSignal = s
	}

	return a.propertiesSignal, nil
}

// Unregister for changes signalling
func (a *GattProfile1) unregisterPropertiesSignal() {
	if a.propertiesSignal != nil {
		a.propertiesSignal <- nil
		a.propertiesSignal = nil
	}
}

// WatchProperties updates on property changes
func (a *GattProfile1) WatchProperties() (chan *bluez.PropertyChanged, error) {
	return bluez.WatchProperties(a)
}

func (a *GattProfile1) UnwatchProperties(ch chan *bluez.PropertyChanged) error {
	return bluez.UnwatchProperties(a, ch)
}




/*
Release 
			This method gets called when the service daemon
			unregisters the profile. The profile can use it to
			do cleanup tasks. There is no need to unregister the
			profile, because when this method gets called it has
			already been unregistered.


*/
func (a *GattProfile1) Release() error {
	
	return a.client.Call("Release", 0, ).Store()
	
}

