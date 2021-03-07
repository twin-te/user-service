import {
  AddAuthenticationResponse,
  CreateUserResponse,
  GetUserResponse,
  UserService,
} from '../../generated'
import { addAuthenticationUseCase } from '../usecase/addAuthentication'
import { createUserUseCase } from '../usecase/createUser'
import { getUserUseCase } from '../usecase/getUser'
import { toGrpcError } from './converter'
import { GrpcServer } from './type'

export const userService: GrpcServer<UserService> = {
  async createUser(_, callback) {
    try {
      const res = await createUserUseCase()
      callback(null, CreateUserResponse.create({ ...res }))
    } catch (e) {
      callback(toGrpcError(e))
    }
  },
  async addAuthentication({ request }, callback) {
    try {
      await addAuthenticationUseCase(
        request.id,
        request.provider,
        request.socialId
      )
      callback(null, AddAuthenticationResponse.create())
    } catch (e) {
      callback(toGrpcError(e))
    }
  },
  async getUser({ request }, callback) {
    try {
      const user = await getUserUseCase(request.id)
      callback(
        null,
        GetUserResponse.create({
          id: request.id,
          authentications: user.authentications,
        })
      )
    } catch (e) {
      callback(toGrpcError(e))
    }
  },
}
