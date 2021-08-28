import {IRideQuery} from "../hooks/useRides";

export default function filterBuilder(query: IRideQuery): string {
  let result = new URLSearchParams()

  if (query.type)
    result.set("type", query.type)
  if (query.direction)
    result.set("direction", query.direction)
  if (query.startDate)
    result.set("start_date", query.startDate)
  if (query.endDate)
    result.set("end_date", query.endDate)

  return result.toString()
}