import { cache, mutate } from 'swr'

export default async function mutateWithQueries(base: string) {
  const pattern = new RegExp(`^(${base})(\\?.*|)$`)
  const links = cache.keys().filter((key) => pattern.test(key))
  for (const link of links) {
    await mutate(link)
  }
}
