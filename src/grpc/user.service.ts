import {
  AddAuthenticationResponse,
  GetOrCreateUserResponse,
  GetUserResponse,
  Provider as GProvider,
  UserService,
} from '../../generated'
import { Provider } from '../database/model/userAuthentications'
import { addAuthenticationUseCase } from '../usecase/addAuthentication'
import { getOrCreateUserUseCase } from '../usecase/getOrCreateUser'
import { getUserUseCase } from '../usecase/getUser'
import { toGrpcError } from './converter'
import { GrpcServer } from './type'

export const userService: GrpcServer<UserService> = {
  async getOrCreateUser({ request }, callback) {
    try {
      const res = await getOrCreateUserUseCase(
        toDBProvider(request.provider),
        request.socialId
      )
      callback(null, GetOrCreateUserResponse.create({ ...res }))
    } catch (e) {
      callback(toGrpcError(e))
    }
  },
  async addAuthentication({ request }, callback) {
    try {
      await addAuthenticationUseCase(
        request.id,
        toDBProvider(request.provider),
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
          authentications: user.authentications.map((a) => ({
            provider: toGrpcProvider(a.provider),
            socialId: a.socialId,
          })),
        })
      )
    } catch (e) {
      callback(toGrpcError(e))
    }
  },
}

function toDBProvider(provider: GProvider) {
  switch (provider) {
    case GProvider.Google:
      return Provider.Google
    case GProvider.Twitter:
      return Provider.Twitter
    case GProvider.Apple:
      return Provider.Apple
  }
}

function toGrpcProvider(provider: Provider) {
  switch (provider) {
    case Provider.Google:
      return GProvider.Google
    case Provider.Twitter:
      return GProvider.Twitter
    case Provider.Apple:
      return GProvider.Apple
  }
}
