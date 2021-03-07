/**
 * Grpcサービスの型定義を生成する
 */

import * as grpc from '@grpc/grpc-js'
import * as protobuf from 'protobufjs'

export type DeepRequired<T> = {
  [K in keyof T]-?: NonNullable<DeepRequired<T[K]>>
}

type FilteredKeys<T, U> = {
  [P in keyof T]: T[P] extends U ? P : never
}[keyof T]

type ProtobufFn = (request: {}) => PromiseLike<{}>
type RequestArgType<T extends ProtobufFn> = T extends (
  request: infer U
) => PromiseLike<any>
  ? U
  : never
type ResponseType<T extends ProtobufFn> = T extends (
  request: any
) => PromiseLike<infer U>
  ? U
  : never

export type GrpcServer<T extends protobuf.rpc.Service> = {
  [K in FilteredKeys<T, ProtobufFn>]: grpc.handleUnaryCall<
    DeepRequired<RequestArgType<T[K]>>,
    ResponseType<T[K]>
  >
}

export type GrpcClient<T extends protobuf.rpc.Service> = grpc.Client &
  {
    [K in FilteredKeys<T, ProtobufFn>]: (
      req: RequestArgType<T[K]>,
      callback: grpc.requestCallback<ResponseType<T[K]>>
    ) => void
  }

declare module '@grpc/grpc-js' {
  // eslint-disable-next-line no-unused-vars
  class Server {
    addService<T extends grpc.UntypedServiceImplementation>(
      service: grpc.ServiceDefinition,
      implementation: T
    ): void
  }
}
