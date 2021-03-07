import { connectDatabase, disconnectDatabase } from '../../src/database'
import { createUserUseCase } from '../../src/usecase/createUser'
import { clearDB } from '../_cleardb'

beforeAll(async () => {
  await connectDatabase()
  await clearDB()
})

test('create user', () => {
  return expect(createUserUseCase()).resolves.toEqual({
    id: expect.any(String),
  })
})

afterAll(disconnectDatabase)
