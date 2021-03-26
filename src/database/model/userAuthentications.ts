import { Column, Entity, ManyToOne, PrimaryGeneratedColumn } from 'typeorm'
import { User } from './user'

export enum Provider {
  Google,
  Twitter,
  Apple,
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
