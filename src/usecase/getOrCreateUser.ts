import { getConnection, getRepository } from 'typeorm'
import { v4 } from 'uuid'
import { User } from '../database/model/user'
import {
  Provider,
  UserAuthentication,
} from '../database/model/userAuthentications'

export async function getOrCreateUserUseCase(
  provider: Provider,
  socialId: string
): Promise<{ id: string }> {
  const auth = await getRepository(UserAuthentication).findOne({
    where: { provider, socialId },
    relations: ['user'],
  })
  if (auth) return getRepository(User).findOneOrFail({ id: auth.user.id })
  else
    return getRepository(User).save({
      id: v4(),
      authentications: [{ socialId, provider }],
    })
}
