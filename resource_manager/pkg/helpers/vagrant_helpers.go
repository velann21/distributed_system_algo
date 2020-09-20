package helper

import (
	"github.com/bmatcuk/go-vagrant"
)


type Vagrant struct {
	client *vagrant.VagrantClient

}

func (v *Vagrant) NewVagrantClient(vagFileDir string)(*vagrant.VagrantClient, error){
	client, err := vagrant.NewVagrantClient(vagFileDir)
	if err != nil {
		return nil, err
	}
	v.client = client
	return client, nil
}

func (v *Vagrant) CreateVms(verbose bool)(*vagrant.UpCommand, error){
	upcmd := v.client.Up()
	upcmd.Verbose = verbose
	if err := upcmd.Start(); err != nil {
		return nil, err
	}
	return upcmd, nil
}

func (v *Vagrant) GetVagrantVersion()(*vagrant.VersionCommand, error){
	vercmd := v.client.Version()
	if err := vercmd.Run(); err != nil {
		return nil, err
	}
	return vercmd, nil
}

func (v *Vagrant) WaitTaskToComplete(upcmd *vagrant.UpCommand)error{
	if err := upcmd.Wait(); err != nil {
		return err
	}
	return nil
}

func (v *Vagrant) GetVMStatus()(*vagrant.StatusCommand, error){
	statuscmd := v.client.Status()
	if err := statuscmd.Run(); err != nil {
		return nil, err
	}
	return statuscmd, nil
}

func (v *Vagrant) Destroy()error{
	if err := v.client.Destroy().Run(); err != nil {
		return err
	}
	return nil
}


//func Example(vagFileDir string) {
//	client, err := vagrant.NewVagrantClient(vagFileDir)
//	if err != nil {
//		fmt.Println("Got error while creating client:", err)
//		os.Exit(-1)
//	}
//
//	// Let's start bringing up the vm
//	upcmd := client.Up()
//	upcmd.Verbose = true
//	fmt.Println("Bringing up the vm")
//	if err := upcmd.Start(); err != nil {
//		fmt.Println("Error bringing up vm:", err)
//		os.Exit(-1)
//	}
//
//	// while we're waiting, let's get version info
//	vercmd := client.Version()
//	if err := vercmd.Run(); err != nil {
//		fmt.Println("Error retrieving version info:", err)
//	}
//
//	// now wait for up to finish
//	if err := upcmd.Wait(); err != nil {
//		fmt.Println("Error waiting for up:", err)
//		os.Exit(-1)
//	}
//
//	fmt.Println("\n\nInstalled Vagrant version:", vercmd.InstalledVersion)
//
//	// Get vm status
//	statuscmd := client.Status()
//	if err := statuscmd.Run(); err != nil {
//		fmt.Println("Error getting status:", err)
//	} else {
//		for vm, status := range statuscmd.Status {
//			fmt.Printf("%v: %v\n", vm, status)
//		}
//	}
//
//	// Destroy vm
//	if err := client.Destroy().Run(); err != nil {
//		fmt.Println("Error destroying vm:", err)
//		os.Exit(-1)
//	}
//}
//
//func main() {
//	Example("./conf")
//}