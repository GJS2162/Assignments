# Project Flow Overview

## Step 1: Compile the eBPF Program

- **Command**: `clang -O2 -target bpf -c -o bpf_program.o bpf_program.c`
- **Description**: This command compiles the eBPF program written in C (bpf_program.c) into an object file (bpf_program.o).
- **Reason**: The eBPF program needs to be compiled into an object file before it can be loaded into the kernel.

## Step 2: Run Go Mod Tidy

- **Command**: `go mod tidy`
- **Description**: This command ensures that the Go module's dependencies are in sync with the go.mod file.
- **Reason**: It helps maintain a clean and consistent dependency tree for the project.

## Step 3: Generate Code

- **Command**: `go generate`
- **Description**: This command runs any Go code generation directives specified in the project.
- **Reason**: It generates any necessary code or files required for the project's functionality.

## Step 4: Build the Go Program

- **Command**: `go build .`
- **Description**: This command compiles the main Go program into an executable binary.
- **Reason**: It produces the executable binary that will be run to load the eBPF program into the kernel and configure it.

## Step 5: Run the Go Program

- **Command**: `sudo go run . {port}`
- **Description**: This command runs the main Go program (main.go) with elevated privileges (sudo) and passes the desired port number as a command-line argument.
- **Reason**: The Go program is responsible for loading the compiled eBPF program into the kernel and configuring it to run on the desired network interface. It needs the port number as input to configure the eBPF program to drop packets targeting that port.

## Step 6: Monitor Network Traffic

- **Command**: `sudo tcpdump -i $INTERFACE_NAME$ port {port}`
- **Description**: This command uses tcpdump to monitor network traffic on the specified port.
- **Reason**: It helps in verifying the behavior of the eBPF program by observing incoming packets on the specified port.

## Step 7: Test the Connection

- **Command**: `nc localhost {port}`
- **Description**: This command attempts to establish a TCP connection to the specified port on localhost (127.0.0.1) using netcat (nc).
- **Reason**: It is used to test the behavior of the eBPF program by attempting to connect to the port and observing whether the connection is successful.

## Overall Flow

1. Compile the eBPF program (`bpf_program.c`) into an object file (`bpf_program.o`).
2. Run Go Mod Tidy to ensure dependency consistency.
3. Generate any necessary code or files using `go generate`.
4. Build the main Go program into an executable binary.
5. Run the main Go program with elevated privileges, passing the desired port number as input.
6. Monitor network traffic on the specified port using tcpdump to verify the behavior of the eBPF program.
7. Test the connection to the specified port using netcat to observe the impact of the eBPF program on incoming packets.
