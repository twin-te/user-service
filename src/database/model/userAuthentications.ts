import {
  Column,
  Entity,
  JoinColumn,
  ManyToOne,
  PrimaryGeneratedColumn,
} from 'typeorm'
import { User } from './user'

export enum Provider {
  Google = 'Google',
  Twitter = 'Twitter',
  Apple = 'Apple',
}

@Entity({
  name: 'user_authentications',
})
export class UserAuthentication {
  @PrimaryGeneratedColumn({
    name: 'id',
    type: 'integer',
  })
  id!: number

  @ManyToOne((type) => User, (user) => user.id)
  @JoinColumn({ name: 'user_id' })
  user!: User

  @Column({
    name: 'provider',
    type: 'enum',
    enum: Provider,
  })
  provider!: Provider

  @Column({
    name: 'social_id',
  })
  socialId!: string
}
