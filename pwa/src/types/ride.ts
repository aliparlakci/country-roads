import ILocation from './location'
import RideType from './rideType'

export default interface IRide {
    id: string
    type: RideType
    date: number
    destination: ILocation
    direction: string
    createdAt: number
}
