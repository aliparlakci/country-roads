import ILocation from './location'
import RideType from './rideType'
import { IForeignUser } from './user'

export default interface IRide {
  id: string
  type: RideType
  date: number
  destination: ILocation
  direction: string
  createdAt: number
  owner: IForeignUser
}
