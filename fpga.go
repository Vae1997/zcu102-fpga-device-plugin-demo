package main

import (
	"fmt"
	"io/ioutil"
	pluginapi "k8s.io/kubernetes/pkg/kubelet/apis/deviceplugin/v1beta1"
	//	"os"
	"path"
	"strconv"
	"strings"
)

const (
	SysfsDevices    = "/sys/bus/pci/devices"
	EmbeddedDevices = "/sys/bus/platform/devices" //renderD128
	MgmtPrefix      = "/dev/xclmgmt"
	UserPrefix      = "/dev/dri" //renderD128
	QdmaPrefix      = "/dev/xfpga"
	QDMASTR         = "dma.qdma.u"
	UserPFKeyword   = "drm"     //renderD128
	DRMSTR          = "renderD" //renderD128
	ROMSTR          = "rom"
	DSAverFile      = "VBNV"
	DSAtsFile       = "timestamp"
	InstanceFile    = "instance"
	MgmtFile        = "mgmt_pf"
	UserFile        = "user_pf"
	VendorFile      = "vendor"
	DeviceFile      = "device"
	DevFile         = "dev"    //renderD128--226:128
	UeventFile      = "uevent" //renderD128--MAJOR=226 MINOR=128
	//	    --DEVNAME=dri/renderD128
	//	    --DEVTYPE=drm_minor
	//	UeventFile     = "uevent"			///sys/bus/platform/devices/XXX/uevent
	//OF_NAME=zyxclmm_drm
	ModaliasFile = "modalias"
	MODALIAS     = "of:Nzyxclmm_drmT<NULL>Cxlnx,zocl"

	XilinxVendorID = "0x10ee"
	ADVANTECH_ID   = "0x13fe"
	AWS_ID         = "0x1d0f"
)

type Pairs struct {
	Mgmt string
	User string //renderD128
	Qdma string
}

type Device struct {
	index     string
	shellVer  string
	timestamp string
	DBDF      string // this is for user pf
	deviceID  string //devid of the user pf
	Healthy   string
	Nodes     *Pairs //renderD128
}

func GetFileNameFromPrefix(dir string, prefix string) (string, error) {
	userFiles, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("Can't read folder %s", dir)
	}
	for _, userFile := range userFiles {
		fname := userFile.Name()

		if !strings.HasPrefix(fname, prefix) {
			continue
		}
		return fname, nil
	}
	return "", nil
}

func GetFileContent(file string) (string, error) {
	if buf, err := ioutil.ReadFile(file); err != nil {
		return "", fmt.Errorf("Can't read file %s", file)
	} else {
		ret := strings.Trim(string(buf), "\n")
		return ret, nil
	}
}

func GetDevices() ([]Device, error) {
	var devices []Device
	pairMap := make(map[string]*Pairs)
	//pciFiles, err := ioutil.ReadDir(SysfsDevices)
	//if err != nil {
	//	return nil, fmt.Errorf("Can't read folder %s", SysfsDevices)
	//}
	//fmt.Printf("The pciFiles is %v \n",pciFiles)
	//获取platform/devices目录下的文件
	platformFiles, err := ioutil.ReadDir(EmbeddedDevices)
	if err != nil {
		return nil, fmt.Errorf("Can't read folder %s", EmbeddedDevices)
	}
	//fmt.Printf("The platformFiles is %v \n",platformFiles)
	//fmt.Printf("The platformFiles's count is %v \n",len(platformFiles))
	//仿照for循环
	for _, platformFile := range platformFiles {
		//当前文件
		fileName := platformFile.Name()
		//当前文件下的modalias文件
		mFile := path.Join(EmbeddedDevices, fileName, ModaliasFile)
		mFileContent, err := GetFileContent(mFile)
		if err != nil {
			return nil, err
		}
		//定位到zocl所在文件夹
		if strings.EqualFold(mFileContent, MODALIAS) != true {
			continue
		}
		//fmt.Printf("The fileName is %v\n",fileName)
		//fmt.Printf("The mFileContent is %v\n",mFileContent)
		//截取a0000000
		DBD := fileName[:8]
		//fmt.Printf("The DBD is %v\n",DBD)
		//初始化
		if _, ok := pairMap[DBD]; !ok {
			pairMap[DBD] = &Pairs{
				Mgmt: "",
				User: "",
				Qdma: "",
			}
		}
		//文件名设为device的DBDF字段
		userDBDF := fileName
		//设置device的Nodes字段，先定位到drm文件夹，再通过renderD前缀定位
		userpf, err := GetFileNameFromPrefix(path.Join(EmbeddedDevices, fileName, UserPFKeyword), DRMSTR)
		if err != nil {
			return nil, err
		}
		//定位/dev/dri目录，userpf为renderD128
		userNode := path.Join(UserPrefix, userpf)
		//仅将User赋值，其他为""
		pairMap[DBD].User = userNode
		//for k,v := range pairMap{
		//	fmt.Printf("key %v is:%v\n",k,*v)
		//}
		//设置device的ID,通过renderD128下的dev文件
		dev := path.Join(EmbeddedDevices, fileName, UserPFKeyword, userpf, DevFile)
		content, err := GetFileContent(dev)
		if err != nil {
			return nil, err
		}
		//将226:128作为deviceID
		devid := content
		//通过renderD128下的uevent设置shellVer和timestamp
		uevent := path.Join(EmbeddedDevices, fileName, UserPFKeyword, userpf, UeventFile)
		content, err = GetFileContent(uevent)
		if err != nil {
			return nil, err
		}
		//fmt.Printf("content is:\n%v\n",content)
		//对content处理，按换行分割，得到dsaVer
		contents := strings.Split(content, "\n")
		//DEVTYPE=XXX
		dsaVer := strings.Split(contents[3], "=")[1]
		//将Unix时间戳作为dsaTs(会反复创建sock，不可取)
		//dsaTs := strconv.FormatInt(time.Now().Unix(), 10)
		dsaTs := "20191217"
		healthy := pluginapi.Healthy
		devices = append(devices, Device{
			index:     strconv.Itoa(len(devices) + 1),
			shellVer:  dsaVer,
			timestamp: dsaTs,
			DBDF:      userDBDF,
			deviceID:  devid,
			Healthy:   healthy,
			Nodes:     pairMap[DBD],
		})
	}
	return devices, nil
}

/*
func main() {
	devices, err := GetDevices()
	if err != nil {
		fmt.Printf("%s !!!\n", err)
		return
	}
 //fmt.Printf("The devices is %v\n", devices)
	for _, device := range devices {
		fmt.Printf("%v\n", device)
	}
}
*/
