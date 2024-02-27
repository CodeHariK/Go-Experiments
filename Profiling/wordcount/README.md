# Profiling

* [GopherCon 2019: Dave Cheney - Two Go Programs, Three Different Profiling Techniques](https://www.youtube.com/watch?v=nok0aYiGiYA&ab_channel=GopherAcademy)

* [High Performance Go Workshop](https://dave.cheney.net/high-performance-go-workshop/dotgo-paris.html)

* time wc README.md
* time go run wordcount.go README.md 
* go tool pprof -http=:8080 cpu.pprof
* go tool pprof -http=:8080 mem.pprof

## WordCount1
real    0m8.119s
user    0m2.399s
sys     0m5.193s

## WordCount2
real    0m1.129s
user    0m1.131s
sys     0m0.280s

## WordCount3
real    0m0.450s
user    0m0.433s
sys     0m0.259s