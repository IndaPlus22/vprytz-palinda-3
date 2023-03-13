# Written answers for the questions in the README.md file

## Task 1

Take a look at the program [matching.go](src/matching.go). Explain what happens and why it happens if you make the following changes. Try first to reason about it, and then test your hypothesis by changing and running the program.

- What happens if you remove the `go-command` from the `Seek` call in the `main` function?
  - **Answer**: The program will not run in parallel, but sequentially. This means that `Anna` will send a message to `Bob`, `Cody` to `Dave` but `Eva`'s message will never be recieved by anyone. Since `Eva` is the last person to send a message in an odd number of people, the messages send by `Eva` will never be recieved by anyone.
- What happens if you switch the declaration `wg := new(sync.WaitGroup)` to `var wg sync.WaitGroup` and the parameter `wg *sync.WaitGroup` to `wg sync.WaitGroup`?
  - **Answer**: If we had just replaced `wg := new(sync.WaitGroup)` with `var wg sync.WaitGroup` and added `&` before passing `wg` to the function, nothing would have changed since we're still passing the pointer to our function `Seek`. However, since we're also meant to replace the parameter, to indicate that the function should no longer take a pointer/reference as a parameter, a copy of the `WaitGroup` will be passed to the function. This means that the `WaitGroup` will not be shared between the `main` function and the `Seek` function. This means that the `WaitGroup` will not be able to keep track of the number of goroutines that are still running, and the program will not wait for the goroutines to finish before exiting.
- What happens if you remove the buffer on the channel match?
  - **Answer**: If you remove the buffer on the channel then the channel becomes unbuffered, meaning that the channel will be blocked in `Seek` until a value is read. This will cause a deadlock since the `Seek` function will be blocked on reading from the channel when the other goroutines already have sent a message to the channel. This means that the program will never exit.
- What happens if you remove the default-case from the case-statement in the `main` function?
  - **Answer**: In theory, nothing would happen, since we have an odd number of people the default case will never be executed. However, if we had an even number of people, the case-statement will keep trying to read from the channel even when all the goroutines have finished. This means that the program will never exit, causing a deadlock.

Hint: Think about the order of the instructions and what happens with arrays of different lengths.
