package main

type ServerInfo struct{
	var Addr byte[22]
	var Memory int
}

type ContPoint struct{
	var Addr byte[22]
	var Flag int
}

func RingGen(ServerList []ServerInfo){
	Continum := make([]ContPoint, len(ServerList) * 160)
	ServerN := len(ServerList)
	TotalMemory := 0
	for i := 0;i < ServerN;i++{
		TotalMemory += ServerList[i].Memory
	}
	Idx := 0
	var GroupCount int
	for i := 0;i < len(ServerList);i++{
		GroupCount = math.Floor(ServerList[i].Memory * 40 * ServerN / TotalMemory)
		for j := 0;j < GroupCount;j++{
			SrcStr := fmt.Sprintf("%s-%s", ServerList[i].Addr, GroupCount)
			Digest := md5(SrcStr)
			for k := 0;k < 4;k++{
				Flag := (Digest[3 + k * 4]) << 24
					   |(Digest[2 + k * 4]) << 16
					   |(Digest[1 + k * 4]) << 8
					   |(Digest[0 + k * 4])
				copy(Continum[Idx].Addr, ServerList[i].Addr)
				Continum[Idx].Flag = Flag
				Idx++
			}
		}
	}
	sort.SliceStable(Continum, func(i, j int){
		return Continum[i].Flag <= Continum[j].Flag
	})
}

func Find(Key []byte){
	ContN := len(Continum)
	Front := ContN
	Rear := 0
	var Mid int
	Digest := md5(Key)
	Digest = (Digest[3] << 24)
			|(Digest[2] << 16)
			|(Digest[1] << 8)
			|(Digest[0])
	for{
		if Front == Rear{
			return Continum[Front].Addr
		}
		Mid = math.Floor((Front + Rear)/2)
		if Continum[Mid].Flag == Key{
			return Continum[Mid].Addr
		}else if Continum[Mid].Flag > Digest{
			Front = Mid - 1
		}else{
			Rear = Mid + 1
		}
	}
}