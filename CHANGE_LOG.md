# Change Log

## [v1.2](https://github.com/thewizardplusplus/go-sync-utils/tree/v1.2) (2021-11-08)

- wrapper for an abstract handler allowed to call it concurrently:
  - starting:
    - start in the caller goroutine;
    - start in a goroutine pool;
  - stopping:
    - can be called after both kinds of the starting;
    - blocks the execution flow until the stopping will be completed.

## [v1.1](https://github.com/thewizardplusplus/go-sync-utils/tree/v1.1) (2020-09-26)

- sending to a channel without blocking even if the channel is busy;
- adding the package comment.

## [v1.0](https://github.com/thewizardplusplus/go-sync-utils/tree/v1.0) (2020-09-16)
