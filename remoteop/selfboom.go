package remoteop

import "os"

func DeleteMyself(){
	_ = os.Remove(os.Args[0])
}
