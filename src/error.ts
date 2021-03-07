class BaseError extends Error {
  constructor(message?: string, error?: Error) {
    super(message)
    Object.defineProperty(this, 'name', {
      configurable: true,
      enumerable: false,
      value: this.constructor.name,
      writable: true,
    })
    Error.captureStackTrace(this, this.constructor)
    if (error?.stack) this.stack = this.mergeStack(this.stack!, error.stack)
  }

  private mergeStack(newStack: string, oldStack: string) {
    return (
      newStack
        .split('\n')
        .filter((line) => !oldStack.includes(line))
        .join('\n') +
      '\n' +
      oldStack
    )
  }
}

export class NotFoundError extends BaseError {
  constructor(message?: string, error?: Error, readonly resources?: string[]) {
    super(message, error)
  }
}
export class InvalidArgumentError extends BaseError {
  constructor(
    message?: string,
    error?: Error,
    readonly args?: { name: string; expect: string }[]
  ) {
    super(message, error)
  }
}
export class AlreadyExistError extends BaseError {}
