import { getConnection } from 'typeorm'
import { User } from '../database/model/user'
import { Provider } from '../database/model/userAuthentications'
import { NotFoundError } from '../error'

type Result = {
  id: string
  authentications: {
    provider: Provider
    socialId: string
  }[]
}

export async function getUserUseCase(id: string): Promise<Result> {
  const res = await getConnection()
    .getRepository(User)
    .findOne({ id }, { relations: ['authentications'] })
  if (!res) throw new NotFoundError('指定されたユーザーが見つかりませんでした')
  return res
}
