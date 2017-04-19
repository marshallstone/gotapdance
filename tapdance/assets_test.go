package tapdance

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path"
	"reflect"
	"testing"
)

func TestAssets_Decoys(t *testing.T) {
	oldpath := Assets().path
	Assets().saveClientConf()
	dir1, err := ioutil.TempDir("/tmp/", "decoy1")
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
	dir2, err := ioutil.TempDir("/tmp/", "decoy2")
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}


	var testDecoys1 = []*TLSDecoySpec{
		initTLSDecoySpec("4.8.15.16", "ericw.us"),
		initTLSDecoySpec("19.21.23.42", "sergeyfrolov.github.io"),
	}

	var testDecoys2 = []*TLSDecoySpec{
		initTLSDecoySpec("0.1.2.3", "whatever.cn"),
		initTLSDecoySpec("255.254.253.252", "particular.ir"),
		initTLSDecoySpec("11.22.33.44", "what.is.up"),
		initTLSDecoySpec("8.255.255.8", "heh.meh"),
	}

	Assets().SetAssetsDir(dir1)
	err = Assets().SetDecoys(testDecoys1)
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
	Assets().SetAssetsDir(dir2)
	err = Assets().SetDecoys(testDecoys2)
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
	if !reflect.DeepEqual(Assets().config.DecoyList.TlsDecoys, testDecoys2) {
		fmt.Println("Assets are not equal!")
		fmt.Println("Assets().decoys:", Assets().config.DecoyList.TlsDecoys)
		fmt.Println("testDecoys2:", testDecoys2)
		t.Fail()
	}

	decoyInList := func(d *TLSDecoySpec, decoyList []*TLSDecoySpec) bool {
		for _, elem := range decoyList {
			if reflect.DeepEqual(elem, d) {
				return true
			}
		}
		return false
	}

	for i := 0; i < 10; i++ {
		_sni, addr := Assets().GetDecoyAddress()
		host_addr, _, err := net.SplitHostPort(addr)
		if err != nil {
			fmt.Println("Corrupted addr:", addr, ". Error:", err.Error())
			t.Fail()
		}
		decoyServ := initTLSDecoySpec(host_addr, _sni)
		if !decoyInList(decoyServ, Assets().config.DecoyList.TlsDecoys) {
			fmt.Println("decoyServ not in List!")
			fmt.Println("decoyServ:", decoyServ)
			fmt.Println("Assets().decoys:", Assets().config.DecoyList.TlsDecoys)
			t.Fail()
		}
	}
	Assets().SetAssetsDir(dir1)
	if !reflect.DeepEqual(Assets().config.DecoyList.TlsDecoys, testDecoys1) {
		fmt.Println("Assets are not equal!")
		fmt.Println("Assets().decoys:", Assets().config.DecoyList.TlsDecoys)
		fmt.Println("testDecoys1:", testDecoys1)
		t.Fail()
	}
	for i := 0; i < 10; i++ {
		_sni, addr := Assets().GetDecoyAddress()
		host_addr, _, err := net.SplitHostPort(addr)
		if err != nil {
			fmt.Println("Corrupted addr:", addr, ". Error:", err.Error())
			t.Fail()
		}
		decoyServ := initTLSDecoySpec(host_addr, _sni)
		if !decoyInList(decoyServ, Assets().config.DecoyList.TlsDecoys) {
			fmt.Println("decoyServ not in List!")
			fmt.Println("decoyServ:", decoyServ)
			fmt.Println("Assets().decoys:", Assets().config.DecoyList.TlsDecoys)
			t.Fail()
		}
	}
	os.Remove(path.Join(dir1, Assets().filenameClientConf))
	os.Remove(path.Join(dir2, Assets().filenameClientConf))
	os.Remove(dir1)
	os.Remove(dir2)
	Assets().SetAssetsDir(oldpath)
	fmt.Println("TestAssets_Decoys OK")
}

func TestAssets_Pubkey(t *testing.T) {
	initPubKey := func(defaultKey []byte) PubKey {
		defualtKeyType := KeyType_AES_GCM_128
		return PubKey{Key: defaultKey, Type: &defualtKeyType}
	}

	oldpath := Assets().path
	Assets().saveClientConf()
	dir1, err := ioutil.TempDir("/tmp/", "pubkey1")
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
	dir2, err := ioutil.TempDir("/tmp/", "pubkey2")
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}

	var pubkey1 = initPubKey([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
		12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
		27, 28, 29, 30, 31})
	var pubkey2 = initPubKey([]byte{200, 201, 202, 203, 204, 205, 206, 207, 208, 209, 210, 211,
		212, 213, 214, 215, 216, 217, 218, 219, 220, 221, 222, 223, 224, 225, 226,
		227, 228, 229, 230, 231})

	Assets().SetAssetsDir(dir1)
	err = Assets().SetPubkey(pubkey1)
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
	Assets().SetAssetsDir(dir2)
	err = Assets().SetPubkey(pubkey2)
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
	if !bytes.Equal(Assets().config.DefaultPubkey.Key[:], pubkey2.Key[:]) {
		fmt.Println("Pubkeys are not equal!")
		fmt.Println("Assets().stationPubkey:", Assets().config.DefaultPubkey.Key[:])
		fmt.Println("pubkey2:", pubkey2)
		t.Fail()
	}

	Assets().SetAssetsDir(dir1)
	if !bytes.Equal(Assets().config.DefaultPubkey.Key[:], pubkey1.Key[:]) {
		fmt.Println("Pubkeys are not equal!")
		fmt.Println("Assets().stationPubkey:", Assets().config.DefaultPubkey.Key[:])
		fmt.Println("pubkey1:", pubkey1)
		t.Fail()
	}
	os.Remove(path.Join(dir1, Assets().filenameStationPubkey))
	os.Remove(path.Join(dir2, Assets().filenameStationPubkey))
	os.Remove(dir1)
	os.Remove(dir2)
	Assets().SetAssetsDir(oldpath)
	fmt.Println("TestAssets_Pubkey OK")
}