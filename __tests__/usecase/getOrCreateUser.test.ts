import { connectDatabase, disconnectDatabase } from '../../src/database'
import { Provider } from '../../src/database/model/userAuthentications'
import { getOrCreateUserUseCase } from '../../src/usecase/getOrCreateUser'
import { clearDB } from '../_cleardb'

beforeAll(async () => {
  await connectDatabase()
  await clearDB()
})

let userId = ''
test('初回作成', async () => {
  const res = await getOrCreateUserUseCase(
    Provider.Google,
    '100000000000000000'
  )
  userId = res.id
  return expect(res).toEqual(
    expect.objectContaining({
      id: expect.any(String),
    })
  )
})

test('get user', () =>
  expect(
    getOrCreateUserUseCase(Provider.Google, '100000000000000000')
  ).resolves.toEqual(
    expect.objectContaining({
      id: userId,
    })
  ))

afterAll(disconnectDatabase)
