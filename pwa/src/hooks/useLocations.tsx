import useSWR from 'swr'

import ILocation from '../types/location'
import CONSTANTS from '../constants'

export interface ILocationResponse {
    results?: ILocation[]
    error?: string
}

export default function useLocations() {
    const { data, error } = useSWR<ILocationResponse>(CONSTANTS.API().LOCATIONS)
    if (error) console.error(error)
    return { data, error }
}
