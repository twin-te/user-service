import { ServerUnaryCall, sendUnaryData, Metadata } from '@grpc/grpc-js'
import {
  ServerErrorResponse,
  ServerStatusResponse,
} from '@grpc/grpc-js/build/src/server-call'
import { grpcLogger as logger } from '../logger'
import { GrpcServer } from './type'

/**
 * grpc通信のログを記録する
 * @param i 任意のService実装
 */
export function applyLogger<T extends GrpcServer<any>>(i: T): T {
  const impl = i
  Object.getOwnPropertyNames(impl)
    .filter((k) => typeof impl[k] === 'function')
    .forEach((k) => {
      const originalImpl = impl[k]

      // ログ出力を追加するラップ関数
      // @ts-ignore
      impl[k] = function (
        call: ServerUnaryCall<any, any>,
        callback: sendUnaryData<any>
      ) {
        if (logger.isTraceEnabled())
          logger.trace('REQUEST', originalImpl.name, call.request)
        else logger.info('REQUEST', originalImpl.name)

        const originalCallback = callback
        callback = function (
          error: ServerErrorResponse | ServerStatusResponse | null,
          value?: any | null,
          trailer?: Metadata,
          flags?: number
        ) {
          if (error) logger.error('RESPONSE', originalImpl.name, error)
          else if (logger.isTraceEnabled())
            logger.trace('RESPONSE', originalImpl.name, {
              error,
              value: JSON.stringify(value),
              trailer,
              flags,
            })
          else logger.info('RESPONSE', originalImpl.name, 'ok')

          originalCallback(error, value, trailer, flags)
        }
        // @ts-ignore
        originalImpl(call, callback)
      }
    })
  // @ts-ignore
  return impl
}
