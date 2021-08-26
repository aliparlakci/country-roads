import { cache, mutate } from "swr";

export default function mutateWithQueries(base: string) {
  const pattern = new RegExp(`^(${base})(\\?.*|)$`);
  cache
    .keys()
    .filter((key) => pattern.test(key))
    .forEach((key) => mutate(key));
}
