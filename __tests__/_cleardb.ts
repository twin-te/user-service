import { getConnection } from 'typeorm'
import { User } from '../src/database/model/user'
import { UserAuthentication } from '../src/database/model/userAuthentications'

export async function clearDB() {
  await getConnection().getRepository(UserAuthentication).delete({})
  await getConnection().getRepository(User).delete({})
}
