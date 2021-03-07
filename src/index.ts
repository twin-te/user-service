import { connectDatabase } from './database'
import { startGrpcServer } from './grpc'
import { logger } from './logger'

async function main() {
  logger.info('starting...')
  await connectDatabase()
  await startGrpcServer()
}

main()
