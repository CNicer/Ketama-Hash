package main

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"crypto/md5"
	"encoding/binary"
)

type ServerInfo struct{
	Addr string
	Memory int
}

type ContPoint struct{
	Addr string
	Flag uint32
}

func RingGen(ServerList []ServerInfo) []ContPoint{
	Continum := make([]ContPoint, len(ServerList) * 160)
	ServerN := len(ServerList)
	TotalMemory := 0
	for i := 0;i < ServerN;i++{
		TotalMemory += ServerList[i].Memory
	}
	Idx := 0
	var GroupCount int
	for i := 0;i < len(ServerList);i++{
		GroupCount = ServerList[i].Memory * 40 * ServerN / TotalMemory
		for j := 0;j < GroupCount;j++{
			SrcStr := fmt.Sprintf("%s-%s", ServerList[i].Addr, strconv.Itoa(j))
			Digest := md5.Sum([]byte(SrcStr))
			for k := 0;k < 4;k++{
				Flag := binary.BigEndian.Uint32(Digest[k*4:k*4 + 4])
				Continum[Idx].Addr = ServerList[i].Addr
				Continum[Idx].Flag = Flag
				Idx++
			}
		}
	}
	Continum = Continum[:Idx]
	sort.SliceStable(Continum, func(i, j int) bool{
		return Continum[i].Flag <= Continum[j].Flag
	})
	return Continum
}

func Find(Continum []ContPoint, Key string) string{
	ContN := len(Continum)
	Front := ContN - 1
	Rear := 0
	var Mid int
	Digest := md5.Sum([]byte(Key))
	Flag := binary.BigEndian.Uint32(Digest[:4])
	for{
		if Front <= Rear{
			return Continum[Rear].Addr
		}
		Mid = (Front + Rear)/2
		if Continum[Mid].Flag == Flag{
			return Continum[Mid].Addr
		}else if Continum[Mid].Flag > Flag{
			Front = Mid - 1
		}else{
			Rear = Mid + 1
		}
	}
}

func main(){
	ServerList := [6]ServerInfo{
		{"127.0.0.1:9000", 1},
		{"192.168.100.230:5001", 2},
		{"201.130.79.3:3000", 3},
		{"89.50.10.32:4500", 5},
		{"168.0.3.23:9900", 3},
		{"172.30.23.13:99", 2},
	}
	Continum := RingGen(ServerList[:])
	Check := make(map[string]int)
	for i := 0;i < 6;i++{
		Check[ServerList[i].Addr] = 0
	}
	Count := 100000
	for i := 0;i < Count;i++{
		Check[Find(Continum, strconv.Itoa(rand.Intn(1000000)))]++
	}
	for i := 0;i < 6;i++{
		fmt.Println(ServerList[i].Memory, "Wight: ", float64(ServerList[i].Memory) / 16.0, "rate: ", float64(Check[ServerList[i].Addr]) / float64(Count))
	}
}