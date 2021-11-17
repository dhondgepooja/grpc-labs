# grpc-labs

## Types of APIs

![types of APIs](/diagrams/gRPC.png)

## Error Handling:

### References:
https://grpc.io/docs/guides/error/

https://avi.im/grpc-errors/

### Client side:
```
res,err := client.MyFunc(context,req)

if err != nil {
    respErr,ok := status.FromError(err)
    if ok {
        //actual error from gRPC (user error)
    }else{
        // something went wrong while calling myFunc
    }
}
```

## Deadlines:
- Allow client to specify how long they are willing to wait for an RPC to complete before RPC is terminated with DEADLINE_EXCEEDED error 
- gRPC recommends setting a deadline for all client RPC calls
- Server should check if deadline exceeded and cancel the remaining work

### References:
https://grpc.io/blog/deadlines/

## Reflections:

Commands:
```
evans -p 50051 -r

//evans commands
show package
show service
package default
service CalculatorService
call Sum
call ComputeAverage
// to exit from streaming client ctrl + D
```
