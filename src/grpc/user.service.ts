import { Status } from '@grpc/grpc-js/build/src/constants'
import { GreetResponse, UserService } from '../../generated'
import { GrpcServer } from './type'

export const userService: GrpcServer<UserService> = {
  greet({ request }, callback) {
    if (!request.name)
      callback({ code: Status.INVALID_ARGUMENT, details: '名前が空です' })
    else
      callback(null, GreetResponse.create({ text: `hello! ${request.name}` }))
  },
}
