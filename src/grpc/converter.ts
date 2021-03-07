import { StatusObject, Metadata } from '@grpc/grpc-js'
import { Status } from '@grpc/grpc-js/build/src/constants'
import { ServerErrorResponse } from '@grpc/grpc-js/build/src/server-call'
import {
  NotFoundError,
  InvalidArgumentError,
  AlreadyExistError,
} from '../error'

export function toGrpcError(
  e: Error
): Partial<StatusObject> | ServerErrorResponse {
  if (e instanceof NotFoundError)
    return Object.assign(e, {
      code: Status.NOT_FOUND,
      metadata: makeMetadata({ resources: e.resources }),
    })
  else if (e instanceof InvalidArgumentError)
    return Object.assign(e, {
      code: Status.INVALID_ARGUMENT,
      metadata: makeMetadata({ args: e.args }),
    })
  else if (e instanceof AlreadyExistError)
    return Object.assign(e, {
      code: Status.ALREADY_EXISTS,
    })
  else return Object.assign(e, { code: Status.UNKNOWN })
}

function makeMetadata(obj: any): Metadata {
  const metadata = new Metadata()
  Object.keys(obj).forEach((k) => metadata.add(k, obj[k]))
  return metadata
}
