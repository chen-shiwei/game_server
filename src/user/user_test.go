package user

import (
	"fmt"
	"testing"
)

func TestGetUserLass(T *testing.T) {
	u, err := GetUserLass("b245648bf655f3adfd7a92590fc0b9f7")
	if err != nil {
		T.Error(err.Error())
		return
	}
	fmt.Print(u)
}
