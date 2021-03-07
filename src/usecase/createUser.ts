import { getConnection } from 'typeorm'
import { v4 } from 'uuid'
import { User } from '../database/model/user'

export function createUserUseCase(): Promise<{ id: string }> {
  return getConnection().getRepository(User).save({
    id: v4(),
  })
}
