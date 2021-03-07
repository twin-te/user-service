import { getConnection } from 'typeorm'
import { v4 } from 'uuid'
import { connectDatabase, disconnectDatabase } from '../../src/database'
import { User } from '../../src/database/model/user'
import { Provider } from '../../src/database/model/userAuthentications'
import { NotFoundError } from '../../src/error'
import { getUserUseCase } from '../../src/usecase/getUser'
import { clearDB } from '../_cleardb'

const userId = v4()

const data = {
  id: userId,
  authentications: [
    { provider: Provider.Google, socialId: '100000000000000000' },
  ],
}

beforeAll(async () => {
  await connectDatabase()
  await clearDB()
  await getConnection().getRepository(User).save(data)
})

test('add authentication', () => {
  return expect(getUserUseCase(userId)).resolves.toEqual(data)
})

test('not exist user', () => {
  return expect(getUserUseCase(v4())).rejects.toThrow(NotFoundError)
})

afterAll(disconnectDatabase)
