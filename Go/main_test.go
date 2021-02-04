package main

import "testing"



func Benchmark(b *testing.B) {

for i:=0; i<b.N; i++{
	GruppenMaster := loadData("D:\\go\\src\\marketdata\\MDL-WPTS_20210113.txt")
	GruppenMasterNoTime := loadDataNoTime("D:\\go\\src\\marketdata\\MDL-WP_20210113.txt")
	request := loadDataRequest("D:\\go\\src\\marketdata\\INST01_MDA_20210114_1244.txt")

	wg.Add(2)
	go GetMatchesNoTIme(request,GruppenMasterNoTime)
	go GetMatches(request,GruppenMaster)
	wg.Wait()
}
}