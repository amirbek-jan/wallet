package main

import "github.com/amirbek-jan/wallet/pkg/wallet"

func main(){
	svc := &wallet.Service{}

	svc.RegisterAccount("+992000000001")
}