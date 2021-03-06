package connman

import (
	"fmt"
	"log"
	"os"

	"github.com/godbus/dbus"
)

type Agent struct {
	Name       string
	Path       dbus.ObjectPath
	Interface  string
	Passphrase string
}

func NewAgent(psk string) *Agent {
	agent := &Agent{
		Name:       "com.develboard.webadmin",
		Path:       "/com/develboard/webadmin/Agent",
		Interface:  "net.connman.Agent",
		Passphrase: psk,
	}

	conn, err := dbus.SystemBus()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	reply, err := conn.RequestName(agent.Name, dbus.NameFlagDoNotQueue)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if reply != dbus.RequestNameReplyPrimaryOwner {
		fmt.Fprintln(os.Stderr, "Name already taken")
		return nil
	}

	conn.Export(agent, agent.Path, agent.Interface)
	return agent
}

func (a *Agent) Destroy() error {
	conn, err := dbus.SystemBus()
	if err != nil {
		return err
	}

	reply, err := conn.ReleaseName(a.Name)
	if err != nil {
		return err
	}
	if reply != dbus.ReleaseNameReplyReleased {
		return fmt.Errorf("Could not release the name\n")
	}

	conn.Export(nil, a.Path, a.Interface)
	return nil
}

func (a *Agent) RequestInput(service dbus.ObjectPath, rq map[string]dbus.Variant) (map[string]dbus.Variant, *dbus.Error) {
	return map[string]dbus.Variant{
		"Passphrase": dbus.MakeVariant(a.Passphrase),
	}, nil
}

func (a *Agent) ReportError(service dbus.ObjectPath, err string) *dbus.Error {
	log.Printf("%s: %s\n", service, err)
	return nil
}
