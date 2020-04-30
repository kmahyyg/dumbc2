package remoteop

import (
	"fmt"
	"os"
)

func DeleteMyself() {
	_ = os.Remove(os.Args[0])
	fmt.Println("Destroy Called done!")
}
