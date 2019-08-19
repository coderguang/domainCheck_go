package domainCheckGenTxt

import (
	"os"

	"github.com/coderguang/GameEngine_go/sgthread"

	"github.com/coderguang/GameEngine_go/sgfile"

	"github.com/coderguang/GameEngine_go/sgalgorithm"
	"github.com/coderguang/GameEngine_go/sglog"
)

func init() {
	numlist = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	charlist = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}

}

var numlist []string
var charlist []string

func CreateDominFile(fileName string) {
	path, err := sgfile.GetPath(fileName)
	if err != nil {
		sglog.Error("get path error,err=%s", err)
		sgthread.DelayExit(2)
	}
	sgfile.AutoMkDir(path)
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0766)
	if err != nil {
		sglog.Error("open file error,err=%s", err)
		sgthread.DelayExit(2)
	}
	allList := []string{}
	allList = append(allList, numlist...)
	allList = append(allList, charlist...)
	//create all
	sum := 0
	tmpSum := createDomainFileAndWrite(allList, 1, file)
	sum += tmpSum
	tmpSum = createDomainFileAndWrite(allList, 2, file)
	sum += tmpSum
	tmpSum = createDomainFileAndWrite(allList, 3, file)
	sum += tmpSum
	//create only num
	tmpSum = createDomainFileAndWrite(numlist, 4, file)
	sum += tmpSum
	tmpSum = createDomainFileAndWrite(numlist, 5, file)
	sum += tmpSum

	//create only char
	tmpSum = createDomainFileAndWrite(charlist, 4, file)
	sum += tmpSum

	defer file.Close()
	sglog.Info("sum is %d", sum)
}

func createDomainFileAndWrite(srcList []string, num int, file *os.File) (sum int) {
	sglog.Info("now gen %s,num=%d", srcList, num)
	sum = 0
	result := []string{}
	sgalgorithm.GenPermutation(srcList, num, &result)
	zonelist := []string{"com", "cn", "net"} //modify domain
	for _, n := range result {
		for _, k := range zonelist {
			str := n + "." + k + "\n"
			file.Write([]byte(str))
			sum++
		}
	}
	sglog.Info("gen success total num=%d", sum)
	return
}
