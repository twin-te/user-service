import { startGrpcServer, stopGrpcServer } from '../../src/grpc'
import * as protoLoader from '@grpc/proto-loader'
import path from 'path'
import * as grpc from '@grpc/grpc-js'
import { Provider, UserService } from '../../generated'
import { ServiceClientConstructor } from '@grpc/grpc-js/build/src/make-client'
import { GrpcClient } from '../../src/grpc/type'
import { Status } from '@grpc/grpc-js/build/src/constants'
import { getOrCreateUserUseCase } from '../../src/usecase/getOrCreateUser'
import { mocked } from 'ts-jest/utils'
import { v4 } from 'uuid'
import { addAuthenticationUseCase } from '../../src/usecase/addAuthentication'
import { AlreadyExistError, NotFoundError } from '../../src/error'
import { getUserUseCase } from '../../src/usecase/getUser'
import { deepContaining } from '../_deepContaining'

const def = protoLoader.loadSync(
  path.resolve(__dirname, `../../protos/UserService.proto`)
)
const pkg = grpc.loadPackageDefinition(def)
const ClientConstructor = pkg.UserService as ServiceClientConstructor
let client: GrpcClient<UserService>

jest.mock('../../src/usecase/getOrCreateUser')
jest.mock('../../src/usecase/addAuthentication')
jest.mock('../../src/usecase/getUser')

beforeAll(async () => {
  await startGrpcServer()
  client = (new ClientConstructor(
    'localhost:50051',
    grpc.ChannelCredentials.createInsecure()
  ) as unknown) as GrpcClient<UserService>
})

describe('createUser', () => {
  test('success', (done) => {
    const id = v4()
    mocked(getOrCreateUserUseCase).mockImplementation(async () => ({
      id,
    }))
    client.getOrCreateUser({}, (err, res) => {
      expect(err).toBeNull()
      expect(res?.id).toEqual(id)
      done()
    })
  })

  test('unexpected error', (done) => {
    mocked(getOrCreateUserUseCase).mockImplementation(() => {
      throw new Error('Unexpected Error')
    })
    client.getOrCreateUser({}, (err, res) => {
      expect(err?.code).toEqual(Status.UNKNOWN)
      done()
    })
  })
})

describe('addAuthentication', () => {
  const data = {
    id: v4(),
    provider: Provider.Google,
    socialId: '100000000000000000',
  }
  test('success', (done) => {
    mocked(addAuthenticationUseCase).mockImplementation(
      async (id, provider, socialId) => {
        expect(id).toEqual(data.id)
        expect(provider).toEqual('Google') // 後で綺麗にする
        expect(socialId).toEqual(data.socialId)
      }
    )
    client.addAuthentication(data, (err, res) => {
      expect(err).toBeNull()
      expect(res).toBeTruthy()
      done()
    })
  })

  test('already exist', (done) => {
    mocked(addAuthenticationUseCase).mockImplementation(() => {
      throw new AlreadyExistError()
    })
    client.addAuthentication(data, (err, res) => {
      expect(err?.code).toBe(Status.ALREADY_EXISTS)
      done()
    })
  })

  test('not found', (done) => {
    mocked(addAuthenticationUseCase).mockImplementation(() => {
      throw new NotFoundError()
    })
    client.addAuthentication(data, (err, res) => {
      expect(err?.code).toBe(Status.NOT_FOUND)
      done()
    })
  })
})

describe('getUser', () => {
  const data = {
    id: v4(),
    authentications: [
      { provider: Provider.Google, socialId: '100000000000000000' },
    ],
  }
  test('success', (done) => {
    // 後で綺麗にする
    // @ts-ignore
    mocked(getUserUseCase).mockImplementation(async (id) => {
      expect(id).toEqual(data.id)
      return {
        id,
        authentications: [
          { provider: 'Google', socialId: '100000000000000000' },
        ],
      }
    })
    client.getUser(data, (err, res) => {
      expect(err).toBeNull()
      expect(res).toEqual(deepContaining(data))
      done()
    })
  })

  test('not found', (done) => {
    mocked(getUserUseCase).mockImplementation(() => {
      throw new NotFoundError()
    })
    client.getUser(data, (err, res) => {
      expect(err?.code).toBe(Status.NOT_FOUND)
      done()
    })
  })
})

afterAll(stopGrpcServer)
