# Concurrency Patterns in Go
# Concurrency is not parallelism

Parellelism is about doing lot of things at once
Concurrency is about dealing with lot of things at once

Notes taken from Rob Pike's presentation '[Google I/O 2012 - Go Concurrency Patterns](https://youtu.be/f6kdp27TYZs)'.

## 1. Generator: function that returns a channel
The code is as follows:

``` Go
c := boring("boring!")  // function returning a channel

for i := 0; i < 5; i++ {
    fmt.Printf("You say: %q\n", <-c)
}

fmt.Println("You're boring; I'm leaving.")

func boring(msg string) <-chan string {  // returns receive-only channel of strings
    c := make(chan string)

    go func() {  // we launch the goroutine from inside the function
        for i := 0; ; i++ {
            c <- fmt.Sprintf("%s, %d", msg, i)
            time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
        }
    }()

    return c  // return the channel to the caller
}
```

The generator function act as a service, returning a channel that lets us communicate with the underlying service provided. We can even create multiple instances/routines, each operating concurrently and returning values independently like so.

``` Go
func main() {
    joe := boring("Joe")
    ann := boring("Ann")

    for i := 0; i < 5; i++ {
        fmt.Println(<-joe)
        fmt.Println(<-ann)
    }

    fmt.Println("You're both boring; I'm leaving.")
}
```

The above approach blocks both Joe and Ann until the other has returned and released the routine. The answer to that: multiplexing with a fan-in function.

## 2. Multiplexing: fan-in function to stitch two channels together
We've already established our `func boring(string) <-chan string`, now we multiplex two channels together with a fan-in.

``` Go
func fanIn(input1, input2 <-chan string) <-chan string {
    c := make(chan string)

    go func() { for { c <- <-input1 } }()  // todo: why do we double the channel operator?
    go func() { for { c <- <-input2 } }()

    return c
}

func main() {
    c := fanIn(boring("Joe"), boring("Ann"))

    for i := 0; i < 10; i++ {
        fmt.Println(<-c)
    }

    fmt.Println("You're both boring; I'm leaving")
}
```

Now, the `boring` services for Ann and Joe aren't tied together and we retrieve data as it's available -- the services can independently execute.

## 3. Restore Sequencing: construct message types with blocking fields
Sometimes we need to re-sync concurrent, uncoupled goroutines. To achieve this, we can pass a `wait` channel on another channel, shown below with the `type Message struct` that contains a message string and a boolean channel for blocking. So the goroutine (`boring`) creates a channel within the channel that enables us to tell it to wait its turn. Receive all messages, then enable them again.

``` Go
func boring(msg string) <-chan Message {
    c := make(chan Message)
    waitForIt := make(chan bool)

    go func() {
        for i := 0; ; i++ {
            c <- Message{ fmt.Sprintf("%s, %d", msg, i), waitForIt}
            time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
            <-waitForIt  // boring goroutines now block here until we set the Message#wait boolean field
        }
    }()

    return c  // return the channel to the caller
}

func fanIn(input1, input2 <-chan Message) <-chan Message {
    c := make(chan Message)

    go func() { for { c <- <-input1 } }()
    go func() { for { c <- <-input2 } }()

    return c
}

type Message struct {
    str string
    wait chan bool
}

func main() {
    c := fanIn(boring("Joe"), boring("Ann"))

    for i := 0; i < 5; i++ {
        msg1 := <- c; fmt.Println(msg1.str)
        msg2 := <- c; fmt.Println(msg2.str)

        msg1.wait <- true  // Setting this unblocks `boring` to now iterate through the `for` loop again
        msg2.wait <- true
    }

    fmt.Println("You're both boring; I'm leaving")
}
```

## 4. The Select Statement: switch-like control structure for concurrency that enables us to proceed based on what communication is available at any given moment
- Like a switch, but each case is a communication -- a channel operation (send/receive).
- Select blocks until a communication can proceed, then does. *Note: pseudo-random when multiple can proceed*.
- The default proceeds if non other can.

``` Go
select {
    case v1 := <- c1:
        fmt.Printf("Received %v from c1\n", v1)
    case v2 := <- c2:
        fmt.Printf("Received %v from c2\n", v1)
    case c3 <- 23:
        fmt.Printf("Sent %v from c3\n", 23)
    default:
        fmt.Printf("No one ready")
}
```

## 5. Fan-in Select: reduce the number of goroutines from the prior example
Refactoring our fan-in from #2, we can take:

``` Go
func fanIn(input1, input2 <-chan string) <-chan string {
    c := make(chan string)

    go func() { for { c <- <-input1 } }()  // todo: why do we double the channel operator?
    go func() { for { c <- <-input2 } }()

    return c
}
```

And make:

``` Go
func fanIn(input1, input2 <-chan string) <-chan string {
    c := make(chan string)

    go func() { 
        for {
            select {
                case s:= <- input1: c <- s
                case s:= <- input2: c <- s                
            }
        }        
    }()

    return c
}
```

Now, rather than create two separate goroutined, we create a single goroutine and use the select to multiplex the received channels into the single returned channel.

##  6. Timeout Select
We can use the `time.After(...)` function, which returns a channel, to terminate iterations of the for-select loop.

``` Go
func main() {
    c := boring("Joe")

    for {
        select {
            case s := <- c:
                fmt.Println(s)

            case <- time.After(1 * time.Second):  // We kill the loop if we don't receive a message from the boring("Joe") within 1sec per iteration
                fmt.Println("You're too slow")
                return                
        }
    }
}
```

We can perform a similar operation but with a total of 5sec time elapsed over the entire loop, rather than the 1sec per message in the prior example.

``` Go
func main() {
    c := boring("Joe")
    timeout := time.After(5 * time.Second):  // We kill the loop after 5sec total

    for {
        select {
            case s := <- c:
                fmt.Println(s)

            case <- timeout
                fmt.Println("You talk too much")
                return                
        }
    }
}
```

## 7. Quit Channel: passed to the goroutine
We can pass a quit channel to out routine which tells the routine to finish.

``` Go
func boring(msg string, quit chan bool) <-chan String {
    c := make(chan string)

    go func() {  // we launch the goroutine from inside the function
        for i := 0; ; i++ {
            select {
                case c <- fmt.Sprintf("%s, %d", msg, i):
                    // Do nothing
                case <- quite:
                    return  // Parent routine tells us to finish, so we return from the goroutine
            }
        }
    }()

    return c  // return the channel to the caller
}

func main() {
    quit := make(chan bool)
    c := boring("Joe", quit)
    for i := rand.Intn(10); i >= 0; i-- { fmt.Println(<-c) }
    quite <- true  // Tell the routine to finish
}
```

## 8. Quite Channel: two-way 'wrap it up' communication
The example above ends the goroutine rather abruptly the second we pass a value to the `quit` channel. What it we wanted to allow the routine to wrap things up? Well, the channel offers two-way communication, so we can wait for confirmation from the routine like so.

``` Go
func boring(msg string, quit chan string) <-chan String {
    c := make(chan string)

    go func() {  // we launch the goroutine from inside the function
        for i := 0; ; i++ {
            select {
                case c <- fmt.Sprintf("%s, %d", msg, i):
                    // Do nothing
                case <- quite:
                    cleanup()
                    quit <- "See you!"
                    return
            }
        }
    }()

    return c  // return the channel to the caller
}

func main() {
    quit := make(chan string)
    c := boring("Joe", quit)
    for i := rand.Intn(10); i >= 0; i-- { fmt.Println(<-c) }
    quite <- true  // Tell the routine to finish
    fmt.Printf("Joe says: %q\n", <-quit) // This blocks main until the goroutine confirms it's done
}
```

## 9. Google Search: A fake framework
Lets simulate a search function.

The below runs sequentially. There's no element of concurrency present in this iteration. We construct 3 search functions that conform to the `type Search` form.

``` Go
var (
    Web = fakeSearch("web")
    Image = fakeSearch("image")
    Video = fakeSearch("video")
)

type Search func(query string) Result

func fakeSearch(kind string) Search {
    return func(query string) Result {
        time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
        return Result(fmt.Sprintf("%s result for %q\n", kind, query))
    }
}

func Google(query string) (results []Result) {
    results = append(results, Web(query))  // Currently, we block in each of these search queries
    results = append(results, Image(query))
    results = append(results, Video(query))

    return
}

func main() {
    rand.Seed(time.Now().UnixNano())
    start := time.Now()
    results := Google("golang")
    elapsed := time.Since(start)
    fmt.Println(results)
    fmt.Println(elapsed)
}
```

Alternatively, we could refactor the `Google` function to run in a goroutine and run web, image and video searches concurrently, then wait on all results.

``` Go
func Google(query string) (results []Result) {
    c := make(chan Result)

    // Spin each of the searches off into their own goroutine and pipe all results back to the channel created above
    go func() { c <- Web(query)) }()
    go func() { c <- Image(query)) }()
    go func() { c <- Video(query)) }()

    // Pull each result out of the channel as it becomes available and append it to the results splice returned
    for i := 0; i < 3; i++ {
        result := <- c
        results = append(results, result)
    }

    return
}
```

Here we only wait for the slowest search to complete. But what if we want to timeout communication if a server isn't returning quick enough?... Apply timeout and select and refactor the `Google(...)` function again.

``` Go
func Google(query string) (results []Result) {
    c := make(chan Result)

    // Spin each of the searches off into their own goroutine and pipe all results back to the channel created above
    go func() { c <- Web(query)) }()
    go func() { c <- Image(query)) }()
    go func() { c <- Video(query)) }()

    // Now we pull each Result from the channel, but timeout if searches take longer than 80ms
    timeout := time.After(80 * time.Millisecond)
    for i := 0; i < 3; i++ {
        select {
            case result <- c:
                results = append(results, result)
            case <- timeout:
                fmt.Println("Timed out")
                return
        }
    }

    return
}
```

This could be annoying -- what if we keep killing slow servers before results are returned? We can avoid timeout by replicating servers and returning the first response.

``` Go
// Returns the first response received from one of the many replica Search functions received
func First(query string, replicas ...Search) Result {
    c := make(chan Result)
    searchReplica := func(i int) { c <- replicas[i](query) }  // A function that abstract calling a given Search function with the given query, then piping the result into channel c

    // Iterate over all the given replica Search functions and spin off into their own goroutines
    for i := range replicas {
        go searchReplica(i)
    }

    // Return the first result received only
    return <- c
}
```

Lets add all this together to create a new, final `Google` function that combines multiplexing multiple channels, timeouts, and replicated backend searches

``` Go
func Google(query string) (results []Result) {
    c := make(chan Result)

    go func() { c <- First(query, Web1, Web2) }()
    go func() { c <- First(query, Image1, Image2) }()
    go func() { c <- First(query, Video1, Video2) }()

    timeout := time.After(80 * time.Millisecond)
    for i := 0; i < 3; i++ {
        select {
            case result <- c:
                results = append(results, result)
            case <- timeout:
                fmt.Println("Timed out")
                return
        }
    }

    return
}
```