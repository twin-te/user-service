import { getConnection } from 'typeorm'
import { v4 } from 'uuid'
import { connectDatabase, disconnectDatabase } from '../../src/database'
import { User } from '../../src/database/model/user'
import { Provider } from '../../src/database/model/userAuthentications'
import { AlreadyExistError, NotFoundError } from '../../src/error'
import { addAuthenticationUseCase } from '../../src/usecase/addAuthentication'
import { clearDB } from '../_cleardb'

const userId = v4()

beforeAll(async () => {
  await connectDatabase()
  await clearDB()
  await getConnection().getRepository(User).save({ id: userId })
})

test('add authentication', () => {
  return expect(
    addAuthenticationUseCase(userId, Provider.Google, '100000000000000000')
  ).resolves.toBeUndefined()
})

test('already exist authentication', () => {
  return expect(
    addAuthenticationUseCase(userId, Provider.Google, '100000000000000000')
  ).rejects.toThrow(AlreadyExistError)
})

test('not exist user', () => {
  return expect(
    addAuthenticationUseCase(v4(), Provider.Google, '100000000000000000')
  ).rejects.toThrow(NotFoundError)
})

afterAll(disconnectDatabase)
