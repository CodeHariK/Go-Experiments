# Mandelbrot

go tool trace trace.out

# time go run .
real    0m3.951s
user    0m3.523s
sys     0m0.246s

# time go run . -mode px
real    0m2.184s
user    0m6.951s
sys     0m0.775s

# time go run . -mode row
real    0m1.159s
user    0m4.246s
sys     0m0.253s

# time go run . -mode workers -workers=8
real    0m3.910s
user    0m6.718s
sys     0m2.094s