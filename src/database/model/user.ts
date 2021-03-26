import { Entity, OneToMany, PrimaryColumn } from 'typeorm'
import { UserAuthentication } from './userAuthentications'

@Entity({
  name: 'users',
})
export class User {
  @PrimaryColumn({
    name: 'id',
    type: 'uuid',
  })
  id!: string

  @OneToMany((type) => UserAuthentication, (auth) => auth.user, {
    cascade: true,
  })
  authentications!: UserAuthentication[]
}
