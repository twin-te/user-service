import { getConnection } from 'typeorm'
import { User } from '../database/model/user'
import {
  Provider,
  UserAuthentication,
} from '../database/model/userAuthentications'
import { AlreadyExistError, NotFoundError } from '../error'

export async function addAuthenticationUseCase(
  id: string,
  provider: Provider,
  socialId: string
): Promise<void> {
  const userRepo = getConnection().getRepository(User)
  if (!(await userRepo.findOne({ id })))
    throw new NotFoundError('指定されたユーザーが見つかりません')

  const authRepo = getConnection().getRepository(UserAuthentication)
  if (await authRepo.findOne({ provider, user: { id } }))
    throw new AlreadyExistError('認証情報が既に登録されています')

  await authRepo.save({
    provider,
    socialId,
    user: { id },
  })
}
