# hconcurrent

```
//
//                     |-      -|                    |-      -|                    |-      -|
//                     | handle |                    | handle |                    | handle |
//                     | handle |                    | handle |                    | handle |
//                     | handle |                    | handle |                    | handle |
// >>> input chan >>> -| handle |- >>> out chan >>> -| handle |- >>> out chan >>> -| handle |
//    (channelSize)    | handle |    (channelSize)   | handle |    (channelSize)   | handle |
//                     | handle |                    | handle |                    | handle |
//                     |   .    |                    |   .    |                    |   .    |
//                     |   .    |                    |   .    |                    |   .    |
//                     |-  .   -|                    |-  .   -|                    |-  .   -|
//               goroutineCount x handle       goroutineCount x handle       goroutineCount x handle
//
```

import:

```go
go get -v -u github.com/hqpko/hconcurrent
```

example:

```go
	c := hconcurrent.NewConcurrentWithOptions(
		NewOption(channSize, goroutineCount, handle1),
		NewOption(channSize, goroutineCount, handle2),
		NewOption(channSize, goroutineCount, handle3),
	)
	c.Start()

	// ...
	inputSuccess := c.Input(v)
	// c.MustInput(v)
	// inputSuccess := c.InputWithTimeout(v, time.Second)

	// ...

	c.Stop()
```

simple:

```go
	c := hconcurrent.NewConcurrent(16, 1, handle)
	c.Start()
	
	// ...
	inputSuccess := c.Input(v)
	// c.MustInput(v)
	// inputSuccess := c.InputWithTimeout(v, time.Second)

	// ...
	
	c.Stop()
```
