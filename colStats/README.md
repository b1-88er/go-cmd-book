```
go test -bench . -benchtime=10x -run ^$ -cpuprofile cpu01.prof
go test -bench . -benchtime=10x -run ^$ -benchmem | tee benchresult01m.txt
go test -bench . -benchtime=10x -run ^$ -trace trace02.out
```